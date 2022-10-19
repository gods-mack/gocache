package main
import (
	"fmt"
	"io"
	//"io/ioutil"
	"net/http"
	"strconv"
	"github.com/syndtr/goleveldb/leveldb"
	"flag"
	"strings"
	//"log"
	"sync"
	//"time

)

type GoCache struct {
	db *leveldb.DB
	 
	mutex sync.Mutex
	key_lock_map map[string]bool
	// params
	port int
	replicas int
	md5sum bool
	volumes []string
	db_path string
}


func getPingPong(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Pong")
}




func main() {
	fmt.Println("====== GoCache is Running... =======")
	
	port := flag.Int("port", 3699, "Port to serve")
	db_path := flag.String("db_path", "", "Path to key-value db may levelDb")
	replicas := flag.Int("replicas", 3, "Amount of replicas")
	pvolumes := flag.String("volumes", "", "voulmes to use to store , comma seprated")
	md5sum := flag.Bool("md5sum", true, "calculate and store md5 checksum")
	volumes := strings.Split(*pvolumes, ",")
	flag.Parse()

	fmt.Printf("Server running at 0:%d\n",*port)
	

	db, db_err := leveldb.OpenFile(*db_path, nil)
	if db_err != nil {
		panic(fmt.Sprintf("LevelDB open failed %s", db_err))
	}
	defer  db.Close()  // close db later


	//http.HandleFunc("/ping", getPingPong)


	client := GoCache{
			db: db,
			port: *port,
			db_path: *db_path,
			replicas: *replicas,
			volumes: volumes,
			md5sum: *md5sum,
	}

	err := http.ListenAndServe(":"+strconv.Itoa(*port), &client)
	if err != nil{
		fmt.Println("Unexpected error occured.")
	}
}
