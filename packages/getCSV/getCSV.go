package getCSV

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
)

func getCSVfromRequest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		log.Fatal("Error parsing form data")
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Fatal("Error retrieving file")
	}
	defer file.Close()

	reader := csv.NewReader(file)

	record, err := reader.Read()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(record)

}
