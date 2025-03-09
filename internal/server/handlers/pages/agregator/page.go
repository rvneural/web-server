package text2speech

import (
	formHandler "WebServer/internal/server/handlers/forms/agregator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Page struct {
	handler *formHandler.Agregator
}

func New(handler *formHandler.Agregator) *Page {
	return &Page{
		handler: handler,
	}
}

func (p *Page) GetPage(c *gin.Context) {

	meta := gin.H{}
	meta["OG_TITLE"] = "Нейросетевая служба новостей"
	meta["OG_DESCRIPTION"] = "Нейросети, которые помогают понимать новости"
	meta["OG_IMAGE"] = "https://neuron-nexus.ru/web/static/img/news.jpg"

	c.HTML(http.StatusOK, "agregator.html", meta)
}

func (p *Page) GetPageWithID(c *gin.Context) {
	id := c.Param("id")
	news, err := p.handler.GetNewsData(id)
	meta := gin.H{}
	if err != nil {
		meta["OG_TITLE"] = "Нейросетевая служба новостей"
		meta["OG_DESCRIPTION"] = "Нейросети, которые помогают понимать новости"
		meta["OG_IMAGE"] = "https://neuron-nexus.ru/web/static/img/news.jpg"
	} else {
		meta["OG_TITLE"] = news.Title
		meta["OG_DESCRIPTION"] = news.Description
		enclosure := ""
		for _, enter := range news.Sources {
			if enter.Enclosure != "" {
				enclosure = enter.Enclosure
				break
			}
		}
		if enclosure == "" {
			enclosure = "https://neuron-nexus.ru/web/static/img/news.jpg"
		}
		meta["OG_IMAGE"] = enclosure
	}
	c.HTML(http.StatusOK, "agregator.html", meta)
}
