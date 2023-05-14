package lit

import "net/http"

// Response is the output of a Lit handler function.
type Response func(writer http.ResponseWriter)
