package logging

import (
	model "core-ledger/model/core-ledger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/jobs"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	dispatcher queue.Dispatcher
}

func NewService(db *gorm.DB, dispatcher queue.Dispatcher) *Service {
	return &Service{
		db:         db,
		dispatcher: dispatcher,
	}
}

// -------------------- MANUAL CREATE LOG --------------------
func (s *Service) CreateLog(obj Loggable, action string, oldVal, newVal any, createdBy *uint64, metadata any) error {
	oldJSON, _ := json.Marshal(oldVal)
	newJSON, _ := json.Marshal(newVal)
	metaJSON, _ := json.Marshal(metadata)

	log := model.Log{
		LoggableID:   obj.GetLoggableID(),
		LoggableType: obj.GetLoggableType(),
		Action:       action,
		OldValue:     datatypes.JSON(oldJSON),
		NewValue:     datatypes.JSON(newJSON),
		Metadata:     datatypes.JSON(metaJSON),
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
	}

	return s.db.Create(&log).Error
}

// -------------------- REGISTER CALLBACKS --------------------
func RegisterCallbacks(db *gorm.DB, dispatcher queue.Dispatcher) {
	// Store dispatcher in a way that callbacks can access it
	// We'll use a global variable or pass it through context
	// For now, we'll use a package-level variable
	globalDispatcher = dispatcher
	
	db.Callback().Create().After("gorm:after_create").Register("auto_log_create", autoLogCreate)
	db.Callback().Update().Before("gorm:before_update").Register("auto_log_before_update", autoLogBeforeUpdate)
	db.Callback().Update().After("gorm:after_update").Register("auto_log_after_update", autoLogAfterUpdate)
	db.Callback().Delete().After("gorm:after_delete").Register("auto_log_delete", autoLogDelete)
}

// Global dispatcher for callbacks (since callbacks don't have access to service instance)
var globalDispatcher queue.Dispatcher

// -------------------- CREATE CALLBACK --------------------
func autoLogCreate(db *gorm.DB) {
	fmt.Printf("🔔 [DEBUG] autoLogCreate called, dest type: %T\n", db.Statement.Dest)
	
	// Track processed IDs to avoid duplicates in batch operations
	processed := make(map[string]bool)
	
	iterateLoggable(db.Statement.Dest, func(obj Loggable) {
		fmt.Printf("🔔 [DEBUG] Processing loggable: %s, ID: %d\n", obj.GetLoggableType(), obj.GetLoggableID())
		if !shouldLog(obj) {
			fmt.Printf("⚠️  [DEBUG] Skipping log for %s (shouldLog returned false)\n", obj.GetLoggableType())
			return
		}

		id := obj.GetLoggableID()
		if id == 0 {
			fmt.Printf("⚠️  [DEBUG] Skipping log for %s (ID is 0)\n", obj.GetLoggableType())
			return
		}

		// Create unique key to track duplicates
		key := fmt.Sprintf("%s:%d", obj.GetLoggableType(), id)
		if processed[key] {
			fmt.Printf("⚠️  [DEBUG] Skipping duplicate log for %s (already processed)\n", key)
			return
		}
		processed[key] = true

		newJSON, _ := json.Marshal(obj)
		log := model.Log{
			LoggableID:   id,
			LoggableType: obj.GetLoggableType(),
			Action:       "CREATED",
			NewValue:     datatypes.JSON(newJSON),
			CreatedAt:    time.Now(),
		}

		// Dispatch log job to queue instead of inserting directly
		if err := dispatchLogJob(&log); err != nil {
			fmt.Printf("❌ autoLogCreate: failed to dispatch log job for %s:%d: %v\n", obj.GetLoggableType(), id, err)
		} else {
			fmt.Printf("✅ autoLogCreate: successfully dispatched log job for %s:%d\n", obj.GetLoggableType(), id)
		}
	})
}

