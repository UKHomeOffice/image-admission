package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/stretchr/testify/assert"
)

func TestPutImage(t *testing.T) {
	r, db, err := testSetup()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/images", strings.NewReader(`{"id": "123", "name": "foo", "tags": ["foo", "v1.0"]}`))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestPutMalformedPayload(t *testing.T) {
	r, db, err := testSetup()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/images", strings.NewReader(`{"": "123", "name": "foo", "tags": `))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutUpdate(t *testing.T) {
	r, db, err := testSetup()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/images", strings.NewReader(`{"id": "123", "name": "foo"}`))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.HeaderMap.Get("Content-Type"), "application/json")

	// only tags should be updated
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/images", strings.NewReader(`{"id": "123", "name": "foo/bar", "tags": ["latest"]}`))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.HeaderMap.Get("Content-Type"), "application/json")

	var image Image
	if err := json.Unmarshal(w.Body.Bytes(), &image); err != nil {
		t.Error(err)
	}

	assert.Len(t, image.Tags, 1, "failed to update tags")
}

func testSetup() (*gin.Engine, *gorm.DB, error) {
	// We use sqlite for testing, since we aren't using postgres specific features.
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return &gin.Engine{}, db, err
	}

	db.LogMode(true)

	if err := db.AutoMigrate(&Image{}).Error; err != nil {
		return &gin.Engine{}, db, err
	}

	return newRouter("", db), db, nil
}
