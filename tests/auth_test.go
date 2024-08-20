package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shin/src/app/models"

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

	It("Should return status 200", func() {
		w := httptest.NewRecorder()
		reqBody, _ := json.Marshal(gin.H{"email": usersData[0]["email"]})
		req, _ := http.NewRequest("POST", "/auth/otp", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(200))
	})

	It("Should return status 200", func() {
		//Get OTP
		otp := new(models.OTP)
		db.Get(otp, "SELECT * FROM otps LIMIT 1")

		w := httptest.NewRecorder()
		reqBody, _ := json.Marshal(gin.H{"email": usersData[0]["email"], "code": otp.Code})
		req, _ := http.NewRequest("POST", "/auth/otp/verify", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(200))
	})
}
