package rembg

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

type BackgroundRemover struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *BackgroundRemover {
	return &BackgroundRemover{
		logger: logger,
	}
}

func (i *BackgroundRemover) HandleForm(c *gin.Context) {
	imgFile, header, err := c.Request.FormFile("image")

	if err != nil {
		i.logger.Error("Getting file from form", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imgBytes, err := io.ReadAll(imgFile)
	if err != nil {
		i.logger.Error("Reading file from form", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	remover_url := os.Getenv("REMOVER_URL")
	if remover_url == "" {
		i.logger.Error("NO REMOVER URL")
		c.JSON(http.StatusBadRequest, gin.H{"error": "NO REMOVER URL"})
		return
	}

	type Request struct {
		Image []byte `json:"image"`
	}

	type Response struct {
		Image []byte `json:"image"`
	}

	req := &Request{
		Image: imgBytes,
	}

	byteReq, err := json.Marshal(req)
	if err != nil {
		i.logger.Error("Marshalling request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.Post(remover_url, "application/json", bytes.NewBuffer(byteReq))
	if err != nil {
		i.logger.Error("Sending request to remover", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		i.logger.Error("Decoding response", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imgData := res.Image
	fileName := "../../web/uploads/bg-" + strings.TrimSpace(header.Filename) + ".png"
	file, err := os.Create(fileName)
	if err != nil {
		i.logger.Error("Creating file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	_, err = file.Write(imgData)
	if err != nil {
		i.logger.Error("Writing file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": fileName})
	go deleteImage(fileName, time.After(time.Hour*1))

}

func deleteImage(url string, delay <-chan time.Time) {
	select {
	case <-delay:
		os.Remove(url)
	}
}
