package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Функция обработки запроса в корень
func mainPage(w http.ResponseWriter, r *http.Request) {

}

// Функция обработки запроса на страницу переписывания текста
func rewritePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection to rewritePage from:", r.RemoteAddr)
	t, err := template.ParseFiles("pages/html/page.html", "pages/base.html", "pages/styles/style.html", "pages/scripts/style.html")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}

// Функция обработки запроса на страницу обработки текста
func processTextPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection to precessTextPage from:", r.RemoteAddr)
	t, err := template.ParseFiles("pages/html/page.html", "pages/base.html", "pages/styles/style.html", "pages/scripts/style.html")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}

// Функция обработки запроса на страницу обработки изображений
func imagePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection to imagePage from:", r.RemoteAddr)
	t, err := template.ParseFiles("pages/html/page.html", "pages/base.html", "pages/styles/style.html", "pages/scripts/style.html")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}

func main() {

	// Register handlers for routes
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/rewrite", rewritePage)
	http.HandleFunc("/text", processTextPage)
	http.HandleFunc("/image", imagePage)

	// Register handlers for Main Server
	http.HandleFunc("/recognize", recognizeFromWeb)
	http.HandleFunc("/rewriteFromWeb", rewriteFromWeb)
	http.HandleFunc("/processTextFromWeb", processTextFromWeb)
	http.HandleFunc("/generateImage", generateImageFromWeb)

	// Serve static files
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./pages"))))
}
