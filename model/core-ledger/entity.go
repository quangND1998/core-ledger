package model

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Entity struct {
}

func (e *Entity) ScopeSort(sortStr string, model interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		allowed := map[string]bool{}

		t := reflect.TypeOf(model)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			gormTag := f.Tag.Get("gorm")
			col := ""
			if strings.Contains(gormTag, "column:") {
				parts := strings.Split(gormTag, ";")
				for _, p := range parts {
					if strings.HasPrefix(p, "column:") {
						col = strings.TrimPrefix(p, "column:")
						break
					}
				}
			}
			if col == "" {
				col = toSnakeCase(f.Name)
			}
			allowed[col] = true
		}

		if strings.TrimSpace(sortStr) == "" {
			return db
		}

		orders := []string{}
		for _, pair := range strings.Split(sortStr, ",") {
			parts := strings.Split(pair, ":")
			if len(parts) != 2 {
				continue
			}
			col := strings.TrimSpace(parts[0])
			dir := strings.TrimSpace(parts[1])
			if !allowed[col] {
				continue
			}
			order := col
			if dir == "-1" {
				order += " DESC"
			} else {
				order += " ASC"
			}
			orders = append(orders, order)
		}

		if len(orders) > 0 {
			db = db.Order(strings.Join(orders, ", "))
		}
		return db
	}
}

// helper snake_case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		if r >= 'A' && r <= 'Z' {
			r = r + ('a' - 'A')
		}
		result = append(result, r)
	}
	return string(result)
}
