# Employee Data Viewer for HR Management System
# View and export employee data from the API

param(
    [string]$Action = "list",
    [int]$EmployeeId = 0,
    [string]$Department = "",
    [string]$Search = "",
    [string]$Status = "",
    [switch]$Export
)

$baseUrl = "http://localhost:8080"
$apiUrl = "$baseUrl/api/v1"

# Login function
function Get-AuthToken {
    try {
        $loginBody = '{"email":"admin@company.com","password":"password123"}'
        $login = Invoke-RestMethod -Uri "$apiUrl/auth/login" -Method POST -Body $loginBody -ContentType "application/json" -ErrorAction Stop
        return $login.data.token
    }
    catch {
        Write-Host "Error: Unable to authenticate. Is the server running?" -ForegroundColor Red
        exit 1
    }
}

# Get all employees
function Get-AllEmployees {
    param($Headers)
    
    $url = "$apiUrl/employees"
    $employees = Invoke-RestMethod -Uri $url -Method GET -Headers $Headers
    return $employees.data
}

# Get single employee
function Get-SingleEmployee {
    param($Headers, $Id)
    
    $url = "$apiUrl/employees/$Id"
    try {
        $employee = Invoke-RestMethod -Uri $url -Method GET -Headers $Headers
        return $employee.data
    }
    catch {
        Write-Host "Error: Employee with ID $Id not found" -ForegroundColor Red
        return $null
    }
}

# Search employees
function Search-Employees {
    param($Headers, $SearchTerm, $DeptId, $EmpStatus)
    
    $url = "$apiUrl/employees?"
    $params = @()
    
    if ($SearchTerm) { $params += "search=$SearchTerm" }
    if ($DeptId) { $params += "department_id=$DeptId" }
    if ($EmpStatus) { $params += "status=$EmpStatus" }
    
    $url += ($params -join "&")
    
    $employees = Invoke-RestMethod -Uri $url -Method GET -Headers $Headers
    return $employees.data
}

# Display employee table
function Show-EmployeeTable {
    param($Employees)
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "EMPLOYEE LIST" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
    
    $Employees | Format-Table @(
        @{Label="ID"; Expression={$_.id}; Width=5}
        @{Label="Code"; Expression={$_.employee_code}; Width=10}
        @{Label="Name"; Expression={"$($_.first_name) $($_.last_name)"}; Width=25}
        @{Label="Position"; Expression={$_.position}; Width=25}
        @{Label="Department"; Expression={$_.department.name}; Width=20}
        @{Label="Status"; Expression={$_.employment_status}; Width=10}
        @{Label="Salary"; Expression={"$" + $_.salary}; Width=10}
    ) -AutoSize
    
    Write-Host "Total: $($Employees.Count) employees" -ForegroundColor Green
}

# Display detailed employee info
function Show-EmployeeDetails {
    param($Employee)
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "EMPLOYEE DETAILS" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan
    
    Write-Host "Personal Information" -ForegroundColor Yellow
    Write-Host "  ID:              $($Employee.id)" -ForegroundColor White
    Write-Host "  Employee Code:   $($Employee.employee_code)" -ForegroundColor White
    Write-Host "  Full Name:       $($Employee.first_name) $($Employee.last_name)" -ForegroundColor White
    Write-Host "  Gender:          $($Employee.gender)" -ForegroundColor White
    Write-Host "  Phone:           $($Employee.phone)" -ForegroundColor White
    Write-Host "  Address:         $($Employee.address)" -ForegroundColor White
    
    Write-Host "`nEmployment Information" -ForegroundColor Yellow
    Write-Host "  Position:        $($Employee.position)" -ForegroundColor White
    Write-Host "  Department:      $($Employee.department.name)" -ForegroundColor White
    Write-Host "  Status:          $($Employee.employment_status)" -ForegroundColor $(if ($Employee.employment_status -eq 'active') { 'Green' } else { 'Red' })
    Write-Host "  Hire Date:       $($Employee.hire_date.Substring(0,10))" -ForegroundColor White
    Write-Host "  Salary:          `$$($Employee.salary)" -ForegroundColor White
    
    Write-Host "`nAccount Information" -ForegroundColor Yellow
    Write-Host "  Email:           $($Employee.user.email)" -ForegroundColor White
    Write-Host "  Role:            $($Employee.user.role)" -ForegroundColor White
    Write-Host "  Account Active:  $($Employee.user.is_active)" -ForegroundColor $(if ($Employee.user.is_active) { 'Green' } else { 'Red' })
    Write-Host "  Created:         $($Employee.created_at.Substring(0,19))" -ForegroundColor White
    Write-Host "  Last Updated:    $($Employee.updated_at.Substring(0,19))" -ForegroundColor White
    
    Write-Host ""
}

# Export to CSV
function Export-EmployeesToCSV {
    param($Employees, $FilePath)
    
    $exportData = $Employees | Select-Object @(
        @{Name="ID"; Expression={$_.id}}
        @{Name="EmployeeCode"; Expression={$_.employee_code}}
        @{Name="FirstName"; Expression={$_.first_name}}
        @{Name="LastName"; Expression={$_.last_name}}
        @{Name="Email"; Expression={$_.user.email}}
        @{Name="Phone"; Expression={$_.phone}}
        @{Name="Gender"; Expression={$_.gender}}
        @{Name="Position"; Expression={$_.position}}
        @{Name="Department"; Expression={$_.department.name}}
        @{Name="Status"; Expression={$_.employment_status}}
        @{Name="Salary"; Expression={$_.salary}}
        @{Name="HireDate"; Expression={$_.hire_date.Substring(0,10)}}
        @{Name="Role"; Expression={$_.user.role}}
    )
    
    $exportData | Export-Csv -Path $FilePath -NoTypeInformation -Encoding UTF8
    Write-Host "`nData exported to: $FilePath" -ForegroundColor Green
}

