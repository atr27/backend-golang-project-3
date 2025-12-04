# PowerShell script to reset and reseed the database
# Run this from the backend directory: .\reset-and-seed.ps1

Write-Host "🔄 Resetting and reseeding the database..." -ForegroundColor Cyan

# Database credentials (update these if different)
$env:PGPASSWORD = "secure_password"
$DB_USER = "hrms_user"
$DB_NAME = "hrms_db"
$DB_HOST = "localhost"
$DB_PORT = "5432"

Write-Host "📦 Dropping and recreating database..." -ForegroundColor Yellow

# Drop and recreate database
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U $DB_USER -h $DB_HOST -p $DB_PORT -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -U $DB_USER -h $DB_HOST -p $DB_PORT -d postgres -c "CREATE DATABASE $DB_NAME;"

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Database reset successfully!" -ForegroundColor Green
    
    Write-Host "🌱 Starting backend with seeding..." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "The backend will now:" -ForegroundColor Cyan
    Write-Host "  1. Connect to the database" -ForegroundColor White
    Write-Host "  2. Run all migrations" -ForegroundColor White
    Write-Host "  3. Seed 16 employees across 4 departments" -ForegroundColor White
    Write-Host "  4. Start the API server" -ForegroundColor White
    Write-Host ""
    
    # Run backend with seed flag
    go run cmd/api/main.go --seed
} else {
    Write-Host "❌ Failed to reset database. Please check your PostgreSQL connection." -ForegroundColor Red
}
