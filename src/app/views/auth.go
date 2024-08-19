package views

import (
	"context"
	"math/rand/v2"
	"net/http"
	"shin/src/app/auth"
	"shin/src/app/models"
	"shin/src/app/services"
	"shin/src/config"
	"shin/src/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func authGroup(router *gin.Engine) {
	g := router.Group("auth")

	g.POST("/login", func(c *gin.Context) {
		form := new(auth.LoginForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, err := models.GetUserByEmail(form.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email/password not match"})
			return
		}
		if err := auth.CheckPasswordHash(form.Password, *u.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email/password not match"})
			return
		}
		accessToken, err := auth.GenerateToken(u.ID.String(), false)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
		})
	})

	g.POST("/register", func(c *gin.Context) {
		form := new(auth.RegisterForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u := new(models.User)
		utils.Copy(form, u)
		password, _ := auth.HashPassword(form.Password)
		u.Password = &password
		ctx, _ := c.Get("ctx")
		if err := u.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		accessToken, err := auth.GenerateToken(u.ID.String(), false)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
		})
	})

	g.POST("/otp", auth.LoginRequired(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx, _ := c.Get("ctx")
		otp := models.OTP{
			UserID: uuid.MustParse(userID.(string)),
			Code:   int(100000 + rand.Float64()*900000),
		}
		otp.Create(ctx.(context.Context))

		//Sending Email
		u, err := models.GetUser(uuid.MustParse(userID.(string)))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "User does not found",
			})
			return
		}
		items := map[string]string{"name": *u.FirstName, "code": strconv.Itoa(otp.Code)}
		err = services.SendGridClient.SendWithTemplate(u.Email, "OTP Code", services.SendGridTemplates["otp"], items)
		if err != nil && config.Config.Env != "test" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Couldn't send OTP Code to mailbox",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	g.POST("/otp/resend", auth.LoginRequired(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		otp, err := models.GetOTPByUserID(uuid.MustParse(userID.(string)))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		//Sending Email
		u, err := models.GetUser(uuid.MustParse(userID.(string)))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "User does not found",
			})
			return
		}
		items := map[string]string{"name": *u.FirstName, "code": strconv.Itoa(otp.Code)}
		err = services.SendGridClient.SendWithTemplate(u.Email, "OTP Code", services.SendGridTemplates["otp"], items)
		if err != nil && config.Config.Env != "test" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Couldn't send OTP Code to mailbox",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	g.POST("/otp/verify", auth.LoginRequired(), func(c *gin.Context) {

		form := new(auth.OTPConfirmForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("user_id")
		ctx, _ := c.Get("ctx")

		//Verifying OTP
		otp := models.OTP{
			UserID: uuid.MustParse(userID.(string)),
			Code:   form.Code,
		}

		err := otp.Verify(ctx.(context.Context))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "A problem occured when trying to verify the code",
			})
			return
		} else if otp.IsVerified == false {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   nil,
				"message": "Code does not found or it is wrong",
			})
			return
		}

		//Verifying User
		u := models.User{
			ID:     uuid.MustParse(userID.(string)),
			Status: "ACTIVE",
		}
		if err := u.Verify(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})

	})

}
