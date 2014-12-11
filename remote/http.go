package remote

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/huangml/dispatch/msg"
)

type HTTPResolver interface {
	ResolveRequest(r *http.Request) msg.Request
	ResolveResponse(r *http.Response) msg.Response
	BuildRequest(r msg.Request, url, method string) *http.Request
	WriteResponse(w http.ResponseWriter, r msg.Response)
}

type defaultResolver struct{}

func (d defaultResolver) ResolveRequest(r *http.Request) msg.Request {
	return DefaultRequestResolveFunc(r)
}

func (d defaultResolver) ResolveResponse(r *http.Response) msg.Response {
	return DefaultResponseResolveFunc(r)
}

func (d defaultResolver) BuildRequest(r msg.Request, url, method string) *http.Request {
	return DefaultRequestBuildFunc(r, url, method)
}

func (d defaultResolver) WriteResponse(w http.ResponseWriter, r msg.Response) {
	DefaultResponseWriteFunc(w, r)
}

var DefaultHTTPResolver HTTPResolver = &defaultResolver{}

func DefaultRequestResolveFunc(r *http.Request) msg.Request {
	return &Request{
		Proto:   r.RequestURI,
		Sink:    parseSink(r.Body, r.Header),
		Auth:    parseAuth(r),
		TimeOut: parseTimeOut(r),
	}
}

func DefaultResponseResolveFunc(r *http.Response) msg.Response {
	var err *statusError
	if r.StatusCode != http.StatusOK {
		e := Err(r.StatusCode)
		err = &e
	}
	s := parseSink(r.Body, r.Header)
	return &Response{
		statusError: err,
		Sink:        s,
	}
}

func DefaultRequestBuildFunc(r msg.Request, url, method string) *http.Request {
	sink := r.Body()
	if sink == nil {
		return nil
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(sink.Bytes()))
	if err != nil {
		return nil
	}

	req.Header.Set(ContentTypeKey, contentToHTTP(sink.ContentType))
	if r, ok := r.(*Request); ok {
		if r.Auth != nil {
			req.SetBasicAuth(r.Auth.UserName, r.Auth.Password)
		}
		if r.TimeOut > 0 {
			req.Header.Set(TimeOutKey, r.TimeOut.String())
		}
	}

	return req
}

func DefaultResponseWriteFunc(w http.ResponseWriter, r msg.Response) {
	if sink := r.Body(); sink != nil {
		defer w.Write(sink.Bytes())
		w.Header().Set(ContentTypeKey, contentToHTTP(sink.ContentType))
	}

	statusCode := http.StatusOK
	if r.Error() != nil {
		if r, ok := r.(*Response); ok {
			statusCode = r.statusCode
		} else {
			statusCode = http.StatusInternalServerError
		}
	}

	w.WriteHeader(statusCode)
}

func parseSink(body io.ReadCloser, header http.Header) *msg.Sink {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}
	c := contentTypeFromHTTP(header.Get(ContentTypeKey))
	s := msg.NewBytesSink(b)
	s.ContentType = c
	return s
}

const (
	OctetStream = "application/octet-stream"
	XProtoBuf   = "application/x-protobuf"
	TextPlain   = "text/plain"

	TimeOutKey     = "X-Dispatch-Timeout"
	ContentTypeKey = "Content-Type"
)

func contentTypeFromHTTP(v string) msg.ContentType {
	switch v {
	case OctetStream:
		return msg.Bytes
	case TextPlain:
		return msg.Text
	case XProtoBuf:
		return msg.ProtoBuf
	default:
		return msg.Bytes
	}
}

func contentToHTTP(c msg.ContentType) string {
	switch c {
	case msg.Bytes:
		return OctetStream
	case msg.Text:
		return TextPlain
	case msg.ProtoBuf:
		return XProtoBuf
	default:
		return OctetStream
	}
}

func parseAuth(r *http.Request) *Auth {
	if username, password, ok := BasicAuth(r); ok {
		return &Auth{
			UserName: username,
			Password: password,
		}
	}
	return nil
}

func parseTimeOut(r *http.Request) time.Duration {
	if t, err := time.ParseDuration(r.Header.Get(TimeOutKey)); err == nil {
		return t
	}

	return 0
}
