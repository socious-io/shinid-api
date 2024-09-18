package tests_test

import (
	"fmt"
	"log"
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
			s, err := shortner.New(url)
			if err != nil {
				log.Fatalf("Shortner error : %v", err)
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
