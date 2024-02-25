package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
	"time"
)

const contentType = "application/json"

func TestServiceSuite(t *testing.T) {
	suite.Run(t, &ServiceSuite{
		ctx: context.Background(),
	})
}

type ServiceSuite struct {
	suite.Suite
	ctx        context.Context
	httpClient *http.Client
}

func (s *ServiceSuite) SetupSuite() {
	s.httpClient = &http.Client{
		Timeout: time.Second,
	}
	go func() {
		if err := Run(s.ctx); err != nil {
			s.Require().Error(err)
		}
	}()

	for it := 0; it < 5; it++ {
		time.Sleep(time.Millisecond * 100)
		_, err := s.httpClient.Get(makeURL("/ping"))
		if err == nil {
			break
		}
	}
}

func (s *ServiceSuite) Test() {
	var shortID string
	s.Run("Shorten", func() {
		Body := s.buildBody("https://google.com", 1)
		resp, err := s.httpClient.Post(makeURL("/shorten"), contentType, Body)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		data, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)
		var shortURL ShortURL
		s.Require().NoError(json.Unmarshal(data, &shortURL))

		s.Require().Equal("https://google.com", shortURL.URL)
		timeDiff := shortURL.ExpireAt.Sub(time.Now())
		s.Require().InDelta(time.Hour*24, timeDiff, float64(time.Minute))

		shortID = shortURL.ID
	})

	s.Run("ResolveShortLink", func() {
		resp, err := s.httpClient.Get(makeURL(fmt.Sprintf("/%s", shortID)))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)

		s.Require().Equal(`"https://google.com"`, string(body))
	})

	s.Run("UpdateLink", func() {
		body := s.buildBody("http://ya.ru", 0)
		resp, err := s.httpClient.Post(makeURL(fmt.Sprintf("/update/%s", shortID)), contentType, body)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		data, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)
		var shortURL ShortURL
		s.Require().NoError(json.Unmarshal(data, &shortURL))

		s.Require().Equal("http://ya.ru", shortURL.URL)
		s.Require().Empty(shortURL.ExpireAt)
	})

	s.Run("Delete", func() {
		req, err := http.NewRequest(http.MethodDelete, makeURL(fmt.Sprintf("/%s", shortID)), nil)
		s.Require().NoError(err)

		resp, err := s.httpClient.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		resp, err = s.httpClient.Get(makeURL(fmt.Sprintf("/%s", shortID)))
		s.Require().NoError(err)
		s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	})
}

func (s *ServiceSuite) buildBody(url string, ttlDays int) io.Reader {
	data, err := json.Marshal(&URLRequest{
		URL:     url,
		TTLDays: ttlDays,
	})
	s.Require().NoError(err)
	return bytes.NewReader(data)
}

func makeURL(path string) string {
	return "http://localhost:8000" + path
}
