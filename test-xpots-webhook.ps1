# XPOTS Webhook Test Script
# This script tests the XPOTS webhook integration

Write-Host "`n=== XPOTS Webhook Integration Test ===" -ForegroundColor Cyan

# Configuration
$baseUrl = "http://localhost:8082/api/licenseplate"
$webhookKey = "test-webhook-key-12345"

Write-Host "`nNote: Make sure to set WEBHOOK_API_KEY=$webhookKey in your .env file" -ForegroundColor Yellow
Write-Host "Press any key to continue..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# 1. Get webhook info
Write-Host "`n1. Getting webhook configuration info..." -ForegroundColor Yellow
curl.exe -X GET "$baseUrl/webhook/info"

# 2. Test webhook - entry event
Write-Host "`n`n2. Sending entry event (vehicle arriving)..." -ForegroundColor Yellow
$entryPayload = @{
    event_type = "entry"
    plate_number = "XPOTS-001"
    timestamp = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
    location = "Main Gate"
    confidence = 0.98
    camera_id = "CAM-001"
    direction = "in"
    vehicle_type = "car"
} | ConvertTo-Json

curl.exe -X POST "$baseUrl/webhook/xpots" `
    -H "Authorization: Bearer $webhookKey" `
    -H "Content-Type: application/json" `
    -d $entryPayload

# 3. Verify the record was created
Write-Host "`n`n3. Verifying record was created..." -ForegroundColor Yellow
Start-Sleep -Seconds 1
curl.exe -X GET "$baseUrl/records/XPOTS-001"

# 4. Test webhook - another entry event (different plate)
Write-Host "`n`n4. Sending another entry event (different vehicle)..." -ForegroundColor Yellow
$entry2Payload = @{
    event_type = "entry"
    plate_number = "XPOTS-002"
    timestamp = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
    location = "Side Gate"
    confidence = 0.95
    camera_id = "CAM-002"
    direction = "in"
} | ConvertTo-Json

curl.exe -X POST "$baseUrl/webhook/xpots" `
    -H "Authorization: Bearer $webhookKey" `
    -H "Content-Type: application/json" `
    -d $entry2Payload

# 5. Test webhook - exit event
Write-Host "`n`n5. Sending exit event (vehicle leaving)..." -ForegroundColor Yellow
Start-Sleep -Seconds 2
$exitPayload = @{
    event_type = "exit"
    plate_number = "XPOTS-001"
    timestamp = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
    location = "Main Gate"
    confidence = 0.99
    camera_id = "CAM-001"
    direction = "out"
} | ConvertTo-Json

curl.exe -X POST "$baseUrl/webhook/xpots" `
    -H "Authorization: Bearer $webhookKey" `
    -H "Content-Type: application/json" `
    -d $exitPayload

# 6. Verify check-out was recorded
Write-Host "`n`n6. Verifying check-out was recorded..." -ForegroundColor Yellow
Start-Sleep -Seconds 1
curl.exe -X GET "$baseUrl/records/XPOTS-001"

# 7. Get all records
Write-Host "`n`n7. Getting all records..." -ForegroundColor Yellow
curl.exe -X GET "$baseUrl/records"

# 8. Test invalid API key
Write-Host "`n`n8. Testing webhook with invalid API key (should fail)..." -ForegroundColor Yellow
curl.exe -X POST "$baseUrl/webhook/xpots" `
    -H "Authorization: Bearer wrong-key" `
    -H "Content-Type: application/json" `
    -d $entryPayload

Write-Host "`n`n=== Test Complete ===" -ForegroundColor Green
Write-Host "Review the responses above to verify:" -ForegroundColor Cyan
Write-Host "  ✓ Entry events create new records" -ForegroundColor White
Write-Host "  ✓ Exit events update check-out times" -ForegroundColor White
Write-Host "  ✓ Invalid API keys are rejected" -ForegroundColor White
Write-Host "  ✓ All data is stored in the database" -ForegroundColor White
