package lib

import "testing"

func TestLog(t *testing.T) {
	mapData := map[string]string{
		"message": "hello",
	}
	Info("test", mapData, "test")
}