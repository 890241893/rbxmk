package dump

import (
	"sort"

	"github.com/anaminus/rbxmk/dump/dt"
)

// Enums maps a name to an enum.
type Enums map[string]Enum

// Enum describes the API of an enum.
type Enum struct {
	// Summary is a fragment reference pointing to a short summary of the enum.
	Summary string `json:",omitempty"`
	// Description is a fragment reference pointing to a detailed description of
	// the enum.
	Description string `json:",omitempty"`

	// Items are the items that exist on the enum.
	Items EnumItems
}

const V_Enum = "Enum"

func (v Enum) v() {}

func (v Enum) Kind() string { return V_Enum }

// Type implements Value by returning the Enum primitive.
func (v Enum) Type() dt.Type {
	return dt.Prim("Enum")
}

func (v Enum) Index(path []string, name string) ([]string, Value) {
	return append(path, "Enum", "Items", name), v.Items[name]
}

func (v Enum) Indices() []string {
	l := make([]string, 0, len(v.Items))
	for k := range v.Items {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}

// EnumItems maps a name to an enum.
type EnumItems map[string]EnumItem

// EnumItem describes the API of an enum item.
type EnumItem struct {
	// Value is the value of the item.
	Value int

	// Summary is a fragment reference pointing to a short summary of the enum
	// item.
	Summary string `json:",omitempty"`
	// Description is a fragment reference pointing to a detailed description of
	// the enum item.
	Description string `json:",omitempty"`
}

const V_EnumItem = "EnumItem"

func (v EnumItem) v() {}

func (v EnumItem) Kind() string { return V_EnumItem }

// Type implements Value by returning the EnumItem primitive.
func (v EnumItem) Type() dt.Type {
	return dt.Prim("EnumItem")
}

func (v EnumItem) Index(path []string, name string) ([]string, Value) { return nil, nil }

func (v EnumItem) Indices() []string { return nil }
