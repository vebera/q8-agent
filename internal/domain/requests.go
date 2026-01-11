package domain

// MongoDBUserCreateRequest represents a request to create a MongoDB user/db
type MongoDBUserCreateRequest struct {
	Host          string `json:"host"`
	Port          string `json:"port"`
	AdminUser     string `json:"admin_user"`
	AdminPassword string `json:"admin_password"`
	DatabaseName  string `json:"database_name"`
	NewUser       string `json:"new_user"`
	NewPassword   string `json:"new_password"`
}
