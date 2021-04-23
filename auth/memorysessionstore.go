package auth

type MemorySessionStore struct {
	sessions map[string]User
}

func CreateMemorySessionStore() *MemorySessionStore {
	mss := MemorySessionStore{
		sessions: make(map[string]User),
	}
	return &mss
}

func (mss *MemorySessionStore) AddUser(u *User) string {
	cstr := makeSessionCookieString()
	for _,ok := mss.sessions[cstr]; ok; {
		cstr = makeSessionCookieString()
	}

	mss.sessions[cstr] = *u
	return cstr
}

func (mss *MemorySessionStore) RemoveUser(cookie string) {
	delete(mss.sessions, cookie)
}

func (mss *MemorySessionStore) GetUser(cookie string) *User {
	u, ok := mss.sessions[cookie]
	if ok {
		return &u
	} else {
		return nil
	}
}
