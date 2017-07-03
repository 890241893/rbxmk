package format

import (
	"github.com/anaminus/rbxmk"
	"github.com/robloxapi/rbxapi"
	"github.com/robloxapi/rbxfile"
	"github.com/robloxapi/rbxfile/xml"
	"io"
)

func init() {
	Formats.Register(rbxmk.Format{
		Name: "XML Properties",
		Ext:  "properties.xml",
		Codec: func(opt rbxmk.Options, ctx interface{}) rbxmk.FormatCodec {
			return &XMLCodec{API: opt.Config.API}
		},
		InputDrills:  []rbxmk.Drill{DrillProperty},
		OutputDrills: []rbxmk.Drill{DrillProperty},
		Merger:       MergeProperties,
	})
}

type XMLCodec struct {
	API *rbxapi.API
}

func (c *XMLCodec) Decode(r io.Reader, data *rbxmk.Data) (err error) {
	doc := &xml.Document{}
	if _, err = doc.ReadFrom(r); err != nil {
		return err
	}
	if doc.Root == nil || doc.Root.StartName != "Properties" {
		*data = nil
		return nil
	}
	inst := &rbxfile.Instance{Properties: make(map[string]rbxfile.Value, len(doc.Root.Tags))}
	xml.RobloxCodec{API: c.API}.DecodeProperties(doc.Root.Tags, inst, nil)
	*data = inst.Properties
	return nil
}

func (c *XMLCodec) Encode(w io.Writer, data rbxmk.Data) (err error) {
	switch v := data.(type) {
	case *[]*rbxfile.Instance:
		if len(*v) > 0 {
			data = (*v)[0].Properties
		}
	case *rbxfile.Instance:
		data = v.Properties
	case Property:
		data = map[string]rbxfile.Value{v.Name: v.Properties[v.Name]}
	case nil:
		data = map[string]rbxfile.Value{}
	}

	props, ok := data.(map[string]rbxfile.Value)
	if !ok {
		return NewDataTypeError(data)
	}

	doc := &xml.Document{Indent: "\t"}
	root := &xml.Tag{StartName: "Properties"}
	doc.Root = root

	inst := &rbxfile.Instance{Properties: props}
	root.Tags = xml.RobloxCodec{API: c.API}.EncodeProperties(inst)
	_, err = doc.WriteTo(w)
	return err
}
