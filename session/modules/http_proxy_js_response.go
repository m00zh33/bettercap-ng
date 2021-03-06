package session_modules

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

type JSResponse struct {
	Status      int
	ContentType string
	Headers     string
	Body        string

	wasUpdated bool
	resp       *http.Response
}

func NewJSResponse(res *http.Response) *JSResponse {
	cType := ""
	headers := ""

	for name, values := range res.Header {
		for _, value := range values {
			if name == "Content-Type" {
				cType = value
			}
			headers += name + ": " + value + "\r\n"
		}
	}

	return &JSResponse{
		Status:      res.StatusCode,
		ContentType: cType,
		Headers:     headers,
		resp:        res,
	}
}

func (j *JSResponse) Updated() {
	j.wasUpdated = true
}

func (j *JSResponse) ToResponse(req *http.Request) (resp *http.Response) {
	resp = goproxy.NewResponse(req, j.ContentType, j.Status, j.Body)
	if j.Headers != "" {
		for _, header := range strings.Split(j.Headers, "\n") {
			header = strings.Trim(header, "\n\r\t ")
			if header == "" {
				continue
			}
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				resp.Header.Add(parts[0], parts[1])
			} else {
				log.Warningf("Unexpected header '%s'", header)
			}
		}
	}
	return
}

func (j *JSResponse) ReadBody() string {
	defer j.resp.Body.Close()

	raw, err := ioutil.ReadAll(j.resp.Body)
	if err != nil {
		log.Errorf("Could not read response body: %s", err)
		return ""
	}

	j.Body = string(raw)

	return j.Body
}
