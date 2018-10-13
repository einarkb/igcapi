package main

import (
	"strconv"
	"testing"
)

func Test_Init(t *testing.T) {
	db := TrackURLsDB{}
	db.Init()
	if db.nextID != 1 {
		t.Error("nextId not initialized to 1")
		return
	}
	if db.urls == nil {
		t.Error("map not initialized")
		return
	}
}

func Test_Add(t *testing.T) {
	db := TrackURLsDB{}
	db.nextID = 1
	db.urls = make(map[int]string)
	db.Add("test1.com")
	url, exists := db.urls[db.nextID-1]
	if !exists {
		t.Error("test1.com was not added")
		return
	}
	if url != "test1.com" {
		t.Error("expected 'test1.com' but got '" + url + "'")
		return
	}

	db.nextID = 1
	db.urls = make(map[int]string)
	for i := 1; i < 10; i++ {
		db.Add("test" + strconv.Itoa(i) + ".com")
	}
	for i := 1; i < 10; i++ {
		url, exists := db.urls[i]
		if !exists {
			t.Error("No track found with id: " + strconv.Itoa(i))
			return
		}
		if url != ("test" + strconv.Itoa(i) + ".com") {
			t.Error("Expected test" + strconv.Itoa(i) + ".com, got " + url)
		}
	}
}

func Test_Get(t *testing.T) {
	db := TrackURLsDB{}
	db.nextID = 1
	db.urls = make(map[int]string)
	db.urls[1] = "test.com"
	url, exists := db.Get(1)
	if !exists {
		t.Error("Failed to get track from index 1")
		return
	}
	if url != "test.com" {
		t.Error("Expected 'test.com' but got '" + url + "'")
	}
}
