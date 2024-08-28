package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func orgGroup() {
	if focused {
		authGroup()
	}

	It("should create organizations", func() {
		for i, data := range organizationsData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/organizations", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(201))
			organizationsData[i]["id"] = body["id"]
		}
	})

	It("should get organizations", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/organizations", nil)
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(200))
		Expect(len(body["results"].([]interface{}))).To(Equal(1))
	})

	It("should get organization", func() {
		for _, data := range organizationsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/organizations/%s", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(200))
			Expect(body["id"]).To(Equal(data["id"]))
		}
	})
}
