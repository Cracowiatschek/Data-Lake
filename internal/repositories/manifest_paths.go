package repositories

import (
	"fmt"
)

func ManifestPath(layer, entity, dt string) string {
	return fmt.Sprintf("%s/%s/dt=%s/_MANIFEST.json", layer, entity, dt)
}

func SuccessPath(layer, entity, dt string) string {
	return fmt.Sprintf("%s/%s/dt=%s/_SUCCESS", layer, entity, dt)
}

func FailedPath(layer, entity, dt string) string {
	return fmt.Sprintf("%s/%s/dt=%s/_FAILED", layer, entity, dt)
}

func InProgressPath(layer, entity, dt string) string {
	return fmt.Sprintf("%s/%s/dt=%s/_INPROGRESS", layer, entity, dt)
}
