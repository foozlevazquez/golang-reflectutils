package reflectutils

import (
	"reflect"
	"strconv"
	"strings"
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

var stCache = map[reflect.Type]map[string]map[string]string{}

var fidxCache = map[reflect.Type]map[string]int{}

// StructTags returns a map[string]string of the tags and values of the struct.
func StructTags(o interface{}) map[string]map[string]string {
	var oVal = reflect.ValueOf(o).Elem()
	var oType = oVal.Type()

	stMap, ok := stCache[oType]
	if !ok {
		//         fName	   tagKey tagVal
		stMap = map[string]map[string]string{}

		//				 fName    fIdx
		var fidxMap = map[string]int{}
		for i := 0; i < oVal.NumField(); i++ {
			field := oType.Field(i)
			stMap[field.Name] = parseTags(field.Tag)
			fidxMap[field.Name] = i
		}
		stCache[oType] = stMap
	}
	return stMap
}

func parseTags(tag reflect.StructTag) map[string]string {
	// Copied ruthlessly from src/reflect/type.go
	var stMap = map[string]string{}

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax
		// error.  Strictly speaking, control chars include the range [0x7f,
		// 0x9f], not just [0x00, 0x1f], but in practice, we ignore the
		// multi-byte control characters as it is simpler to inspect the tag's
		// bytes than the tag's runes.
		i = 0
		for (i < len(tag) && tag[i] > ' ' && tag[i] != ':' &&
			tag[i] != '"' && tag[i] != 0x7f ){
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		stMap[name] = value
	}
	return stMap
}

// GetTagNameToFieldIndexMap returns a map the value of each field's <tagName>
// value to the index of the field it was found in.
//
// Tag values are stripped of anything past the comma (including the comma).
//
// type Foo struct {
//	  MyThing string `json:"mything,omitempty"`
//    OThing  string `json:"othing"`
// }
//		->
//	map[string]int{
//		"mything": 0,
//		"othing": 1,
//	}
//
func GetTagNameToFieldIndexMap(o interface{}, tagName string) map[string]int {
	var oFidxMap = fidxCache[reflect.TypeOf(o)]
	var fidxMap = map[string]int{}

	for fName, tagMap := range StructTags(o) {
		tVal, ok := tagMap[tagName]
		if ok {
			if idx := strings.Index(tVal, ","); idx != -1 {
				tVal = tVal[:idx]
			}
			fidxMap[tVal] = oFidxMap[fName]
		}
	}
	return fidxMap
}
