package images

import (
	modelImage "WebServer/internal/models/db/results/image"
	"WebServer/internal/server/handlers/interfaces"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Image struct {
	URI      string
	Prompt   string
	Data     string
	Date     string
	Name     string
	LastName string
	USER_ID  int
}

type AdminImages struct {
	dbWorker interfaces.DBWorker
}

func New(dbWorker interfaces.DBWorker) *AdminImages {
	return &AdminImages{dbWorker: dbWorker}
}

func (a *AdminImages) GetPage(c *gin.Context) {
	limit := c.DefaultQuery("limit", "100")
	operations, err := a.dbWorker.GetAllOperations(limit, "image", "")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	id, err := c.Cookie("user_id")
	if err != nil || id == "" {
		id = "-1"
	}
	var user_id = -1
	if id != "-1" {
		user_id, err = strconv.Atoi(id)
		if err != nil {
			user_id = -1
		}
	}
	var user_status = -1
	if user_id != -1 {
		current_user, err := a.dbWorker.GetUserByID(user_id)
		if err != nil {
			user_status = -1
		} else {
			user_status = current_user.USER_STATUS
		}
	}

	images := make([]Image, 0, len(operations))
	for _, operation := range operations {
		if user_status != -1 && operation.USER_STATUS > user_status {
			continue
		}
		img := modelImage.DBResult{}
		err = json.Unmarshal(operation.DATA, &img)
		if err != nil {
			continue
		}

		if img.B64string == "" || operation.IN_PROGRESS {
			continue
		}
		images = append(images, Image{
			URI:      "/operation/" + operation.OPERATION_ID,
			Prompt:   img.Prompt,
			Data:     img.B64string,
			Date:     operation.CREATION_DATE.Format("02.01.2006 15:04"),
			Name:     operation.FIRST_NAME,
			LastName: operation.LAST_NAME,
			USER_ID:  operation.USER_ID,
		})
	}

	style := "/web/styles/admin-images-style.css"
	c.HTML(200, "admin-images.html", gin.H{
		"title":  "Cозданные изображения",
		"Images": images,
		"style":  style,
	})
}
