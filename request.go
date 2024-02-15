package mango

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

type DataMap map[string][]string
type FileMap map[string][]*multipart.FileHeader

type Request struct {
	Method  string
	Path    string
	Args    DataMap
	Headers DataMap
	Data    []byte
	Form    DataMap
	Files   FileMap
}

func ParseRequest(req *http.Request) Request {

	req.ParseForm()
	argsMap := req.Form
	bodyForm := req.PostForm

	for key := range bodyForm {
		delete(argsMap, key)
	}

	var data []byte
	mediaType, params, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))

	var fileMap FileMap = make(FileMap)
	if req.Body != nil {
		if strings.HasPrefix(mediaType, "multipart/") {
			mr := multipart.NewReader(req.Body, params["boundary"])
			form, _ := mr.ReadForm(32 << 20)

			for key, values := range form.Value {
				for _, value := range values {
					bodyForm.Add(key, value)
				}
			}
			for key, values := range form.File {
				fileMap[key] = values
			}
		} else {
			data, _ = io.ReadAll(req.Body)
		}

		defer req.Body.Close()
	}

	return Request{
		Method:  req.Method,
		Path:    req.URL.Path,
		Headers: DataMap(req.Header),
		Args:    DataMap(argsMap),
		Form:    DataMap(bodyForm),
		Data:    data,
		Files:   fileMap,
	}
}
