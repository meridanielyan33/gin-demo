package utils

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func BuildProjection(fieldsToInclude, fieldsToExclude string) bson.M {
	fieldsToInclude = strings.TrimSpace(fieldsToInclude)
	fieldsToExclude = strings.TrimSpace(fieldsToExclude)

	if fieldsToInclude != "" {
		proj := bson.M{}
		for _, f := range strings.Split(fieldsToInclude, ",") {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			proj[f] = 1
		}
		return proj
	}

	if fieldsToExclude != "" {
		proj := bson.M{}
		for _, f := range strings.Split(fieldsToExclude, ",") {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			proj[f] = 0
		}
		return proj
	}

	return nil
}

func FieldIncluded(projection bson.M, field string) bool {
	if projection == nil || len(projection) == 0 {
		return true
	}
	if v, ok := projection[field]; ok && v == 0 {
		return false
	}

	return true
}

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
