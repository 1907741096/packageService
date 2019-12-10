package lib

import (
	"testing"
	"strings"
	"log"
)

func TestPath(t *testing.T) {
	str := strings.Title("hello")
	log.Print(str)
	t.Error(str)
}