// -------------------- BEFORE UPDATE CALLBACK --------------------
func autoLogBeforeUpdate(db *gorm.DB) {
	fmt.Printf("🔔 [DEBUG] autoLogBeforeUpdate called, dest type: %T\n", db.Statement.Dest)
	oldMap := map[uint64]json.RawMessage{}

	// Handle Model().Updates() case - get ID from Model
	if db.Statement.Model != nil {
		if modelObj, ok := db.Statement.Model.(Loggable); ok {
			id := modelObj.GetLoggableID()
			if id != 0 && shouldLog(modelObj) {
				fmt.Printf("🔔 [DEBUG] Found loggable model in Statement.Model: %s, ID: %d\n", modelObj.GetLoggableType(), id)
				oldStruct := cloneEmpty(modelObj)
				if oldStruct != nil {
					// Use Statement.DB to get the base DB instance and create a fresh query
					baseDB := db.Statement.DB
					if baseDB != nil {
						if err := baseDB.WithContext(context.Background()).First(oldStruct, id).Error; err == nil {
							oldJSON, _ := json.Marshal(oldStruct)
							oldMap[id] = oldJSON
							fmt.Printf("🔔 [DEBUG] Saved old value for %s:%d\n", modelObj.GetLoggableType(), id)
						} else {
							fmt.Printf("⚠️  [DEBUG] Failed to get old value for %s:%d: %v\n", modelObj.GetLoggableType(), id, err)
						}
					}
				}
			}
		}
	}

	// Struct update
	iterateLoggable(db.Statement.Dest, func(obj Loggable) {
		if !shouldLog(obj) {
			return
		}
		id := obj.GetLoggableID()
		if id == 0 {
			return
		}

		oldStruct := cloneEmpty(obj)
		if oldStruct == nil {
			return
		}

		// Use Statement.DB to get the base DB instance and create a fresh query
		baseDB := db.Statement.DB
		if baseDB != nil {
			if err := baseDB.WithContext(context.Background()).First(oldStruct, id).Error; err != nil {
				fmt.Printf("⚠️  [DEBUG] Failed to get old value for %s:%d (from Dest): %v\n", obj.GetLoggableType(), id, err)
				return
			}

			oldJSON, _ := json.Marshal(oldStruct)
			oldMap[id] = oldJSON
			fmt.Printf("🔔 [DEBUG] Saved old value for %s:%d (from Dest)\n", obj.GetLoggableType(), id)
		}
	})

	// Map update
	if m, ok := db.Statement.Dest.(map[string]interface{}); ok {
		if idVal, exists := m["id"]; exists {
			if id := convertToUint64(idVal); id != 0 {
				var old model.CoaAccount
				// Use Statement.DB to get the base DB instance and create a fresh query
				baseDB := db.Statement.DB
				if baseDB != nil {
					if err := baseDB.WithContext(context.Background()).First(&old, id).Error; err == nil {
						oldJSON, _ := json.Marshal(old)
						oldMap[id] = oldJSON
						fmt.Printf("🔔 [DEBUG] Saved old value for coa_accounts:%d (from map)\n", id)
					} else {
						fmt.Printf("⚠️  [DEBUG] Failed to get old value for coa_accounts:%d (from map): %v\n", id, err)
					}
				}
			}
		}
	}

	db.InstanceSet("old_values", oldMap)
	
	// Also save update fields if available (for Updates() method)
	if db.Statement.Dest != nil {
		if fieldsMap, ok := db.Statement.Dest.(map[string]interface{}); ok {
			db.InstanceSet("update_fields", fieldsMap)
			fmt.Printf("🔔 [DEBUG] Saved update fields: %v\n", fieldsMap)
		}
	}
	
	fmt.Printf("🔔 [DEBUG] Saved %d old values\n", len(oldMap))
}

