package helpers

import (
	"encoding/json"
)

func ToJSON(value any) string {
	data, _ := json.Marshal(value)
	return string(data)
}

func ToPrettyJSON(value any) string {
	data, _ := json.MarshalIndent(value, "", "  ")
	return string(data)
}
