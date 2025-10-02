package models

type Picture struct{
	URL string`json:"url"`
	PostID int `json:"post_id"`
}