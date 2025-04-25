package service

import (
	"errors"
	"myproject/dao"
	"myproject/model"
	"strconv"
)

func GetComment(productID uint) ([]model.Comment, error) {
	db := dao.DB
	var comments []model.Comment
	result := db.Where("product_id = ?", productID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}
func AddComment(userID uint, productID uint, content string) (string, error) {
	db := dao.DB
	//
	var product model.Product
	if result := db.First(&product, "ID=?", productID); result.Error != nil {
		return "", errors.New("product not found")
	}
	comment := model.Comment{
		ProductID: productID,
		UserID:    userID,
		Content:   content,
	}
	result := db.Create(&comment)
	if result.Error != nil {
		return " ", result.Error
	}
	return strconv.Itoa(int(comment.ID)), nil
}
func DeleteComment(commentID string) error {
	db := dao.DB
	var comment = model.Comment{}
	// 使用 Where 删除指定 ID 的记录
	result := db.Where("id = ?", commentID).Delete(&comment)
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
	db := dao.DB
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
func PraiseComment(commentID uint, Model int) error {
	db := dao.DB

	comment := model.Comment{
		CommentModel: Model,
	}
	//查询评论是否存在,并将数据库内容返回到comments中
	result := db.Where("id = ?", commentID).First(&comment)
	if result.Error != nil {
		return result.Error
	}
	//更新model
	if Model == 1 {
		comment.PraiseCount++
		comment.BadModel--
		comment.CommentModel = 1
	}
	if Model == 2 {
		comment.PraiseCount--
		comment.BadModel++
		comment.CommentModel = 2
	}
	//保存更新后的评论
	result = db.Save(&comment)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
