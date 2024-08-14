package constants

type PrivilegeDto struct {
	Name        string
	Description string
}

var DefaultPrivileges = []PrivilegeDto{
	{
		Name:        "CRUD_USER",
		Description: "Access Update Or Create User",
	},

	{
		Name:        "GET_USERS",
		Description: "Get Users",
	},

	{
		Name:        "UPDATE_ROLE_USERS",
		Description: "Update role users",
	},

	{
		Name:        "GET_ROLE_USERS",
		Description: "Get role users",
	},
}
