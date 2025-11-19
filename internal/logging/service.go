package logging

import (
	model "core-ledger/model/core-ledger"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
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
func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Create().After("gorm:after_create").Register("auto_log_create", autoLogCreate)
	db.Callback().Update().Before("gorm:before_update").Register("auto_log_before_update", autoLogBeforeUpdate)
	db.Callback().Update().After("gorm:after_update").Register("auto_log_after_update", autoLogAfterUpdate)
	db.Callback().Delete().After("gorm:after_delete").Register("auto_log_delete", autoLogDelete)
}

// -------------------- CREATE CALLBACK --------------------
func autoLogCreate(db *gorm.DB) {
	iterateLoggable(db.Statement.Dest, func(obj Loggable) {
		if !shouldLog(obj) {
			return
		}

		id := obj.GetLoggableID()
		if id == 0 {
			return
		}

		newJSON, _ := json.Marshal(obj)
		log := model.Log{
			LoggableID:   id,
			LoggableType: obj.GetLoggableType(),
			Action:       "CREATED",
			NewValue:     datatypes.JSON(newJSON),
			CreatedAt:    time.Now(),
		}

		// Use a separate DB session to avoid invalid transaction
		tx := db.Session(&gorm.Session{NewDB: true})
		if err := tx.Create(&log).Error; err != nil {
			fmt.Printf("autoLogCreate: failed to create log for %s:%d: %v\n", obj.GetLoggableType(), id, err)
		}
	})
}

// -------------------- BEFORE UPDATE CALLBACK --------------------
func autoLogBeforeUpdate(db *gorm.DB) {
	oldMap := map[uint64]json.RawMessage{}

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

		if err := db.First(oldStruct, id).Error; err != nil {
			return
		}

		oldJSON, _ := json.Marshal(oldStruct)
		oldMap[id] = oldJSON
	})

	// Map update
	if m, ok := db.Statement.Dest.(map[string]interface{}); ok {
		if idVal, exists := m["id"]; exists {
			if id := convertToUint64(idVal); id != 0 {
				var old model.CoaAccount
				if err := db.First(&old, id).Error; err == nil {
					oldJSON, _ := json.Marshal(old)
					oldMap[id] = oldJSON
				}
			}
		}
	}

	db.InstanceSet("old_values", oldMap)
}

// -------------------- AFTER UPDATE CALLBACK --------------------
func autoLogAfterUpdate(db *gorm.DB) {
	oldMap, _ := db.InstanceGet("old_values")
	oldMapTyped, _ := oldMap.(map[uint64]json.RawMessage)

	iterateLoggable(db.Statement.Dest, func(obj Loggable) {
		if !shouldLog(obj) {
			return
		}
		id := obj.GetLoggableID()
		if id == 0 {
			return
		}

		newJSON, _ := json.Marshal(obj)
		var oldJSON json.RawMessage
		if oldMapTyped != nil {
			oldJSON = oldMapTyped[id]
			delete(oldMapTyped, id)
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

		tx := db.Session(&gorm.Session{NewDB: true})
		if err := tx.Create(&log).Error; err != nil {
			fmt.Printf("autoLogAfterUpdate: failed to create log for %s:%d: %v\n", obj.GetLoggableType(), id, err)
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
			tx := db.Session(&gorm.Session{NewDB: true})
			_ = tx.Create(&log).Error
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

		tx := db.Session(&gorm.Session{NewDB: true})
		_ = tx.Create(&log).Error
	})
}

// -------------------- HELPERS --------------------
func iterateLoggable(dest interface{}, fn func(Loggable)) {
	switch v := dest.(type) {
	case Loggable:
		fn(v)
	case []Loggable:
		for _, obj := range v {
			fn(obj)
		}
	case []*model.CoaAccount:
		for _, obj := range v {
			fn(obj)
		}
	}
}

func shouldLog(obj Loggable) bool {
	switch obj.(type) {
	case *model.CoaAccount:
		return true
	default:
		return false
	}
}

func cloneEmpty(obj Loggable) interface{} {
	switch obj.(type) {
	case *model.CoaAccount:
		return &model.CoaAccount{}
	default:
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
