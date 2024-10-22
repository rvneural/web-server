package upscale

import (
	"encoding/json"
	"io"
	"log"
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
}

func New() *ImageUpscaler {
	return &ImageUpscaler{}
}

// [x] Image Upscale Handler
func (i *ImageUpscaler) HandleForm(c *gin.Context) {
	imgFile, header, err := c.Request.FormFile("image")
	log.Println("New image for upscaling:", header.Filename)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imgBytes, err := io.ReadAll(imgFile)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := os.Create(header.Filename)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer os.Remove(header.Filename)
	defer file.Close()

	_, err = file.Write(imgBytes)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	upscaler := imageupscaler.New()
	err = upscaler.SetImage(header.Filename)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	renderFile, err := os.Create("../../web/uploads/upscaled-" + strings.ReplaceAll(header.Filename, " ", "-") + ".jpg")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	upscaler.Upscale(3, 2)
	upscaler.Render(imageupscaler.JPG, renderFile, nil)
	file.Close()
	resp := &Response{URL: "/web/uploads/upscaled-" + strings.ReplaceAll(header.Filename, " ", "-") + ".jpg"}
	d, _ := json.Marshal(resp)
	c.JSON(http.StatusOK, d)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go deleteImage("../../web/uploads/upscaled-"+strings.ReplaceAll(header.Filename, " ", "-")+".jpg", time.After(1*time.Hour))
}

func deleteImage(url string, delay <-chan time.Time) {
	select {
	case <-delay:
		os.Remove(url)
	}
}
