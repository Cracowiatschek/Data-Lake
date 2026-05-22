package main

import (
	"DataLake/internal/services"
	"fmt"
	"time"
)

func main() {
	dt := time.Now().Format("2006/01/02")

	job := services.NewFetchSensorsService(dt)
	err := job.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Success!")
}
