package mocker

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MethodEndpoint represents a single endpoint
type MethodEndpoint struct {
	Preset    string      `yaml:"preset"` // random, ratio
	Responses []*Response `yaml:"responses"`
	Headers   []Header    `yaml:"headers"`
}

// PickResponse picks a random response according to the ratio defined in the
// responses
func (e MethodEndpoint) PickResponse() *Response {
	var sum int
	var out *Response
	for _, c := range e.Responses {
		sum += c.Ratio
	}
	o := rand.Intn(sum)
	for _, r := range e.Responses {
		o -= r.Ratio
		if o < 0 {
			out = r
			break
		}
	}
	return out
}

// CalcRatios computes the ratios for every response that doesn't have one
func (e *MethodEndpoint) CalcRatios() {
	var tot int
	var allocated int
	for _, r := range e.Responses {
		if r.Ratio != 0 {
			allocated += r.Ratio
		} else {
			tot++
		}
	}
	remaining := 100 - allocated
	o := float64(remaining) / float64(tot)
	for _, r := range e.Responses {
		if r.Ratio == 0 {
			r.Ratio = int(o)
		}
	}
}

// ToHandler generates a handler to apply on the router
func (e MethodEndpoint) ToHandler() func(c *gin.Context) {
	e.CalcRatios()
	return func(c *gin.Context) {
		for _, h := range e.Headers {
			if h.Required && c.GetHeader(h.Name) == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "header is required", "missing": h.Name})
				return
			}
		}
		if e.Responses == nil {
			c.Status(http.StatusOK)
			return
		}

		r := e.PickResponse()
		switch r.Preset {
		case "json":
			c.Header("Content-Type", "application/json; charset=utf-8")
		}
		if r.Body != "" {
			c.String(r.Code, r.Body)
		} else {
			c.Status(r.Code)
		}
	}
}

// EndpointGenerator is an interface that allows to generate endpoints
type EndpointGenerator interface {
	Generate(string, *gin.Engine)
}

// GetEndpoint implements the EndpointGenerator interface
type GetEndpoint struct {
	MethodEndpoint `yaml:",inline"`
}

// Generate generates the endpoint
func (e GetEndpoint) Generate(path string, r *gin.Engine) { r.GET(path, e.ToHandler()) }

// PostEndpoint implements the EndpointGenerator interface
type PostEndpoint struct {
	MethodEndpoint `yaml:",inline"`
}

// Generate generates the endpoint
func (e PostEndpoint) Generate(path string, r *gin.Engine) { r.POST(path, e.ToHandler()) }

// PutEndpoint implements the EndpointGenerator interface
type PutEndpoint struct {
	MethodEndpoint `yaml:",inline"`
}

// Generate generates the endpoint
func (e PutEndpoint) Generate(path string, r *gin.Engine) { r.PUT(path, e.ToHandler()) }

// PatchEndpoint implements the EndpointGenerator interface
type PatchEndpoint struct {
	MethodEndpoint `yaml:",inline"`
}

// Generate generates the endpoint
func (e PatchEndpoint) Generate(path string, r *gin.Engine) { r.PATCH(path, e.ToHandler()) }

// DeleteEndpoint implements the EndpointGenerator interface
type DeleteEndpoint struct {
	MethodEndpoint `yaml:",inline"`
}

// Generate generates the endpoint
func (e DeleteEndpoint) Generate(path string, r *gin.Engine) { r.DELETE(path, e.ToHandler()) }

// HeadEndpoint implements the EndpointGenerator interface
type HeadEndpoint struct {
	MethodEndpoint `yaml:",inline"`
}

// Generate generates the endpoint
func (e HeadEndpoint) Generate(path string, r *gin.Engine) { r.HEAD(path, e.ToHandler()) }

// OptionsEndpoint implements the EndpointGenerator interface
type OptionsEndpoint struct {
	MethodEndpoint `yaml:",inline"`
}

// Generate generates the endpoint
func (e OptionsEndpoint) Generate(path string, r *gin.Engine) { r.OPTIONS(path, e.ToHandler()) }