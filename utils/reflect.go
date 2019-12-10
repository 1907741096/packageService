package utils

import (
	"fmt"
	"reflect"
	"runtime"
)

func setField(obj interface{}, name string, value interface{}) error {
	structData := reflect.ValueOf(obj).Elem()
	fieldVal := structData.FieldByName(name)

	if !fieldVal.IsValid() {
		return fmt.Errorf("utils.setField() No such field: %s in obj ", name)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Cannot set %s field value ", name)
	}

	val := reflect.ValueOf(value)

	if val.Kind() == fieldVal.Kind() ||
		val.Kind() == reflect.Float64 && fieldVal.Kind() == reflect.Int {
		fieldVal.Set(val.Convert(fieldVal.Type()))
	} else {
		valTypeStr := val.Type().String()
		fieldTypeStr := fieldVal.Type().String()
		return fmt.Errorf("Provided value type " + valTypeStr + " didn't match obj field type " + fieldTypeStr)
	}
	return nil
}

// SetStructByJSON 由json对象生成 struct
func ParseJsonToStruct(obj interface{}, mapData map[string]interface{}) error {
	for key, value := range mapData {
		if err := setField(obj, key, value); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func GetFunctionName() string {
	pc, _, _, ok := runtime.Caller(2)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()

		return funcName
	}

	return "default"
}

