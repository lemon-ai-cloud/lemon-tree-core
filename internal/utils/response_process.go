package utils

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"reflect"
	"time"
)

// ConvertToMapAndReplaceTime 递归转换任意对象为 map，并将 time.Time 替换为 Unix 时间戳
func ConvertToMapAndReplaceTime(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return nil
	}

	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		return ConvertToMapAndReplaceTime(val.Elem().Interface())

	case reflect.Struct:
		// Special case: time.Time
		if t, ok := v.(time.Time); ok {
			return t.Unix()
		}

		// Special case: uuid.UUID
		if id, ok := v.(uuid.UUID); ok {
			return id.String()
		}

		// Special case: sql.NullTime
		if nt, ok := v.(sql.NullTime); ok {
			if nt.Valid {
				return nt.Time.Unix()
			}
			return nil
		}

		// Special case: gorm.DeletedAt
		if da, ok := v.(gorm.DeletedAt); ok {
			if da.Valid {
				return da.Time.Unix()
			}
			return nil
		}

		out := map[string]interface{}{}
		t := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" { // unexported field
				continue
			}

			fieldVal := val.Field(i).Interface()

			// Flatten embedded fields
			if field.Anonymous {
				embedded := ConvertToMapAndReplaceTime(fieldVal)
				if m, ok := embedded.(map[string]interface{}); ok {
					for k, v := range m {
						out[k] = v
					}
				}
				continue
			}

			jsonKey := field.Tag.Get("json")
			if jsonKey == "" || jsonKey == "-" {
				jsonKey = field.Name
			}

			out[jsonKey] = ConvertToMapAndReplaceTime(fieldVal)
		}
		return out

	case reflect.Slice, reflect.Array:
		out := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			out[i] = ConvertToMapAndReplaceTime(val.Index(i).Interface())
		}
		return out

	case reflect.Map:
		out := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			strKey := fmt.Sprintf("%v", key.Interface())
			out[strKey] = ConvertToMapAndReplaceTime(val.MapIndex(key).Interface())
		}
		return out

	default:
		return v
	}
}

func JsonResponse(c *gin.Context, code int, obj interface{}) {
	//result := gin.H{"body": obj}
	//fmt.Println("yyy: ", result)
	//converted := ConvertToMapAndReplaceTime(result)
	//converted := ConvertToMapAndReplaceTime(obj)
	c.JSON(code, obj)
}
