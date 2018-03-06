package swagger_handler

import (
	"net/http"
	"io"
	"bytes"
	"github.com/go-openapi/loads"
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
)

type Opts struct {
	BasePath string
	Path     string
	SpecUrl  string
}

func NewHandler(opts Opts, file []byte) (handler http.Handler, err error) {
	specDoc, err := loads.Analyzed(file, "")
	if err != nil {
		return
	}

	b, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
	if err != nil {
		return
	}

	handler = http.NotFoundHandler()
	handler = middleware.Redoc(middleware.RedocOpts{
		BasePath: opts.BasePath,
		SpecURL:  opts.SpecUrl,
		Path:     opts.Path,
	}, handler)
	handler = middleware.Spec(opts.BasePath, b, handler)

	return
}

func NewAssetHandler(file []byte) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		reader := bytes.NewBuffer(file)
		io.Copy(rw, reader)
	})
}