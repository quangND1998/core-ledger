package repo

import (
	"core-ledger/model/dto"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Scoper interface{}

// ApplyFilterScopeDynamic applies dynamic scopes from filters to any model with ScopeXxx(...) methods.
// ApplyFilterScopeDynamic áp dụng các filter scope từ map,
// includeSort=true → áp dụng cả ScopeSort
// includeSort=false → bỏ ScopeSort (dùng cho COUNT)
func ApplyFilterScopeDynamic[T any](db *gorm.DB, filters map[string]any, includeSort bool) *gorm.DB {
	var t T
	tType := reflect.TypeOf(t)

	var model reflect.Value
	if tType.Kind() == reflect.Pointer {
		model = reflect.New(tType.Elem()) // *CoaAccount
	} else {
		model = reflect.New(tType) // CoaAccount → &CoaAccount{}
	}

	modelPtr := model
	if modelPtr.Kind() != reflect.Pointer {
		modelPtr = modelPtr.Addr()
	}

	for key, rawParam := range filters {
		if key == "page" || key == "limit" || key == "offset" {
			continue
		}

		scopeName := jsonKeyToScopeName(key)
		if !includeSort && scopeName == "ScopeSort" {
			continue // bỏ sort khi includeSort=false
		}

		log.Printf("[DynamicScope] Checking method: %s with param %+v", scopeName, rawParam)

		method, ok := modelPtr.Type().MethodByName(scopeName)
		if !ok {
			log.Printf("[DynamicScope] Method %s not found", scopeName)
			continue
		}

		converted, err := convertDynamicValue(rawParam, method.Type.In(1))
		if err != nil {
			log.Printf("[DynamicScope] Cannot convert param for %s: %v", scopeName, err)
			continue
		}

		out := method.Func.Call([]reflect.Value{modelPtr, converted})
		if len(out) == 1 {
			if scopeFn, ok := out[0].Interface().(func(*gorm.DB) *gorm.DB); ok {
				db = db.Scopes(scopeFn)
				log.Println("[DynamicScope] Scope applied")
			} else {
				log.Printf("[DynamicScope] Method %s return not func(*gorm.DB)*gorm.DB", scopeName)
			}
		}
	}

	return db
}

// convertDynamicValue: convert interface{} sang reflect.Value theo target type
func convertDynamicValue(val any, targetType reflect.Type) (reflect.Value, error) {
	v := reflect.ValueOf(val)

	switch targetType.Kind() {
	case reflect.String:
		str, ok := val.(string)
		if !ok {
			return reflect.Value{}, ErrInvalidType
		}
		return reflect.ValueOf(str), nil

	case reflect.Slice:
		sliceVal := reflect.MakeSlice(targetType, 0, 0)
		arr, ok := val.([]any)
		if !ok {
			return reflect.Value{}, ErrInvalidType
		}

		for _, item := range arr {
			sliceVal = reflect.Append(sliceVal, reflect.ValueOf(item).Convert(targetType.Elem()))
		}
		return sliceVal, nil

	default:
		return v.Convert(targetType), nil
	}
}

var ErrInvalidType = fmt.Errorf("invalid type conversion")

func jsonKeyToScopeName(key string) string {
	key = strings.Split(key, ",")[0] // remove omitempty

	// Convert snake_case → camelCase
	parts := strings.Split(key, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}

	return "Scope" + strings.Join(parts, "")
}

// convertDynamicValue converts JSON values to the method's expected type.

// BuildParamsFromFilter converts filter struct to map[string]any
func BuildParamsFromFilter(f interface{}) map[string]any {
	result := map[string]any{}
	if f == nil {
		return result
	}

	v := reflect.ValueOf(f)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return result
		}
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// handle embedded structs
		if field.Kind() == reflect.Struct && fieldType.Anonymous {
			embedded := BuildParamsFromFilter(field.Addr().Interface())
			for k, val := range embedded {
				result[k] = val
			}
			continue
		}

		key := fieldType.Tag.Get("json")
		if key == "" {
			key = fieldType.Name
		} else {
			key = strings.Split(key, ",")[0] // <--- important fix
		}

		switch field.Kind() {
		case reflect.Ptr:
			if !field.IsNil() {
				result[key] = field.Elem().Interface()
			}
		case reflect.Slice:
			if field.Len() > 0 {
				var arr []any
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					if elem.Kind() == reflect.Ptr && !elem.IsNil() {
						arr = append(arr, elem.Elem().Interface())
					} else {
						arr = append(arr, elem.Interface())
					}
				}
				result[key] = arr
			}
		case reflect.String:
			if field.String() != "" {
				result[key] = field.String()
			}
		case reflect.Int, reflect.Int64:
			if field.Int() != 0 {
				result[key] = field.Int()
			}
		case reflect.Uint, reflect.Uint64:
			if field.Uint() != 0 {
				result[key] = field.Uint()
			}
		case reflect.Bool:
			result[key] = field.Bool()
		}
	}

	return result
}

// Paginate generic helper
func CustomPaginate[T any](db *gorm.DB, filters map[string]any, page, limit int64, out *[]*T) (*dto.PaginationResponse[*T], error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if limit <= 0 {
		limit = 25
	}
	if page <= 0 {
		page = 1
	}

	// COUNT tổng số (bỏ sort)
	countDB := ApplyFilterScopeDynamic[T](db.Session(&gorm.Session{}), filters, false)
	var total int64
	if err := countDB.Count(&total).Error; err != nil {
		log.Println("COUNT SQL:", countDB.Statement.SQL.String(), countDB.Statement.Vars)
		return nil, fmt.Errorf("failed to count: %w", err)
	}

	// Query dữ liệu thực tế (filter + sort)
	dataDB := ApplyFilterScopeDynamic[T](db.Session(&gorm.Session{}), filters, true)
	offset := int((page - 1) * limit)
	if err := dataDB.Limit(int(limit)).Offset(offset).Find(out).Error; err != nil {
		log.Println("FIND SQL:", dataDB.Statement.SQL.String(), dataDB.Statement.Vars)
		return nil, fmt.Errorf("failed to find: %w", err)
	}

	// Build pagination response
	totalPage := (total + limit - 1) / limit
	var nextPage, prevPage *int64
	if page < totalPage {
		n := page + 1
		nextPage = &n
	}
	if page > 1 {
		p := page - 1
		prevPage = &p
	}

	return &dto.PaginationResponse[*T]{
		Items:     *out,
		Total:     total,
		Limit:     limit,
		Page:      page,
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}
