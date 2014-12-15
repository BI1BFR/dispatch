package dispatch

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type RemoveResolver struct {
	RequestResolver  func(*http.Request) Request
	ResponseResolver func(r *http.Response) Response
	RequestBuilder   func(r Request, url, method string) *http.Request
	ResponseWriter   func(w http.ResponseWriter, r Response)
}

var DefaultRemoteResolver = RemoveResolver{
	RequestResolver:  DefaultRemoteRequestResolveFunc,
	ResponseResolver: DefaultRemoteResponseResolveFunc,
	RequestBuilder:   DefaultRemoteRequestBuildFunc,
	ResponseWriter:   DefaultRemoteResponseWriteFunc,
}

func DefaultRemoteRequestResolveFunc(r *http.Request) Request {
	return &RemoteRequest{
		Proto:   r.RequestURI,
		Sink:    parseSinkFromHTTP(r.Body, r.Header),
		Auth:    parseAuthFromHTTP(r),
		TimeOut: parseTimeOutFromHTTP(r),
	}
}

func DefaultRemoteResponseResolveFunc(r *http.Response) Response {
	var err *statusError
	if r.StatusCode != http.StatusOK {
		e := Err(r.StatusCode)
		err = &e
	}
	return &RemoteResponse{
		statusError: err,
		Sink:        parseSinkFromHTTP(r.Body, r.Header),
	}
}

func DefaultRemoteRequestBuildFunc(r Request, url, method string) *http.Request {
	sink := r.Body()
	if sink == nil {
		return nil
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(sink.Bytes()))
	if err != nil {
		return nil
	}

	req.Header.Set(ContentTypeKey, contentTypeToHTTP(sink.ContentType))
	if r, ok := r.(*RemoteRequest); ok {
		if r.Auth != nil {
			req.SetBasicAuth(r.Auth.UserName, r.Auth.Password)
		}
		if r.TimeOut > 0 {
			req.Header.Set(TimeOutKey, r.TimeOut.String())
		}
	}

	return req
}

func DefaultRemoteResponseWriteFunc(w http.ResponseWriter, r Response) {
	if sink := r.Body(); sink != nil {
		defer w.Write(sink.Bytes())
		w.Header().Set(ContentTypeKey, contentTypeToHTTP(sink.ContentType))
	}

	statusCode := http.StatusOK
	if r.Error() != nil {
		if r, ok := r.(*RemoteResponse); ok {
			statusCode = r.statusCode
		} else {
			statusCode = http.StatusInternalServerError
		}
	}

	w.WriteHeader(statusCode)
}

func parseSinkFromHTTP(body io.ReadCloser, header http.Header) *Sink {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}
	c := contentTypeFromHTTP(header.Get(ContentTypeKey))
	s := NewBytesSink(b)
	s.ContentType = c
	return s
}

func parseAuthFromHTTP(r *http.Request) *Auth {
	if username, password, ok := r.BasicAuth(); ok {
		return &Auth{
			UserName: username,
			Password: password,
		}
	}
	return nil
}

func parseTimeOutFromHTTP(r *http.Request) time.Duration {
	if t, err := time.ParseDuration(r.Header.Get(TimeOutKey)); err == nil {
		return t
	}

	return 0
}

const (
	OctetStream = "application/octet-stream"
	XProtoBuf   = "application/x-protobuf"
	TextPlain   = "text/plain"

	TimeOutKey     = "X-Dispatch-Timeout"
	ContentTypeKey = "Content-Type"
)

func contentTypeFromHTTP(v string) ContentType {
	switch v {
	case OctetStream:
		return Bytes
	case TextPlain:
		return Text
	case XProtoBuf:
		return ProtoBuf
	default:
		return Bytes
	}
}

func contentTypeToHTTP(c ContentType) string {
	switch c {
	case Bytes:
		return OctetStream
	case Text:
		return TextPlain
	case ProtoBuf:
		return XProtoBuf
	default:
		return OctetStream
	}
}
