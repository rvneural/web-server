package audio

import (
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	modelAudio "WebServer/internal/models/db/results/audio"
	modelImage "WebServer/internal/models/db/results/image"
	modelText "WebServer/internal/models/db/results/text"

	"github.com/gin-gonic/gin"
)

type Result struct {
	notFoundOperation interfaces.NoResultPage
	progressOperation interfaces.NoResultPage
	dbWorker          interfaces.DBWorker
	logger            *slog.Logger
}

func New(notFoundOperation, progressOperation interfaces.NoResultPage, dbWorker interfaces.DBWorker, logger *slog.Logger) *Result {
	return &Result{
		notFoundOperation: notFoundOperation,
		progressOperation: progressOperation,
		dbWorker:          dbWorker,
		logger:            logger,
	}
}

func (r *Result) GetPage(c *gin.Context) {

	id := c.Param("id")

	res, err := r.dbWorker.GetResult(id)

	if err != nil {
		r.logger.Error("Getting result from DB", "error", err)
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
			r.logger.Error("Unmarshalling result", "error", err)
			r.notFoundOperation.GetPage(c, id)
			return
		}

		if result.Name == "" {
			result.Name = "image.jpg"
		}

		c.HTML(http.StatusOK, "image-generation-result.html", gin.H{
			"title":  "Результаты генерации",
			"style":  style,
			"prompt": strings.ReplaceAll(result.Prompt, "&#34;", "\""),
			"seed":   result.Seed,
			"image":  result.B64string,
			"name":   result.Name,
		})
		return
	case "audio":
		style := "/web/styles/results/recognition-style.css"
		result := modelAudio.DBResult{}
		err = json.Unmarshal(res.DATA, &result)

		if err != nil {
			r.logger.Error("Unmarshalling result", "error", err)
			r.notFoundOperation.GetPage(c, id)
			return
		}

		c.HTML(http.StatusOK, "recognition-result.html", gin.H{
			"title":     "Результаты расшифровки",
			"style":     style,
			"filename":  result.FileName,
			"raw_text":  strings.ReplaceAll(result.RawText, "&#34;", "\""),
			"norm_text": strings.ReplaceAll(result.NormText, "&#34;", "\""),
			"version":   res.VERSION,
			"script":    "/web/scripts/results/recognition-result-script.js",
			"id":        res.OPERATION_ID,
		})
		return
	case "text":
		style := "/web/styles/results/text-processing-style.css"
		result := modelText.DBResult{}

		err = json.Unmarshal(res.DATA, &result)

		if err != nil {
			r.logger.Error("Unmarshalling result", "error", err)
			r.notFoundOperation.GetPage(c, id)
			return
		}

		c.HTML(http.StatusOK, "text-processing-result.html", gin.H{
			"title":    "Результаты обработки",
			"style":    style,
			"old_text": strings.ReplaceAll(result.OldText, "&#34;", "\""),
			"new_text": strings.ReplaceAll(result.NewText, "&#34;", "\""),
			"prompt":   strings.ReplaceAll(result.Prompt, "&#34;", "\""),
			"version":  res.VERSION,
			"script":   "/web/scripts/results/text-processing-result-script.js",
			"id":       res.OPERATION_ID,
		})
		return
	default:
		r.notFoundOperation.GetPage(c, id)
	}

}
