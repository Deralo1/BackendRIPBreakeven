package ds

type User struct {
	UserId      int    `gorm:"primaryKey;column:UserID"`
	Login       string `gorm:"type varchar(100);unique;not null;column:Login"`
	Password    string `gorm:"type varchar(100);not null;column:Password"`
	IsModerator bool   `gorm:"type boolean;not null;column:IsModerator"`
}
