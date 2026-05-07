package repositories

import (
	"fmt"
)

func manifestPath(layer, entity, dt string) string {
	return fmt.Sprintf("%s/%s/dt=%s/_MANIFEST.json", layer, entity, dt)
}

func successPath(layer, entity, dt string) string {
	return fmt.Sprintf("%s/%s/dt=%s/_SUCCESS", layer, entity, dt)
}

func failedPath(layer, entity, dt string) string {
	return fmt.Sprintf("%s/%s/dt=%s/_FAILED", layer, entity, dt)
}

func inProgressPath(layer, entity, dt string) string {
	return fmt.Sprintf("%s/%s/dt=%s/_INPROGRESS", layer, entity, dt)
}
