package post

import (
	"encoding/json"
	"fmt"
	"soa-hw-ilyaleshchyk/api"
	acc "soa-hw-ilyaleshchyk/internal/account"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID         uuid.UUID `gorm:"primaryKey"`
	Name       string
	Desc       string
	OwnerID    acc.AccountID `gorm:"index"`
	Private    bool
	Tags       json.RawMessage `gorm:"type:jsonb"`
	ParsedTags []string        `gorm:"-"`
}

func (p *Post) ForApi() *api.Post {
	return &api.Post{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Desc,
		OwnerID:     int64(p.OwnerID),
		Tags:        p.ParsedTags,
		DateCreate:  p.CreatedAt.Unix(),
	}
}

func (p *Post) ParseTags() error {
	if p.Tags == nil {
		p.ParsedTags = make([]string, 0)
		return nil
	}
	err := json.Unmarshal(p.Tags, &p.ParsedTags)
	if err != nil {
		return err
	}
	return nil
}

func CreatePost(db *gorm.DB, post *Post) error {
	err := db.Create(&post).Error
	return err
}

func GetPostByID(db *gorm.DB, postID uuid.UUID) (*Post, error) {
	var res *Post
	fres := db.Model(&Post{}).Where("id = ?", postID).Find(&res)

	if fres.Error != nil {
		return nil, fres.Error
	}

	if fres.RowsAffected == 0 {
		return nil, nil
	}

	if err := res.ParseTags(); err != nil {
		return nil, err
	}

	return res, nil
}

func DeletePost(db *gorm.DB, post *Post) error {
	return db.Unscoped().Delete(&post).Error
}

func UpdatePost(db *gorm.DB, postID uuid.UUID, newData map[string]interface{}) error {
	return db.Model(&Post{}).Where("id = ?", postID).Updates(newData).Error
}

func GetPostsForViewer(db *gorm.DB, lastID uuid.UUID, viewer *acc.Account, ownerID acc.AccountID, limit int) ([]*Post, error) {
	query := "id > ? AND owner_id = ?"
	queryArgs := make([]interface{}, 0)
	queryArgs = append(queryArgs, lastID)
	queryArgs = append(queryArgs, ownerID)
	if viewer.ID != acc.AccountID(ownerID) {
		query += " AND private = ?"
		queryArgs = append(queryArgs, false)
	}

	var posts []*Post
	fmt.Println("limit:", limit)

	err := db.Model(&Post{}).Where(query, queryArgs...).Limit(limit).Find(&posts).Order("id ASC").Error
	if err != nil {
		return nil, err
	}

	for _, p := range posts {
		if err := p.ParseTags(); err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func GetTotalPostsCountForViewer(db *gorm.DB, viewerID, ownerID acc.AccountID) (int64, error) {

	var cnt int64

	query := "owner_id = ?"
	qArgs := []interface{}{ownerID}
	if viewerID != ownerID {
		query += " AND private = ?"
		qArgs = append(qArgs, false)
	}

	err := db.Model(&Post{}).Where(query, qArgs...).Count(&cnt).Error
	if err != nil {
		return 0, err
	}

	return cnt, nil
}
