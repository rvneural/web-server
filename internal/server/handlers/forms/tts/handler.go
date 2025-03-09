package tts

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Text       string `json:"text" form:"text" binding:"required"`
	Voice      string `json:"voice" form:"voice" binding:"required"`
	Role       string `json:"role" form:"role"`
	Speed      string `json:"speed" form:"speed"`
	PitchShift string `json:"pitchShift" form:"pitchShift"`
}

type Response struct {
	Audio []byte `json:"audio"`
	Error string `json:"error"`
}

type Handler struct {
	url string
}

func New() *Handler {
	return &Handler{
		url: os.Getenv("TEXT_2_SPEECH_URL"),
	}
}

func (h *Handler) Handler(c *gin.Context) {
	var req Request
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Text = strings.ReplaceAll(req.Text, "&lt;", "<")
	req.Text = strings.ReplaceAll(req.Text, "&gt;", ">")
	req.Text = strings.ReplaceAll(req.Text, "&#34;", "\"")

	data, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := http.Post(h.url, "application/json", bytes.NewReader(data))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
