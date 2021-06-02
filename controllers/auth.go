package controllers

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/prastamaha/auth-basic/database"
	"github.com/prastamaha/auth-basic/models"
	"github.com/prastamaha/auth-basic/repository"
	"github.com/prastamaha/auth-basic/utils"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {

	// Parse body request into map form
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Create new user model
	user := &models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: data["password"],
	}

	// Insert into database
	db := repository.NewUserRepository(database.DB)
	response, err := db.Insert(context.Background(), user)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// return success response
	return c.JSON(fiber.Map{
		"message": response,
		"data":    user,
	})
}

func Login(c *fiber.Ctx) error {

	// Parse body request into map form
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Query user
	db := repository.NewUserRepository(database.DB)
	user, err := db.QueryByEmail(context.Background(), data["email"])
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "invalid1 email or password",
		})
	}
	if ok := utils.VerifyPassword(user.Password, data["password"]); !ok {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "invalid2 email or password",
		})
	}
	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "invalid3 email or password",
		})
	}

	// Create new claim
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Name,
		ExpiresAt: jwt.NewTime(float64(time.Now().Add(time.Minute * 30).Unix())),
	})

	// Generate token from claim
	token, err := claim.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	// Create new cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Minute * 30),
		Path:     "/api/",
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "login success",
	})
}

func Logout(c *fiber.Ctx) error {
	// check if user login
	ck := c.Cookies("jwt")
	if ck == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "you are not logged in",
		})
	}

	// create new cookie with that invalid in last 3s ago
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "deleted",
		Path:     "/api/",
		Expires:  time.Now().Add(-time.Second * 3),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	// return logout success
	return c.JSON(fiber.Map{
		"message": "logout success",
	})
}

func UserProfile(c *fiber.Ctx) error {

	// get jwt cookies
	cookie := c.Cookies("jwt")

	// parse jwt
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(SecretKey), nil
	})

	// if there have error
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// convert token.Claims into jwt standartclaims
	claims := token.Claims.(*jwt.StandardClaims)
	response := "hello " + claims.Issuer

	// check if valid
	if token.Valid {
		return c.JSON(fiber.Map{
			"message": response,
		})
	} else {
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}
}
