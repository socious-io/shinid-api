package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func authGroup() {
	It("should return status 200 with jwt tokens", func() {
		w := httptest.NewRecorder()
		reqBody, _ := json.Marshal(usersData[0])
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(200))
		bodyExpect(body, gin.H{"access_token": "<ANY>"})
		authTokens = append(authTokens, body["access_token"].(string))
	})
}
