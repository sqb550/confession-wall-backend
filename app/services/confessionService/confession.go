package confessionService

import (
	"confession-wall-backend/app/models"
	"confession-wall-backend/config/database"

	"gorm.io/gorm"
)

func ReleaseConfession(confession *models.Confession) error {
	result := database.DB.Create(confession)
	return result.Error
}

func ShowConfession(offset, pageSize int) ([]models.Confession, error) {
	confessions := []models.Confession{}
	result := database.DB.Where("invisible=?", true).Limit(pageSize).Offset(offset).Find(&confessions)
	if result.Error != nil {
		return nil, result.Error
	}
	return confessions, nil
}

func ShowMyConfession(user_id int) ([]models.Confession, error) {
	confessions := []models.Confession{}
	result := database.DB.Where("user_id=?", user_id).Find(&confessions)
	if result.Error != nil {
		return nil, result.Error
	}
	return confessions, nil
}

func Delete(id int) error {
	result := database.DB.Where("id=?", id).Delete(&models.Confession{})
	return result.Error
}

func Update(confession_id int, content string) error {
	result := database.DB.Model(&models.Confession{}).Where("id=?", confession_id).Update("content", content)
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

func ShowComments(ConfessionID int) ([]models.Comment, error) {
	comments := []models.Comment{}
	result := database.DB.Where("confession_id=?", ConfessionID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}
func IncrComments(ConfessionID int)error{
	result:=database.DB.Model(&models.Confession{}).Where("confession_id=?",ConfessionID).Update("comments",gorm.Expr("comments+?",1))
	return result.Error
}