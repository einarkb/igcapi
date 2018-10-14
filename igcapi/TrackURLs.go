package igcapi

// global varible of the storage system
var globalTracksDb TrackURLsDB

// TrackURLsDB stores urls and and an id for the each url
type TrackURLsDB struct {
	urls   map[int]string
	nextID int
}

// Init initializes the TrackURLsD object. call once after creating object
func (db *TrackURLsDB) Init() {
	db.nextID = 1
	db.urls = make(map[int]string)
}

// Add inserts and stores a new url
// returns the id of inserted item or -1 if it alreayd existed
func (db *TrackURLsDB) Add(url string) (int, bool) {
	for k, v := range db.urls {
		if v == url {
			return k, false
		}
	}
	db.urls[db.nextID] = url
	db.nextID++
	return db.nextID - 1, true
}

// Get return the url and true if it was found/false if not
func (db *TrackURLsDB) Get(id int) (string, bool) {
	url, exists := db.urls[id]
	return url, exists
}

// GetIDs returns an array with the id of every track
func (db *TrackURLsDB) GetIDs() []int {
	var arr []int
	for id := range db.urls {
		arr = append(arr, id)
	}
	return arr
}
