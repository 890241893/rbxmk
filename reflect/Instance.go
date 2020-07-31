package reflect

import (
	"fmt"

	. "github.com/anaminus/rbxmk"
	"github.com/anaminus/rbxmk/rtypes"
	"github.com/robloxapi/rbxdump"
	"github.com/robloxapi/types"
	"github.com/yuin/gopher-lua"
)

// pushPropertyTo behaves like PushVariantTo, except that exprims types are
// reflected as userdata.
func pushPropertyTo(s State, v types.Value) (lv lua.LValue, err error) {
	switch v.(type) {
	case types.Numberlike:
	case types.Intlike:
	case types.Stringlike:
	default:
		return PushVariantTo(s, v)
	}
	typ := s.Type(v.Type())
	if typ.Name == "" {
		return nil, fmt.Errorf("unknown type %q", string(v.Type()))
	}
	if typ.PushTo == nil {
		return nil, fmt.Errorf("unable to cast %s to Variant", typ.Name)
	}
	u := s.L.NewUserData()
	u.Value = v
	s.L.SetMetatable(u, s.L.GetTypeMetatable(typ.Name))
	return u, nil
}

// convertType tries to convert v to t.
func convertType(s State, t string, v types.Value) (nv types.Value, ok bool) {
	if v.Type() == t {
		return v, true
	}
	if s.Type(t).Name == "" {
		return v, false
	}
	switch t {
	case "int":
		switch v := v.(type) {
		case types.Intlike:
			return types.Int(v.Intlike()), true
		case types.Numberlike:
			return types.Int(v.Numberlike()), true
		}
	case "int64":
		switch v := v.(type) {
		case types.Intlike:
			return types.Int64(v.Intlike()), true
		case types.Numberlike:
			return types.Int64(v.Numberlike()), true
		}
	case "float":
		switch v := v.(type) {
		case types.Numberlike:
			return types.Float(v.Numberlike()), true
		case types.Intlike:
			return types.Float(v.Intlike()), true
		}
	case "double":
		switch v := v.(type) {
		case types.Numberlike:
			return types.Double(v.Numberlike()), true
		case types.Intlike:
			return types.Double(v.Intlike()), true
		}
	case "string":
		if v, ok := v.(types.Stringlike); ok {
			return types.String(v.Stringlike()), true
		}
	case "BinaryString":
		if v, ok := v.(types.Stringlike); ok {
			return types.BinaryString(v.Stringlike()), true
		}
	case "ProtectedString":
		if v, ok := v.(types.Stringlike); ok {
			return types.ProtectedString(v.Stringlike()), true
		}
	case "Content":
		if v, ok := v.(types.Stringlike); ok {
			return types.Content(v.Stringlike()), true
		}
	case "SharedString":
		if v, ok := v.(types.Stringlike); ok {
			return types.SharedString(v.Stringlike()), true
		}
	case "Color3":
		if v, ok := v.(rtypes.Color3uint8); ok {
			return types.Color3(v), true
		}
	case "Color3uint8":
		if v, ok := v.(types.Color3); ok {
			return rtypes.Color3uint8(v), true
		}
	}
	return v, false
}

// getPropDesc gets a property descriptor from a class, or any class it inherits
// from.
func getPropDesc(root *rtypes.RootDesc, class *rbxdump.Class, name string) (prop *rbxdump.Property) {
	for class != nil {
		prop, _ = class.Members[name].(*rbxdump.Property)
		if prop != nil {
			return prop
		}
		class = root.Classes[class.Superclass]
	}
	return nil
}

func checkEnumDesc(s State, desc *rtypes.RootDesc, enum, class, prop string) *rtypes.Enum {
	enumValue := desc.EnumTypes.Enum(enum)
	if enumValue == nil {
		if desc.Enums[enum] == nil {
			s.L.RaiseError(
				"no enum descriptor %q for property descriptor %s.%s",
				enum,
				class,
				prop,
			)
			return nil
		}
		s.L.RaiseError(
			"no enum value %q generated for property descriptor %s.%s",
			enum,
			class,
			prop,
		)
		return nil
	}
	return enumValue
}

func checkClassDesc(s State, desc *rtypes.RootDesc, typ, class, prop string) *rbxdump.Class {
	classDesc := desc.Classes[typ]
	if classDesc == nil {
		s.L.RaiseError(
			"no class descriptor %q for property descriptor %s.%s",
			typ,
			class,
			prop,
		)
		return nil
	}
	return classDesc
}

