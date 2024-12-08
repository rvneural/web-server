package generator

import (
	config "WebServer/internal/config/services/text2image-service"
	models "WebServer/internal/models/img"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"WebServer/internal/server/handlers/interfaces"

	"github.com/gin-gonic/gin"
)

type ImageGenerationHandler struct {
	dbWorker interfaces.DBWorker
	logger   *slog.Logger
}

func New(dbWorker interfaces.DBWorker, logger *slog.Logger) *ImageGenerationHandler {
	return &ImageGenerationHandler{
		dbWorker: dbWorker,
		logger:   logger,
	}
}

func (n *ImageGenerationHandler) HandleForm(c *gin.Context) {

	id := c.Request.FormValue("id")
	id = strings.TrimSpace(id)
	str_id, err := c.Cookie("user_id")
	var user_id int
	if err != nil {
		user_id = 0
	} else {
		user_id, err = strconv.Atoi(str_id)
		if err != nil {
			user_id = 0
		}
	}

	prompt := c.Request.FormValue("prompt")
	seed := c.Request.FormValue("seed")
	widthRatio := c.Request.FormValue("widthRatio")
	heightRatio := c.Request.FormValue("heightRatio")

	var request models.Request
	request.Prompt = prompt
	request.Seed = seed
	request.WidthRatio = widthRatio
	request.HeightRatio = heightRatio
	request.UserID = user_id
	request.OperationId = id

	byteRequets, err := json.Marshal(request)

	if err != nil {
		n.logger.Error("Marshalling request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	httpRequest, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(byteRequets))
	if err != nil {
		n.logger.Error("Creating request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		n.logger.Error("Sending request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	byteResp, err := io.ReadAll(resp.Body)

	if err != nil {
		n.logger.Error("Reading response", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model := models.Response{}
	err = json.Unmarshal(byteResp, &model)

	if err != nil {
		n.logger.Error("Unmarshalling response", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model)
}
