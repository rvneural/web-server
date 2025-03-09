package user

import (
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"WebServer/internal/models/db/model"
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
	str_id, err := c.Cookie("user_id")
	var user_id int
	if err != nil {
		p.logger.Error("Getting user id", "error", err)
		user_id = 0
	} else {
		user_id, err = strconv.Atoi(str_id)
		if err != nil {
			p.logger.Error("Converting user id", "error", err)
			user_id = 0
		}
	}
	user, err := p.worker.GetUserByID(user_id)
	if err != nil {
		p.logger.Error("Getting user", "error", err)
		user.EMAIL = "неизвестно"
		user.FIRSTNAME = "Неизвестный"
		user.LASTNAME = "Пользователь"
	}
	p.logger.Info("New request to user page")
	var image_operations []model.DBResult
	var text_operations []model.DBResult
	var audio_operations []model.DBResult

	wgGet := sync.WaitGroup{}
	wgGet.Add(3)

	go func() {
		defer wgGet.Done()
		image_operations, err = p.worker.GetUserOperations(user_id, 20, "image")
		if err != nil {
			p.logger.Error("Getting images operations", "error", err)
		}
	}()

	go func() {
		defer wgGet.Done()
		text_operations, err = p.worker.GetUserOperations(user_id, 20, "text")
		if err != nil {
			p.logger.Error("Getting text operations", "error", err)
		}
	}()

	go func() {
		defer wgGet.Done()
		audio_operations, err = p.worker.GetUserOperations(user_id, 20, "audio")
		if err != nil {
			p.logger.Error("Getting audio operations", "error", err)
		}
	}()

	wgGet.Wait()

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

	style := "/web/styles/user-style.css"
	script := "/web/scripts/user-script.js"
	c.HTML(http.StatusOK, "user.html", gin.H{
		"style":        style,
		"script":       script,
		"title":        "Страница пользователя",
		"FirstName":    user.FIRSTNAME,
		"LastName":     user.LASTNAME,
		"Email":        user.EMAIL,
		"UserID":       str_id,
		"Images":       Images,
		"Texts":        Texts,
		"Recognitions": Audios,
	})
}
