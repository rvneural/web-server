package generator

import (
	config "WebServer/internal/config/services/text2image-service"
	models "WebServer/internal/models/img"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

type ImageGenerationHandler struct {
}

func New() *ImageGenerationHandler {
	return &ImageGenerationHandler{}
}

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

	httpRequest, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(byteRequets))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+config.KEY)

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

	model := models.Response{}
	err = json.Unmarshal(byteResp, &model)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model)

	store := ginsession.FromContext(c)
	store.Set("generation-image", model.Image.B64String)
	store.Set("generation-seed", model.Image.Seed)
	store.Set("generation-prompt", prompt)

	err = store.Save()
	if err != nil {
		log.Panicln("Can't save session: " + err.Error())
	}
}
