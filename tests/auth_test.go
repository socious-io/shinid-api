package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"shin/src/app/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Group", func() {
	var (
		router *gin.Engine
		db     *sqlx.DB
	)

	BeforeEach(func() {
		db, router = setupTestEnvironment()
	})

	Context("POST /auth/register", func() {
		It("should return status 200 with jwt tokens", func() {
			register(router)
		})
	})
	Context("POST /auth/otp", func() {
		It("should return status 200", func() {
			sendOTP(router)
		})
	})
	Context("POST /auth/otp/resend", func() {
		It("should return status 200", func() {
			resendOTP(router)
		})
	})
	Context("POST /auth/otp/verify", func() {
		It("should return status 200", func() {
			verifyOTP(db, router)
		})
	})

	AfterEach(func() {
		// teardownTestEnvironment(db)
	})
})

func register(r *gin.Engine) {
	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(users_data[0])
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body := decodeBody(w.Body)

	Expect(w.Code).To(Equal(200))
	bodyExpect(body, gin.H{"access_token": "<ANY>"})
	users_data[0]["access_token"] = body["access_token"]
}

func sendOTP(r *gin.Engine) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/otp", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", users_data[0]["access_token"]))
	r.ServeHTTP(w, req)

	Expect(w.Code).To(Equal(200))
}

func resendOTP(r *gin.Engine) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/otp/resend", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", users_data[0]["access_token"]))
	r.ServeHTTP(w, req)

	Expect(w.Code).To(Equal(200))
}

func verifyOTP(db *sqlx.DB, r *gin.Engine) {

	//Get OTP
	otp := new(models.OTP)
	db.Get(otp, "SELECT * FROM otps LIMIT 1")

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(gin.H{"code": otp.Code})
	req, _ := http.NewRequest("POST", "/auth/otp/verify", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", users_data[0]["access_token"]))
	r.ServeHTTP(w, req)

	Expect(w.Code).To(Equal(200))
}
