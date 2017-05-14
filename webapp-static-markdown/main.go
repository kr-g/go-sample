package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world!")
}

func handlerShowFile(w http.ResponseWriter, r *http.Request) {

	var resp string

	defer func() {
		fmt.Fprintf(w, "%s", resp)
	}()

	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// change url to folder structure
	mp := strings.Replace(r.URL.Path, "/show/", "/markdown/", 1)

	fullpath := path.Join(ex, mp)
	log.Println("using path ", fullpath)

	fileinfo, err := os.Stat(fullpath)
	if err != nil {
		log.Printf("stat error %s path %s", err, fullpath)
		http.NotFound(w, r)
		return
	}

	file, err := os.Open(fullpath)
	if err != nil {
		log.Printf("open error %s path %s", err, fullpath)
		return
	}

	data := make([]byte, fileinfo.Size())
	count, err := io.ReadFull(file, data)
	if int64(count) != fileinfo.Size() {
		log.Printf("read error %s path %s", err, fullpath)
		return
	}

	resp = resp + fmt.Sprintf("%s", data)
}

func handlerMarkdown(w http.ResponseWriter, r *http.Request) {

	var resp string

	defer func() {
		fmt.Fprintf(w, "%s", resp)
	}()

	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fullpath := path.Join(ex, r.URL.Path)
	log.Println("using path ", fullpath)

	fileinfo, err := os.Stat(fullpath)
	if err != nil || fileinfo.IsDir() {
		log.Printf("stat error %s path %s", err, fullpath)
		http.NotFound(w, r)
		return
	}

	data, err := exec.Command("markdown", fullpath).Output()
	if err != nil {
		log.Printf("exec error %s path %s", err, fullpath)
		return
	}

	resp = resp + fmt.Sprintf("%s", data)
}

func main() {
	http.Handle("/stat/", http.StripPrefix("/stat/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/show/", handlerShowFile)
	http.HandleFunc("/markdown/", handlerMarkdown)
	log.Fatal(http.ListenAndServe(":8081", nil))
	log.Println("up and running")
}
