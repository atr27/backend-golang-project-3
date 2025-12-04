package database

import (
	"hr-backend/internal/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed populates the database with initial test data
func Seed() error {
	// Check if data already exists
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	
	if userCount > 0 {
		log.Println("Database already contains data. Skipping seeding.")
		return nil
	}

	log.Println("Seeding database...")

	// Seed departments first
	if err := seedDepartments(); err != nil {
		return err
	}

	// Seed users and employees
	if err := seedUsersAndEmployees(); err != nil {
		return err
	}

	// Seed attendance records
	if err := seedAttendance(); err != nil {
		return err
	}

	// Seed leave balances
	if err := seedLeaveBalances(); err != nil {
		return err
	}

	// Seed leave requests
	if err := seedLeaves(); err != nil {
		return err
	}

	// Seed payroll records
	if err := seedPayroll(); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func seedDepartments() error {
	departments := []models.Department{
		{
			BaseModel:   models.BaseModel{ID: 1},
			Name:        "Engineering",
			Description: "Product development and engineering",
		},
		{
			BaseModel:   models.BaseModel{ID: 2},
			Name:        "Human Resources",
			Description: "HR management and recruitment",
		},
		{
			BaseModel:   models.BaseModel{ID: 3},
			Name:        "Sales",
			Description: "Sales and business development",
		},
		{
			BaseModel:   models.BaseModel{ID: 4},
			Name:        "Marketing",
			Description: "Marketing and brand management",
		},
	}

	for _, dept := range departments {
		if err := DB.Create(&dept).Error; err != nil {
			log.Printf("Failed to seed department %s: %v", dept.Name, err)
			return err
		}
		log.Printf("Seeded department: %s", dept.Name)
	}

	return nil
}

func seedUsersAndEmployees() error {
	// Default password for all test accounts
	defaultPassword := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	hireDate := now.AddDate(-2, 0, 0) // 2 years ago

	// Seed data
	type seedData struct {
		User     models.User
		Employee models.Employee
	}

	seeds := []seedData{
		// Admin user
		{
			User: models.User{
				Email:        "budi.santoso@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "admin",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP001",
				FirstName:        "Budi",
				LastName:         "Santoso",
				Gender:           "Male",
				Phone:            "+628123456789",
				Address:          "Jl. Sudirman No. 123, Jakarta",
				DepartmentID:     uintPtr(2), // HR
				Position:         "System Administrator",
				HireDate:         hireDate,
				EmploymentStatus: "active",
				Salary:           12000000, // Rp 12 juta
			},
		},
		// HR Manager
		{
			User: models.User{
				Email:        "sarah.wijaya@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "hr_manager",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP002",
				FirstName:        "Sarah",
				LastName:         "Wijaya",
				Gender:           "Female",
				Phone:            "+628123456790",
				Address:          "Jl. Thamrin No. 456, Jakarta",
				DepartmentID:     uintPtr(2), // HR
				Position:         "HR Manager",
				HireDate:         hireDate,
				EmploymentStatus: "active",
				Salary:           11000000, // Rp 11 juta
			},
		},
		// Engineering Department Manager
		{
			User: models.User{
				Email:        "ahmad.pratama@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "department_manager",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP003",
				FirstName:        "Ahmad",
				LastName:         "Pratama",
				Gender:           "Male",
				Phone:            "+628123456791",
				Address:          "Jl. Gatot Subroto No. 789, Jakarta",
				DepartmentID:     uintPtr(1), // Engineering
				Position:         "Engineering Manager",
				HireDate:         hireDate,
				EmploymentStatus: "active",
				Salary:           14000000, // Rp 14 juta
			},
		},
		// Software Engineer
		{
			User: models.User{
				Email:        "joko.widodo@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "employee",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP004",
				FirstName:        "Joko",
				LastName:         "Widodo",
				Gender:           "Male",
				Phone:            "+628123456792",
				Address:          "Jl. Kuningan No. 101, Jakarta",
				DepartmentID:     uintPtr(1), // Engineering
				Position:         "Software Engineer",
				HireDate:         hireDate.AddDate(0, 6, 0), // 1.5 years ago
				EmploymentStatus: "active",
				Salary:           9500000, // Rp 9.5 juta
			},
		},
		// Sales Representative
		{
			User: models.User{
				Email:        "siti.nurhaliza@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "employee",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP005",
				FirstName:        "Siti",
				LastName:         "Nurhaliza",
				Gender:           "Female",
				Phone:            "+628123456793",
				Address:          "Jl. Merdeka No. 202, Jakarta",
				DepartmentID:     uintPtr(3), // Sales
				Position:         "Sales Representative",
				HireDate:         hireDate.AddDate(0, 8, 0), // 1 year 4 months ago
				EmploymentStatus: "active",
				Salary:           7500000, // Rp 7.5 juta
			},
		},
		// Marketing Specialist
		{
			User: models.User{
				Email:        "dewi.lestari@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "employee",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP006",
				FirstName:        "Dewi",
				LastName:         "Lestari",
				Gender:           "Female",
				Phone:            "+628123456794",
				Address:          "Jl. Permata No. 303, Jakarta",
				DepartmentID:     uintPtr(4), // Marketing
				Position:         "Marketing Specialist",
				HireDate:         hireDate.AddDate(1, 0, 0), // 1 year ago
				EmploymentStatus: "active",
				Salary:           8000000, // Rp 8 juta
			},
		},
		// Senior Software Engineer
		{
			User: models.User{
				Email:        "andi.susanto@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "employee",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP007",
				FirstName:        "Andi",
				LastName:         "Susanto",
				Gender:           "Male",
				Phone:            "+628123456795",
				Address:          "Jl. Senopati No. 104, Jakarta",
				DepartmentID:     uintPtr(1), // Engineering
				Position:         "Senior Software Engineer",
				HireDate:         hireDate.AddDate(-1, 0, 0), // 3 years ago
				EmploymentStatus: "active",
				Salary:           15000000, // Rp 15 juta (tertinggi)
			},
		},
		// QA Engineer
		{
			User: models.User{
				Email:        "rina.kusuma@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "employee",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP008",
				FirstName:        "Rina",
				LastName:         "Kusuma",
				Gender:           "Female",
				Phone:            "+628123456796",
				Address:          "Jl. Menteng No. 105, Jakarta",
				DepartmentID:     uintPtr(1), // Engineering
				Position:         "QA Engineer",
				HireDate:         hireDate.AddDate(0, 3, 0), // 1 year 9 months ago
				EmploymentStatus: "active",
				Salary:           8500000, // Rp 8.5 juta
			},
		},
		// Junior Developer
		{
			User: models.User{
				Email:        "agus.setiawan@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "employee",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP009",
				FirstName:        "Agus",
				LastName:         "Setiawan",
				Gender:           "Male",
				Phone:            "+628123456797",
				Address:          "Jl. Kemang No. 106, Jakarta",
				DepartmentID:     uintPtr(1), // Engineering
				Position:         "Junior Developer",
				HireDate:         hireDate.AddDate(1, 3, 0), // 9 months ago
				EmploymentStatus: "active",
				Salary:           5000000, // Rp 5 juta (terendah)
			},
		},
		// HR Specialist
		{
			User: models.User{
				Email:        "ayu.kartika@company.com",
				PasswordHash: string(hashedPassword),
				Role:         "employee",
				IsActive:     true,
			},
			Employee: models.Employee{
				EmployeeCode:     "EMP010",
				FirstName:        "Ayu",
				LastName:         "Kartika",
				Gender:           "Female",
				Phone:            "+628123456798",
				Address:          "Jl. Cikini No. 457, Jakarta",
				DepartmentID:     uintPtr(2), // HR
				Position:         "HR Specialist",
				HireDate:         hireDate.AddDate(0, 10, 0), // 1 year 2 months ago
				EmploymentStatus: "active",
				Salary:           7000000, // Rp 7 juta
			},
		},
	}

	// Create users and employees in a transaction
	return DB.Transaction(func(tx *gorm.DB) error {
		for _, seed := range seeds {
			// Create user
			if err := tx.Create(&seed.User).Error; err != nil {
				log.Printf("Failed to create user %s: %v", seed.User.Email, err)
				return err
			}

			// Link employee to user
			seed.Employee.UserID = &seed.User.ID

			// Create employee
			if err := tx.Create(&seed.Employee).Error; err != nil {
				log.Printf("Failed to create employee %s: %v", seed.Employee.EmployeeCode, err)
				return err
			}

			log.Printf("Seeded user: %s (%s)", seed.User.Email, seed.User.Role)
		}

		// Update department managers
		updates := map[uint]uint{
			2: 2, // HR department manager is Sarah Wijaya (EMP002)
			1: 3, // Engineering department manager is Ahmad Pratama (EMP003)
		}

		for deptID, empID := range updates {
			if err := tx.Model(&models.Department{}).Where("id = ?", deptID).Update("manager_id", empID).Error; err != nil {
				log.Printf("Failed to update department manager: %v", err)
				return err
			}
		}

		return nil
	})
}

func uintPtr(n uint) *uint {
	return &n
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// seedAttendance creates sample attendance records for the last 30 days
func seedAttendance() error {
	log.Println("Seeding attendance records...")

	now := time.Now()
	startDate := now.AddDate(0, -1, 0) // Last 30 days

	// Get all active employees
	var employees []models.Employee
	if err := DB.Where("employment_status = ?", "active").Find(&employees).Error; err != nil {
		return err
	}

	attendances := []models.Attendance{}

	// Create attendance for each employee for the last 30 days
	for _, emp := range employees {
		for d := 0; d < 30; d++ {
			date := startDate.AddDate(0, 0, d)
			
			// Skip weekends
			if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
				continue
			}

			// Random attendance patterns
			clockInHour := 8 + (d % 2) // 8 AM or 9 AM
			clockInMinute := (d * 7) % 60
			clockIn := time.Date(date.Year(), date.Month(), date.Day(), clockInHour, clockInMinute, 0, 0, date.Location())
			
			clockOutHour := 17 + (d % 2) // 5 PM or 6 PM
			clockOutMinute := (d * 11) % 60
			clockOut := time.Date(date.Year(), date.Month(), date.Day(), clockOutHour, clockOutMinute, 0, 0, date.Location())

			workingHours := clockOut.Sub(clockIn).Hours()
			overtimeHours := 0.0
			if workingHours > 8 {
				overtimeHours = workingHours - 8
				workingHours = 8
			}

			status := "present"
			notes := ""
			if clockInHour > 8 {
				notes = "Late arrival"
			}

			attendance := models.Attendance{
				EmployeeID:    emp.ID,
				Date:          date,
				ClockIn:       &clockIn,
				ClockOut:      &clockOut,
				WorkingHours:  workingHours,
				OvertimeHours: overtimeHours,
				Status:        status,
				Notes:         notes,
			}
			attendances = append(attendances, attendance)
		}
	}

	// Bulk insert
	if err := DB.Create(&attendances).Error; err != nil {
		log.Printf("Failed to seed attendance: %v", err)
		return err
	}

	log.Printf("Seeded %d attendance records", len(attendances))
	return nil
}

// seedLeaveBalances creates leave balances for all employees
func seedLeaveBalances() error {
	log.Println("Seeding leave balances...")

	// Get all employees
	var employees []models.Employee
	if err := DB.Find(&employees).Error; err != nil {
		return err
	}

	currentYear := time.Now().Year()
	leaveTypes := []struct {
		Type  string
		Total int
	}{
		{"annual", 12},
		{"sick", 12},
		{"casual", 5},
		{"unpaid", 0},
	}

	balances := []models.LeaveBalance{}

	for _, emp := range employees {
		for _, lt := range leaveTypes {
			balance := models.LeaveBalance{
				EmployeeID:    emp.ID,
				LeaveType:     lt.Type,
				TotalDays:     lt.Total,
				UsedDays:      0,
				RemainingDays: lt.Total,
				Year:          currentYear,
			}
			balances = append(balances, balance)
		}
	}

	if err := DB.Create(&balances).Error; err != nil {
		log.Printf("Failed to seed leave balances: %v", err)
		return err
	}

	log.Printf("Seeded %d leave balance records", len(balances))
	return nil
}

// seedLeaves creates sample leave requests
func seedLeaves() error {
	log.Println("Seeding leave requests...")

	now := time.Now()
	
	leaves := []models.Leave{
		// Pending leave - Future
		{
			EmployeeID: 4, // Joko Widodo
			LeaveType:  "annual",
			StartDate:  now.AddDate(0, 0, 10),
			EndDate:    now.AddDate(0, 0, 12),
			TotalDays:  3,
			Reason:     "Liburan keluarga",
			Status:     "pending",
		},
		{
			EmployeeID: 8, // Rina Kusuma
			LeaveType:  "sick",
			StartDate:  now.AddDate(0, 0, 5),
			EndDate:    now.AddDate(0, 0, 6),
			TotalDays:  2,
			Reason:     "Sakit demam",
			Status:     "pending",
		},
		// Approved leaves - Past
		{
			EmployeeID: 7, // Andi Susanto
			LeaveType:  "annual",
			StartDate:  now.AddDate(0, 0, -15),
			EndDate:    now.AddDate(0, 0, -13),
			TotalDays:  3,
			Reason:     "Urusan keluarga",
			Status:     "approved",
			ApprovedBy: uintPtr(3), // Ahmad Pratama (Engineering Manager)
			ApprovedAt: timePtr(now.AddDate(0, 0, -20)),
		},
		{
			EmployeeID: 5, // Siti Nurhaliza
			LeaveType:  "annual",
			StartDate:  now.AddDate(0, 0, -25),
			EndDate:    now.AddDate(0, 0, -23),
			TotalDays:  3,
			Reason:     "Acara keluarga",
			Status:     "approved",
			ApprovedBy: uintPtr(1), // Budi Santoso (Admin)
			ApprovedAt: timePtr(now.AddDate(0, 0, -30)),
		},
		// Rejected leave
		{
			EmployeeID: 9, // Agus Setiawan
			LeaveType:  "annual",
			StartDate:  now.AddDate(0, 0, -5),
			EndDate:    now.AddDate(0, 0, -3),
			TotalDays:  3,
			Reason:     "Liburan mendadak",
			Status:     "rejected",
			ApprovedBy: uintPtr(3), // Ahmad Pratama
			ApprovedAt: timePtr(now.AddDate(0, 0, -10)),
		},
		// More pending leaves
		{
			EmployeeID: 6, // Dewi Lestari
			LeaveType:  "casual",
			StartDate:  now.AddDate(0, 0, 7),
			EndDate:    now.AddDate(0, 0, 7),
			TotalDays:  1,
			Reason:     "Urusan pribadi",
			Status:     "pending",
		},
		{
			EmployeeID: 10, // Ayu Kartika
			LeaveType:  "sick",
			StartDate:  now.AddDate(0, 0, 2),
			EndDate:    now.AddDate(0, 0, 3),
			TotalDays:  2,
			Reason:     "Pemeriksaan kesehatan",
			Status:     "pending",
		},
		// Approved sick leave
		{
			EmployeeID: 4, // Joko Widodo
			LeaveType:  "sick",
			StartDate:  now.AddDate(0, 0, -7),
			EndDate:    now.AddDate(0, 0, -7),
			TotalDays:  1,
			Reason:     "Sakit",
			Status:     "approved",
			ApprovedBy: uintPtr(3), // Ahmad Pratama
			ApprovedAt: timePtr(now.AddDate(0, 0, -8)),
		},
	}

	if err := DB.Create(&leaves).Error; err != nil {
		log.Printf("Failed to seed leaves: %v", err)
		return err
	}

	log.Printf("Seeded %d leave requests", len(leaves))
	return nil
}

// seedPayroll creates sample payroll records
func seedPayroll() error {
	log.Println("Seeding payroll records...")

	// Get all active employees
	var employees []models.Employee
	if err := DB.Where("employment_status = ?", "active").Find(&employees).Error; err != nil {
		return err
	}

	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()
	lastMonth := currentMonth - 1
	lastYear := currentYear
	if lastMonth < 1 {
		lastMonth = 12
		lastYear--
	}

	payrolls := []models.Payroll{}

	// Create payroll for last month (paid)
	for _, emp := range employees {
		basicSalary := emp.Salary
		allowances := basicSalary * 0.10  // 10% tunjangan
		tax := basicSalary * 0.05          // 5% pajak
		deductions := basicSalary * 0.02   // 2% potongan (BPJS, dll)
		netSalary := basicSalary + allowances - tax - deductions

		paymentDate := time.Date(lastYear, time.Month(lastMonth), 25, 10, 0, 0, 0, now.Location())

		payroll := models.Payroll{
			EmployeeID:  emp.ID,
			Month:       lastMonth,
			Year:        lastYear,
			BasicSalary: basicSalary,
			Allowances:  allowances,
			Deductions:  deductions,
			Tax:         tax,
			NetSalary:   netSalary,
			PaymentDate: &paymentDate,
			Status:      "paid",
		}
		payrolls = append(payrolls, payroll)
	}

	// Create payroll for current month (pending)
	for _, emp := range employees {
		basicSalary := emp.Salary
		allowances := basicSalary * 0.10
		tax := basicSalary * 0.05
		deductions := basicSalary * 0.02
		netSalary := basicSalary + allowances - tax - deductions

		payroll := models.Payroll{
			EmployeeID:  emp.ID,
			Month:       currentMonth,
			Year:        currentYear,
			BasicSalary: basicSalary,
			Allowances:  allowances,
			Deductions:  deductions,
			Tax:         tax,
			NetSalary:   netSalary,
			PaymentDate: nil,
			Status:      "pending",
		}
		payrolls = append(payrolls, payroll)
	}

	if err := DB.Create(&payrolls).Error; err != nil {
		log.Printf("Failed to seed payroll: %v", err)
		return err
	}

	log.Printf("Seeded %d payroll records", len(payrolls))
	return nil
}
