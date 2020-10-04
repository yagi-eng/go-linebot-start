package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// ハンドラの登録
	http.HandleFunc("/", helloHandler)

	fmt.Println("http://localhost:8080 で起動中...")
	// HTTPサーバを起動
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello World!!!!"
	fmt.Fprintf(w, msg)
}
