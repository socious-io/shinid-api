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

func schemaGroup() {
	if focused {
		authGroup()
	}
	It("it should create schema", func() {
		for i, data := range schemasData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/credentials/schemas", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			schemasData[i]["id"] = body["id"]
			Expect(w.Code).To(Equal(201))
		}
	})

	It("it should get schema", func() {
		for _, data := range schemasData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/credentials/schemas/%s", data["id"]), nil)
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			fmt.Println(body, "-----@@@@@------")
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(w.Code).To(Equal(200))
		}
	})
}
