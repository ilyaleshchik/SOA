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

// @Summary Create a new post
// @Description Creates a new post with the provided details
// @Tags posts
// @Accept json
// @Produce json
// @Param request body api.CreatePost true "Post data"
// @Success 200 {object} api.Post
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /posts/create [post]
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

// @Summary Delete a post
// @Description Deletes a post by its ID
// @Tags posts
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /posts/{post_id} [delete]
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

// @Summary Update a post
// @Description Updates the details of an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param request body api.PostUpdate true "Updated post data"
// @Success 200 {object} api.Post
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /posts/{post_id} [patch]
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

// @Summary Get a single post
// @Description Retrieves details of a specific post
// @Tags posts
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} api.Post
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /posts/{post_id} [get]
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

// @Summary Get multiple posts
// @Description Retrieves a paginated list of posts
// @Tags posts
// @Produce json
// @Param owner_id path string true "Owner ID"
// @Param accountID header string true "Viewer Account ID"
// @Param limit query int false "Limit"
// @Param prev_id query string false "Previous Post ID"
// @Success 200 {object} map[string]string "{total_count: int, posts: []api.Post}"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /account/posts/{owner_id} [get]
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
