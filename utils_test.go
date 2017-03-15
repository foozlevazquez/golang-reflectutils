package reflectutils

import (
	"testing"
	"reflect"
	"os"
)

type foo struct {
	Slot0 int
	Slot1 float32
}

var (
	slotNames = []string{ "Slot0", "Slot1" }
)


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

	for _, slotName := range slotNames {
		if !hasField(fm, slotName) {
			t.Errorf("Missing slot %s", slotName)
		}
	}
}

func TestStructFieldValue(t *testing.T) {
	mFoo := &foo{ Slot0: 99, Slot1: 123.9 }

	if StructFieldValue(mFoo, "Slot0").(int) != mFoo.Slot0 {
		t.Errorf("slot0 incorrect %v != %v",
			StructFieldValue(mFoo, "Slot0"), mFoo.Slot0)
	}
	if StructFieldValue(mFoo, "Slot1").(float32) != mFoo.Slot1 {
		t.Errorf("slot1 incorrect %v != %v",
			StructFieldValue(mFoo, "Slot1"), mFoo.Slot1)
	}
}
