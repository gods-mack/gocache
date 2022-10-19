package main
import (
	"io/ioutil"
	"log"
	//"github.com/gods-mack/gocache/src/main"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"

)





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

	if r.Method == "PUT" {
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
