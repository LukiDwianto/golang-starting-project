package models

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Role struct {
	ID         int64        `gorm:"primaryKey;column:id" json:"id"`
	Name       string       `gorm:"varchar(233);column:name;unique:true" json:"name" validate:"required" message:"nama role harus ada"`
	Privileges *[]Privilege `gorm:"many2many:role_privileges;" json:"privileges"`
}

type SearchRole struct {
	ID     int64  `json:"id"`
	Name   string `json:"role"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

const TABLENAME_ROLE string = "roles"

func GetRolesPagination(filterRole SearchRole) ([]Role, int64, error) {
	roleList := []Role{}
	var totalSize int64

	dbRole := db.Model(&Role{}).Where("name != ?", "ADMIN")

	if filterRole.Name != "" {
		dbRole = dbRole.Where("name LIKE ?", "%"+filterRole.Name+"%")
	}

	err := dbRole.Count(&totalSize).Debug().Error

	if err != nil {
		log.Println(err.Error())
		return roleList, 0, err
	}

	if filterRole.Limit != 0 {
		dbRole = dbRole.Limit(int(filterRole.Limit))
	}

	if filterRole.Offset != 0 {
		dbRole = dbRole.Offset(int(filterRole.Offset))
	}

	err = dbRole.Find(&roleList).Debug().Error

	if err != nil {
		log.Println(err.Error())
		return roleList, 0, err
	}

	return roleList, totalSize, nil
}

func GetRoles() ([]Role, error) {
	var roleList []Role

	err := db.Table(TABLENAME_ROLE).Where("name != ?", "ADMIN").Find(&roleList).Debug().Error

	if err != nil {
		return roleList, err
	}

	return roleList, nil
}

func CreateRole(role Role, privileges []uint) (Role, error) {

	tx := db.Begin()

	err := tx.Table(TABLENAME_ROLE).Omit(clause.Associations).Create(&role).Error
	if err != nil {
		tx.Rollback()
		return role, err
	}

	var privilegesNew []Privilege

	if err := tx.Model(Privilege{}).Where("id IN ?", privileges).Find(&privilegesNew).Error; err != nil {
		tx.Rollback()
		fmt.Println("Failed get privilege :", err)
		return role, err
	}

	if err := tx.Model(&role).Association("Privileges").Append(&privilegesNew); err != nil {
		tx.Rollback()
		fmt.Println("Failed append privilege :", err)
		return role, err
	}

	err = tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return role, err
	}

	return role, nil
}

func FindOneRole(id int64) (Role, error) {
	var role Role
	err := db.Table(TABLENAME_ROLE).Where("id = ?", id).Preload("Privileges").Find(&role).Error
	if err != nil {
		fmt.Println("error get role : ", err)
		return role, err
	}
	return role, nil
}

func UpdateRole(id int64, role Role) (Role, error) {
	err := db.Model(&role).Omit(clause.Associations).Where("id =?", id).Updates(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func UpdateRolePrivilages(idRole int64, IdPrivileges []int64) error {
	var role Role

	tx := db.Begin()
	if err := tx.Where("id=? and name != ?", idRole, "ADMIN").First(&role).Error; err != nil {
		fmt.Println("Failed get role :", err)
		tx.Rollback()
		return err
	}

	if err := tx.Model(&role).Association("Privileges").Clear(); err != nil {
		fmt.Println("Failed clear privilege :", err)
		tx.Rollback()
		return err
	}

	var privileges []Privilege

	if err := tx.Model(Privilege{}).Where("id IN ?", IdPrivileges).Find(&privileges).Error; err != nil {
		fmt.Println("Failed get privilege :", err)
		tx.Rollback()
		return err
	}

	if err := tx.Model(&role).Association("Privileges").Append(&privileges); err != nil {
		fmt.Println("Failed append privilege :", err)
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func DeleteRole(idRole int64) error {

	var role Role
	var user User

	tx := db.Begin()

	if err := tx.Where("role_id =?", idRole).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			fmt.Println("Failed get user :", err)
			tx.Rollback()
			return err
		}
	}

	if user.ID != 0 {
		tx.Rollback()
		return errors.New("masih ada user dengan role terkait")
	}

	if err := tx.Where("id=? and name != ?", idRole, "ADMIN").First(&role).Error; err != nil {
		fmt.Println("Failed get role :", err)
		tx.Rollback()
		return err
	}

	if err := tx.Model(&role).Association("Privileges").Clear(); err != nil {
		fmt.Println("Failed clear privilege :", err)
		tx.Rollback()
		return err
	}

	if err := tx.Model(&role).Delete(&role).Error; err != nil {
		fmt.Println("Failed clear privilege :", err)
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
