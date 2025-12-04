package models

type User struct {
	BaseModel
	Email        string    `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"not null" json:"role" binding:"required,oneof=admin hr_manager department_manager employee"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	Employee     *Employee `gorm:"constraint:OnDelete:CASCADE;" json:"employee,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
