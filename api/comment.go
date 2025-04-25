package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"myproject/service"
	"net/http"
	"strconv"
)

func GetComment(ctx context.Context, c *app.RequestContext) {
	productStrID := c.Param("product_id")
	productID, err := strconv.ParseUint(productStrID, 10, 64)
	if err != nil {
		c.JSON(400, utils.H{"message": err})
		return
	}
	comments, err := service.GetComment(uint(productID))
	if err != nil {
		c.JSON(400, utils.H{"info": err.Error()})
		return
	}
	c.JSON(200, utils.H{"status": 10000, "info": "success", "comments": comments})
}

// AddComment 添加评论
func AddComment(ctx context.Context, c *app.RequestContext) {
	productIDStr := c.Param("product_id")
	var commentData struct {
		Content string `json:"content"`
		UserID  uint   `json:"user_id"`
	}

	if err := c.BindJSON(&commentData); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"info": "invalid request"})
		return
	}
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(400, utils.H{"message": err})
		return
	}
	commentID, err := service.AddComment(commentData.UserID, uint(productID), commentData.Content)
	if err != nil {
		c.JSON(400, utils.H{"info": "wrong comment"})
		return
	}
	c.JSON(200, utils.H{"info": "success", "status": 10000, "data": commentID})
}

// DeleteComment 删除评论
func DeleteComment(ctx context.Context, c *app.RequestContext) {
	commentID := c.Param("comment_id")
	err := service.DeleteComment(commentID)
	if err != nil {
		c.JSON(400, utils.H{"info": "wrong"})
		return
	}
	c.JSON(200, utils.H{"info": "success", "status": 10000})
}

func UpdateComment(ctx context.Context, c *app.RequestContext) {
	commentID := c.Param("comment_id")
	var commentData struct {
		Content string `json:"content"`
	}
	if err := c.BindJSON(&commentData); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"info": "invalid request"})
		return
	}
	if err := service.UpdateComment(commentID, commentData.Content); err != nil {
		c.JSON(400, utils.H{"info": "failed"})
		return
	}
	c.JSON(200, utils.H{"info": "success", "status": 10000})
}
func PraiseComment(ctx context.Context, c *app.RequestContext) {

	var model struct {
		UserID       uint `json:"user_id"`
		CommentID    uint `json:"comment_id"`
		CommentModel int  `json:"model"`
	}
	// 绑定请求体数据到model
	if err := c.BindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, utils.H{"info": "invalid request"})
		return
	}

	//点赞点踩操作是否成功
	if err := service.PraiseComment(model.CommentID, model.CommentModel); err != nil {
		c.JSON(400, utils.H{"info": err.Error()})
		return
	}

	c.JSON(200, utils.H{"info": "success", "status": 10000})
}
