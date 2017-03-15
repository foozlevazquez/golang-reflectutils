package reflectutils

import (
	"testing"
	"reflect"
	"os"
)

type foo struct {
	slot0 int
	slot1 float32
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}


func hasField(mp map[string]reflect.Type, name string) bool {
	_, ok := mp[name]
	return ok
}

func TestStructFieldData(t *testing.T) {
	fm := StructFieldData(&foo{})

	slotNames := []string{ "slot0", "slot1" }

	for _, slotName := range slotNames {
		if !hasField(fm, slotName) {
			t.Errorf("Missing slot %s", slotName)
		}
	}
}
