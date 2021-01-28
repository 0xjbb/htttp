package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var(
	serverDirectory *string

)

func main(){
	dir,err  := os.Getwd()

	if err != nil{
		log.Fatal(err)
	}

	port := flag.String("p", "8080", "Listening port.")
	serverDirectory = flag.String("d", dir, "Serve directory.")

	flag.Parse()

	fmt.Println("[+] Listening on port", *port, "...")
	fmt.Println("[+] Serving directory:",*serverDirectory)
	fmt.Println("[+] Uploads directory:",*serverDirectory)//todo

	http.HandleFunc("/upload", uploadHandler)

	http.Handle("/", http.FileServer(http.Dir(*serverDirectory)))
	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, logRequest(http.DefaultServeMux)))
}

//curl -F file=@test.txt http://localhost:8080/upload
func uploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		fmt.Fprintf(w, "Send it as post.")
		//todo add an upload form, for an easier time in RDP Sessions.
		return
	}

	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")

	if err != nil {
		log.Println(err)
		return
	}

	fileName := fmt.Sprintf("%s/%s", *serverDirectory, handler.Filename)

	fh, err := os.Create(fileName)

	if err != nil{
		log.Println(err)//?fatal???
		return
	}
	defer fh.Close()


	io.Copy(fh, file)
	ipAddr := strings.Split(r.RemoteAddr, ":")[0]
	fmt.Printf("%s %s %s\n", ipAddr, "UPLOAD", handler.Filename)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// maybe ignore /uploads?
		ipAddr := strings.Split(r.RemoteAddr, ":")[0]//remote port.
		fmt.Printf("%s %s %s\n", ipAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}