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
	"log"
	"sync"
	//"time"
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


func (clnt *GoCache) LockKey(key []byte) bool {
	clnt.mutex.Lock()  // ON mutex

	defer clnt.mutex.Unlock()  // Unlock the mutex at end

	_, is_lock := clnt.key_lock_map[string(key)]
	fmt.Println("LOCK ", is_lock)
	if is_lock { // if lock enabled on this given key
		return false
	}
	clnt.key_lock_map = make(map[string]bool)
	clnt.key_lock_map[string(key)] = true
	return true
}


func (clnt *GoCache) UnlockKey(key []byte)  {
	clnt.mutex.Lock()
	delete(clnt.key_lock_map, string(key))
	clnt.mutex.Unlock()
}


func getPingPong(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Pong")
}

func  (clnt *GoCache) GetRecord(key []byte) string {
	data, err := clnt.db.Get(key, nil)
	if err != leveldb.ErrNotFound {
	//	rec = toRecord(data)
	}
	return string(data)
}

func (clnt *GoCache) PutRecord(key []byte, data []byte) bool {
	err := clnt.db.Put(key, data, nil)
	if err != nil{
		return false 
	}
	//time.Sleep(1*time.Second)
	return true
}



func (clnt *GoCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := []byte(r.URL.Path)


	if  r.Method == "POST" || r.Method == "DELETE" || r.Method == "PUT" {
		enabled := clnt.LockKey(key)
		if !enabled {
			w.WriteHeader(409)
			w.Write([]byte("key (thread LOCKED)"))
			log.Println(r.Method, string(key), 409)
			return
		}
		defer clnt.UnlockKey(key)
	}




	if r.Method == "GET" {
		rec := clnt.GetRecord(key)
		//io.WriteString(w, rec)
		w.WriteHeader(200)
		w.Write([]byte(rec))
		log.Println("GET ",string(key)," [200 OK]")

	}

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil{
			panic("IO readAll didn't work.")
		
		}
		code := 201
		msg := "success"
		if clnt.PutRecord(key, body) != true{
			code = 400
			msg = "failed"
		}
		//io.WriteString(w, "success")
		w.WriteHeader(code)
        w.Write([]byte(msg))
        log.Println("PUT ",string(key), code)

	}
	
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


	http.HandleFunc("/ping", getPingPong)


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
