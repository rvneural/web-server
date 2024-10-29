package image

import (
	model "WebServer/internal/models/db/results/image"
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecognitionResult struct {
	notFoundOperation interfaces.NoResultPage
	progressOperation interfaces.NoResultPage
	dbWorker          interfaces.DBWorker
}

func New(notFoundOperation, progressOperation interfaces.NoResultPage, dbWorker interfaces.DBWorker) *RecognitionResult {
	return &RecognitionResult{
		notFoundOperation: notFoundOperation,
		progressOperation: progressOperation,
		dbWorker:          dbWorker,
	}
}

func (r *RecognitionResult) GetPage(c *gin.Context) {

	style := "/web/styles/results/image-generation-style.css"

	id := c.Param("id")

	res, err := r.dbWorker.GetResult(id)

	if err != nil {
		r.notFoundOperation.GetPage(c, id)
		return
	} else if res.OPERATION_TYPE != "image" {
		c.Redirect(http.StatusMovedPermanently, "/operation/"+res.OPERATION_TYPE+"/"+res.OPERATION_ID)
		return
	} else if res.IN_PROGRESS {
		r.progressOperation.GetPage(c, id)
		return
	}

	result := model.DBResult{}

	err = json.Unmarshal(res.DATA, &result)

	if err != nil {
		r.notFoundOperation.GetPage(c, id)
		return
	}

	c.HTML(http.StatusOK, "image-generation-result.html", gin.H{
		"title":  "Результаты генерации",
		"style":  style,
		"prompt": result.Prompt,
		"seed":   result.Seed,
		"image":  result.B64string,
	})
}
