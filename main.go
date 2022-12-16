package main

import (
	"fmt"
	"grass-control-go/elements"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		div := elements.Div{}

		text := elements.Text{Value: "TestText"}
		div.Add(text)

		result := div.Render()
		fmt.Fprintf(w, result)
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":8888", nil))

}
