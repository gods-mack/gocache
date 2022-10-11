package main
import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"github.com/syndtr/goleveldb/leveldb"
	"flag"
	"strings"
)

type GoCache struct {
	db *leveldb.DB
	
	// params
	port int
	replicas int
	md5sum bool
	volumes []string
	db_path string
}


func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Server is Running...")
}

func getPingPong(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Pong")
}

func  (clnt *GoCache) GetRecord(key []byte) string {
	data, err := clnt.db.Get(key, nil)
	//rec := Record{[]string{}, HARD, ""}
	fmt.Println("DATA ", data)
	if err != leveldb.ErrNotFound {
	//	rec = toRecord(data)
	}
	return string(data)
}

func (clnt *GoCache) PutRecord(key []byte, data []byte) {
	fmt.Println("PurRecord ", key, data)
	err := clnt.db.Put(key, data, nil)
	if err != nil{
		panic("Put didn't work")
	}
}



func (clnt *GoCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := []byte(r.URL.Path)
	log_key := string(key[1:])
	new_key := []byte(log_key)


	fmt.Println("GET_key ", log_key, new_key)

	if r.Method == "GET" {
		rec := clnt.GetRecord(key)
		io.WriteString(w, rec)

	}

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		fmt.Println("r.BODY ", body, err)
		clnt.PutRecord(key, body)
		io.WriteString(w, "success")
	}
	
}




func main() {
	fmt.Println("GoCache is Running...")
	fmt.Println("============================")
	
	port := flag.Int("port", 3699, "Port to serve")
	db_path := flag.String("db_path", "", "Path to key-value db may levelDb")
	replicas := flag.Int("replicas", 3, "Amount of replicas")
	pvolumes := flag.String("volumes", "", "voulmes to use to store , comma seprated")
	md5sum := flag.Bool("md5sum", true, "calculate and store md5 checksum")
	volumes := strings.Split(*pvolumes, ",")
	flag.Parse()

	_a, _b, _c , _d, _e := port, db_path, replicas, volumes, md5sum
	fmt.Println("cmdline params ", *_a, *_b,*_c, _d, *_e)
	fmt.Printf("Server running at 0:%d\n",*port)
	

	db, db_err := leveldb.OpenFile(*db_path, nil)
	if db_err != nil {
		panic(fmt.Sprintf("LevelDB open failed %s", db_err))
	}
	defer  db.Close()  // close db later


	//http.HandleFunc("/", HttpHandler)
	http.HandleFunc("/ping", getPingPong)
	//http.HandleFunc("/get", GETfunc)


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
