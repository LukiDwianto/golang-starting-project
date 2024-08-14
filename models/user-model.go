package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SearchUser struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	RoleID   int64  `json:"role_id"`
	RoleName string `json:"role_name"`
	Limit    int64  `json:"limit"`
	Offset   int64  `json:"offset"`
}

type User struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"id"`
	UserName    string    `gorm:"varchar(233);column:username;unique:true" json:"username"`
	Password    string    `gorm:"varchar(233);column:password;omitempty" json:"password"`
	IsActive    bool      `gorm:"column:is_active" json:"is_active"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	LastUpdated time.Time `gorm:"column:last_updated" json:"last_updated"`
	RoleID      int64     `gorm:"column:role_id" json:"role_id"`
	Role        *Role     `gorm:"foreignKey:RoleID" json:"role"`
	CreatedBy   string    `gorm:"varchar(100);column:created_by" json:"created_by"`
	UpdatedBy   string    `gorm:"varchar(100);column:updated_by" json:"updated_by"`
}

const TABLENAME_USER string = "users"

func GetUsers(filterUser SearchUser) ([]User, int64, error) {
	userList := []User{}
	var totalSize int64

	dbUser := db.Table(TABLENAME_USER)

	dbUser = dbUser.Where("username != ?", "admin")

	if filterUser.UserName != "" {
		dbUser = dbUser.Where("username LIKE ?", "%"+filterUser.UserName+"%")
	}

	if filterUser.RoleID != 0 {
		dbUser = dbUser.Where("role_id =?", filterUser.RoleID)
	}

	err := dbUser.Count(&totalSize).Debug().Error

	if err != nil {
		log.Println(err.Error())
		return userList, 0, err
	}

	if filterUser.Limit != 0 {
		dbUser = dbUser.Limit(int(filterUser.Limit))
	}

	if filterUser.Offset != 0 {
		dbUser = dbUser.Offset(int(filterUser.Offset))
	}

	err = dbUser.Preload("Role").Omit("password").Find(&userList).Debug().Error

	if err != nil {
		log.Println(err.Error())
		return userList, 0, err
	}

	return userList, totalSize, nil
}

func CheckUpdateUSerExist(username string, id int64) (User, error) {
	var user User
	err := db.Table(TABLENAME_USER).Where("username=? AND id != ?", username, id).Find(&user).Error
	if err != nil {
		log.Println(err.Error())
		if err != gorm.ErrRecordNotFound {
			return user, err
		}
	}
	return user, nil
}

func CreateUser(user User) (User, error) {
	user.CreatedAt = time.Now()
	err := db.Table(TABLENAME_USER).Omit(clause.Associations).Create(&user).Debug().Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func CheckUserExist(username string) (User, error) {
	var user User
	err := db.Table(TABLENAME_USER).Where("username=?", username).Find(&user).Error
	if err != nil {
		log.Println(err.Error())
		if err != gorm.ErrRecordNotFound {
			return user, err
		}
	}
	return user, nil
}


func UpdateUser(id int64, user User) (User, error) {
	user.LastUpdated = time.Now()
	fmt.Println(user)
	err := db.Debug().Model(&user).Where("id =? AND username != ? ", id, "ADMIN").Omit(clause.Associations).
		Updates(
			map[string]interface{}{
				"username":  user.UserName,
				"role_id":   user.RoleID,
				"is_active": user.IsActive}).
		Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func UpdatePasswordUser(id int64, user User) (User, error) {
	user.LastUpdated = time.Now()
	err := db.Model(&user).Where("id =? AND id != 1 ", id).Omit(clause.Associations).Updates(
		User{
			Password: user.Password,
		}).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func UpdateRoleUser(user User) error {
	err := db.Table(TABLENAME_USER).Where("username=?", user.UserName).Updates(User{RoleID: user.RoleID, LastUpdated: time.Now(), UpdatedBy: user.UpdatedBy}).Debug().Error
	if err != nil {
		return err
	}
	return nil
}


func GetUserByUsername(username string) (User, error) {
	var user User
	err := db.Table(TABLENAME_USER).Where("username=?", username).Preload("Role").Preload("Role.Privileges").Find(&user).Debug().Error
	if err != nil {
		return user, err
	}
	return user, nil
}
