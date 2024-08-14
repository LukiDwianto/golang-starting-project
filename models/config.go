package models

import (
	"fmt"
	"golang-starting-project/constants"
	"golang-starting-project/middleware"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func UpdateAdminRoles() {
	var role Role = Role{Name: constants.DefaultRole}

	if err := db.Where("name=?", role.Name).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				fmt.Println("Failed to create role:", err)
			} else {
				fmt.Println("Created Role:", role.Name)
			}
		}
	}

	if err := db.Model(&role).Association("Privileges").Clear(); err != nil {
		fmt.Println("Failed clear privilege :", err)
	}

	var privileges []Privilege

	if err := db.Model(Privilege{}).Find(&privileges).Error; err != nil {
		fmt.Println("Failed get privilege :", err)
	}

	if err := db.Model(&role).Association("Privileges").Append(&privileges); err != nil {
		fmt.Println("Failed append privilege :", err)
	}

	log.Println("Success update role admin \n")
}

func UpdateUserAdmin() {

	username := "admin"

	user, err := CheckUserExist(username)

	if err != nil {
		fmt.Println("Err on check user :", err)
	}

	if user.ID == 0 {

		var role Role = Role{Name: constants.DefaultRole}

		if err := db.Where("name=?", role.Name).First(&role).Error; err != nil {
			fmt.Println("Error get role :", err)
		}

		var user User
		user.UserName = "admin"
		user.RoleID = role.ID
		user.IsActive = true
		hashedPassword, err := middleware.HashPassword(os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			log.Fatal(err.Error())
		}

		user.Password = hashedPassword
		_, err = CreateUser(user)

		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Println("Success update useradmin")
}

func InitDB() {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	log.Println("Connecting DB")

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Fatal("Failed to connect to database!", err.Error())
	}

	db.AutoMigrate(
		&User{},
		&Role{},
		&Privilege{},
	)

	sqlDb, err := db.DB()

	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	sqlDb.SetMaxOpenConns(8)
	sqlDb.SetMaxIdleConns(4)
	sqlDb.SetConnMaxLifetime(time.Hour)

	UpdatePrivileges()
	UpdateAdminRoles()
	UpdateUserAdmin()

	fmt.Println("Database connected successfully!")

}

func UpdatePrivileges() {
	for _, privilegeData := range constants.DefaultPrivileges {
		var privilege Privilege

		if err := db.Where("name = ?", privilegeData.Name).First(&privilege).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				newPrivilege := Privilege{
					Name:        privilegeData.Name,
					Description: privilegeData.Description,
				}
				if err := db.Create(&newPrivilege).Error; err != nil {
					fmt.Println("Failed to create privilege:", err)
				} else {
					fmt.Println("Created privilege:", newPrivilege.Name)
				}
			} else {
				fmt.Println("Failed to query privilege:", err)
			}
		} else {
			fmt.Println("Privilege already exists:", privilegeData.Name)
		}
	}
}
