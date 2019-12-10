package main
import (
	"fmt"
	"github.com/gohouse/converter"
)
func main() {
	err := converter.NewTable2Struct().
		SavePath("/Users/mac/Go/src/AS/models/table.go").
		Dsn("root:root@tcp(localhost:3306)/app_package?charset=utf8").
		Run()
	fmt.Println(err)
}