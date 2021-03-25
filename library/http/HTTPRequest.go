package reflect

import (
	"github.com/anaminus/rbxmk"
	"github.com/anaminus/rbxmk/dump"
	"github.com/anaminus/rbxmk/dump/dt"
	"github.com/robloxapi/types"
)

func init() { register(HTTPRequest) }
func HTTPRequest() rbxmk.Reflector {
	return rbxmk.Reflector{
		Name:     "HTTPRequest",
		PushTo:   rbxmk.PushPtrTypeTo("HTTPRequest"),
		PullFrom: rbxmk.PullTypeFrom("HTTPRequest"),
		Methods: rbxmk.Methods{
			"Resolve": {
				Func: func(s rbxmk.State, v types.Value) int {
					req := v.(*rbxmk.HTTPRequest)
					resp, err := req.Resolve()
					if err != nil {
						return s.RaiseError("%s", err)
					}
					return s.Push(*resp)
				},
				Dump: func() dump.Function {
					return dump.Function{
						Returns: dump.Parameters{
							{Name: "resp", Type: dt.Prim("HTTPResponse")},
						},
						CanError:    true,
						Summary:     "libraries/http/types/HTTPRequest:Methods/Resolve/Summary",
						Description: "libraries/http/types/HTTPRequest:Methods/Resolve/Description",
					}
				},
			},
			"Cancel": {
				Func: func(s rbxmk.State, v types.Value) int {
					req := v.(*rbxmk.HTTPRequest)
					req.Cancel()
					return 0
				},
				Dump: func() dump.Function {
					return dump.Function{
						Summary:     "libraries/http/types/HTTPRequest:Methods/Cancel/Summary",
						Description: "libraries/http/types/HTTPRequest:Methods/Cancel/Description",
					}
				},
			},
		},
		Dump: func() dump.TypeDef {
			return dump.TypeDef{
				Summary:     "libraries/http/types/HTTPRequest:Summary",
				Description: "libraries/http/types/HTTPRequest:Description",
			}
		},
	}
}
