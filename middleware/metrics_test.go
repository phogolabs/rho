package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/go-chi/chi"
	"github.com/phogolabs/rest/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var _ = Describe("Metrics", func() {
	It("writes the metrics", func() {
		router := chi.NewMux()
		router.Use(middleware.Metrics)

		handler := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello")
		}

		router.Mount("/", http.HandlerFunc(handler))
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com/", nil))

		data, err := prometheus.DefaultGatherer.Gather()
		Expect(err).To(BeNil())
		Expect(data).NotTo(HaveLen(0))
	})
})
