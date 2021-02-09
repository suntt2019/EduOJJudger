package web_test

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/suntt2019/EduOJJudger/base"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func R() *resty.Request {
	return base.HC.R()
}

func checkFile(t *testing.T, path, content string) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	assert.Nil(t, err)
	b := make([]byte, len(content))
	_, err = file.Read(b)
	assert.Nil(t, err)
	assert.Equal(t, content, string(b))
}

func testServerRoute(wr http.ResponseWriter, r *http.Request) {
	index := strings.Index(r.RequestURI[1:], "/")
	if index == -1 {
		panic("could not find the second '/' to find out service name")
	}
	switch r.RequestURI[1 : index+1] {
	case "echoURI":
		echoURI(wr, r)
	case "backendErr":
		backendError(wr, r)
	case "fileURI":
		fileURI(wr, r)
	case "script":
		script(wr, r)

	default:
		panic(`invalid service for test server: "` + r.RequestURI[1:index+1] + `"`)
	}

}

func echoURI(wr http.ResponseWriter, r *http.Request) {
	if _, err := wr.Write([]byte(r.RequestURI[9:])); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	config := ``
	if err := viper.ReadConfig(strings.NewReader(config)); err != nil {
		panic(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(testServerRoute))
	base.HC = resty.New().SetHostURL(ts.URL)
	base.ScriptPath = "../test_file/scripts"
	base.RunPath = "../test_file/runs"
	if err := base.RemoveBuffer(); err != nil {
		panic(err)
	}
	ret := m.Run()
	if err := base.RemoveBuffer(); err != nil {
		panic(err)
	}
	os.Exit(ret)
}
