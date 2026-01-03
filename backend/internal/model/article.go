package model

type Article struct {
	ID		int 	`gorm:"primerykey" json:"id"`
	Title	string	`json:"title"`
	Body	string	`json:"body"`
	IsPinned bool	`json:"is_pinned"`
}