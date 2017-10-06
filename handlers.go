package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// Image is a representation of an image entry.
type Image struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	ID        string         `json:"id" binding:"required" gorm:"primary_key;type:varchar(64)"`
	Name      string         `json:"name" binding:"required" form:"name"`
	Tags      pq.StringArray `json:"tags,omitempty" gorm:"type:varchar(100)[]" form:"tags"`
}

func getImages(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var image Image
		var images []Image

		// Return a specific image entry.
		if id := c.Param("id"); id != "" {
			r := db.Where(&Image{ID: id}).First(&image)
			if r.RecordNotFound() {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			if r.Error != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.JSON(http.StatusOK, image)
			return
		}

		// Figure out sorting. Default is sorting by updated_at.
		validSortBy := []string{"created_at", "updated_at", "name", "id"}
		sortBy := "updated_at"
		if s := c.Query("sort"); s != "" {
			if func() bool {
				for _, i := range validSortBy {
					if i == s {
						return true
					}
				}
				return false
			}() {
				sortBy = s
			}
		}

		if name := c.Query("name"); name != "" {
			image.Name = c.Query("name")
			if err := db.Where(&image).Order(sortBy + " desc").Find(&images).Error; err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.JSON(http.StatusOK, images)
			return
		}

		// Return a list of all images.
		if err := db.Order(sortBy + " desc").Find(&images).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, images)
	}

	return gin.HandlerFunc(fn)
}

func putImage(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var image Image

		if c.BindJSON(&image) != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Create a new image entry or update tags field of an existing entry.
		if err := db.Where(Image{ID: image.ID}).Assign(Image{Tags: image.Tags}).FirstOrCreate(&image).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Return created or updated record from the database.
		r := Image{}
		db.Where(&Image{ID: image.ID}).First(&r)
		c.JSON(http.StatusOK, r)
	}

	return gin.HandlerFunc(fn)
}

func deleteImage(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		image := Image{ID: c.Param("id")}
		if err := db.Delete(&image).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	}

	return gin.HandlerFunc(fn)
}
