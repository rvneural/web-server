package user

import (
	"WebServer/internal/server/handlers/interfaces"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
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
	user, err := p.worker.GetUserByID(user_id)
	if err != nil {
		p.logger.Error("Getting user", "error", err)
		user.EMAIL = "неизвестно"
		user.FIRSTNAME = "Неизвестный"
		user.LASTNAME = "Пользователь"
	}

	image_operations, err := p.worker.GetUserOperations(user_id, 0, "image")
	if err != nil {
		p.logger.Error("Getting images operations", "error", err)
	}
	text_operations, err := p.worker.GetUserOperations(user_id, 0, "text")
	if err != nil {
		p.logger.Error("Getting text operations", "error", err)
	}
	audio_operations, err := p.worker.GetUserOperations(user_id, 0, "audio")
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
			Images[i].ID = operation.OPERATION_ID
			Images[i].Placeholder = "Генерация изображения"
			Images[i].Date = operation.CREATION_DATE.Format("02.01.2006 15:04:05")
		}
	}()
	go func() {
		defer wg.Done()
		for i, operation := range text_operations {
			Texts[i].ID = operation.OPERATION_ID
			Texts[i].Placeholder = "Генерация текста"
			Texts[i].Date = operation.CREATION_DATE.Format("02.01.2006 15:04:05")
		}
	}()
	go func() {
		defer wg.Done()
		for i, operation := range audio_operations {
			Audios[i].ID = operation.OPERATION_ID
			Audios[i].Placeholder = "Генерация аудио"
			Audios[i].Date = operation.CREATION_DATE.Format("02.0101.2006 15:04:05")
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
