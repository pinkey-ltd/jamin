package main

import (
	"fmt"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/file"
	"io"
	"mime/multipart"
)

func main() {
	// initialise gofr object
	app := gofr.New()

	// register route upload
	app.POST("/upload", func(ctx *gofr.Context) (interface{}, error) {

		return "Hello World!", nil
	})

	// Runs the server, it will listen on the default port 8000.
	// it can be over-ridden through configs
	app.Run()
}

// Payload is the struct that we are trying to bind files to
type Payload struct {
	// Name represents the non-file field in the struct
	Name string `form:"name"`

	// The FileHeader determines the generic file format that we can get
	// from the multipart form that gets parsed by the incoming HTTP request
	FileHeader *multipart.FileHeader `file:"upload"`
}

func fileUploadHandler(ctx *gofr.Context) (resp interface{}, err error) {
	var p Payload

	// bind the multipart data into the variable p
	err = ctx.Bind(&p)
	if err != nil {
		return nil, err
	}

	// Retrieve the file from form data
	f, err := p.FileHeader.Open()
	if err != nil {
		return nil, err
	}

	defer func(file multipart.File) {
		if closeErr := file.Close(); closeErr != nil {
			err = closeErr
		}
	}(f)

	fmt.Printf("Uploaded File: %s\n", p.FileHeader.Filename)
	fmt.Printf("File Size: %d\n", p.FileHeader.Size)
	fmt.Printf("MIME Header: %v\n", p.FileHeader.Header)

	// Now let’s save it locally
	err = ctx.File.ChDir("upload")
	if err != nil {
		return nil, err
	}

	dst, err := ctx.File.Create(p.FileHeader.Filename)
	if err != nil {
		return nil, err
	}

	defer func(dst file.File) {
		if createErr := dst.Close(); createErr != nil {
			err = createErr
		}
	}(dst)

	// read the file content
	content, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}

	// Copy the uploaded file to the destination file
	if _, err := dst.Write(content); err != nil {
		return nil, err
	}
	return fmt.Sprintf("saved files: %s, len of file `a`: %d", p.FileHeader.Filename, p.FileHeader.Size), nil
}
