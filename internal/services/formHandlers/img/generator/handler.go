package generator

import (
	models "WebServer/internal/models/img"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ImageGenerationHandler struct {
}

func New() *ImageGenerationHandler {
	return &ImageGenerationHandler{}
}

// [ ] Image Generation handler
func (n *ImageGenerationHandler) HandleForm(w http.ResponseWriter, r *http.Request) {
	log.Println("New generate image from web request from", r.RemoteAddr)

	prompt := r.FormValue("prompt")
	seed := r.FormValue("seed")
	widthRatio := r.FormValue("widthRatio")
	heightRatio := r.FormValue("heightRatio")

	var request models.Request
	request.Prompt = prompt
	request.Seed = seed
	request.WidthRatio = widthRatio
	request.HeightRatio = heightRatio

	byteRequets, err := json.Marshal(request)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": ` + err.Error()))
		return
	}

	httpRequest, err := http.NewRequest("POST", "http://127.0.0.1:8083/", bytes.NewBuffer(byteRequets))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": ` + err.Error()))
		return
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer KU2WTZBzFWn4Ko9lJ7TlpmUXwkHc8Y")

	client := http.Client{}

	resp, err := client.Do(httpRequest)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": ` + err.Error()))
		return
	}

	defer resp.Body.Close()

	byteResp, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ` + err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(byteResp)
}
