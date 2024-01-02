// Package render contains implementations of [lit.Response], suitable for responding requests.
//
// # JSON responses
//
// For simple JSON responses, one can use the function [JSON] or their "shortcuts", such as [OK], [BadRequest] or
// [Unauthorized]. All of them are constructors for [JSONResponse].
//
// # 204 No Content responses
//
// In a special case, the constructor [NoContent] returns a [NoContentResponse], that does not contain a body.
//
// # Redirections
//
// The constructor [Redirect] can be used to redirect a request. If one wants more granularity, the specific redirection
// functions are also available, such as [Found] or [PermanentRedirect].
//
// # Serving files
//
// The constructor [File] can be used to serve files. It uses internally the [http.ServeFile] function.
//
// # Streaming
//
// The constructor [Stream] can be used to serve streams. It uses internally the [http.ServeContent] function.
//
// # Custom responses
//
// In order to create new responses not mapped in this package (or to use Facades), such as a YAML response or a new
// way of serving streams, one can:
//   - Create a type that implements the [lit.Response] interface and use it as the return of a handler;
//   - Return a [lit.ResponseFunc], writing to the [http.ResponseWriter] directly.
package render
