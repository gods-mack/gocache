package main
import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"github.com/syndtr/goleveldb/leveldb"
)

type gocache struct {
	db *leveldb.DB
	
	// params
	port int
	replicas int
	md5sum bool
	volumes []string
}


func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Server is Running...")
}

func main() {
	PORT := 3699
	fmt.Println("GoCache is Running on Port:", PORT)
	fmt.Println("============================")
	
	http.HandleFunc("/", getRoot)
	err := http.ListenAndServe(":"+strconv.Itoa(PORT), nil)
	if err != nil{
		fmt.Println("Unexpected error occured.")
	}
}
