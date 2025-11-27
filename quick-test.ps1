# Quick Test - License Plate Plugin

Write-Host ""
Write-Host "=== Testing License Plate Plugin ===" -ForegroundColor Cyan
Write-Host ""

# Test 1: Scan a license plate
Write-Host "Test 1: Scanning license plate ABC123..." -ForegroundColor Yellow
$response1 = Invoke-RestMethod -Uri "http://localhost:8082/api/licenseplate/scan" -Method POST -ContentType "application/json" -Body '{"plate_number":"ABC123","guest_name":"John Doe","room_number":"101","vehicle_make":"Toyota","vehicle_model":"Camry"}'
Write-Host "Response:" -ForegroundColor Green
$response1 | ConvertTo-Json -Depth 5
Write-Host ""

# Test 2: Scan another license plate
Write-Host "Test 2: Scanning license plate XYZ789..." -ForegroundColor Yellow
$response2 = Invoke-RestMethod -Uri "http://localhost:8082/api/licenseplate/scan" -Method POST -ContentType "application/json" -Body '{"plate_number":"XYZ789","guest_name":"Jane Smith","room_number":"205","vehicle_make":"Honda","vehicle_model":"Civic"}'
Write-Host "Response:" -ForegroundColor Green
$response2 | ConvertTo-Json -Depth 5
Write-Host ""

# Test 3: Get all records
Write-Host "Test 3: Getting all records..." -ForegroundColor Yellow
$response3 = Invoke-RestMethod -Uri "http://localhost:8082/api/licenseplate/records" -Method GET
Write-Host "Response:" -ForegroundColor Green
$response3 | ConvertTo-Json -Depth 5
Write-Host ""

# Test 4: Get specific record
Write-Host "Test 4: Getting record for ABC123..." -ForegroundColor Yellow
$response4 = Invoke-RestMethod -Uri "http://localhost:8082/api/licenseplate/records/ABC123" -Method GET
Write-Host "Response:" -ForegroundColor Green
$response4 | ConvertTo-Json -Depth 5
Write-Host ""

# Test 5: Delete a record
Write-Host "Test 5: Deleting record XYZ789..." -ForegroundColor Yellow
$response5 = Invoke-RestMethod -Uri "http://localhost:8082/api/licenseplate/records/XYZ789" -Method DELETE
Write-Host "Response:" -ForegroundColor Green
$response5 | ConvertTo-Json -Depth 5
Write-Host ""

Write-Host "=== All Tests Complete! ===" -ForegroundColor Green
Write-Host ""
