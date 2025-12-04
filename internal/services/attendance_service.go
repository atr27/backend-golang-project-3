package services

import (
	"errors"
	"hr-backend/internal/models"
	"hr-backend/internal/repositories"
	"time"
)

type AttendanceService struct {
	attendanceRepo *repositories.AttendanceRepository
	employeeRepo   *repositories.EmployeeRepository
}

func NewAttendanceService(attendanceRepo *repositories.AttendanceRepository, employeeRepo *repositories.EmployeeRepository) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
	}
}

func (s *AttendanceService) ClockIn(req *models.ClockInRequest) (*models.Attendance, error) {
	// Check if employee exists
	_, err := s.employeeRepo.FindByID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Check if already clocked in today
	existing, err := s.attendanceRepo.FindByEmployeeAndDate(req.EmployeeID, req.Date.Time)
	if err == nil && existing.ClockIn != nil {
		return nil, errors.New("already clocked in today")
	}

	attendance := &models.Attendance{
		EmployeeID: req.EmployeeID,
		Date:       req.Date.Time,
		ClockIn:    &req.ClockIn,
		Status:     "present",
	}

	if err := s.attendanceRepo.Create(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *AttendanceService) ClockOut(req *models.ClockOutRequest) (*models.Attendance, error) {
	attendance, err := s.attendanceRepo.FindByEmployeeAndDate(req.EmployeeID, req.Date.Time)
	if err != nil {
		return nil, errors.New("no clock-in record found for today")
	}

	if attendance.ClockOut != nil {
		return nil, errors.New("already clocked out")
	}

	attendance.ClockOut = &req.ClockOut

	// Calculate working hours
	if attendance.ClockIn != nil {
		duration := req.ClockOut.Sub(*attendance.ClockIn)
		hours := duration.Hours()
		attendance.WorkingHours = hours

		// Calculate overtime (assuming 8 hours is standard)
		if hours > 8 {
			attendance.OvertimeHours = hours - 8
		}
	}

	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *AttendanceService) GetAttendanceByEmployee(employeeID uint, month, year int) ([]models.Attendance, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	return s.attendanceRepo.FindByEmployee(employeeID, startDate, endDate)
}

func (s *AttendanceService) GetAttendanceReport(startDate, endDate time.Time) ([]models.Attendance, error) {
	return s.attendanceRepo.FindByDateRange(startDate, endDate)
}

func (s *AttendanceService) CreateManualAttendance(attendance *models.Attendance) (*models.Attendance, error) {
	// Check if employee exists
	_, err := s.employeeRepo.FindByID(attendance.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Check if attendance already exists
	existing, err := s.attendanceRepo.FindByEmployeeAndDate(attendance.EmployeeID, attendance.Date)
	if err == nil && existing.ID > 0 {
		return nil, errors.New("attendance record already exists for this date")
	}

	if err := s.attendanceRepo.Create(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}
