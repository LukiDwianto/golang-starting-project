package models

type Privilege struct {
	ID          int64   `gorm:"primaryKey;column:id" json:"id"`
	Name        string  `gorm:"varchar(233);column:name;unique:true" json:"name" `
	Description string  `gorm:"varchar(233);column:description" json:"description"`
	Roles       *[]Role `gorm:"many2many:role_privileges;"`
}

func GetPrivileges() ([]Privilege, error) {
	var privilegeList []Privilege

	err := db.Model(Privilege{}).Find(&privilegeList).Debug().Error

	if err != nil {
		return privilegeList, err
	}

	return privilegeList, nil
}

func CreatePrivilege(privilege Privilege) (Privilege, error) {
	err := db.Save(&privilege).Error
	if err != nil {
		return privilege, err
	}
	return privilege, nil
}

func FindOnePrivilege(id int64) (Privilege, error) {
	var privilege Privilege
	err := db.Model(Privilege{}).Where("id = ?", id).Find(privilege).Error
	if err != nil {
		return privilege, err
	}
	return privilege, nil
}
