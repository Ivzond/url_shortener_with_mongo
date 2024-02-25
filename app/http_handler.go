package app

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"io"
	"net/http"
	"net/url"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type URLRequest struct {
	URL     string `json:"url"`
	TTLDays int    `json:"ttlDays"`
}

func (h *Handler) Shorten(rw web.ResponseWriter, req *web.Request) (interface{}, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var rJson URLRequest
	if err := json.Unmarshal(body, &rJson); err != nil {
		return nil, err
	}
	if _, err := url.ParseRequestURI(rJson.URL); err != nil {
		return nil, ErrBadRequest
	}
	return h.service.Shorten(req.Context(), rJson.URL, rJson.TTLDays)
}

func (h *Handler) Update(rw web.ResponseWriter, req *web.Request) (interface{}, error) {
	id := req.PathParams["shortUrl"]
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var rJson URLRequest
	if err := json.Unmarshal(body, &rJson); err != nil {
		return nil, err
	}

	if _, err := url.ParseRequestURI(rJson.URL); err != nil {
		return nil, ErrBadRequest
	}

	return h.service.Update(req.Context(), id, rJson.URL, rJson.TTLDays)
}

func (h *Handler) Delete(rw web.ResponseWriter, req *web.Request) (interface{}, error) {
	id := req.PathParams["shortUrl"]
	return nil, h.service.Delete(req.Context(), id)
}

func (h *Handler) GetFullURL(rw web.ResponseWriter, req *web.Request) (interface{}, error) {
	val := req.PathParams["shortUrl"]
	return h.service.GetFullURL(req.Context(), val)
}

func (h *Handler) Ping(rw web.ResponseWriter, req *web.Request) (interface{}, error) {
	return nil, nil
}

type EndpointHandler func(rw web.ResponseWriter, req *web.Request) (interface{}, error)

func WrapEndpoint(h EndpointHandler) interface{} {
	fn := func(rw web.ResponseWriter, req *web.Request, h EndpointHandler) error {
		result, err := h(rw, req)
		if err != nil {
			return err
		}

		data, err := json.Marshal(result)
		if err != nil {
			return err
		}
		_, err = rw.Write(data)
		return err
	}
	return func(rw web.ResponseWriter, req *web.Request) {
		err := fn(rw, req, h)
		if err != nil {
			fmt.Println(err.Error())
			writeHttpCode(rw, err)
		}
	}
}

func writeHttpCode(rw http.ResponseWriter, err error) {
	switch err {
	case ErrNotFound:
		rw.WriteHeader(http.StatusNotFound)
	case ErrBadRequest:
		rw.WriteHeader(http.StatusBadRequest)
	default:
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
