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

// ------------------ Public CreateLog ------------------
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
	}

	return s.db.Create(&log).Error
}

// ------------------ Register Callbacks ------------------
func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Create().After("gorm:after_create").Register("auto_log_create", autoLogCreate)
	db.Callback().Update().Before("gorm:before_update").Register("auto_log_before_update", autoLogBeforeUpdate)
	db.Callback().Update().After("gorm:after_update").Register("auto_log_after_update", autoLogAfterUpdate)
	db.Callback().Delete().After("gorm:after_delete").Register("auto_log_delete", autoLogDelete)
}

// ------------------ CREATE ------------------
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

		// Tạo log trên session mới, skip hooks để tránh recursion
		if err := db.Session(&gorm.Session{NewDB: true, SkipHooks: true}).Create(&log).Error; err != nil {
			fmt.Printf("autoLogCreate: failed to create log for %s:%d: %v\n", obj.GetLoggableType(), id, err)
		}
	})
}

// ------------------ BEFORE UPDATE ------------------
func autoLogBeforeUpdate(db *gorm.DB) {
	iterateLoggable(db.Statement.Dest, func(obj Loggable) {
		if !shouldLog(obj) {
			return
		}

		id := obj.GetLoggableID()
		if id == 0 {
			return
		}

		var old interface{}
		switch obj.(type) {
		case *model.CoaAccount:
			old = &model.CoaAccount{}
		default:
			return
		}

		if err := db.First(old, id).Error; err != nil {
			return
		}

		oldJSON, _ := json.Marshal(old)
		var oldMap map[uint64]json.RawMessage
		if v, ok := db.InstanceGet("old_values"); ok {
			if m, ok2 := v.(map[uint64]json.RawMessage); ok2 {
				oldMap = m
			}
		}
		if oldMap == nil {
			oldMap = map[uint64]json.RawMessage{}
		}
		oldMap[id] = oldJSON
		db.InstanceSet("old_values", oldMap)
	})

	// Map update handling
	if m, ok := db.Statement.Dest.(map[string]interface{}); ok {
		if idVal, hasID := m["id"]; hasID {
			var id uint64
			switch v := idVal.(type) {
			case uint64:
				id = v
			case int:
				id = uint64(v)
			case int64:
				id = uint64(v)
			default:
				return
			}
			if id == 0 {
				return
			}

			var old model.CoaAccount
			if err := db.First(&old, id).Error; err != nil {
				return
			}

			oldJSON, _ := json.Marshal(old)
			var oldMap map[uint64]json.RawMessage
			if v, ok := db.InstanceGet("old_values"); ok {
				if m, ok2 := v.(map[uint64]json.RawMessage); ok2 {
					oldMap = m
				}
			}
			if oldMap == nil {
				oldMap = map[uint64]json.RawMessage{}
			}
			oldMap[id] = oldJSON
			db.InstanceSet("old_values", oldMap)
		}
	}
}

// ------------------ AFTER UPDATE ------------------
func autoLogAfterUpdate(db *gorm.DB) {
	iterateLoggable(db.Statement.Dest, func(obj Loggable) {
		if !shouldLog(obj) {
			return
		}

		id := obj.GetLoggableID()
		if id == 0 {
			return
		}

		var oldJSON json.RawMessage
		if v, ok := db.InstanceGet("old_values"); ok {
			if m, ok2 := v.(map[uint64]json.RawMessage); ok2 {
				oldJSON = m[id]
				delete(m, id)
				db.InstanceSet("old_values", m)
			}
		}

		newJSON, _ := json.Marshal(obj)
		meta := map[string]any{}
		if v, ok := db.Statement.Context.Value("log_metadata").(map[string]any); ok {
			meta = v
		}
		metaJSON, _ := json.Marshal(meta)

		log := model.Log{
			LoggableID:   id,
			LoggableType: obj.GetLoggableType(),
			Action:       "UPDATED",
			OldValue:     datatypes.JSON(oldJSON),
			NewValue:     datatypes.JSON(newJSON),
			Metadata:     datatypes.JSON(metaJSON),
			CreatedAt:    time.Now(),
		}

		if err := db.Session(&gorm.Session{NewDB: true, SkipHooks: true}).Create(&log).Error; err != nil {
			fmt.Printf("autoLogAfterUpdate: failed to create log for %s:%d: %v\n", obj.GetLoggableType(), id, err)
		}
	})

	// Map update case
	if m, ok := db.Statement.Dest.(map[string]interface{}); ok {
		if idVal, hasID := m["id"]; hasID {
			var id uint64
			switch v := idVal.(type) {
			case uint64:
				id = v
			case int:
				id = uint64(v)
			case int64:
				id = uint64(v)
			default:
				return
			}
			if id == 0 {
				return
			}

			var newObj model.CoaAccount
			if err := db.First(&newObj, id).Error; err != nil {
				return
			}

			var oldJSON json.RawMessage
			if v, ok := db.InstanceGet("old_values"); ok {
				if m, ok2 := v.(map[uint64]json.RawMessage); ok2 {
					oldJSON = m[id]
					delete(m, id)
					db.InstanceSet("old_values", m)
				}
			}

			newJSON, _ := json.Marshal(newObj)
			meta := map[string]any{}
			if v, ok := db.Statement.Context.Value("log_metadata").(map[string]any); ok {
				meta = v
			}
			metaJSON, _ := json.Marshal(meta)

			log := model.Log{
				LoggableID:   id,
				LoggableType: "coa_accounts",
				Action:       "UPDATED",
				OldValue:     datatypes.JSON(oldJSON),
				NewValue:     datatypes.JSON(newJSON),
				Metadata:     datatypes.JSON(metaJSON),
				CreatedAt:    time.Now(),
			}
			if err := db.Session(&gorm.Session{NewDB: true, SkipHooks: true}).Create(&log).Error; err != nil {
				fmt.Printf("autoLogAfterUpdate(map): failed to create log for %s:%d: %v\n", "coa_accounts", id, err)
			}
		}
	}
}

// ------------------ DELETE ------------------
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

		if err := db.Session(&gorm.Session{NewDB: true, SkipHooks: true}).Create(&log).Error; err != nil {
			fmt.Printf("autoLogDelete: failed to create log for %s:%d: %v\n", obj.GetLoggableType(), id, err)
		}
	})
}

// ------------------ Helpers ------------------
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

// chỉ log các model cần thiết, chặn bảng log để tránh recursion
func shouldLog(obj Loggable) bool {
	switch obj.(type) {
	case *model.CoaAccount:
		return true

	default:
		return false
	}
}
