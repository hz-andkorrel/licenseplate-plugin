# Database Setup Script for License Plate Plugin

Write-Host ""
Write-Host "=== Setting up License Plate Database ===" -ForegroundColor Cyan

# Database connection details
$dbUser = "broker"
$dbPassword = "broker123"
$dbName = "broker_db"
$dbHost = "localhost"
$dbPort = "5432"

# Set PGPASSWORD environment variable
$env:PGPASSWORD = $dbPassword

Write-Host ""
Write-Host "Running migration..." -ForegroundColor Yellow

# Run the migration
psql -U $dbUser -h $dbHost -p $dbPort -d $dbName -f migrations/001_create_license_plates_table.sql

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "Database setup complete!" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "Database setup failed." -ForegroundColor Red
    Write-Host "Make sure PostgreSQL is running and credentials are correct." -ForegroundColor Red
}

# Clear password
Remove-Item Env:PGPASSWORD

Write-Host ""


Write-Host ""
