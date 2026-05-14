package ma

import (
	"fmt"
	"log"

	gios "DataLake/internal/infrastructure/gios"
	httpclient "DataLake/internal/infrastructure/http"
)

func main() {
	http := httpclient.New()
	client := gios.New(http)

	resp, err := client.FetchStations()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TotalPages: %d\n", resp.TotalPages)
	fmt.Printf("Stations count (page): %d\n", len(resp.Stations))

	// pokaż pierwszą stację
	if len(resp.Stations) > 0 {
		for i, x := range resp.Stations {
			fmt.Println(i, x.City)
		}
	}
}
