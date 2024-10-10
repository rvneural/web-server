package upscale

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

func (i *ImageUpscaler) HandleForm(w http.ResponseWriter, r *http.Request) {
	imgFile, header, err := r.FormFile("image")
	log.Println("New image for upscaling:", header.Filename)

	if err != nil {
		log.Println(err)
		resp := &Response{Error: err.Error()}
		d, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(d)
		return
	}

	imgBytes, err := io.ReadAll(imgFile)
	if err != nil {
		log.Println(err)
		resp := &Response{Error: err.Error()}
		d, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(d)
		return
	}

	file, err := os.Create(header.Filename)
	if err != nil {
		log.Println(err)
		resp := &Response{Error: err.Error()}
		d, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(d)
		return
	}
	defer os.Remove(header.Filename)
	defer file.Close()

	_, err = file.Write(imgBytes)
	if err != nil {
		log.Println(err)
		resp := &Response{Error: err.Error()}
		d, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(d)
		return
	}

	upscaler := imageupscaler.New()
	err = upscaler.SetImage(header.Filename)
	if err != nil {
		log.Println(err)
		resp := &Response{Error: err.Error()}
		d, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(d)
		return
	}

	renderFile, err := os.Create("../../web/uploads/upscaled-" + strings.ReplaceAll(header.Filename, " ", "-") + ".jpg")
	if err != nil {
		log.Println(err)
		resp := &Response{Error: err.Error()}
		d, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(d)
		return
	}

	upscaler.Upscale(3, 2)
	upscaler.Render(imageupscaler.JPG, renderFile, nil)
	file.Close()
	resp := &Response{URL: "/web/uploads/upscaled-" + strings.ReplaceAll(header.Filename, " ", "-") + ".jpg"}
	d, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(d)
	if err != nil {
		log.Println(err)
		resp := &Response{Error: err.Error()}
		d, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(d)
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
