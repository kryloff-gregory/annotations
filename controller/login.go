package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"main/user"
)

var jwtKey = []byte("SECRET_KEY")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type login struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginController struct {
	userProvider *user.Provider
}

func NewLoginController(provider *user.Provider) *LoginController {
	return &LoginController{userProvider: provider}
}

func (u *LoginController) Login(c *gin.Context) {
	var data login

	// Bind the request body data to var data and check if all details are provided
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required details"})
		c.Abort()
		return
	}

	usr := u.userProvider.GerUserByName(data.UserName)
	if usr == nil {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	if err := validatePassword(data.Password, usr.HashedPassword); err != nil {
		c.JSON(403, gin.H{"message": "Invalid user credentials"})
		c.Abort()
		return
	}

	jwtToken, err := generateToken(usr.Name)
	if err != nil {
		c.JSON(403, gin.H{"message": "There was a problem logging you in, try again later"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Log in success", "token": jwtToken})
}
