package main




func (clnt *GoCache) LockKey(key []byte) bool {
	clnt.mutex.Lock()  // ON mutex

	defer clnt.mutex.Unlock()  // Unlock the mutex at end

	_, is_lock := clnt.key_lock_map[string(key)]
	if is_lock { // if lock enabled on this given key
		return false
	}
	return true
}


func (clnt *GoCache) UnlockKey(key []byte) bool {
	clnt.mutex.Lock()
	delete(clnt.key_lock_map, string(key))
	clnt.mutex.Unlock()
}