# Main execution
Write-Host "`n╔══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║         HR Management System - Employee Data Viewer          ║" -ForegroundColor Cyan
Write-Host "╚══════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan

# Authenticate
Write-Host "`nAuthenticating..." -ForegroundColor Yellow
$token = Get-AuthToken
$headers = @{Authorization="Bearer $token"}
Write-Host "✓ Authentication successful`n" -ForegroundColor Green

# Execute action
switch ($Action.ToLower()) {
    "list" {
        $employees = Get-AllEmployees -Headers $headers
        Show-EmployeeTable -Employees $employees
        
        if ($Export) {
            $exportPath = "employee_list_$(Get-Date -Format 'yyyyMMdd_HHmmss').csv"
            Export-EmployeesToCSV -Employees $employees -FilePath $exportPath
        }
    }
    
    "view" {
        if ($EmployeeId -eq 0) {
            Write-Host "Error: Please specify -EmployeeId parameter" -ForegroundColor Red
            exit 1
        }
        
        $employee = Get-SingleEmployee -Headers $headers -Id $EmployeeId
        if ($employee) {
            Show-EmployeeDetails -Employee $employee
        }
    }
    
    "search" {
        if (-not $Search -and -not $Department -and -not $Status) {
            Write-Host "Error: Please specify at least one search parameter (-Search, -Department, or -Status)" -ForegroundColor Red
            exit 1
        }
        
        Write-Host "Searching employees..." -ForegroundColor Yellow
        $employees = Search-Employees -Headers $headers -SearchTerm $Search -DeptId $Department -EmpStatus $Status
        
        if ($employees.Count -gt 0) {
            Show-EmployeeTable -Employees $employees
            
            if ($Export) {
                $exportPath = "employee_search_$(Get-Date -Format 'yyyyMMdd_HHmmss').csv"
                Export-EmployeesToCSV -Employees $employees -FilePath $exportPath
            }
        }
        else {
            Write-Host "`nNo employees found matching your criteria." -ForegroundColor Yellow
        }
    }
    
    "stats" {
        $employees = Get-AllEmployees -Headers $headers
        
        Write-Host "`n========================================" -ForegroundColor Cyan
        Write-Host "EMPLOYEE STATISTICS" -ForegroundColor Cyan
        Write-Host "========================================`n" -ForegroundColor Cyan
        
        Write-Host "Total Employees: $($employees.Count)" -ForegroundColor White
        
        Write-Host "`nBy Department:" -ForegroundColor Yellow
        $employees | Group-Object { $_.department.name } | Sort-Object Count -Descending | ForEach-Object {
            Write-Host "  $($_.Name): $($_.Count) employees" -ForegroundColor White
        }
        
        Write-Host "`nBy Status:" -ForegroundColor Yellow
        $employees | Group-Object employment_status | ForEach-Object {
            $color = if ($_.Name -eq 'active') { 'Green' } else { 'Red' }
            Write-Host "  $($_.Name): $($_.Count) employees" -ForegroundColor $color
        }
        
        Write-Host "`nBy Role:" -ForegroundColor Yellow
        $employees | Group-Object { $_.user.role } | Sort-Object Count -Descending | ForEach-Object {
            Write-Host "  $($_.Name): $($_.Count) employees" -ForegroundColor White
        }
        
        Write-Host "`nSalary Statistics:" -ForegroundColor Yellow
        $salaries = $employees | Select-Object -ExpandProperty salary
        $avgSalary = ($salaries | Measure-Object -Average).Average
        $minSalary = ($salaries | Measure-Object -Minimum).Minimum
        $maxSalary = ($salaries | Measure-Object -Maximum).Maximum
        $totalPayroll = ($salaries | Measure-Object -Sum).Sum
        
        Write-Host "  Average Salary:  `$$([math]::Round($avgSalary, 2))" -ForegroundColor White
        Write-Host "  Minimum Salary:  `$$minSalary" -ForegroundColor White
        Write-Host "  Maximum Salary:  `$$maxSalary" -ForegroundColor White
        Write-Host "  Total Payroll:   `$$totalPayroll" -ForegroundColor White
        
        Write-Host "`nTop 5 Highest Paid Employees:" -ForegroundColor Yellow
        $employees | Sort-Object salary -Descending | Select-Object -First 5 | ForEach-Object {
            Write-Host "  $($_.first_name) $($_.last_name) - $($_.position) - `$$($_.salary)" -ForegroundColor White
        }
        
        Write-Host ""
    }
    
    default {
        Write-Host "Error: Invalid action '$Action'" -ForegroundColor Red
        Write-Host "`nUsage Examples:" -ForegroundColor Yellow
        Write-Host "  .\view_employee_data.ps1 -Action list" -ForegroundColor White
        Write-Host "  .\view_employee_data.ps1 -Action list -Export" -ForegroundColor White
        Write-Host "  .\view_employee_data.ps1 -Action view -EmployeeId 1" -ForegroundColor White
        Write-Host "  .\view_employee_data.ps1 -Action search -Search 'john'" -ForegroundColor White
        Write-Host "  .\view_employee_data.ps1 -Action search -Department 2" -ForegroundColor White
        Write-Host "  .\view_employee_data.ps1 -Action search -Status active" -ForegroundColor White
        Write-Host "  .\view_employee_data.ps1 -Action stats" -ForegroundColor White
        Write-Host ""
    }
}