// -------------------- AFTER UPDATE CALLBACK --------------------
func autoLogAfterUpdate(db *gorm.DB) {
	fmt.Printf("🔔 [DEBUG] autoLogAfterUpdate called, dest type: %T, model type: %T\n", db.Statement.Dest, db.Statement.Model)
	oldMap, _ := db.InstanceGet("old_values")
	oldMapTyped, _ := oldMap.(map[uint64]json.RawMessage)

	// Handle Model().Updates() case - merge old value with changes
	if db.Statement.Model != nil {
		if modelObj, ok := db.Statement.Model.(Loggable); ok {
			id := modelObj.GetLoggableID()
			if id != 0 && shouldLog(modelObj) {
				fmt.Printf("🔔 [DEBUG] Processing loggable from Statement.Model: %s, ID: %d\n", modelObj.GetLoggableType(), id)
				
				var oldJSON json.RawMessage
				if oldMapTyped != nil {
					oldJSON = oldMapTyped[id]
					delete(oldMapTyped, id)
				}

				// Get new value by merging old value with changes
				// Note: With Updates(), we may not have access to the fields map directly
				// So we'll use old value as base and try to merge with any available changes
				var newJSON json.RawMessage
				if oldJSON != nil {
					// Parse old value
					var oldMap map[string]interface{}
					if err := json.Unmarshal(oldJSON, &oldMap); err == nil {
						// Try to get update fields from InstanceSet (saved in before_update)
						if updateFields, exists := db.InstanceGet("update_fields"); exists {
							if fieldsMap, ok := updateFields.(map[string]interface{}); ok {
								// Merge update fields into old map
								for key, value := range fieldsMap {
									oldMap[key] = value
								}
							}
						}
						// Note: With Updates(), we can't easily get the fields map from statement
						// So we rely on update_fields saved in before_update
						newJSON, _ = json.Marshal(oldMap)
					}
				}

				// If we couldn't merge, use old value as new value (not ideal but prevents errors)
				// The important thing is we have the old value to track what changed
				if newJSON == nil {
					newJSON = oldJSON
				}

				meta := getMetadataFromContext(db)

				log := model.Log{
					LoggableID:   id,
					LoggableType: modelObj.GetLoggableType(),
					Action:       "UPDATED",
					OldValue:     datatypes.JSON(oldJSON),
					NewValue:     datatypes.JSON(newJSON),
					Metadata:     datatypes.JSON(meta),
					CreatedAt:    time.Now(),
				}

				// Dispatch log job to queue instead of inserting directly
				if err := dispatchLogJob(&log); err != nil {
					fmt.Printf("❌ autoLogAfterUpdate: failed to dispatch log job for %s:%d: %v\n", modelObj.GetLoggableType(), id, err)
				} else {
					fmt.Printf("✅ autoLogAfterUpdate: successfully dispatched log job for %s:%d\n", modelObj.GetLoggableType(), id)
				}
			}
		}
	}

	iterateLoggable(db.Statement.Dest, func(obj Loggable) {
		if !shouldLog(obj) {
			return
		}
		id := obj.GetLoggableID()
		if id == 0 {
			return
		}

		var oldJSON json.RawMessage
		if oldMapTyped != nil {
			oldJSON = oldMapTyped[id]
			delete(oldMapTyped, id)
		}

		// Get new value - try to merge old value with update fields
		var newJSON json.RawMessage
		if oldJSON != nil {
			// Try to merge old value with changes
			var oldMap map[string]interface{}
			if err := json.Unmarshal(oldJSON, &oldMap); err == nil {
				// Try to get update fields from InstanceSet (saved in before_update)
				if updateFields, exists := db.InstanceGet("update_fields"); exists {
					if fieldsMap, ok := updateFields.(map[string]interface{}); ok {
						// Merge update fields into old map
						for key, value := range fieldsMap {
							oldMap[key] = value
						}
					}
				}
				newJSON, _ = json.Marshal(oldMap)
			}
		}

		// If we couldn't merge, use the object as-is or old value
		if newJSON == nil {
			newJSON, _ = json.Marshal(obj)
			// If object doesn't have updated values, use old value
			if newJSON == nil || len(newJSON) == 0 {
				newJSON = oldJSON
			}
		}

		meta := getMetadataFromContext(db)

		log := model.Log{
			LoggableID:   id,
			LoggableType: obj.GetLoggableType(),
			Action:       "UPDATED",
			OldValue:     datatypes.JSON(oldJSON),
			NewValue:     datatypes.JSON(newJSON),
			Metadata:     datatypes.JSON(meta),
			CreatedAt:    time.Now(),
		}

		// Dispatch log job to queue instead of inserting directly
		if err := dispatchLogJob(&log); err != nil {
			fmt.Printf("❌ autoLogAfterUpdate: failed to dispatch log job for %s:%d: %v\n", obj.GetLoggableType(), id, err)
		} else {
			fmt.Printf("✅ autoLogAfterUpdate: successfully dispatched log job for %s:%d\n", obj.GetLoggableType(), id)
		}
	})

	// Map update case: log remaining old_values
	if oldMapTyped != nil {
		for id, oldJSON := range oldMapTyped {
			var newObj model.CoaAccount
			if err := db.First(&newObj, id).Error; err != nil {
				continue
			}
			newJSON, _ := json.Marshal(newObj)
			meta := getMetadataFromContext(db)
			log := model.Log{
				LoggableID:   id,
				LoggableType: "coa_accounts",
				Action:       "UPDATED",
				OldValue:     datatypes.JSON(oldJSON),
				NewValue:     datatypes.JSON(newJSON),
				Metadata:     datatypes.JSON(meta),
				CreatedAt:    time.Now(),
			}
			// Dispatch log job to queue instead of inserting directly
			if err := dispatchLogJob(&log); err != nil {
				fmt.Printf("❌ autoLogAfterUpdate (map): failed to dispatch log job for coa_accounts:%d: %v\n", id, err)
			}
		}
		db.InstanceSet("old_values", map[uint64]json.RawMessage{})
	}
}

