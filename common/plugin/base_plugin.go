package plugin

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type BasePlugin struct {
}

func (op *BasePlugin) Name() string {
	return "basePlugin"
}

func (op *BasePlugin) Initialize(db *gorm.DB) (err error) {
	// 创建字段的时候ulid
	db.Callback().Create().Before("gorm:create").Replace("id", func(db *gorm.DB) {
		// 检查是否存在Schema以及ID字段
		if db.Statement.Schema != nil {
			snakeInitials := getSnakeInitials(db.Statement.Schema.Table)
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Struct:
				field := db.Statement.Schema.LookUpField("id")
				// 检查ID字段是否存在以及是否已经设置了值
				if field != nil && field.PrimaryKey {
					if _, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue); isZero {
						// 如果ID字段的值为零值，则生成一个新ID
						db.Statement.SetColumn("id", fmt.Sprintf("%s_%s", snakeInitials, strings.ToLower(ulid.Make().String())))
					}
				}
				break
			case reflect.Slice, reflect.Array:
				// 如果是切片或者数组，遍历每一个元素
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					field := db.Statement.Schema.LookUpField("id")
					// 检查ID字段是否存在以及是否已经设置了值
					if field != nil && field.PrimaryKey {
						if _, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue.Index(i)); isZero {
							// 如果ID字段的值为零值，则生成一个新ID
							field.Set(db.Statement.Context, db.Statement.ReflectValue.Index(i), fmt.Sprintf("%s_%s", snakeInitials, strings.ToLower(ulid.Make().String())))
						}
					}
				}
				break
			}
		}
	})
	return
}

func getSnakeInitials(table string) string {
	words := strings.Split(table, "_")
	initials := ""
	for _, word := range words {
		initials += string(word[0])
	}
	return initials
}
