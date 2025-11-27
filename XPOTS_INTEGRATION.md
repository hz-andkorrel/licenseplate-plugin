# XPOTS Webhook Integration Guide

## Quick Setup

### 1. Generate API Key

Add a secure webhook key to your `.env` file:

```bash
# Windows PowerShell
-join ((48..57) + (65..90) + (97..122) | Get-Random -Count 32 | ForEach-Object {[char]$_})
```

Copy the generated key to `.env`:
```
WEBHOOK_API_KEY=your-generated-key-here
```

### 2. Configure XPOTS System

In your XPOTS parking management system, configure the webhook:

**Webhook URL:**
```
http://your-server-ip:8082/api/licenseplate/webhook/xpots
```

**Method:** `POST`

**Headers:**
```
Authorization: Bearer your-generated-key-here
Content-Type: application/json
```

**Events to Send:**
- Vehicle Entry (when plate detected at entrance)
- Vehicle Exit (when plate detected at exit)

### 3. Start the Plugin

```bash
go run main.go
```

The plugin will:
- Listen for XPOTS webhooks on `/api/licenseplate/webhook/xpots`
- Automatically create records for detected plates
- Update check-in times for entries
- Update check-out times for exits

## How It Works

### Entry Event
When XPOTS detects a vehicle entering:
1. Plugin receives webhook with plate number and metadata
2. If plate exists → updates check-in time
3. If plate is unknown → creates new record with "Unknown Guest (Auto-detected)"
4. Stores XPOTS metadata (location, confidence, camera ID) in notes field

### Exit Event
When XPOTS detects a vehicle leaving:
1. Plugin receives webhook with plate number
2. If plate exists → updates check-out time
3. If plate is unknown → creates exit-only record for logging

### Guest Information
For auto-detected vehicles:
- Initially marked as "Unknown Guest (Auto-detected)"
- Hotel staff can manually update guest info later via API or management system
- Can be linked to reservations through hotel management integration

## Testing

Test the webhook without XPOTS:

```bash
# Run the test script
.\test-xpots-webhook.ps1

# Or manually with curl
curl -X POST http://localhost:8082/api/licenseplate/webhook/xpots \
  -H "Authorization: Bearer your-webhook-key" \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "entry",
    "plate_number": "TEST-123",
    "timestamp": "2025-11-27T10:30:00Z",
    "location": "Main Gate",
    "confidence": 0.98,
    "camera_id": "CAM-001",
    "direction": "in"
  }'
```

## XPOTS Payload Format

The plugin expects this JSON structure from XPOTS:

```json
{
  "event_type": "entry",        // "entry", "exit", "scan", "in", "out"
  "plate_number": "ABC-123",    // License plate detected
  "timestamp": "2025-11-27T10:30:00Z",
  "location": "Main Gate",       // Gate/camera location
  "confidence": 0.98,            // Recognition confidence (0-1)
  "camera_id": "CAM-001",        // Camera identifier
  "direction": "in"              // "in" or "out"
}
```

## Security

- **API Key Authentication:** All webhook requests must include valid API key
- **HTTPS Recommended:** Use HTTPS in production to encrypt webhook data
- **Key Rotation:** Change WEBHOOK_API_KEY regularly
- **IP Whitelisting:** Consider restricting webhook endpoint to XPOTS server IPs

## Troubleshooting

### Webhook Not Receiving Data
1. Check XPOTS webhook configuration
2. Verify plugin is running: `curl http://localhost:8082/health`
3. Check firewall allows traffic on port 8082
4. Review plugin logs for authentication errors

### Authentication Failures
1. Verify API key matches in both `.env` and XPOTS config
2. Check Authorization header format: `Bearer your-key`
3. Look for "Invalid API key" in responses

### Missing Records
1. Check database connection: `docker ps` (PostgreSQL should be running)
2. Verify migrations ran: Check `license_plates` table exists
3. Review plugin logs for database errors

## Integration with Hotel Management

To link auto-detected plates with guest reservations:

1. **Manual Matching:** Hotel staff updates guest info via API
2. **Reservation System:** Hotel management system calls API when guest checks in
3. **Pre-registration:** Guests provide plate numbers during booking

Example API call to update guest info:
```bash
curl -X POST http://localhost:8082/api/licenseplate/scan \
  -H "Content-Type: application/json" \
  -d '{
    "plate_number": "ABC-123",
    "guest_name": "John Smith",
    "room_number": "305",
    "vehicle_make": "Toyota",
    "vehicle_model": "Camry"
  }'
```

## Production Deployment

For production use:

1. **Use HTTPS** with valid SSL certificate
2. **Strong API Key:** Generate 256-bit random key
3. **Database Backup:** Regular backups of PostgreSQL
4. **Monitoring:** Track webhook failures and response times
5. **Rate Limiting:** Protect against webhook spam
6. **Logging:** Archive XPOTS webhook data for auditing

## Support

For XPOTS-specific configuration:
- Contact your XPOTS system administrator
- Refer to XPOTS webhook documentation
- Check XPOTS support portal

For plugin issues:
- Check README.md for general setup
- Review database migrations
- Verify environment configuration
