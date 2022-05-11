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

const T_HTTPHeaders = "HTTPHeaders"

func init() { register(HTTPHeaders) }
func HTTPHeaders() rbxmk.Reflector {
	return rbxmk.Reflector{
		Name: T_HTTPHeaders,
		PushTo: func(c rbxmk.Context, v types.Value) (lv lua.LValue, err error) {
			headers, ok := v.(rtypes.HTTPHeaders)
			if !ok {
				return nil, rbxmk.TypeError{Want: T_HTTPHeaders, Got: v.Type()}
			}
			table := c.CreateTable(0, len(headers))
			for name, values := range headers {
				vs := c.CreateTable(len(values), 0)
				for _, value := range values {
					vs.Append(lua.LString(value))
				}
				table.RawSetString(name, vs)
			}
			return table, nil
		},
		PullFrom: func(c rbxmk.Context, lv lua.LValue) (v types.Value, err error) {
			table, ok := lv.(*lua.LTable)
			if !ok {
				return nil, rbxmk.TypeError{Want: T_Table, Got: lv.Type().String()}
			}
			headers := make(rtypes.HTTPHeaders)
			err = table.ForEach(func(k, lv lua.LValue) error {
				name, ok := k.(lua.LString)
				if !ok {
					return nil
				}
				values, err := pullStringArray(lv)
				if err != nil {
					return fmt.Errorf("header %q: %w", string(name), err)
				}
				headers[string(name)] = values
				return nil
			})
			if err != nil {
				return nil, err
			}
			return headers, nil
		},
		SetTo: func(p interface{}, v types.Value) error {
			switch p := p.(type) {
			case *rtypes.HTTPHeaders:
				*p = v.(rtypes.HTTPHeaders)
			default:
				return setPtrErr(p, v)
			}
			return nil
		},
		Dump: func() dump.TypeDef {
			return dump.TypeDef{
				Underlying:  dt.Map{K: dt.Prim(T_String), V: dt.Or{dt.Prim(T_String), dt.Array{T: dt.Prim(T_String)}}},
				Summary:     "Types/HTTPHeaders:Summary",
				Description: "Types/HTTPHeaders:Description",
			}
		},
	}
}

// Convert a string or an array of strings.
func pullStringArray(v lua.LValue) ([]string, error) {
	switch v := v.(type) {
	case lua.LString:
		return []string{string(v)}, nil
	case *lua.LTable:
		n := v.Len()
		if n == 0 {
			return nil, fmt.Errorf("expected string or array of strings")
		}
		values := make([]string, n)
		for i := 1; i <= n; i++ {
			value, ok := v.RawGetInt(i).(lua.LString)
			if !ok {
				return nil, fmt.Errorf("index %d: expected string, got %s", i, value.Type())
			}
			values[i-1] = string(value)
		}
		return values, nil
	default:
		return nil, fmt.Errorf("expected string or array of strings")
	}
}
