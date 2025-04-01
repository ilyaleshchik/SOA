package main

import (
	"encoding/json"
	"net/http"
	"soa-hw-ilyaleshchyk/api"
	acc "soa-hw-ilyaleshchyk/internal/account"
	"soa-hw-ilyaleshchyk/internal/post"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) createPost(c *gin.Context) error {

	var req api.CreatePost

	accountID, err := strconv.ParseInt(c.GetHeader("accountID"), 10, 64)

	if err != nil {
		return err
	}

	account, err := acc.GetAccountByID(s.db, acc.AccountID(accountID))
	if err != nil {
		return err
	}

	if err := c.BindJSON(&req); err != nil {
		return err
	}

	reqTags, err := json.Marshal(req.Tags)
	if err != nil {
		return err
	}

	newPost := &post.Post{
		ID:      uuid.New(),
		Name:    req.Name,
		Desc:    req.Description,
		OwnerID: account.ID,
		Private: req.Private,
		Tags:    reqTags,
	}

	err = post.CreatePost(s.db, newPost)

	if err != nil {
		return err
	}

	newPost.ParsedTags = req.Tags

	c.JSON(http.StatusOK, newPost.ForApi())
	return nil
}

func (s *Server) deletePost(c *gin.Context) error {

	postID, err := uuid.Parse(c.Param("post_id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wrong id"})
		return nil
	}

	curPost, err := post.GetPostByID(s.db, postID)

	if err != nil {
		return err
	}

	if curPost == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wrong id"})
		return nil
	}

	err = post.DeletePost(s.db, curPost)

	return err
}

func (s *Server) updatePost(c *gin.Context) error {

	req := api.PostUpdate{}
	if err := c.BindJSON(&req); err != nil {
		return err
	}

	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wrong id"})
		return nil
	}

	curPost, err := post.GetPostByID(s.db, postID)
	if err != nil {
		return err
	}

	if curPost == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wrong id"})
		return nil
	}

	dataToUpdate := make(map[string]interface{})
	needUpd := false

	if req.Description != nil {
		dataToUpdate["des"] = *req.Description
		curPost.Desc = *req.Description
		needUpd = true
	}

	if req.Name != nil {
		dataToUpdate["name"] = *req.Name
		curPost.Name = *req.Name
		needUpd = true
	}

	if req.Tags != nil {
		newTags, merr := json.Marshal(*req.Tags)
		if merr != nil {
			return merr
		}
		dataToUpdate["tags"] = newTags
		curPost.ParsedTags = *req.Tags
		needUpd = true
	}

	if req.Private != nil {
		dataToUpdate["private"] = *req.Private
		curPost.Private = *req.Private
		needUpd = true
	}

	if needUpd {
		err := post.UpdatePost(s.db, postID, dataToUpdate)
		if err != nil {
			return err
		}
	}

	c.JSON(http.StatusOK, curPost.ForApi())
	return nil
}

func (s *Server) getPost(c *gin.Context) error {

	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wrong id"})
		return nil
	}

	curPost, err := post.GetPostByID(s.db, postID)
	if err != nil {
		return err
	}

	if curPost == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wrong id"})
		return nil
	}

	c.JSON(http.StatusOK, curPost.ForApi())

	return nil
}

func (s *Server) getPosts(c *gin.Context) error {

	accountID, err := strconv.ParseInt(c.GetHeader("accountID"), 10, 64)

	if err != nil {
		return err
	}

	viewer, err := acc.GetAccountByID(s.db, acc.AccountID(accountID))
	if err != nil {
		return err
	}

	p := api.PaginationParams{}
	if err := c.BindQuery(&p); err != nil {
		return err
	}

	ownerID, err := strconv.ParseInt(c.Param("owner_id"), 10, 64)
	if err != nil {
		return err
	}

	lastID := uuid.Nil

	if p.PrevID != "" {
		lastID, err = uuid.Parse(p.PrevID)
		if err != nil {
			return err
		}
	}

	if p.Limit == 0 {
		p.Limit = 30
	}

	totalCnt, err := post.GetTotalPostsCountForViewer(s.db, viewer.ID, acc.AccountID(ownerID))
	if err != nil {
		return err
	}

	posts, err := post.GetPostsForViewer(s.db, lastID, viewer, acc.AccountID(ownerID), p.Limit)
	if err != nil {
		return err
	}

	apiResult := make([]*api.Post, 0, len(posts))
	for _, p := range posts {
		apiResult = append(apiResult, p.ForApi())
	}

	c.JSON(http.StatusOK, gin.H{
		"total_count": totalCnt,
		"posts":       apiResult,
	})
	return nil
}
