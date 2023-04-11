package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools
	s := testTools.RandomString(10)
	println(len(s), s)
	if len(s) != 10 {
		t.Error("Wrong lenght random string returned")
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{name: "allowed no rename", allowedTypes: []string{"image/jpeg", "image/png"}, renameFile: false, errorExpected: false},
	{name: "allowed rename", allowedTypes: []string{"image/jpeg", "image/png"}, renameFile: true, errorExpected: false},
	{name: "not allowed", allowedTypes: []string{"image/jpeg"}, renameFile: false, errorExpected: true},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, e := range uploadTests {
		// setup a pip to aviod buffering
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			// create form data field "file"
			part, err := writer.CreateFormFile("file", "./testdata/img.png")
			if err != nil {
				t.Error(err)
			}
			f, err := os.Open("./testdata/img.png")
			if err != nil {
				t.Error(err)
			}

			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("Error Decoding Image", err)
			}
			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}
		}()
		// read from the pipe which receives data
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testtools Tools
		testtools.AllowedFileTypes = e.allowedTypes
		uploadedFile, err := testtools.UploadFiles(request, "./testdata/uploads/", e.renameFile)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}
		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFile[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file to exists : %s", e.name, err.Error())
			}

			// clean up
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFile[0].NewFileName))
		}
		if !e.errorExpected && err != nil {
			t.Errorf("%s :error expected but none received", e.name)
		}
		wg.Wait()

	}
}
func TestTools_UploadOneFile(t *testing.T) {
	// setup a pip to aviod buffering
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		// create form data field "file"
		part, err := writer.CreateFormFile("file", "./testdata/img.png")
		if err != nil {
			t.Error(err)
		}
		f, err := os.Open("./testdata/img.png")
		if err != nil {
			t.Error(err)
		}

		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			t.Error("Error Decoding Image", err)
		}
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()
	// read from the pipe which receives data
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	var testtools Tools

	uploadedFile, err := testtools.UploadOneFile(request, "./testdata/uploads/", true)
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName)); os.IsNotExist(err) {
		t.Errorf("expected file to exists : %s", err.Error())
	}

	// clean up
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFile.NewFileName))

}
func TestTools_CreateDirIfNotExist(t *testing.T) {
	var testTools Tools

	err := testTools.CreateDirIfNotExist("./testdata/myDir")
	if err != nil {
		t.Error(err)
	}

	err = testTools.CreateDirIfNotExist("./testdata/myDir")
	if err != nil {
		t.Error(err)
	}
	_ = os.Remove("./testdata/myDir")
}

var slugTest = []struct {
	name          string
	s             string
	expected      string
	errorExpected bool
}{
	{name: "valid string", s: "now is the time", expected: "now-is-the-time", errorExpected: false},
	{name: "empty string", s: "", expected: "", errorExpected: true},
	{name: "complex string", s: "Now is the time for all GOOD men! + fish & such &#^123", expected: "now-is-the-time-for-all-good-men-fish-such-123", errorExpected: true},
}

func TestTools_Slugify(t *testing.T) {
	var testTool Tools
	for _, e := range slugTest {
		slug, err := testTool.Slugify(e.s)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: error received when none expected: %s", e.name, err.Error())
		}
		if !e.errorExpected && slug != e.expected {
			t.Errorf("%s: wrong slug returned; expected %s but got %s", e.name, e.expected, slug)
		}
	}
}
