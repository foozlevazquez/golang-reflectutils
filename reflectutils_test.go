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

var tagGetTests = []struct {
	Tag   reflect.StructTag
	Key   string
	Value string
}{
	{`protobuf:"PB(1,2)"`, `protobuf`, `PB(1,2)`},
	{`protobuf:"PB(1,2)"`, `foo`, ``},
	{`protobuf:"PB(1,2)"`, `rotobuf`, ``},
	{`protobuf:"PB(1,2)" json:"name"`, `json`, `name`},
	{`protobuf:"PB(1,2)" json:"name"`, `protobuf`, `PB(1,2)`},
	{`k0:"values contain spaces" k1:"and\ttabs"`, "k0", "values contain spaces"},
	{`k0:"values contain spaces" k1:"and\ttabs"`, "k1", "and\ttabs"},
}

func TestTagGet(t *testing.T) {
	for _, tt := range tagGetTests {
		stMap := parseTags(tt.Tag)
		if stMap[tt.Key] != tt.Value {
			t.Errorf("stMap = %#v want %#v", stMap, tt)
		}
		t.Logf("%#q -> %#q\n", tt, stMap)
	}
}


func TestGetTagNameToFieldIndexMap(t *testing.T) {
	type Foo struct {
		MyThing string `json:"mything,omitempty"`
		OThing  string `json:"othing"`
	}
	result := map[string]int{
		"mything": 0,
		"othing": 1,
	}
	got := GetTagNameToFieldIndexMap(&Foo{}, "json")
	t.Logf("%#q -> %#q\n", result, got)
}
