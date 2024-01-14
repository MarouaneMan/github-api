package middleware

import (
	"bytes"
	"github.com/MarouaneMan/github-api/kvstore"
	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	"net/http"
)

type responseCachingMiddleware struct {
	storeReader kvstore.Reader
	storeWriter kvstore.Writer
}

type cachedItem struct {
	statusCode int
	body       []byte
	headers    map[string][]string
}

// NewResponseCachingMiddleware creates and returns a new response caching middleware.
// This middleware uses the provided storeReader and storeWriter to cache responses.
func NewResponseCachingMiddleware(storeReader kvstore.Reader, storeWriter kvstore.Writer) handlers.Middleware {
	return &responseCachingMiddleware{
		storeReader: storeReader,
		storeWriter: storeWriter,
	}
}

// Apply is the middleware handler that caches responses based on their status code.
// It serves cached responses if available and caches new responses if the status code is below 300.
func (rcm *responseCachingMiddleware) Apply(next handlers.HandlerFunc) handlers.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		log := logger.Get(r.Context())

		// serve cached response if available
		cachedResponse, ok := rcm.storeReader.Read(r.Context(), r.URL.String()).(*cachedItem)
		if ok {
			log.Debug("Cache hit")
			for key, values := range cachedResponse.headers {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.Header().Add("X-From-Cache", "True")
			w.WriteHeader(cachedResponse.statusCode)
			_, err := w.Write(cachedResponse.body)
			return err
		}

		log.Debug("Cache miss")

		// create a custom writer to capture the response
		customWriter := customResponseWriter{ResponseWriter: w}

		// process the request
		err := next(&customWriter, r, vars)

		// cache the response if status code is below 300
		if err == nil && customWriter.statusCode < 300 {
			// cache response: statusCode, body and headers
			cacheErr := rcm.storeWriter.Write(r.Context(), r.URL.String(), &cachedItem{
				statusCode: customWriter.statusCode,
				body:       customWriter.body.Bytes(),
				headers:    customWriter.Header().Clone(),
			}, kvstore.DefaultExpiration)
			if cacheErr != nil {
				log.WithError(cacheErr).Error("Failed to cache response")
			}
		}
		return err
	}
}

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

func (crw *customResponseWriter) Write(b []byte) (int, error) {
	crw.body.Write(b)
	return crw.ResponseWriter.Write(b)
}