// -------------------- DELETE CALLBACK --------------------
func autoLogDelete(db *gorm.DB) {
	iterateLoggable(db.Statement.Dest, func(obj Loggable) {
		if !shouldLog(obj) {
			return
		}
		id := obj.GetLoggableID()
		if id == 0 {
			return
		}

		oldJSON, _ := json.Marshal(obj)
		log := model.Log{
			LoggableID:   id,
			LoggableType: obj.GetLoggableType(),
			Action:       "DELETED",
			OldValue:     datatypes.JSON(oldJSON),
			CreatedAt:    time.Now(),
		}

		// Dispatch log job to queue instead of inserting directly
		if err := dispatchLogJob(&log); err != nil {
			fmt.Printf("❌ autoLogDelete: failed to dispatch log job for %s:%d: %v\n", obj.GetLoggableType(), id, err)
		} else {
			fmt.Printf("✅ autoLogDelete: successfully dispatched log job for %s:%d\n", obj.GetLoggableType(), id)
		}
	})
}

// -------------------- HELPERS --------------------
// dispatchLogJob dispatches log job to queue instead of inserting directly
func dispatchLogJob(log *model.Log) error {
	if globalDispatcher == nil {
		// Fallback: if dispatcher is not available, skip logging
		fmt.Printf("⚠️  [WARN] Dispatcher not available, skipping log for %s:%d\n", log.LoggableType, log.LoggableID)
		return nil
	}

	job := jobs.NewLogJob(
		log.LoggableID,
		log.LoggableType,
		log.Action,
		log.OldValue,
		log.NewValue,
		log.Metadata,
		log.CreatedBy,
		log.CreatedAt,
	)

	err := globalDispatcher.Dispatch(job)
	if err != nil {
		fmt.Printf("❌ [ERROR] Failed to dispatch log job for %s:%d: %v\n", log.LoggableType, log.LoggableID, err)
		return err
	}
	
	fmt.Printf("✅ [DEBUG] Successfully dispatched log job for %s:%d to queue\n", log.LoggableType, log.LoggableID)
	return nil
}

