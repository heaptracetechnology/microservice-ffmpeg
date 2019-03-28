package conversion

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	//"os"
)

type TestArgsData struct {
	VideoBase64 string `json:"video_base64"`
}

func Encode(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buff), nil
}

var _ = Describe("ffmpeg conversion video to image", func() {

	filepath := "../tmp/videos/sample.3gp"
	base64Data, base64Err := Encode(filepath)
	if base64Err != nil {
		fmt.Println("===base64 err======", base64Err)
	}

	testmessage := TestArgsData{VideoBase64: base64Data}

	requestBody := new(bytes.Buffer)
	errr := json.NewEncoder(requestBody).Encode(testmessage)
	if errr != nil {
		log.Fatal(errr)
	}

	req, err := http.NewRequest("POST", "/convertvideotoimage", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(VideoToImage)

	handler.ServeHTTP(recorder, req)

	Describe("conversion", func() {
		Context("video to image", func() {
			It("Should result http.StatusOK", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
