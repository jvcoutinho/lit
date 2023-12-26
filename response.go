package lit

import (
	"net/http"
)

// Response is the output of a [Handler].
type Response interface {
	// Write responds the request.
	Write(w http.ResponseWriter)
}

// ResponseFunc writes response data into [http.ResponseWriter].
type ResponseFunc func(w http.ResponseWriter)

func (r ResponseFunc) Write(w http.ResponseWriter) {
	r(w)
}

// Recorder is a http.ResponseWriter that keeps track of the response' status code and content length.
type Recorder struct {
	written bool

	// StatusCode of this response.
	StatusCode int

	// Size of this response.
	ContentLength int

	http.ResponseWriter
	http.Hijacker
	http.Flusher
}

// NewRecorder creates a new *Recorder instance from a http.ResponseWriter.
func NewRecorder(w http.ResponseWriter) *Recorder {
	hijacker, _ := w.(http.Hijacker)
	flusher, _ := w.(http.Flusher)

	return &Recorder{
		ResponseWriter: w,
		Hijacker:       hijacker,
		Flusher:        flusher,
		StatusCode:     http.StatusOK,
		ContentLength:  0,
		written:        false,
	}
}

func (r *Recorder) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)

	if statusCode >= 200 && statusCode <= 599 && !r.written {
		r.StatusCode = statusCode
		r.written = true
	}
}

func (r *Recorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.ContentLength += n
	return n, err
}