// insertLogRaw inserts log using raw SQL to completely avoid GORM associations (DEPRECATED: use dispatchLogJob instead)
func insertLogRaw(db *gorm.DB, log *model.Log) error {
	baseDB := db.Statement.DB
	if baseDB == nil {
		return fmt.Errorf("cannot get base DB instance")
	}

	// Get underlying *sql.DB to execute raw SQL
	sqlDB, err := baseDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Prepare JSON values
	var oldValueStr, newValueStr, metadataStr string
	if log.OldValue != nil && len(log.OldValue) > 0 {
		oldValueStr = string(log.OldValue)
	}
	if log.NewValue != nil && len(log.NewValue) > 0 {
		newValueStr = string(log.NewValue)
	}
	if log.Metadata != nil && len(log.Metadata) > 0 {
		metadataStr = string(log.Metadata)
	} else {
		metadataStr = "{}"
	}
	
	// Use PostgreSQL placeholders directly with database/sql
	var createdByVal interface{} = log.CreatedBy
	if log.CreatedBy == nil {
		createdByVal = nil
	}
	
	// Handle NULL JSON values
	var oldVal, newVal interface{}
	if oldValueStr != "" {
		oldVal = oldValueStr
	} else {
		oldVal = nil
	}
	if newValueStr != "" {
		newVal = newValueStr
	} else {
		newVal = nil
	}
	
	sql := `INSERT INTO logs (loggable_id, loggable_type, action, old_value, new_value, metadata, created_by, created_at) 
			VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6::jsonb, $7, $8)`
	
	_, err = sqlDB.Exec(sql,
		log.LoggableID,
		log.LoggableType,
		log.Action,
		oldVal,
		newVal,
		metadataStr,
		createdByVal,
		log.CreatedAt,
	)
	
	return err
}

func iterateLoggable(dest interface{}, fn func(Loggable)) {
	if dest == nil {
		return
	}

	val := reflect.ValueOf(dest)
	
	// Handle pointer - dereference it first
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
		dest = val.Interface()
	}

	// Handle single object
	if loggable, ok := dest.(Loggable); ok {
		fn(loggable)
		return
	}

	// Handle slice of Loggable
	if loggables, ok := dest.([]Loggable); ok {
		for _, obj := range loggables {
			fn(obj)
		}
		return
	}

	// Handle slice or array (including after dereferencing pointer to slice)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			
			// Get the element as interface (could be pointer or value)
			elemInterface := elem.Interface()
			
			// Try to cast as Loggable directly first
			if loggable, ok := elemInterface.(Loggable); ok {
				fn(loggable)
				continue
			}
			
			// If element is a pointer, the Loggable check above should have caught it
			// Only try specific types if the general Loggable check failed
			if elem.Kind() == reflect.Ptr {
				if elem.IsNil() {
					continue
				}
				// These checks are redundant if Loggable check passed, but keep as fallback
				// They should not execute if Loggable check above succeeded
			} else {
				// Element is a value, try to get pointer to it
				if elem.CanAddr() {
					if loggablePtr, ok := elem.Addr().Interface().(Loggable); ok {
						fn(loggablePtr)
					}
				}
			}
		}
		return
	}
}

func shouldLog(obj Loggable) bool {
	// Don't log the Log table itself to avoid infinite loop
	if obj.GetLoggableType() == "logs" {
		return false
	}

	// Log all Loggable objects
	switch obj.(type) {
	case *model.CoaAccount, *model.RuleValue:
		return true
	default:
		// Check if it implements Loggable interface
		return obj != nil
	}
}

func cloneEmpty(obj Loggable) interface{} {
	switch v := obj.(type) {
	case *model.CoaAccount:
		return &model.CoaAccount{}
	case *model.RuleValue:
		return &model.RuleValue{}
	default:
		// Use reflection to create a new instance of the same type
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			elemType := val.Elem().Type()
			return reflect.New(elemType).Interface()
		}
		return nil
	}
}

func convertToUint64(val interface{}) uint64 {
	switch v := val.(type) {
	case uint64:
		return v
	case int:
		return uint64(v)
	case int64:
		return uint64(v)
	default:
		return 0
	}
}

func getMetadataFromContext(db *gorm.DB) json.RawMessage {
	meta := map[string]any{}
	if v := db.Statement.Context.Value("log_metadata"); v != nil {
		if m, ok := v.(map[string]any); ok {
			meta = m
		}
	}
	metaJSON, _ := json.Marshal(meta)
	return metaJSON
}
