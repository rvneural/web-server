package audio

import (
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
	"net/http"

	modelAudio "WebServer/internal/models/db/results/audio"
	modelImage "WebServer/internal/models/db/results/image"
	modelText "WebServer/internal/models/db/results/text"

	"github.com/gin-gonic/gin"
)

type Result struct {
	notFoundOperation interfaces.NoResultPage
	progressOperation interfaces.NoResultPage
	dbWorker          interfaces.DBWorker
}

func New(notFoundOperation, progressOperation interfaces.NoResultPage, dbWorker interfaces.DBWorker) *Result {
	return &Result{
		notFoundOperation: notFoundOperation,
		progressOperation: progressOperation,
		dbWorker:          dbWorker,
	}
}

func (r *Result) GetPage(c *gin.Context) {

	id := c.Param("id")

	res, err := r.dbWorker.GetResult(id)

	if err != nil {
		r.notFoundOperation.GetPage(c, id)
		return
	} else if res.IN_PROGRESS {
		r.progressOperation.GetPage(c, id)
		return
	}

	switch res.OPERATION_TYPE {
	case "image":
		style := "/web/styles/results/image-generation-style.css"
		result := modelImage.DBResult{}

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
		return
	case "audio":
		style := "/web/styles/results/recognition-style.css"
		result := modelAudio.DBResult{}
		err = json.Unmarshal(res.DATA, &result)

		if err != nil {
			r.notFoundOperation.GetPage(c, id)
			return
		}

		c.HTML(http.StatusOK, "recognition-result.html", gin.H{
			"title":     "Результаты расшифровки",
			"style":     style,
			"raw_text":  result.RawText,
			"norm_text": result.NormText,
		})
		return
	case "text":
		style := "/web/styles/results/text-processing-style.css"
		result := modelText.DBResult{}

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
		return
	default:
		r.notFoundOperation.GetPage(c, id)
	}

}
