package generator

import (
	models "WebServer/internal/models/img"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageGenerationHandler struct {
}

func New() *ImageGenerationHandler {
	return &ImageGenerationHandler{}
}

// [x] Image Generation handler
func (n *ImageGenerationHandler) HandleForm(c *gin.Context) {
	prompt := c.Request.FormValue("prompt")
	seed := c.Request.FormValue("seed")
	widthRatio := c.Request.FormValue("widthRatio")
	heightRatio := c.Request.FormValue("heightRatio")

	var request models.Request
	request.Prompt = prompt
	request.Seed = seed
	request.WidthRatio = widthRatio
	request.HeightRatio = heightRatio

	byteRequets, err := json.Marshal(request)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	httpRequest, err := http.NewRequest("POST", "http://127.0.0.1:8083/", bytes.NewBuffer(byteRequets))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer KU2WTZBzFWn4Ko9lJ7TlpmUXwkHc8Y")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	byteResp, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, byteResp)
}
