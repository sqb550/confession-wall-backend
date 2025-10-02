package postService

import (
	"confession-wall-backend/app/models"
	"confession-wall-backend/config/database"

	"gorm.io/gorm"
)

func ReleasePost(post *models.Post) (int,error) {
	result := database.DB.Create(post)
	return int(post.ID),result.Error
}
func ReleasePicture(pictrue *models.Picture) error {
	result := database.DB.Create(pictrue)
	return result.Error
}
func GetPictures(postID int)([]models.Picture,error){
	pictures:=[]models.Picture{}
	result:=database.DB.Where("post_id=?",postID).Find(&pictures)
	if result.Error!=nil{
		return nil,result.Error
	}
	return pictures,nil
}

func ShowPost(offset int, pageSize int,blockedID []int) ([]models.Post, error) {
	posts := []models.Post{}
    db := database.DB.Where("invisible = ? AND release_status = ?", 0, 1)
    if len(blockedID) > 0 {
        db = db.Where("user_id NOT IN (?)", blockedID)
    }
    result := db.Limit(pageSize).Offset(offset).Find(&posts)
    if result.Error != nil {
        return nil, result.Error
    }
    return posts, nil
    
  
}

func ShowMyPost(userID int) ([]models.Post, error) {
	posts := []models.Post{}
	result := database.DB.Where("user_id=?", userID).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func Delete(id int) error {
	result := database.DB.Where("id=?", id).Delete(&models.Post{})
	return result.Error
}

func Update(postID int, content string) error {
	result := database.DB.Model(&models.Post{}).Where("id=?", postID).Update("content", content)
	return result.Error
}

func Block(block *models.Block) error {
	result := database.DB.Create(block)
	return result.Error
}

func ShowBlock(userID int) ([]models.Block, error) {
	blocks := []models.Block{}
	result := database.DB.Where("user_id=?", userID).Find(&blocks)
	if result.Error != nil {
		return nil, result.Error
	}
	return blocks, nil
}

func Comment(comment *models.Comment) error {
	result := database.DB.Create(comment)
	return result.Error
}

func ShowComments(postID int,blocked []int) ([]models.Comment, error) {
	comments := []models.Comment{}
	result := database.DB.Where("post_id=? AND user_id NOT IN (?)", postID,blocked).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}
func IncrComments(postID int) error {
	result := database.DB.Model(&models.Post{}).Where("id=?", postID).Update("comments", gorm.Expr("comments+?", 1))
	return result.Error
}

func SeekPost(postID int)(*models.Post,error){
	post:=models.Post{}
	result:=database.DB.Where("id=?",postID).First(&post)
	err:=result.Error
	if err!=nil{
		return nil,err
	}
	return &post,nil
}
