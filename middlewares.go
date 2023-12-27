package lit

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// Recover is a simple middleware that recovers if h panics, responding a 500 Internal Server Error with
// the panic value as the body and logging the stack trace in os.Stderr.
func Recover(h Handler) Handler {
	return func(r *Request) (res Response) {
		defer func() {
			if value := recover(); value != nil {
				log.Printf("recovering from a panic: %v\n%s", value, debug.Stack())

				res = ResponseFunc(func(w http.ResponseWriter) {
					http.Error(w, fmt.Sprintf("%v", value), http.StatusInternalServerError)
				})
			}
		}()

		return h(r)
	}
}

// Log is a simple middleware that logs information data about the request:
//
//   - The method and path of the request;
//   - The status code of the response;
//   - The time of the request;
//   - The client's remote address;
//   - The duration of the request;
//   - The content length of the response body.
func Log(h Handler) Handler {
	return func(r *Request) Response {
		startTime := time.Now()

		res := h(r)

		return ResponseFunc(func(w http.ResponseWriter) {
			res.Write(w)

			recorder, ok := w.(*Recorder)
			if !ok {
				return
			}

			endTime := time.Now()

			var (
				statusCode      = recorder.StatusCode
				statusCodeText  = http.StatusText(statusCode)
				statusCodeColor = getColorFromStatusCode(statusCode)
				method          = r.Method()
				url             = r.URL()
				remoteAddress   = r.base.RemoteAddr
				duration        = endTime.Sub(startTime)
				responseSize    = recorder.ContentLength
			)

			message := &strings.Builder{}
			fmt.Fprintf(message, "\n%s>> %s %s\u001B[0m", statusCodeColor, method, url)
			fmt.Fprintf(message, "\n> %d %s", statusCode, statusCodeText)
			fmt.Fprintf(message, "\n> Start Time: %s", startTime.Format(time.DateTime))
			fmt.Fprintf(message, "\n> Remote Address: %s", remoteAddress)
			fmt.Fprintf(message, "\n> Duration: %s", duration)
			fmt.Fprintf(message, "\n> Content-Length: %d", responseSize)

			logger := log.New(log.Writer(), "", log.Flags()&^(log.Ldate|log.Ltime))
			logger.Println(message)
		})
	}
}

func getColorFromStatusCode(statusCode int) string {
	const (
		red    = "\u001B[97;1;41m"
		green  = "\u001B[97;1;42m"
		yellow = "\u001B[97;1;43m"
		blue   = "\u001B[97;1;104m"
	)

	if statusCode >= 200 && statusCode < 300 {
		return green
	}

	if statusCode >= 300 && statusCode < 400 {
		return blue
	}

	if statusCode >= 400 && statusCode < 500 {
		return yellow
	}

	return red
}
