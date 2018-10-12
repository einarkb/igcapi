package main

// TrackURLsDB stores urls and and an id for the each url
type TrackURLsDB struct {
	urls   map[int]string
	nextID int
}

func (db *TrackURLsDB) Init() {
	db.nextID = 1
	db.urls = make(map[int]string)
}

// Add inserts and stores a new url
// returns the id of inserted item or -1 if it alreayd existed
func (db *TrackURLsDB) Add(url string) int {
	db.urls[db.nextID] = url
	db.nextID++
	return db.nextID - 1
}

// Get return the url and true if it was found/false if not
func (db *TrackURLsDB) Get(id int) (string, bool) {
	url, exists := db.urls[id]
	return url, exists
}