func Instance() Type {
	return Type{
		Name:     "Instance",
		PushTo:   PushTypeTo,
		PullFrom: PullTypeFrom,
		Metatable: Metatable{
			"__tostring": func(s State) int {
				s.L.Push(lua.LString(s.Pull(1, "Instance").(*rtypes.Instance).String()))
				return 1
			},
			"__eq": func(s State) int {
				op := s.Pull(2, "Instance").(*rtypes.Instance)
				return s.Push(types.Bool(s.Pull(1, "Instance").(*rtypes.Instance) == op))
			},
			"__index": func(s State) int {
				inst := s.Pull(1, "Instance").(*rtypes.Instance)
				name := string(s.Pull(2, "string").(types.String))
				desc := s.Desc(inst)
				var classDesc *rbxdump.Class
				if desc != nil {
					classDesc = desc.Classes[inst.ClassName]
				}

				// Try GetService.
				if inst.IsDataModel() && name == "GetService" {
					s.L.Push(s.L.NewFunction(func(l *lua.LState) int {
						u := l.CheckUserData(1)
						if u.Metatable != l.GetTypeMetatable("Instance") {
							TypeError(l, 1, "Instance")
							return 0
						}
						inst, ok := u.Value.(*rtypes.Instance)
						if !ok {
							TypeError(l, 1, "Instance")
							return 0
						}
						s := State{World: s.World, L: l}
						className := string(s.Pull(2, "string").(types.String))
						if desc != nil {
							classDesc := desc.Classes[className]
							if classDesc == nil || !classDesc.GetTag("Service") {
								s.L.RaiseError("%q is not a valid service", className)
							}
						}
						service := inst.FindFirstChildOfClass(className, false)
						if service == nil {
							service = rtypes.NewInstance(className)
							service.IsService = true
							service.SetName(className)
							service.SetParent(inst)
						}
						return s.Push(service)
					}))
					return 1
				}

				// Try property.
				var lv lua.LValue
				var err error
				value := inst.Get(name)
				if classDesc != nil {
					propDesc := getPropDesc(desc, classDesc, name)
					if propDesc == nil {
						s.L.RaiseError("%s is not a valid member", name)
						return 0
					}
					if value == nil {
						s.L.RaiseError("property %s not initialized", name)
						return 0
					}
					switch propDesc.ValueType.Category {
					case "Class":
						inst, ok := value.(*rtypes.Instance)
						if !ok {
							s.L.RaiseError("stored value type %s is not an instance", value.Type())
							return 0
						}
						class := checkClassDesc(s, desc, propDesc.ValueType.Name, classDesc.Name, propDesc.Name)
						if class == nil {
							return 0
						}
						if inst.ClassName != class.Name {
							s.L.RaiseError("instance of class %s expected, got %s", class.Name, inst.ClassName)
							return 0
						}
						return s.Push(inst)
					case "Enum":
						enum := checkEnumDesc(s, desc, propDesc.ValueType.Name, classDesc.Name, propDesc.Name)
						if enum == nil {
							return 0
						}
						token, ok := value.(types.Token)
						if !ok {
							s.L.RaiseError("stored value type %s is not a token", value.Type())
							return 0
						}
						item := enum.Value(int(token))
						if item == nil {
							s.L.RaiseError("invalid stored value %d for enum %s", value, enum.Name())
							return 0
						}
						return s.Push(item)
					default:
						if a, b := value.Type(), propDesc.ValueType.Name; a != b {
							s.L.RaiseError("stored value type %s does not match property type %s", a, b)
							return 0
						}
					}
					// Push without converting exprims.
					lv, err = PushVariantTo(s, value)
				} else {
					if value == nil {
						// Fallback to nil.
						return s.Push(rtypes.Nil)
					}
					lv, err = pushPropertyTo(s, value)
				}
				if err != nil {
					s.L.RaiseError(err.Error())
					return 0
				}
				s.L.Push(lv)
				return 1
			},
			"__newindex": func(s State) int {
				inst := s.Pull(1, "Instance").(*rtypes.Instance)
				name := string(s.Pull(2, "string").(types.String))

				// Try GetService.
				if inst.IsDataModel() && name == "GetService" {
					s.L.RaiseError("%s cannot be assigned to", name)
					return 0
				}

				// Try property.
				value := PullVariant(s, 3)

				desc := s.Desc(inst)
				var classDesc *rbxdump.Class
				if desc != nil {
					classDesc = desc.Classes[inst.ClassName]
				}
				if classDesc != nil {
					propDesc := getPropDesc(desc, classDesc, name)
					if propDesc == nil {
						s.L.RaiseError("%s is not a valid member", name)
						return 0
					}
					switch propDesc.ValueType.Category {
					case "Class":
						inst, ok := value.(*rtypes.Instance)
						if !ok {
							s.L.RaiseError("Instance expected, got %s", value.Type())
							return 0
						}
						class := checkClassDesc(s, desc, propDesc.ValueType.Name, classDesc.Name, propDesc.Name)
						if class == nil {
							return 0
						}
						if inst.ClassName != class.Name {
							s.L.RaiseError("instance of class %s expected, got %s", class.Name, inst.ClassName)
							return 0
						}
						inst.Set(name, inst)
						return 0
					case "Enum":
						enum := checkEnumDesc(s, desc, propDesc.ValueType.Name, classDesc.Name, propDesc.Name)
						if enum == nil {
							return 0
						}
						switch value := value.(type) {
						case types.Token:
							item := enum.Value(int(value))
							if item == nil {
								s.L.RaiseError("invalid value %d for enum %s", value, enum.Name())
								return 0
							}
							inst.Set(name, value)
							return 0
						case *rtypes.EnumItem:
							item := enum.Value(value.Value())
							if item == nil {
								s.L.RaiseError(
									"invalid value %s (%d) for enum %s",
									value.String(),
									value.Value(),
									enum.String(),
								)
								return 0
							}
							if a, b := enum.Name(), value.Enum().Name(); a != b {
								s.L.RaiseError("expected enum %s, got %s", a, b)
								return 0
							}
							if a, b := item.Name(), value.Name(); a != b {
								s.L.RaiseError("expected enum item %s, got %s", a, b)
								return 0
							}
							inst.Set(name, types.Token(item.Value()))
							return 0
						case types.Intlike:
							v := int(value.Intlike())
							item := enum.Value(v)
							if item == nil {
								s.L.RaiseError("invalid value %d for enum %s", v, enum.Name())
								return 0
							}
							inst.Set(name, types.Token(item.Value()))
							return 0
						case types.Numberlike:
							v := int(value.Numberlike())
							item := enum.Value(v)
							if item == nil {
								s.L.RaiseError("invalid value %d for enum %s", v, enum.Name())
								return 0
							}
							inst.Set(name, types.Token(item.Value()))
							return 0
						case types.Stringlike:
							v := value.Stringlike()
							item := enum.Item(v)
							if item == nil {
								s.L.RaiseError("invalid value %s for enum %s", v, enum.Name())
								return 0
							}
							inst.Set(name, types.Token(item.Value()))
							return 0
						default:
							s.L.RaiseError("invalid value for enum %s", enum.Name())
							return 0
						}
					default:
						var ok bool
						value, ok = convertType(s, propDesc.ValueType.Name, value)
						if !ok {
							s.L.RaiseError("%s expected, got %s", propDesc.ValueType.Name, value.Type())
							return 0
						}
					}
				}
				prop, ok := value.(types.PropValue)
				if !ok {
					s.L.RaiseError("cannot assign %s as property", value.Type())
					return 0
				}
				inst.Set(name, prop)
				return 0
			},
		},
		Members: Members{
			"ClassName": Member{
				Get: func(s State, v types.Value) int {
					return s.Push(types.String(v.(*rtypes.Instance).ClassName))
				},
				// Allowed to be set for convenience.
				Set: func(s State, v types.Value) {
					inst := v.(*rtypes.Instance)
					if inst.IsDataModel() {
						s.L.RaiseError("%s cannot be assigned to", "ClassName")
						return
					}
					inst.ClassName = string(s.Pull(3, "string").(types.String))
				},
			},
			"Name": Member{
				Get: func(s State, v types.Value) int {
					return s.Push(types.String(v.(*rtypes.Instance).Name()))
				},
				Set: func(s State, v types.Value) {
					v.(*rtypes.Instance).SetName(string(s.Pull(3, "string").(types.String)))
				},
			},
			"Parent": Member{
				Get: func(s State, v types.Value) int {
					if parent := v.(*rtypes.Instance).Parent(); parent != nil {
						return s.Push(parent)
					}
					return s.Push(rtypes.Nil)
				},
				Set: func(s State, v types.Value) {
					var err error
					switch parent := s.PullAnyOf(3, "Instance", "nil").(type) {
					case *rtypes.Instance:
						err = v.(*rtypes.Instance).SetParent(parent)
					case nil:
						err = v.(*rtypes.Instance).SetParent(nil)
					}
					if err != nil {
						s.L.RaiseError(err.Error())
					}
				},
			},
			"ClearAllChildren": Member{Method: true, Get: func(s State, v types.Value) int {
				v.(*rtypes.Instance).RemoveAll()
				return 0
			}},
			"Clone": Member{Method: true, Get: func(s State, v types.Value) int {
				return s.Push(v.(*rtypes.Instance).Clone())
			}},
			"Destroy": Member{Method: true, Get: func(s State, v types.Value) int {
				v.(*rtypes.Instance).SetParent(nil)
				return 0
			}},
			"FindFirstAncestor": Member{Method: true, Get: func(s State, v types.Value) int {
				name := string(s.Pull(2, "string").(types.String))
				if ancestor := v.(*rtypes.Instance).FindFirstAncestorOfClass(name); ancestor != nil {
					return s.Push(ancestor)
				}
				return s.Push(rtypes.Nil)
			}},
			"FindFirstAncestorOfClass": Member{Method: true, Get: func(s State, v types.Value) int {
				className := string(s.Pull(2, "string").(types.String))
				if ancestor := v.(*rtypes.Instance).FindFirstAncestorOfClass(className); ancestor != nil {
					return s.Push(ancestor)
				}
				return s.Push(rtypes.Nil)
			}},
			"FindFirstChild": Member{Method: true, Get: func(s State, v types.Value) int {
				name := string(s.Pull(2, "string").(types.String))
				recurse := bool(s.PullOpt(3, "bool", types.False).(types.Bool))
				if child := v.(*rtypes.Instance).FindFirstChild(name, recurse); child != nil {
					return s.Push(child)
				}
				return s.Push(rtypes.Nil)
			}},
			"FindFirstChildOfClass": Member{Method: true, Get: func(s State, v types.Value) int {
				className := string(s.Pull(2, "string").(types.String))
				recurse := bool(s.PullOpt(3, "bool", types.False).(types.Bool))
				if child := v.(*rtypes.Instance).FindFirstChildOfClass(className, recurse); child != nil {
					return s.Push(child)
				}
				return s.Push(rtypes.Nil)
			}},
			"GetChildren": Member{Method: true, Get: func(s State, v types.Value) int {
				t := v.(*rtypes.Instance).Children()
				return s.Push(rtypes.Instances(t))
			}},
			"GetDescendants": Member{Method: true, Get: func(s State, v types.Value) int {
				return s.Push(rtypes.Instances(v.(*rtypes.Instance).Descendants()))
			}},
			"GetFullName": Member{Method: true, Get: func(s State, v types.Value) int {
				return s.Push(types.String(v.(*rtypes.Instance).GetFullName()))
			}},
			"IsAncestorOf": Member{Method: true, Get: func(s State, v types.Value) int {
				descendant := s.Pull(2, "Instance").(*rtypes.Instance)
				return s.Push(types.Bool(v.(*rtypes.Instance).IsAncestorOf(descendant)))
			}},
			"IsDescendantOf": Member{Method: true, Get: func(s State, v types.Value) int {
				ancestor := s.Pull(2, "Instance").(*rtypes.Instance)
				return s.Push(types.Bool(v.(*rtypes.Instance).IsDescendantOf(ancestor)))
			}},
		},
		Constructors: Constructors{
			"new": func(s State) int {
				className := string(s.Pull(1, "string").(types.String))
				inst := rtypes.NewInstance(className)
				return s.Push(inst)
			},
		},
		Environment: func(s State) {
			t := s.L.CreateTable(0, 1)
			t.RawSetString("new", s.L.NewFunction(func(l *lua.LState) int {
				dataModel := rtypes.NewDataModel()
				return s.Push(dataModel)
			}))
			s.L.SetGlobal("DataModel", t)
		},
	}
}
