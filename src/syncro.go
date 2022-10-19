package main

import (
	"fmt"
)

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
