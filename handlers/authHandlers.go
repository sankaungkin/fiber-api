package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/badoux/checkmail"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const jwtSecret = "superSecretKey"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `jsong:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func Login(c *fiber.Ctx) error {

	// key, errKey := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// if errKey != nil {
	// 	log.Fatal(errKey)
	// }

	db := database.DB

	json := new(LoginRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	found := models.User{}
	query := models.User{Email: json.Email}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "Record not found",
		})
	}

	if !comparePasswords(found.Password, []byte(json.Password)) {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Authorization Failed",
		})
	}

	// jwt access token
	atClaims := jwt.MapClaims{
		"id":    found.ID,
		"email": found.Email,
		"admin": true,
		"exp":   time.Now().Add(time.Minute * 3).Unix(),
	}

	// Create token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	// Generate encoded token and send it as response.
	// at, err := accessToken.SignedString([]byte(jwtSecret))
	at, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// refresh token
	rtClaims := jwt.MapClaims{
		"id":    found.ID,
		"email": found.Email,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	rt, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    at,
		Path:     "/",
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    rt,
		Path:     "/",
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "login success",
		"data": LoginResponse{
			AccessToken:  at,
			RefreshToken: rt,
		}})

}

func CreateUser(c *fiber.Ctx) error {

	type CreateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		IsAdmin  bool   `json:"isAdmin"`
	}

	db := database.DB
	json := new(CreateUserRequest)

	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	password := hashAndSalt([]byte(json.Password))
	err := checkmail.ValidateFormat(json.Email)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid Email Address",
		})
	}

	newUser := models.User{
		UserName: json.Username,
		Email:    json.Email,
		Password: password,
		IsAdmin:  json.IsAdmin,
	}

	// err := c.BodyParser(&newUser)
	// if err != nil {
	// 	c.Status(http.StatusUnprocessableEntity).JSON(
	// 		&fiber.Map{
	// 			"message": "request failed",
	// 		})
	// 	return err
	// }

	errors := models.ValidateStruct(newUser)
	if errors != nil {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}

	existing := db.Where("email = ?", json.Email).Find(&newUser)
	if existing.RowsAffected > 0 {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "email is already taken",
		})
	}
	err = db.Create(&newUser).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not create new user",
		})
		return err
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"status":  "SUCCESS",
		"message": "user has been created successfully",
		"data":    newUser,
	})

}

func GetUsers(c *fiber.Ctx) error {
	db := database.DB

	users := []models.User{}

	db.Model(&models.User{}).Order("ID asc").Find(&users)

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    "404",
			"message": "NO RECORD",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": len(users),
		"data":    users,
	})
}

func GetUser(c *fiber.Ctx) error {

	db := database.DB

	id := c.Params("id")

	var user models.User

	result := db.First(&user, "id = ?", id)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "FAIL",
				"message": "No data",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "FAIL",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Record found",
		"data":    user,
	})

}

func hashAndSalt(pwd []byte) string {
	hash, _ := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}

func SessionExpires() time.Time {
	return time.Now().Add(5 * 24 * time.Hour)
}

// TODO logout
func Logout(c *fiber.Ctx) error {

	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		Expires:  expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	// c.Cookie(&fiber.Cookie{
	// 	Name:    "refresh_token",
	// 	Value:   "",
	// 	Expires: time.Now().Add(-time.Hour),
	// })

	// return c.Redirect("/")
}

func GetToken(c *fiber.Ctx) string {
	token := c.Context().UserValue("JWT_TOKEN")
	if token == nil {
		return ""
	}
	return token.(string)
}

func Refresh(c *fiber.Ctx) error {
	// type TokenRequest struct {
	// 	RefreshToken string `json:"refresh_token"`
	// }

	// tokenRequst := TokenRequest{}

	/*
			tokenString := "<YOUR TOKEN STRING>"
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		    return []byte("<YOUR VERIFICATION KEY>"), nil
			})
			// ... error handling

			// do something with decoded claims
			for key, val := range claims {
		    	fmt.Printf("Key: %v, value: %v\n", key, val)
			}
	*/
	tokenString := c.Cookies("refresh_token")
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	for key, val := range claims {
		fmt.Printf("key: %v, value: %v\n", key, val)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in

		db := database.DB
		found := models.User{}
		query := models.User{Email: claims["email"].(string)}
		err := db.First(&found, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    404,
				"message": "Record not found",
			})
		}

		rtClaims := jwt.MapClaims{
			"id":    found.ID,
			"email": found.Email,
			"admin": found.IsAdmin,
			"exp":   time.Now().Add(time.Hour * 1).Unix(),
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

		rt, err := refreshToken.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			Secure:   false,
			HTTPOnly: true,
			Domain:   "localhost",
		})
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    rt,
			Path:     "/",
			Secure:   false,
			HTTPOnly: true,
			Domain:   "localhost",
		})

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":  "SUCCESS",
			"message": "login success",
			"data": LoginResponse{
				AccessToken:  "",
				RefreshToken: rt,
			}})

	}

	return err

}
