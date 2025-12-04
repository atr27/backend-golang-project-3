package utils

import (
	"bytes"
	"fmt"
	"hr-backend/internal/models"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePayrollPDF generates a PDF payroll slip for a given payroll record
func GeneratePayrollPDF(payroll *models.Payroll) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set margins
	pdf.SetMargins(20, 20, 20)

	// Header - Company Name
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(31, 41, 55) // gray-800
	pdf.Cell(0, 10, "SLIP GAJI")
	pdf.Ln(12)

	// Employee Information Section
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(75, 85, 99) // gray-600
	pdf.Cell(0, 8, "Informasi Karyawan")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(31, 41, 55)

	// Check if Employee is loaded
	if payroll.Employee == nil {
		return nil, fmt.Errorf("employee data not loaded for payroll ID %d", payroll.ID)
	}

	// Employee details
	employeeName := fmt.Sprintf("%s %s", payroll.Employee.FirstName, payroll.Employee.LastName)
	pdf.Cell(50, 6, "Nama")
	pdf.Cell(5, 6, ":")
	pdf.Cell(0, 6, employeeName)
	pdf.Ln(6)

	pdf.Cell(50, 6, "Kode Karyawan")
	pdf.Cell(5, 6, ":")
	pdf.Cell(0, 6, payroll.Employee.EmployeeCode)
	pdf.Ln(6)

	if payroll.Employee.Department != nil {
		pdf.Cell(50, 6, "Departemen")
		pdf.Cell(5, 6, ":")
		pdf.Cell(0, 6, payroll.Employee.Department.Name)
		pdf.Ln(6)
	}

	pdf.Cell(50, 6, "Posisi")
	pdf.Cell(5, 6, ":")
	pdf.Cell(0, 6, payroll.Employee.Position)
	pdf.Ln(10)

	// Period Information
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(75, 85, 99)
	pdf.Cell(0, 8, "Periode Gaji")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(31, 41, 55)

	monthNames := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	period := fmt.Sprintf("%s %d", monthNames[payroll.Month], payroll.Year)

	pdf.Cell(50, 6, "Bulan/Tahun")
	pdf.Cell(5, 6, ":")
	pdf.Cell(0, 6, period)
	pdf.Ln(10)

	// Salary Details - Table
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(75, 85, 99)
	pdf.Cell(0, 8, "Rincian Gaji")
	pdf.Ln(8)

	// Table header
	pdf.SetFillColor(249, 250, 251) // gray-50
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(31, 41, 55)
	pdf.CellFormat(100, 8, "Keterangan", "1", 0, "L", true, 0, "")
	pdf.CellFormat(70, 8, "Jumlah (Rp)", "1", 1, "R", true, 0, "")

	// Table body
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)

	// Basic Salary
	pdf.CellFormat(100, 8, "Gaji Pokok", "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 8, formatCurrency(payroll.BasicSalary), "1", 1, "R", false, 0, "")

	// Allowances
	pdf.CellFormat(100, 8, "Tunjangan", "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 8, formatCurrency(payroll.Allowances), "1", 1, "R", false, 0, "")

	// Gross Salary
	grossSalary := payroll.BasicSalary + payroll.Allowances
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(249, 250, 251)
	pdf.CellFormat(100, 8, "Total Gaji Kotor", "1", 0, "L", true, 0, "")
	pdf.CellFormat(70, 8, formatCurrency(grossSalary), "1", 1, "R", true, 0, "")

	// Deductions
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(100, 8, "Potongan", "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 8, formatCurrency(payroll.Deductions), "1", 1, "R", false, 0, "")

	// Tax
	pdf.CellFormat(100, 8, "Pajak", "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 8, formatCurrency(payroll.Tax), "1", 1, "R", false, 0, "")

	// Total Deductions
	totalDeductions := payroll.Deductions + payroll.Tax
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(249, 250, 251)
	pdf.CellFormat(100, 8, "Total Potongan", "1", 0, "L", true, 0, "")
	pdf.CellFormat(70, 8, formatCurrency(totalDeductions), "1", 1, "R", true, 0, "")

	// Net Salary - Highlighted
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(220, 252, 231) // green-100
	pdf.SetTextColor(22, 101, 52)   // green-800
	pdf.CellFormat(100, 10, "GAJI BERSIH", "1", 0, "L", true, 0, "")
	pdf.CellFormat(70, 10, formatCurrency(payroll.NetSalary), "1", 1, "R", true, 0, "")

	pdf.Ln(10)

	// Payment Status
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(75, 85, 99)
	pdf.Cell(0, 8, "Status Pembayaran")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(31, 41, 55)

	status := "Pending"
	if payroll.Status == "paid" {
		status = "Lunas"
	}

	pdf.Cell(50, 6, "Status")
	pdf.Cell(5, 6, ":")
	pdf.Cell(0, 6, status)
	pdf.Ln(6)

	if payroll.PaymentDate != nil {
		pdf.Cell(50, 6, "Tanggal Pembayaran")
		pdf.Cell(5, 6, ":")
		pdf.Cell(0, 6, payroll.PaymentDate.Format("02 January 2006"))
		pdf.Ln(6)
	}

	// Footer
	pdf.Ln(15)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(156, 163, 175) // gray-400
	generatedTime := time.Now().Format("02 January 2006, 15:04")
	pdf.Cell(0, 5, fmt.Sprintf("Dokumen dibuat secara otomatis pada %s", generatedTime))

	// Output PDF to bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// formatCurrency formats a float64 as Indonesian Rupiah
func formatCurrency(amount float64) string {
	// Convert to integer
	intAmount := int64(amount)

	// Format with thousand separators
	str := fmt.Sprintf("%d", intAmount)
	result := ""
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}

	return result
}
