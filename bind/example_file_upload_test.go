package bind_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/bind"
	"github.com/jvcoutinho/lit/render"
	"github.com/jvcoutinho/lit/validate"
)

type UploadFileRequest struct {
	Directory  string                `form:"directory"`
	FileHeader *multipart.FileHeader `file:"file"`
}

func (r *UploadFileRequest) Validate() []validate.Field {
	return []validate.Field{
		Directory(&r.Directory),
		validate.Required(&r.FileHeader),
	}
}

func UploadFile(r *lit.Request) lit.Response {
	req, err := bind.Body[UploadFileRequest](r)
	if err != nil {
		return render.BadRequest(err)
	}

	file, err := req.FileHeader.Open()
	if err != nil {
		return render.InternalServerError(err)
	}
	defer file.Close()

	destination, err := os.Create(path.Join(req.Directory, req.FileHeader.Filename))
	if err != nil {
		return render.InternalServerError(err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, file); err != nil {
		return render.InternalServerError(err)
	}

	return render.NoContent()
}

func Example_fileUpload() {
	r := lit.NewRouter()
	r.POST("/upload", UploadFile)

	f := createTemporaryFile()
	defer os.Remove(f.Name())

	body, contentType := createMultipartBody(f)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", contentType)

	r.ServeHTTP(res, req)

	fmt.Println(res.Code)
	// Output:
	// 204
}

func createMultipartBody(f *os.File) (io.Reader, string) {
	var (
		body   = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
	)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", f.Name())
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(part, f); err != nil {
		log.Fatal(err)
	}

	if err := writer.WriteField("directory", os.TempDir()); err != nil {
		log.Fatal(err)
	}

	return body, writer.FormDataContentType()
}

func createTemporaryFile() *os.File {
	f, err := os.CreateTemp("", "temporary_file")
	if err != nil {
		log.Fatal(err)
	}

	return f
}

// Directory validates that target is a valid directory.
func Directory(target *string) validate.Field {
	validation := validate.Field{
		Valid:   false,
		Message: "{0} should be a valid directory path",
		Fields:  []any{target},
	}

	if target == nil {
		return validation
	}

	fileInfo, err := os.Stat(*target)
	if err != nil {
		return validation
	}

	validation.Valid = fileInfo.IsDir()

	return validation
}
