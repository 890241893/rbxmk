package reflect

import (
	"fmt"

	lua "github.com/anaminus/gopher-lua"
	"github.com/anaminus/rbxmk"
	"github.com/anaminus/rbxmk/dump"
	"github.com/anaminus/rbxmk/dump/dt"
	"github.com/anaminus/rbxmk/rtypes"
	"github.com/robloxapi/types"
)

func init() { register(FormatSelector) }
func FormatSelector() rbxmk.Reflector {
	return rbxmk.Reflector{
		Name: "FormatSelector",
		PushTo: func(s rbxmk.Context, v types.Value) (lv lua.LValue, err error) {
			switch v := v.(type) {
			case types.Stringlike:
				table := s.CreateTable(0, 1)
				table.RawSetString("Format", lua.LString(v.Stringlike()))
				return table, nil
			case rtypes.FormatSelector:
				if s.CycleGuard() {
					defer s.CycleClear()
				}
				if s.CycleMark(&v) {
					return nil, fmt.Errorf("format selectors cannot be cyclic")
				}
				format := s.Format(v.Format)
				if format.Name == "" {
					return nil, fmt.Errorf("unknown format")
				}
				if len(format.Options) == 0 {
					table := s.CreateTable(0, 1)
					table.RawSetString("Format", lua.LString(format.Name))
					return table, nil
				}
				table := s.CreateTable(0, len(format.Options))
				for field, fieldTypes := range format.Options {
					value, ok := v.Options[field]
					if ok {
						for _, fieldType := range fieldTypes {
							if v.Type() == fieldType {
								rfl := s.Reflector(fieldType)
								if rfl.Name == "" {
									return nil, fmt.Errorf("unknown type %q for option %s of format %s", fieldType, field, format.Name)
								}
								v, err := rfl.PushTo(s, value)
								if err != nil {
									return nil, fmt.Errorf("field %s for format %s: %w", field, format.Name, err)
								}
								table.RawSetString(field, v)
							}
						}
						return nil, fmt.Errorf("expected type %s for option %s of format %s, got %s", s.ListTypes(fieldTypes), field, format.Name, value.Type())
					}
				}
				return table, nil
			default:
				return nil, rbxmk.TypeError{Want: "FormatSelector or string", Got: v.Type()}
			}
		},
		PullFrom: func(s rbxmk.Context, lv lua.LValue) (v types.Value, err error) {
			switch v := lv.(type) {
			case lua.LString:
				format := s.Format(string(v))
				if format.Name == "" {
					return nil, fmt.Errorf("unknown format %q", v)
				}
				return rtypes.FormatSelector{Format: format.Name}, nil
			case *lua.LTable:
				if s.CycleGuard() {
					defer s.CycleClear()
				}
				if s.CycleMark(v) {
					return nil, fmt.Errorf("tables cannot be cyclic")
				}
				name, ok := v.RawGetString("Format").(lua.LString)
				if !ok {
					return nil, fmt.Errorf("Format field must be a string")
				}
				format := s.Format(string(name))
				if format.Name == "" {
					return nil, fmt.Errorf("unknown format %q", name)
				}
				if len(format.Options) == 0 {
					return rtypes.FormatSelector{Format: format.Name}, nil
				}
				sel := rtypes.FormatSelector{
					Format:  format.Name,
					Options: make(rtypes.Dictionary),
				}
				for field, fieldTypes := range format.Options {
					v, err := s.PullAnyFromDictionaryOpt(v, field, nil, fieldTypes...)
					if err != nil {
						return nil, fmt.Errorf("field %s for format %s: %w", field, format.Name, err)
					}
					if v != nil {
						sel.Options[field] = v
					}
				}
				return sel, nil
			default:
				return nil, rbxmk.TypeError{Want: "string or table", Got: v.Type().String()}
			}
		},
		SetTo: func(p interface{}, v types.Value) error {
			switch p := p.(type) {
			case *rtypes.FormatSelector:
				*p = v.(rtypes.FormatSelector)
			default:
				return setPtrErr(p, v)
			}
			return nil
		},
		Dump: func() dump.TypeDef {
			return dump.TypeDef{
				Underlying:  dt.Or{dt.Prim("string"), dt.Struct{"Format": dt.Prim("string"), "...": dt.Prim("any")}},
				Summary:     "Types/FormatSelector:Summary",
				Description: "Types/FormatSelector:Description",
			}
		},
	}
}
