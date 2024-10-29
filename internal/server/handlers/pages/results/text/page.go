package text

import (
	model "WebServer/internal/models/db/results/text"
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
	"log"
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
	style := "/web/styles/results/text-processing-style.css"

	id := c.Param("id")

	res, err := r.dbWorker.GetResult(id)

	if err != nil {
		r.notFoundOperation.GetPage(c, id)
		return
	} else if res.OPERATION_TYPE != "text" {
		c.Redirect(http.StatusMovedPermanently, "/operation/"+res.OPERATION_TYPE+"/"+res.OPERATION_ID)
		return
	} else if res.IN_PROGRESS {
		r.progressOperation.GetPage(c, id)
		return
	}

	result := model.DBResult{}

	log.Println(string(res.DATA))
	err = json.Unmarshal(res.DATA, &result)

	if err != nil {
		r.notFoundOperation.GetPage(c, id)
		return
	}

	c.HTML(http.StatusOK, "text-processing-result.html", gin.H{
		"title":    "Результаты обработки",
		"style":    style,
		"old_text": result.OldText,
		"new_text": result.NewText,
		"prompt":   result.Prompt,
	})
}
