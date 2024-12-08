package auth

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	authModels "WebServer/internal/models/user/model"
	"WebServer/internal/server/handlers/interfaces"
)

type RegisterWorker interface {
	CheckForRegistered(email string) bool
	Register(email, hashPassword, FirstName, LastName string) string
	CheckForLogin(email, hashPassword string) (status bool, user_id string)
}

type AuthentificationHandler struct {
	worker RegisterWorker
}

func New(worker RegisterWorker) *AuthentificationHandler {
	return &AuthentificationHandler{
		worker: worker,
	}
}

func (a *AuthentificationHandler) HandleRegistration(c *gin.Context) {
	model := authModels.RegistrationInput{}

	err := c.Bind(&model)
	if err != nil {
		log.Println("Error bindind register body")
		c.JSON(http.StatusOK, gin.H{"Register binding error": err.Error()})
		return
	}
	if a.worker.CheckForRegistered(model.Email) {
		c.JSON(http.StatusOK, gin.H{
			"error":  "User with this email already exists",
			"exists": true,
		})
		return
	}
	log.Println("Registering user")
	id := a.worker.Register(model.Email, model.Password, model.FirstName, model.LastName)
	claims := jwt.MapClaims{
		"user_id":    id,
		"user_email": model.Email,
	}

	if id == "" {
		c.JSON(http.StatusOK, gin.H{"error": "Ошибка при регистрации"})
		return
	}

	log.Println("Signing token")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Println("Error signing token")
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	log.Println("Setting token")
	c.SetCookie("NeuronNexusAuth", ss, 3600*24*30, "/", "", false, true)
	c.SetCookie("user_id", id, 3600*24*30, "/", "", false, true)
	c.SetCookie("user_email", model.Email, 3600*24*30, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered"})
}

func (a *AuthentificationHandler) HandleLogin(c *gin.Context) {
	model := authModels.AuthentificationInput{}

	err := c.Bind(&model)
	if err != nil {
		log.Println("Error binding login body")
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	log.Println("Passowrd:", model.Password)
	status, id := a.worker.CheckForLogin(model.Email, model.Password)
	if !status {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid email or password"})
		return
	}

	log.Println("Signing token")
	claims := jwt.MapClaims{
		"user_id":    id,
		"user_email": model.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Println("Error signing token")
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	log.Println("Setting token for user_id:", id)
	c.SetCookie("NeuronNexusAuth", ss, 3600*24*30, "/", "", false, true)
	c.SetCookie("user_id", id, 3600*24*30, "/", "", false, true)
	c.SetCookie("user_email", model.Email, 3600*24*30, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
}

func (a *AuthentificationHandler) AuthMiddleware(authPath string, minimal_level int, db_worker interfaces.DBWorker) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("NeuronNexusAuth")
		if err != nil || tokenString == "" {
			log.Println("No token found")
			c.Redirect(http.StatusPermanentRedirect, authPath)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid token")
			c.Redirect(http.StatusPermanentRedirect, authPath)
			return
		}

		clms, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Redirect(http.StatusPermanentRedirect, authPath)
			return
		}

		id, err := c.Cookie("user_id")
		if err != nil || id == "" {
			c.Redirect(http.StatusPermanentRedirect, authPath)
			return
		}

		email, err := c.Cookie("user_email")
		if err != nil || email == "" {
			c.Redirect(http.StatusPermanentRedirect, authPath)
			return
		}
		if email != clms["user_email"].(string) || id != clms["user_id"].(string) {
			c.Redirect(http.StatusPermanentRedirect, authPath)
			return
		}

		if minimal_level > 0 {
			int_id, err := strconv.Atoi(id)
			if err != nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			user, err := db_worker.GetUserByID(int_id)
			if err != nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			if user.USER_STATUS < minimal_level {
				c.AbortWithStatus(http.StatusForbidden)
				return
			} else {
				c.Next()
			}
		} else {
			c.Next()
		}
	}
}
