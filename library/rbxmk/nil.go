package reflect

import (
	lua "github.com/anaminus/gopher-lua"
	. "github.com/anaminus/rbxmk"
	"github.com/anaminus/rbxmk/rtypes"
	"github.com/robloxapi/types"
)

func init() { register(Nil) }
func Nil() Reflector {
	return Reflector{
		Name: "nil",
		PushTo: func(s State, r Reflector, v types.Value) (lvs []lua.LValue, err error) {
			return []lua.LValue{lua.LNil}, nil
		},
		PullFrom: func(s State, r Reflector, lvs ...lua.LValue) (v types.Value, err error) {
			if lvs[0] == lua.LNil {
				return rtypes.Nil, nil
			}
			return nil, TypeError(nil, 0, "nil")
		},
	}
}
