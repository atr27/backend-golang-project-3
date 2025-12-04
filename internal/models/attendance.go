package models

import (
	"encoding/json"
	"time"
)

// FlexibleDate is a custom type that can unmarshal both date-only strings and full timestamps
type FlexibleDate struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface to handle multiple date formats
func (fd *FlexibleDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Remove quotes
	s = s[1 : len(s)-1]
	
	// Try parsing as date-only first (2006-01-02)
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		fd.Time = t
		return nil
	}
	
	// Try parsing as RFC3339 (full timestamp)
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		fd.Time = t
		return nil
	}
	
	return err
}

// MarshalJSON implements json.Marshaler interface
func (fd FlexibleDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(fd.Time)
}

type Attendance struct {
	BaseModel
	EmployeeID    uint       `gorm:"not null" json:"employee_id" binding:"required"`
	Employee      *Employee  `gorm:"constraint:OnDelete:CASCADE;" json:"employee,omitempty"`
	Date          time.Time  `gorm:"type:date;not null" json:"date" binding:"required"`
	ClockIn       *time.Time `json:"clock_in"`
	ClockOut      *time.Time `json:"clock_out"`
	WorkingHours  float64    `json:"working_hours"`
	OvertimeHours float64    `json:"overtime_hours"`
	Status        string     `gorm:"default:'present'" json:"status"`
	Notes         string     `json:"notes"`
}

type ClockInRequest struct {
	EmployeeID uint         `json:"employee_id" binding:"required"`
	Date       FlexibleDate `json:"date" binding:"required"`
	ClockIn    time.Time    `json:"clock_in" binding:"required"`
}

type ClockOutRequest struct {
	EmployeeID uint         `json:"employee_id" binding:"required"`
	Date       FlexibleDate `json:"date" binding:"required"`
	ClockOut   time.Time    `json:"clock_out" binding:"required"`
}

type AttendanceReport struct {
	EmployeeID    uint    `json:"employee_id"`
	EmployeeName  string  `json:"employee_name"`
	TotalDays     int     `json:"total_days"`
	PresentDays   int     `json:"present_days"`
	AbsentDays    int     `json:"absent_days"`
	LateDays      int     `json:"late_days"`
	TotalHours    float64 `json:"total_hours"`
	OvertimeHours float64 `json:"overtime_hours"`
}
