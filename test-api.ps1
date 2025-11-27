# License Plate Recognition Plugin API Test Script

$baseUrl = "http://localhost:8082/api/licenseplate"

Write-Host "`n=== Testing License Plate Recognition Plugin ===" -ForegroundColor Cyan

# Test 1: Health Check
Write-Host "`n1. Health Check" -ForegroundColor Yellow
curl.exe -X GET http://localhost:8082/health

# Test 2: Scan a license plate (Create record)
Write-Host "`n`n2. Scan License Plate - ABC123" -ForegroundColor Yellow
$scanData1 = @{
    plate_number = "ABC123"
    guest_name = "John Doe"
    room_number = "101"
    vehicle_make = "Toyota"
    vehicle_model = "Camry"
    notes = "Regular guest"
} | ConvertTo-Json

curl.exe -X POST "$baseUrl/scan" `
  -H "Content-Type: application/json" `
  -d $scanData1

# Test 3: Scan another license plate
Write-Host "`n`n3. Scan License Plate - XYZ789" -ForegroundColor Yellow
$scanData2 = @{
    plate_number = "XYZ789"
    guest_name = "Jane Smith"
    room_number = "205"
    vehicle_make = "Honda"
    vehicle_model = "Civic"
} | ConvertTo-Json

curl.exe -X POST "$baseUrl/scan" `
  -H "Content-Type: application/json" `
  -d $scanData2

# Test 4: Get all records
Write-Host "`n`n4. Get All License Plate Records" -ForegroundColor Yellow
curl.exe -X GET "$baseUrl/records"

# Test 5: Get specific record
Write-Host "`n`n5. Get Specific Record (ABC123)" -ForegroundColor Yellow
curl.exe -X GET "$baseUrl/records/ABC123"

# Test 6: Search by guest name
Write-Host "`n`n6. Search by Guest Name (John)" -ForegroundColor Yellow
curl.exe -X GET "$baseUrl/records?guest_name=John"

# Test 7: Delete a record
Write-Host "`n`n7. Delete Record (XYZ789)" -ForegroundColor Yellow
curl.exe -X DELETE "$baseUrl/records/XYZ789"

# Test 8: Verify deletion
Write-Host "`n`n8. Verify Deletion - Get All Records" -ForegroundColor Yellow
curl.exe -X GET "$baseUrl/records"

Write-Host "`n`n=== Tests Complete ===" -ForegroundColor Green
