package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	SERVER = "158.160.85.88" // Основной рабочий сервер
	PORT   = "443"
)

func redirectToTls(w http.ResponseWriter, r *http.Request) {
	log.Println("redirect to TLS server for ", r.RemoteAddr)
	//http.Redirect(w, r, "https://158.160.85.88:443/", http.StatusMovedPermanently)
	http.Redirect(w, r, "https://127.0.0.1:443/", http.StatusMovedPermanently)
}

// Функция обработки запроса в корень
func mainPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection to homePage from:", r.RemoteAddr)
	t, err := template.ParseFiles("pages/html/index.html", "pages/base.html", "pages/styles/homePage.html", "pages/scripts/homePage.html")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}

// Функция обработки запроса на страницу переписывания текста
func rewritePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection to rewritePage from:", r.RemoteAddr)
	t, err := template.ParseFiles("pages/html/rewrite.html", "pages/base.html", "pages/styles/rewritePage.html", "pages/scripts/rewritePage.html")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}

// Функция обработки запроса на страницу обработки текста
func processTextPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection to precessTextPage from:", r.RemoteAddr)
	t, err := template.ParseFiles("pages/html/processText.html", "pages/base.html", "pages/styles/processTextPage.html", "pages/scripts/processTextPage.html")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}

// Функция обработки запроса на страницу обработки изображений
func imagePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection to imagePage from:", r.RemoteAddr)
	t, err := template.ParseFiles("pages/html/image.html", "pages/base.html", "pages/styles/imagePage.html", "pages/scripts/imagePage.html")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
	}
	t.ExecuteTemplate(w, "base", nil)
}

func main() {

	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTls)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

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
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("./node_modules"))))

	log.Printf("Server started at https://%s:%s/\n", SERVER, PORT)

	// Start the main server
	log.Fatal(http.ListenAndServeTLS(":"+PORT, "./ssl/domain.crt", "./ssl/domain.key", nil))
}
