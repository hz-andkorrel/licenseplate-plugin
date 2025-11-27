# Test License Plate Plugin Without Database

Write-Host ""
Write-Host "=== License Plate Plugin - Quick Test ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "NOTE: This plugin requires PostgreSQL to be running." -ForegroundColor Yellow
Write-Host ""
Write-Host "To start PostgreSQL:" -ForegroundColor White
Write-Host "  1. Start Docker Desktop" -ForegroundColor Gray
Write-Host "  2. Navigate to modulair-achterkantje folder" -ForegroundColor Gray
Write-Host "  3. Run: docker compose -f compose.dev.yml up postgres -d" -ForegroundColor Gray
Write-Host ""
Write-Host "Or install PostgreSQL locally and create the broker_db database." -ForegroundColor White
Write-Host ""
Write-Host "Database connection string:" -ForegroundColor White
Write-Host "  postgres://broker:broker123@localhost:5432/broker_db?sslmode=disable" -ForegroundColor Gray
Write-Host ""
Write-Host "After PostgreSQL is running:" -ForegroundColor White
Write-Host "  1. Run: .\setup-database.ps1" -ForegroundColor Gray
Write-Host "  2. Run: go run main.go" -ForegroundColor Gray
Write-Host "  3. Test API: .\test-api.ps1" -ForegroundColor Gray
Write-Host ""

# Check if PostgreSQL is accessible
Write-Host "Checking PostgreSQL connection..." -ForegroundColor Yellow
$testConnection = Test-NetConnection -ComputerName localhost -Port 5432 -WarningAction SilentlyContinue

if ($testConnection.TcpTestSucceeded) {
    Write-Host "PostgreSQL is running on port 5432!" -ForegroundColor Green
    Write-Host ""
    Write-Host "You can now:" -ForegroundColor White
    Write-Host "  1. Setup database: .\setup-database.ps1" -ForegroundColor Cyan
    Write-Host "  2. Start plugin: go run main.go" -ForegroundColor Cyan
    Write-Host "  3. Test API: .\test-api.ps1" -ForegroundColor Cyan
} else {
    Write-Host "PostgreSQL is not running on port 5432." -ForegroundColor Red
    Write-Host "Please start PostgreSQL first (see instructions above)." -ForegroundColor Yellow
}

Write-Host ""
