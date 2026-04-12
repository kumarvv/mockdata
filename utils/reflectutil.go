package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const typeTime = "time.Time"

// ToString converts interface{} to string
func ToString(v interface{}) (s string) {
	if v == nil {
		return ""
	}

	vt := reflect.TypeOf(v)

	if vt.String() == typeTime {
		s = v.(time.Time).Format(time.RFC3339)
		return
	}

	if vt.Kind() == reflect.String {
		s, _ = v.(string)
	} else {
		s = fmt.Sprintf("%v", v)
	}

	return
}

// ToInt converts interface{} to int64
func ToInt(v interface{}) (i int64) {
	if v == nil {
		return 0
	}

	vt := reflect.TypeOf(v)

	switch vt.Kind() {
	case reflect.String:
		i, _ = strconv.ParseInt(v.(string), 10, 64)
	case reflect.Int64:
		i = v.(int64)
	case reflect.Int32:
		i = int64(v.(int32))
	case reflect.Int16:
		i = int64(v.(int16))
	case reflect.Int8:
		i = int64(v.(int8))
	case reflect.Int:
		i = int64(v.(int))
	case reflect.Uint64:
		i = int64(v.(uint64))
	case reflect.Uint32:
		i = int64(v.(uint32))
	case reflect.Uint16:
		i = int64(v.(uint16))
	case reflect.Uint8:
		i = int64(v.(uint8))
	case reflect.Uint:
		i = int64(v.(uint))
	case reflect.Float64:
		i = int64(v.(float64))
	case reflect.Float32:
		i = int64(v.(float32))
	default:
		i = 0
	}
	return
}

// ToBool converts interface{} to bool
func ToBool(v interface{}) (b bool) {
	if v == nil {
		b = false
	} else {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			s := strings.ToLower(v.(string))
			b = s == "true" || s == "y" || s == "yes"
		case reflect.Bool:
			b = v.(bool)
		default:
		}
	}
	return
}

// ToTime converts interface{} to time
// expects time.Type or String type with RFC3339 format
func ToTime(v interface{}) (t time.Time, err error) {
	if v == nil {
		return
	}

	if reflect.TypeOf(v).String() == typeTime {
		t = v.(time.Time)
	} else if reflect.TypeOf(v).Kind() == reflect.String {
		t, err = time.Parse(time.RFC3339, v.(string))
		//if err != nil {
		//	// try converting with mm/dd/yyyy
		//	t, err = time.Parse(timeutil.DateFormatMDY, v.(string))
		//}
		return
	}
	return
}

// ToFloat converts interface{} to float64
func ToFloat(v interface{}) (f float64) {
	if v == nil {
		return 0
	}

	vt := reflect.TypeOf(v)
	switch vt.Kind() {
	case reflect.String:
		f, _ = strconv.ParseFloat(v.(string), 10)
	case reflect.Float64:
		f = float64(v.(float64))
	case reflect.Float32:
		f = float64(v.(float32))
	default:
		f = 0
	}
	return
}
