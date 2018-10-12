package main

// TrackUrlsDB stores urls and and an id for the each url
type TrackUrlsDB struct {
	urls   map[int]string
	ids    map[string]int
	nextID int
}

// Add inserts and stores a new url
// returns the id of inserted item or -1 if it alreayd existed
func (db *TrackUrlsDB) Add(url string) int {
	_, exists := db.ids[url]
	if exists {
		return -1
	}
	db.urls[db.nextID] = url
	db.ids[url] = db.nextID
	db.nextID++
	return db.nextID - 1
}

// Get return the url and true if it was found/false if not
func (db *TrackUrlsDB) Get(id int) (string, bool) {
	url, exists := db.urls[id]
	return url, exists
}
