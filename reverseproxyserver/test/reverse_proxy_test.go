package test

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"reverseproxy/src/handler"
	"testing"
)

func TestMain(m *testing.M) {
	LoadRedisMock()
	LoadReverseProxyConfiguration()
	m.Run()
}

func TestGivenRunningBackendWhenCallingReverseProxyThenMustReturnBackendResponse(t *testing.T) {
	StartBackend()

	req := httptest.NewRequest(http.MethodGet, "/backend/v1/ping", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ReverseProxyHandler(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "pong", string(data))
}

func TestGivenNoRunningBackendWhenCallingReverseProxyThenMustReturnInternalServerError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/backend/v1/ping", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ReverseProxyHandler(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, "error sending request\n", string(data))
}

func TestGivenRunningBackendWhenCallingForTheSecondTimeTheReverseProxyThenMustReturnCachedResponse(t *testing.T) {
	StartBackend()

	req := httptest.NewRequest(http.MethodGet, "/backend/v1/ping", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ReverseProxyHandler(responseRecorder, req)

	req = httptest.NewRequest(http.MethodGet, "/backend/v1/ping", nil)
	responseRecorder = httptest.NewRecorder()
	handler.ReverseProxyHandler(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "true", res.Header.Get("Cached"))
	assert.Equal(t, "pong", string(data))
}

func TestGivenRunningBackendWhenCallingWrongMethodInReverseProxyThenMustReturnError404(t *testing.T) {
	StartBackend()

	req := httptest.NewRequest(http.MethodGet, "/backend/v1/pong", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ReverseProxyHandler(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGivenRunningBackendWhenCallingNotMappedBackendInReverseProxyThenMustReturnError404(t *testing.T) {
	StartBackend()

	req := httptest.NewRequest(http.MethodGet, "/backend/v2/ping", nil)
	responseRecorder := httptest.NewRecorder()

	handler.ReverseProxyHandler(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGivenRunningBackendWhenCallingPostMethodForTheSecondTimeTheReverseProxyThenMustNotReturnCachedResponse(t *testing.T) {
	StartBackend()

	req := httptest.NewRequest(http.MethodPost, "/backend/v1/ping", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ReverseProxyHandler(responseRecorder, req)

	req = httptest.NewRequest(http.MethodPost, "/backend/v1/ping", nil)
	responseRecorder = httptest.NewRecorder()
	handler.ReverseProxyHandler(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Empty(t, res.Header.Get("Cached"))
	assert.Equal(t, "pong", string(data))
}

func TestGivenRunningBackendWhenCallingForTheSecondTimeTheReverseProxyWithNoCacheInHeaderThenMustNotReturnCachedResponse(t *testing.T) {
	StartBackend()

	req := httptest.NewRequest(http.MethodGet, "/backend/v1/ping", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ReverseProxyHandler(responseRecorder, req)

	req = httptest.NewRequest(http.MethodGet, "/backend/v1/ping", nil)
	req.Header.Add("Cache-Control", "no-store")
	responseRecorder = httptest.NewRecorder()
	handler.ReverseProxyHandler(responseRecorder, req)

	res := responseRecorder.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Empty(t, res.Header.Get("Cached"))
	assert.Equal(t, "pong", string(data))
}
