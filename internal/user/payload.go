package user

import "time"

type UserModel struct {
	Id         int64
	TelegramId int64
	Username   string
	CreatedAt  time.Time
	UpdateAt   time.Time
}
