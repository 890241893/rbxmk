package reflect

import (
	lua "github.com/anaminus/gopher-lua"
	"github.com/anaminus/rbxmk"
	"github.com/anaminus/rbxmk/dump"
	"github.com/anaminus/rbxmk/dump/dt"
	"github.com/anaminus/rbxmk/rtypes"
	"github.com/robloxapi/types"
)

func init() { register(PhysicalProperties) }
func PhysicalProperties() rbxmk.Reflector {
	return rbxmk.Reflector{
		Name: "PhysicalProperties",
		PushTo: func(s rbxmk.Context, v types.Value) (lv lua.LValue, err error) {
			if pp, ok := v.(types.PhysicalProperties); ok && !pp.CustomPhysics {
				return lua.LNil, nil
			}
			u := s.L.NewUserData(v)
			s.L.SetMetatable(u, s.L.GetTypeMetatable("PhysicalProperties"))
			return u, nil
		},
		PullFrom: func(s rbxmk.Context, lv lua.LValue) (v types.Value, err error) {
			switch lv := lv.(type) {
			case *lua.LNilType:
				return types.PhysicalProperties{}, nil
			case *lua.LUserData:
				if lv.Metatable != s.L.GetTypeMetatable("PhysicalProperties") {
					return nil, rbxmk.TypeError{Want: "PhysicalProperties", Got: lv.Type().String()}
				}
				v, ok := lv.Value().(types.Value)
				if !ok {
					return nil, rbxmk.TypeError{Want: "PhysicalProperties", Got: lv.Type().String()}
				}
				return v, nil
			default:
				return nil, rbxmk.TypeError{Want: "PhysicalProperties", Got: lv.Type().String()}
			}
		},
		ConvertFrom: func(v types.Value) types.Value {
			switch v := v.(type) {
			case rtypes.NilType:
				return types.PhysicalProperties{}
			case types.PhysicalProperties:
				return v
			}
			return nil
		},
		Metatable: rbxmk.Metatable{
			"__tostring": func(s rbxmk.State) int {
				v := s.Pull(1, "PhysicalProperties").(types.PhysicalProperties)
				s.L.Push(lua.LString(v.String()))
				return 1
			},
		},
		Properties: rbxmk.Properties{
			"Density": {
				Get: func(s rbxmk.State, v types.Value) int {
					return s.Push(types.Float(v.(types.PhysicalProperties).Density))
				},
				Dump: func() dump.Property {
					return dump.Property{
						ValueType:   dt.Prim("float"),
						ReadOnly:    true,
						Summary:     "Types/PhysicalProperties:Properties/Density/Summary",
						Description: "Types/PhysicalProperties:Properties/Density/Description",
					}
				},
			},
			"Friction": {
				Get: func(s rbxmk.State, v types.Value) int {
					return s.Push(types.Float(v.(types.PhysicalProperties).Friction))
				},
				Dump: func() dump.Property {
					return dump.Property{
						ValueType:   dt.Prim("float"),
						ReadOnly:    true,
						Summary:     "Types/PhysicalProperties:Properties/Friction/Summary",
						Description: "Types/PhysicalProperties:Properties/Friction/Description",
					}
				},
			},
			"Elasticity": {
				Get: func(s rbxmk.State, v types.Value) int {
					return s.Push(types.Float(v.(types.PhysicalProperties).Elasticity))
				},
				Dump: func() dump.Property {
					return dump.Property{
						ValueType:   dt.Prim("float"),
						ReadOnly:    true,
						Summary:     "Types/PhysicalProperties:Properties/Elasticity/Summary",
						Description: "Types/PhysicalProperties:Properties/Elasticity/Description",
					}
				},
			},
			"FrictionWeight": {
				Get: func(s rbxmk.State, v types.Value) int {
					return s.Push(types.Float(v.(types.PhysicalProperties).FrictionWeight))
				},
				Dump: func() dump.Property {
					return dump.Property{
						ValueType:   dt.Prim("float"),
						ReadOnly:    true,
						Summary:     "Types/PhysicalProperties:Properties/FrictionWeight/Summary",
						Description: "Types/PhysicalProperties:Properties/FrictionWeight/Description",
					}
				},
			},
			"ElasticityWeight": {
				Get: func(s rbxmk.State, v types.Value) int {
					return s.Push(types.Float(v.(types.PhysicalProperties).ElasticityWeight))
				},
				Dump: func() dump.Property {
					return dump.Property{
						ValueType:   dt.Prim("float"),
						ReadOnly:    true,
						Summary:     "Types/PhysicalProperties:Properties/ElasticityWeight/Summary",
						Description: "Types/PhysicalProperties:Properties/ElasticityWeight/Description",
					}
				},
			},
		},
		Constructors: rbxmk.Constructors{
			"new": {
				Func: func(s rbxmk.State) int {
					var v types.PhysicalProperties
					switch s.Count() {
					case 3:
						v.Density = float32(s.Pull(1, "float").(types.Float))
						v.Friction = float32(s.Pull(2, "float").(types.Float))
						v.Elasticity = float32(s.Pull(3, "float").(types.Float))
					case 5:
						v.Density = float32(s.Pull(1, "float").(types.Float))
						v.Friction = float32(s.Pull(2, "float").(types.Float))
						v.Elasticity = float32(s.Pull(3, "float").(types.Float))
						v.FrictionWeight = float32(s.Pull(4, "float").(types.Float))
						v.ElasticityWeight = float32(s.Pull(5, "float").(types.Float))
					default:
						return s.RaiseError("expected 3 or 5 arguments")
					}
					return s.Push(v)
				},
				Dump: func() dump.MultiFunction {
					return dump.MultiFunction{
						{
							Parameters: dump.Parameters{
								{Name: "density", Type: dt.Prim("float")},
								{Name: "friction", Type: dt.Prim("float")},
								{Name: "elasticity", Type: dt.Prim("float")},
							},
							Returns: dump.Parameters{
								{Type: dt.Prim("PhysicalProperties")},
							},
							Summary:     "Types/PhysicalProperties:Constructors/new/Components/Summary",
							Description: "Types/PhysicalProperties:Constructors/new/Components/Description",
						},
						{
							Parameters: dump.Parameters{
								{Name: "density", Type: dt.Prim("float")},
								{Name: "friction", Type: dt.Prim("float")},
								{Name: "elasticity", Type: dt.Prim("float")},
								{Name: "frictionWeight", Type: dt.Prim("float")},
								{Name: "elasticityWeight", Type: dt.Prim("float")},
							},
							Returns: dump.Parameters{
								{Type: dt.Prim("PhysicalProperties")},
							},
							Summary:     "Types/PhysicalProperties:Constructors/new/Weights/Summary",
							Description: "Types/PhysicalProperties:Constructors/new/Weights/Description",
						},
					}
				},
			},
		},
		Dump: func() dump.TypeDef {
			return dump.TypeDef{
				Summary:     "Types/PhysicalProperties:Summary",
				Description: "Types/PhysicalProperties:Description",
			}
		},
		Types: []func() rbxmk.Reflector{
			Float,
		},
	}
}
