package remote

import (
	"net/http"

	"github.com/huangml/dispatch/msg"
)

type Dest struct {
	addr string
	s    HTTPResolver
}

func (d *Dest) Call(ctx *msg.Context, r msg.Request) msg.Response {
	req := d.s.BuildRequest(r, d.addr+r.Protocol(), "POST")
	var c http.Client
	rsp, _ := c.Do(req)
	return d.s.ResolveResponse(rsp)
}
