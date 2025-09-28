package postService

import (
	"confession-wall-backend/app/models"
	"confession-wall-backend/config/database"

	"gorm.io/gorm"
)

func Releasepost(post *models.Post) error {
	result := database.DB.Create(post)
	return result.Error
}

func Showpost(offset, pageSize int) ([]models.Post, error) {
	posts := []models.Post{}
	result := database.DB.Where("invisible=? AND release_status=?", true,true).Limit(pageSize).Offset(offset).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func ShowMypost(user_id int) ([]models.Post, error) {
	posts := []models.Post{}
	result := database.DB.Where("user_id=?", user_id).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func Delete(id int) error {
	result := database.DB.Where("id=?", id).Delete(&models.Post{})
	return result.Error
}

func Update(post_id int, content string) error {
	result := database.DB.Model(&models.Post{}).Where("id=?", post_id).Update("content", content)
	return result.Error
}

func Block(block *models.Block) error {
	result := database.DB.Create(block)
	return result.Error
}

func ShowBlock(user_id int) ([]models.Block, error) {
	blocks := []models.Block{}
	result := database.DB.Where("user_id=?", user_id).Find(&blocks)
	if result.Error != nil {
		return nil, result.Error
	}
	return blocks, nil
}

func Comment(comment *models.Comment) error {
	result := database.DB.Create(comment)
	return result.Error
}

func ShowComments(postID int) ([]models.Comment, error) {
	comments := []models.Comment{}
	result := database.DB.Where("post_id=?", postID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}
func IncrComments(postID int) error {
	result := database.DB.Model(&models.Post{}).Where("post_id=?", postID).Update("comments", gorm.Expr("comments+?", 1))
	return result.Error
}
