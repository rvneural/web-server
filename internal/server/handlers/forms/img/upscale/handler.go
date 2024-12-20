package upscale

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	imageupscaler "github.com/neuron-nexus/go-image-upscaler"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Response struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

type ImageUpscaler struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *ImageUpscaler {
	return &ImageUpscaler{
		logger: logger,
	}
}

func (i *ImageUpscaler) HandleForm(c *gin.Context) {
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

	file, err := os.Create(header.Filename)
	if err != nil {
		i.logger.Error("Creating file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer os.Remove(header.Filename)
	defer file.Close()

	_, err = file.Write(imgBytes)
	if err != nil {
		i.logger.Error("Writing file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	upscaler := imageupscaler.New()
	err = upscaler.SetImage(header.Filename)
	if err != nil {
		i.logger.Error("Setting image", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	renderFile, err := os.Create("../../web/uploads/upscaled-" + strings.ReplaceAll(header.Filename, " ", "-") + ".png")
	if err != nil {
		i.logger.Error("Creating file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	upscaler.Upscale(3, 2)
	upscaler.Render(imageupscaler.PNG, renderFile, nil)
	file.Close()
	resp := &Response{URL: "/web/uploads/upscaled-" + strings.ReplaceAll(header.Filename, " ", "-") + ".png"}
	c.JSON(http.StatusOK, resp)

	go deleteImage("../../web/uploads/upscaled-"+strings.ReplaceAll(header.Filename, " ", "-")+".png", time.After(1*time.Hour))
}

func deleteImage(url string, delay <-chan time.Time) {
	select {
	case <-delay:
		os.Remove(url)
	}
}
