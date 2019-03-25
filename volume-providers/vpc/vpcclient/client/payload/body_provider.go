/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package payload

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
)

// JSONBodyProvider ...
type JSONBodyProvider struct{ payload interface{} }

// NewJSONBodyProvider ...
func NewJSONBodyProvider(p interface{}) *JSONBodyProvider {
	return &JSONBodyProvider{payload: p}
}

// ContentType ...
func (p *JSONBodyProvider) ContentType() string {
	return "application/json"
}

// Body ...
func (p *JSONBodyProvider) Body() (io.Reader, error) {
	buf := &bytes.Buffer{}

	err := json.NewEncoder(buf).Encode(p.payload)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// MultipartFileBody ...
type MultipartFileBody struct {
	name            string
	contents        io.Reader
	multipartWriter *multipart.Writer
	pipeReader      *io.PipeReader
	pipeWriter      *io.PipeWriter
}

// NewMultipartFileBody ...
func NewMultipartFileBody(name string, contents io.Reader) *MultipartFileBody {
	pr, pw := io.Pipe()
	return &MultipartFileBody{
		name:            name,
		contents:        contents,
		pipeReader:      pr,
		pipeWriter:      pw,
		multipartWriter: multipart.NewWriter(pw),
	}
}

// ContentType ...
func (p *MultipartFileBody) ContentType() string {
	return p.multipartWriter.FormDataContentType()
}

// Body ...
func (p *MultipartFileBody) Body() (io.Reader, error) {
	go p.copyBody()
	return p.pipeReader, nil
}

func (p *MultipartFileBody) copyBody() {
	defer p.Close()

	fileWriter, err := p.multipartWriter.CreateFormFile(p.name, "image")
	if err != nil {
		p.pipeWriter.CloseWithError(err)
	}

	_, err = io.Copy(fileWriter, p.contents)
	if err != nil {
		p.pipeWriter.CloseWithError(err)
	}
}

// Close ...
func (p *MultipartFileBody) Close() {
	p.multipartWriter.Close()
	p.pipeWriter.Close()
}
