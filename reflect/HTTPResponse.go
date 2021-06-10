package reflect

import (
	lua "github.com/anaminus/gopher-lua"
	"github.com/anaminus/rbxmk"
	"github.com/anaminus/rbxmk/dump"
	"github.com/anaminus/rbxmk/dump/dt"
	"github.com/anaminus/rbxmk/rtypes"
	"github.com/robloxapi/types"
)

func init() { register(HTTPResponse) }
func HTTPResponse() rbxmk.Reflector {
	return rbxmk.Reflector{
		Name: "HTTPResponse",
		PushTo: func(s rbxmk.State, v types.Value) (lv lua.LValue, err error) {
			resp, ok := v.(rtypes.HTTPResponse)
			if !ok {
				return nil, rbxmk.TypeError{Want: "HTTPResponse", Got: v.Type()}
			}
			table := s.L.CreateTable(0, 5)
			s.PushToTable(table, lua.LString("Success"), types.Bool(resp.Success))
			s.PushToTable(table, lua.LString("StatusCode"), types.Int(resp.StatusCode))
			s.PushToTable(table, lua.LString("StatusMessage"), types.String(resp.StatusMessage))
			s.PushToTable(table, lua.LString("Headers"), resp.Headers)
			s.PushToTable(table, lua.LString("Cookies"), resp.Cookies)
			s.PushToTable(table, lua.LString("Body"), resp.Body)
			return table, nil
		},
		PullFrom: func(s rbxmk.State, lv lua.LValue) (v types.Value, err error) {
			table, ok := lv.(*lua.LTable)
			if !ok {
				return nil, rbxmk.TypeError{Want: "table", Got: lv.Type().String()}
			}
			resp := rtypes.HTTPResponse{
				Success:       bool(s.PullFromTable(table, lua.LString("Success"), "bool").(types.Bool)),
				StatusCode:    int(s.PullFromTable(table, lua.LString("StatusCode"), "int").(types.Int)),
				StatusMessage: string(s.PullFromTable(table, lua.LString("StatusMessage"), "string").(types.String)),
				Headers:       s.PullFromTableOpt(table, lua.LString("Headers"), rtypes.HTTPHeaders(nil), "HTTPHeaders").(rtypes.HTTPHeaders),
				Cookies:       s.PullFromTableOpt(table, lua.LString("Cookies"), rtypes.Cookies(nil), "Cookies").(rtypes.Cookies),
				Body:          s.PullFromTableOpt(table, lua.LString("Body"), nil, "Variant"),
			}
			return resp, nil
		},
		Dump: func() dump.TypeDef {
			return dump.TypeDef{
				Underlying: dt.Struct{
					"Success":       dt.Prim("bool"),
					"StatusCode":    dt.Prim("int"),
					"StatusMessage": dt.Prim("string"),
					"Headers":       dt.Prim("HTTPHeaders"),
					"Cookies":       dt.Prim("Cookies"),
					"Body":          dt.Optional{T: dt.Prim("any")},
				},
				Summary:     "Types/HTTPResponse:Summary",
				Description: "Types/HTTPResponse:Description",
			}
		},
		Types: []func() rbxmk.Reflector{
			Bool,
			Cookies,
			HTTPHeaders,
			Int,
			String,
			Variant,
		},
	}
}
