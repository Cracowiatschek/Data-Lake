package main

import (
	"DataLake/internal/services"
	"fmt"
)

func main() {

	s := services.NewLoaderService()
	fmt.Println(s.LoadLeatestMeasurementsFromBronze(0))
}
