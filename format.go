package rbxmk

import (
	"io"

	"github.com/robloxapi/types"
)

// Format defines a format for encoding between a sequence of bytes and a
// types.Value. The format can be registered with a World.
type Format struct {
	// Name is the name that identifies the format. The name matches a file
	// extension.
	Name string

	// MediaTypes is a list of media types that are associated with the format,
	// to be used by sources as needed.
	MediaTypes []string

	// CanDecode returns whether the format decodes into the given type.
	CanDecode func(typeName string) bool

	// Encode receives a value of one of a number of types and encodes it as a
	// sequence of bytes written to w.
	Encode func(opt FormatOptions, w io.Writer, v types.Value) error

	// Decode receives a sequence of bytes read from r, and decodes it into a
	// value of a single type.
	Decode func(opt FormatOptions, r io.Reader) (types.Value, error)
}

// FormatOptions contains options to be passed to Format.Encode and
// Format.Decode.
type FormatOptions struct {
}
