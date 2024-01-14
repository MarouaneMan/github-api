package middleware

import (
	"context"
	"github.com/MarouaneMan/github-api/kvstore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseCachingMiddlewareCacheMiss(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		responseCode int
		responseBody string
		shouldCache  bool
	}{
		{
			name:         "StatusCodeGreaterThan300",
			path:         "/?foo=bar",
			responseCode: http.StatusBadRequest, // 400
			responseBody: "bad boys",
			shouldCache:  false,
		},
		{
			name:         "StatusCode200",
			path:         "/?foo=baz",
			responseCode: http.StatusOK, // 200
			responseBody: "good boys",
			shouldCache:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			store := kvstore.NewInMemoryStore(kvstore.DefaultExpiration, kvstore.DefaultExpiration)
			middleware := NewResponseCachingMiddleware(store, store)

			// mock handler
			mockHandler := func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
				w.WriteHeader(tc.responseCode)
				w.Write([]byte(tc.responseBody))
				return nil
			}

			req, _ := http.NewRequest("GET", tc.path, nil)
			rr := httptest.NewRecorder()
			handler := middleware.Apply(mockHandler)
			handler.ServeHTTP(rr, req, nil)

			// check cache
			_, ok := store.Read(context.Background(), tc.path).(*cachedItem)
			if ok != tc.shouldCache {
				t.Errorf("unexpected caching behavior for %s: got %v want %v", tc.path, ok, tc.shouldCache)
			}
		})
	}
}

func TestResponseCachingMiddlewareCacheHit(t *testing.T) {
	testCases := []struct {
		path         string
		responseCode int
		responseBody string
	}{
		{
			path:         "/?foo=bar",
			responseCode: http.StatusOK,
			responseBody: "sweet bytes",
		},
		{
			path:         "/?foo=baz",
			responseCode: http.StatusOK,
			responseBody: "solid bytes",
		},
	}

	for _, tc := range testCases {
		t.Run("CacheHit_"+tc.path, func(t *testing.T) {

			store := kvstore.NewInMemoryStore(kvstore.DefaultExpiration, kvstore.DefaultExpiration)
			middleware := NewResponseCachingMiddleware(store, store)

			// make first call to populate cache
			{
				populateCacheHandler := func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
					w.WriteHeader(tc.responseCode)
					w.Write([]byte(tc.responseBody))
					return nil
				}
				req, _ := http.NewRequest("GET", tc.path, nil)
				rr := httptest.NewRecorder()
				handler := middleware.Apply(populateCacheHandler)
				handler.ServeHTTP(rr, req, nil)
			}

			// test cache hit
			{
				cacheHitHandler := func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
					t.Error("handler should not be called on cache hit")
					return nil
				}

				req, _ := http.NewRequest("GET", tc.path, nil)
				rr := httptest.NewRecorder()
				handler := middleware.Apply(cacheHitHandler)
				handler.ServeHTTP(rr, req, nil)

				// Check if the response is served from cache
				if rr.Body.String() != tc.responseBody || rr.Code != tc.responseCode {
					t.Errorf("expected body %q with status code %v; got body %q with status code %v",
						tc.responseBody, tc.responseCode, rr.Body.String(), rr.Code)
				}
			}
		})
	}
}
