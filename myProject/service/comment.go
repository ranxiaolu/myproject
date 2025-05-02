package service

import (
	"errors"
	"myproject/dao"
	"myproject/model"
	"strconv"
)

func GetComment(productID uint) ([]model.Comment, error) {
	db := dao.GetDB()
	var comments []model.Comment
	result := db.Where("product_id = ?", productID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}
func AddComment(userID uint, productID uint, content string) (string, error) {
	db := dao.GetDB()

	comment := model.Comment{
		ProductID: productID,
		UserID:    userID,
		Content:   content,
	}
	result := db.Create(&comment)
	if result.Error != nil {
		return "wrong", result.Error
	}
	return strconv.Itoa(int(comment.ID)), nil
}
func DeleteComment(commentID string) error {
	db := dao.GetDB()
	comment := model.Comment{}
	// 使用 Where 删除指定 ID 的记录
	result := db.Where("comment_id = ?", commentID).Delete(&comment)
	// 检查是否有错误发生
	if result.Error != nil {
		return result.Error
	}
	// 检查是否存在
	if result.RowsAffected == 0 {
		return errors.New("comment not found")
	}
	return nil
}
func UpdateComment(commentID string, content string) error {
	db := dao.GetDB()
	comments := model.Comment{}
	//查询评论是否存在,并将数据库内容返回到comments中
	result := db.First(&comments, "id = ?", commentID)
	if result.Error != nil {
		return result.Error
	}
	//更新评论内容
	comments.Content = content
	//同步到数据库
	result = db.Save(&comments)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func PraiseComment(commentID string, Model int) error {
	db := dao.GetDB()
	comments := model.Comment{}
	//查询评论是否存在,并将数据库内容返回到comments中
	result := db.First(&comments, "id = ?", commentID)
	if result.Error != nil {
		return result.Error
	}
	//更新点赞点踩操作
	comments.CommentModel = Model

	result = db.Save(&comments)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
