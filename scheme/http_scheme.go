package scheme

import (
	"bytes"
	"errors"
	"github.com/anaminus/rbxmk"
	"net/http"
)

func init() {
	input := rbxmk.InputScheme{
		Handler: httpInputSchemeHandler,
	}
	output := rbxmk.OutputScheme{
		Handler:   httpOutputSchemeHandler,
		Finalizer: httpOutputFinalizer,
	}

	registerInput("http", input)
	registerOutput("http", output)

	registerInput("https", input)
	registerOutput("https", output)
}

func httpInputSchemeHandler(opt rbxmk.Options, node *rbxmk.InputNode, inref []string) (ext string, outref []string, data rbxmk.Data, err error) {
	ext = node.Format
	if !opt.Formats.Registered(ext) {
		return "", nil, nil, errors.New("format is not registered")
	}

	resp, err := http.Get(node.Reference[0])
	if err != nil {
		return "", nil, nil, err
	}
	defer resp.Body.Close()
	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return "", nil, nil, errors.New(resp.Status)
	}

	if err := opt.Formats.Decode(ext, opt, nil, resp.Body, &data); err != nil {
		return "", nil, nil, err
	}
	return ext, inref[1:], data, err
}

func httpOutputSchemeHandler(opt rbxmk.Options, node *rbxmk.OutputNode, inref []string) (ext string, outref []string, data rbxmk.Data, err error) {
	return node.Format, inref[1:], nil, nil
}

func httpOutputFinalizer(opt rbxmk.Options, node *rbxmk.OutputNode, inref []string, ext string, outdata rbxmk.Data) (err error) {
	if !opt.Formats.Registered(ext) {
		return errors.New("format is not registered")
	}
	var buf bytes.Buffer
	if err = opt.Formats.Encode(ext, opt, nil, &buf, outdata); err != nil {
		return err
	}
	resp, err := http.Post(node.Reference[0], "", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return errors.New(resp.Status)
	}
	return nil
}
