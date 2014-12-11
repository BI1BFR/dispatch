package dispatch

// import "net/http"

// type RemoteRequestResolver interface {
// 	Resolve(r *http.Request) Request
// }

// type RemoteResponseWriter interface {
// 	Write(w http.ResponseWriter, r Response)
// }

// type RemoteRequest struct {
// }

// type RemoteDest struct {
// 	addr string
// }

// func (d *RemoteDest) Call(ctx *Context, r Request) Response {
// 	// var c http.Client
// 	// req, err := http.NewRequest("PUT", d.addr+r.Protocol(), bytes.NewBuffer(r.Body().Bytes()))
// 	// rsp, err := c.Do(req)
// 	// if err != nil {
// 	// 	return newErrResponseWithText(http.StatusServiceUnavailable, err.Error())
// 	// }
// 	return nil
// }
