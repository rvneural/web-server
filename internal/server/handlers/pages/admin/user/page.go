package user

import (
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	modelAudio "WebServer/internal/models/db/results/audio"
	modelImage "WebServer/internal/models/db/results/image"
)

type Page struct {
	logger *slog.Logger
	worker interfaces.DBWorker
}

func New(worker interfaces.DBWorker, logger *slog.Logger) *Page {
	return &Page{
		logger: logger,
		worker: worker,
	}
}

func (p *Page) GetPage(c *gin.Context) {
	user_id_str := c.Param("id")
	var user_id int
	user_id, err := strconv.Atoi(user_id_str)
	if err != nil {
		user_id = 0
	}
	p.logger.Info("Getting user", "user_id", user_id)
	user, err := p.worker.GetUserByID(user_id)
	if err != nil {
		p.logger.Error("Getting user", "error", err)
		user.EMAIL = "неизвестно"
		user.FIRSTNAME = "Неизвестный"
		user.LASTNAME = "Пользователь"
	}

	limit_str := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limit_str)
	if err != nil {
		limit = 100
	}

	image_operations, err := p.worker.GetUserOperations(user_id, limit, "image")
	if err != nil {
		p.logger.Error("Getting images operations", "error", err)
	}
	text_operations, err := p.worker.GetUserOperations(user_id, limit, "text")
	if err != nil {
		p.logger.Error("Getting text operations", "error", err)
	}
	audio_operations, err := p.worker.GetUserOperations(user_id, limit, "audio")
	if err != nil {
		p.logger.Error("Getting audio operations", "error", err)
	}

	type Operation struct {
		ID          string
		Placeholder string
		Date        string
	}

	wg := sync.WaitGroup{}
	wg.Add(3)
	Images := make([]Operation, len(image_operations))
	Texts := make([]Operation, len(text_operations))
	Audios := make([]Operation, len(audio_operations))

	go func() {
		defer wg.Done()
		for i, operation := range image_operations {
			result := modelImage.DBResult{}
			var placeholder string
			err := json.Unmarshal(operation.DATA, &result)
			if err != nil {
				p.logger.Error("Unmarshalling image operation", "error", err)
				placeholder = "Генерация изображения"
			} else {
				placeholder = result.Prompt
			}
			Images[i].ID = operation.OPERATION_ID
			Images[i].Placeholder = placeholder
			Images[i].Date = operation.CREATION_DATE.Format("02.01.2006 15:04:05")
		}
	}()
	go func() {
		defer wg.Done()
		for i, operation := range text_operations {
			Texts[i].ID = operation.OPERATION_ID
			Texts[i].Placeholder = "Обработка текста"
			Texts[i].Date = operation.CREATION_DATE.Format("02.01.2006 15:04:05")
		}
	}()
	go func() {
		defer wg.Done()
		for i, operation := range audio_operations {
			result := modelAudio.DBResult{}
			var placeholder string
			err := json.Unmarshal(operation.DATA, &result)
			if err != nil {
				p.logger.Error("Unmarshalling image operation", "error", err)
				placeholder = "Расшифровка фгвшщ"
			} else {
				placeholder = "Расшифровка: " + result.FileName
			}
			Audios[i].ID = operation.OPERATION_ID
			Audios[i].Placeholder = placeholder
			Audios[i].Date = operation.CREATION_DATE.Format("02.01.2006 15:04:05")
		}
	}()
	wg.Wait()

	style := "/web/styles/admin-user.css"
	c.HTML(http.StatusOK, "admin-user.html", gin.H{
		"style":        style,
		"title":        "Страница пользователя",
		"FirstName":    user.FIRSTNAME,
		"LastName":     user.LASTNAME,
		"Email":        user.EMAIL,
		"UserID":       user_id_str,
		"Images":       Images,
		"Texts":        Texts,
		"Recognitions": Audios,
	})
}
