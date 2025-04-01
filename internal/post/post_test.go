package post_test

import (
	"encoding/json"
	"testing"

	"soa-hw-ilyaleshchyk/internal/post"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&post.Post{})
	return db
}

func TestCreatePost(t *testing.T) {
	db := setupTestDB()
	p := &post.Post{
		ID:      uuid.New(),
		Name:    "Test Post",
		Desc:    "Description",
		OwnerID: 1,
		Private: false,
		Tags:    json.RawMessage(`["tag1", "tag2"]`),
	}

	err := post.CreatePost(db, p)
	assert.Nil(t, err)

	var found post.Post
	db.First(&found, "id = ?", p.ID)
	assert.Equal(t, p.Name, found.Name)
}

func TestGetPostByID(t *testing.T) {
	db := setupTestDB()
	p := &post.Post{ID: uuid.New(), Name: "Test", OwnerID: 1, Private: false}
	db.Create(p)

	res, err := post.GetPostByID(db, p.ID)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, p.Name, res.Name)
}

func TestDeletePost(t *testing.T) {
	db := setupTestDB()
	p := &post.Post{ID: uuid.New(), Name: "ToDelete", OwnerID: 1}
	db.Create(p)

	err := post.DeletePost(db, p)
	assert.Nil(t, err)

	var found post.Post
	db.First(&found, "id = ?", p.ID)
	assert.Empty(t, found.ID)
}

func TestUpdatePost(t *testing.T) {
	db := setupTestDB()
	p := &post.Post{ID: uuid.New(), Name: "Old Name", OwnerID: 1}
	db.Create(p)

	newData := map[string]interface{}{"Name": "New Name"}
	err := post.UpdatePost(db, p.ID, newData)
	assert.Nil(t, err)

	var updated post.Post
	db.First(&updated, "id = ?", p.ID)
	assert.Equal(t, "New Name", updated.Name)
}

func TestGetTotalPostsCountForViewer(t *testing.T) {
	db := setupTestDB()
	p1 := &post.Post{ID: uuid.New(), OwnerID: 1, Private: false}
	p2 := &post.Post{ID: uuid.New(), OwnerID: 1, Private: true}
	db.Create(p1)
	db.Create(p2)

	cnt, err := post.GetTotalPostsCountForViewer(db, 2, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), cnt)
}
