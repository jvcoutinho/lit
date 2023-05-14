package render

import "net/http"

type CustomResponse struct {
	render func(http.ResponseWriter)
}

func Custom(render func(http.ResponseWriter)) *CustomResponse {
	return &CustomResponse{render}
}

func (r *CustomResponse) Render(writer http.ResponseWriter) {
	r.Render(writer)
}
