package models

type Department struct {
	BaseModel
	Name        string     `gorm:"uniqueIndex;not null" json:"name" binding:"required"`
	Description string     `json:"description"`
	ManagerID   *uint      `json:"manager_id"`
	Manager     *Employee  `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Employees   []Employee `gorm:"foreignKey:DepartmentID" json:"employees,omitempty"`
}
