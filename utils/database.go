package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

func GenerateUniqueSlug(db *gorm.DB, model interface{}, fieldName string) string {
	val := reflect.ValueOf(model).Elem()
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		panic(fmt.Sprintf("Field %s not found on model %s", fieldName, reflect.TypeOf(model).Name()))
	}

	name := field.String()
	slug := strings.ToLower(name)
	slug = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")

	originalSlug := slug
	counter := 1
	for {
		var count int64
		db.Model(model).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++
	}

	return slug
}
