package tests_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"shin/src/shortner"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func shortnerGroup() {
	var shortLinks []*shortner.ShortnerURL
	It("should create shortner", func() {
		for _, url := range shorteningURLs {
			s := &shortner.ShortnerURL{LongURL: url}
			if err := s.Create(); err != nil {
				continue
			}
			shortLinks = append(shortLinks, s)
		}
	})

	It("should fetch shortner", func() {
		for _, short := range shortLinks {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/fetch", short.ShortID), nil)
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(200))
			Expect(body["long_url"].(string)).To(Equal(short.LongURL))
		}
	})

	It("should  fetch and redirect shortner", func() {
		for _, short := range shortLinks {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", short.ShortID), nil)
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(308))
		}
	})
}
