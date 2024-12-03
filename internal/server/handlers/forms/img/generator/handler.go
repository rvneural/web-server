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

	dbModel "WebServer/internal/models/db/results/image"
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
	var dbError error
	if len(id) != 0 {
		var u_id int
		user_id, err := c.Cookie("user_id")
		if err != nil {
			u_id = 0
		} else {
			u_id, err = strconv.Atoi(user_id)
			if err != nil {
				u_id = 0
			}
		}
		dbError = n.dbWorker.RegisterOperation(id, "image", u_id)
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

	byteRequets, err := json.Marshal(request)

	if err != nil {
		n.logger.Error("Marshalling request", "error", err)
		go n.saveErrorToDB(id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	httpRequest, err := http.NewRequest("POST", config.URL, bytes.NewBuffer(byteRequets))
	if err != nil {
		n.logger.Error("Creating request", "error", err)
		go n.saveErrorToDB(id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		n.logger.Error("Sending request", "error", err)
		go n.saveErrorToDB(id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	byteResp, err := io.ReadAll(resp.Body)

	if err != nil {
		n.logger.Error("Reading response", "error", err)
		go n.saveErrorToDB(id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model := models.Response{}
	err = json.Unmarshal(byteResp, &model)

	if err != nil {
		n.logger.Error("Unmarshalling response", "error", err)
		go n.saveErrorToDB(id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(id string, dbError error) {
		if dbError == nil && len(id) != 0 {
			dbResult := dbModel.DBResult{
				Prompt:    prompt,
				Seed:      model.Image.Seed,
				B64string: model.Image.B64String,
				Name:      strings.ReplaceAll(strings.TrimSpace(prompt), " ", "-") + ".jpg",
			}
			dbByteRes, err := json.Marshal(dbResult)
			if err == nil {
				n.dbWorker.SetResult(id, dbByteRes)
			}
		}
	}(id, dbError)

	c.JSON(http.StatusOK, model)
}

func (n *ImageGenerationHandler) saveErrorToDB(id string, dbError string) {
	if len(id) == 0 {
		return
	}
	dbResult := dbModel.DBResult{
		Prompt:    dbError,
		Seed:      "",
		B64string: "",
	}
	dbByteRes, err := json.Marshal(dbResult)
	if err == nil {
		n.dbWorker.SetResult(id, dbByteRes)
	}
}
