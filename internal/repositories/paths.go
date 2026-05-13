package repositories

import "fmt"

func PageBasedPathJSON(layer, entity, dt string, page int) string {
	return fmt.Sprintf("%s/%s/dt=%s/page=%d.json", layer, entity, dt, page)
}

func AttributeBasedPathJSON(layer, entity, dt, name, value string) string {
	return fmt.Sprintf("%s/%s/dt=%s/%s=%s.json", layer, entity, dt, name, value)
}

func BatchPathJSON(layer, entity, dt string) string {
	return fmt.Sprintf("%s/%s/dt=%s/data.json.gz", layer, entity, dt)
}

func PathParquet(layer, entity, dt, filename string) string {
	return fmt.Sprintf("%s/%s/dt=%s/%s.parquet", layer, entity, dt, filename)
}

func PathJson(layer, entity, dt, filename string) string {
	return fmt.Sprintf("%s/%s/dt=%s/%s.json", layer, entity, dt, filename)
}
