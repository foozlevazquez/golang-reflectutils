package reflectutils

import (
	"reflect"
)

var (
	sfdCache = map[reflect.Type]map[string]reflect.Type{}
)


func StructFieldData(o interface{}) map[string]reflect.Type {
	var oVal = reflect.ValueOf(o).Elem()
	var oType = oVal.Type()

	// Check cache
	fieldMap, ok := sfdCache[oType]
	if !ok {
		fieldMap = map[string]reflect.Type{}
		sfdCache[oType] = fieldMap
	}

	for i := 0; i < oVal.NumField(); i++ {
		var tf = oType.Field(i)
		var vf = oVal.Field(i)
		fieldMap[tf.Name] = vf.Type()
	}
	return fieldMap
}


// StructFieldValue returns the value of slotName of struct o.
func StructFieldValue(o interface{}, slotName string) interface{} {
	var oVal = reflect.ValueOf(o).Elem()

	var fieldVal = oVal.FieldByName(slotName)

	if ! fieldVal.IsValid() {
		panic("No such slot name" + slotName)
	}

	return fieldVal.Interface()
}
