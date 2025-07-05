# Testing Guide for Run Goals Application

## Overview

This guide documents the testing methodology and setup process for the Run Goals application, including how to start services, initialize the database, generate authentication tokens, and test API endpoints.

## Service Architecture

The application consists of three main services:

- **Database**: PostgreSQL container (`run-goals-db`)
- **Backend**: Go application (`run-goals-backend`)
- **Frontend**: Angular application (`run-goals-frontend`)

## Starting the Application

### 1. Start All Services

```zsh
cd /path/to/run-goals
docker-compose up -d
```

### 2. Verify Services are Running

```zsh
docker ps
```

Expected output should show all three containers running:

- `run-goals-db` (PostgreSQL)
- `run-goals-backend` (Go backend)
- `run-goals-frontend` (Angular frontend)

## Database Setup and Initialization

### 1. Check Database Status

```zsh
# Check if database exists and has tables
docker exec run-goals-db psql -U postgres -d run_goals -c "\dt"
```

### 2. Manual Database Initialization (if needed)

If tables don't exist, run initialization scripts:

```zsh
# Create database
docker exec run-goals-db psql -U postgres -d postgres -c "CREATE DATABASE run_goals;"

# Check what initialization files are available
docker exec run-goals-db find /docker-entrypoint-initdb.d -name "*.sql"

# Run initialization scripts manually if needed
docker exec run-goals-db psql -U postgres -d run_goals -f /docker-entrypoint-initdb.d/10_users.sql
docker exec run-goals-db psql -U postgres -d run_goals -f /docker-entrypoint-initdb.d/20_peaks.sql
# ... continue for other table scripts
```

### 3. Verify Sample Data

```zsh
# Check for users
docker exec run-goals-db psql -U postgres -d run_goals -c "SELECT id, strava_athlete_id FROM users LIMIT 5;"

# Check for groups
docker exec run-goals-db psql -U postgres -d run_goals -c "SELECT id, name FROM groups LIMIT 5;"

# Check for group members
docker exec run-goals-db psql -U postgres -d run_goals -c "SELECT id, group_id, user_id, role FROM group_members LIMIT 5;"
```

## Authentication and API Testing

### 1. Generate JWT Token

The backend uses JWT authentication with HMAC-SHA256 signing. Generate a test token using Python:

```python
python3 -c "
import json
import base64
import hmac
import hashlib
import time

# JWT Header
header = {'alg': 'HS256', 'typ': 'JWT'}

# JWT Payload (user ID 1 exists in test data)
payload = {
    'sub': 1.0,  # Backend expects float64 for userID
    'exp': int(time.time()) + 3600,  # 1 hour from now
    'iat': int(time.time())
}

# Encode header and payload
header_encoded = base64.urlsafe_b64encode(json.dumps(header, separators=(',', ':')).encode()).decode().rstrip('=')
payload_encoded = base64.urlsafe_b64encode(json.dumps(payload, separators=(',', ':')).encode()).decode().rstrip('=')

# Create signature using backend secret
message = f'{header_encoded}.{payload_encoded}'
secret = 'secret'  # From backend/.env JWT_SECRET
signature = hmac.new(secret.encode(), message.encode(), hashlib.sha256).digest()
signature_encoded = base64.urlsafe_b64encode(signature).decode().rstrip('=')

# Complete JWT
jwt_token = f'{header_encoded}.{payload_encoded}.{signature_encoded}'
print(jwt_token)
"
```

### 2. Test API Endpoints

#### Test Authentication

```zsh
# Test without token (should fail)
curl "http://localhost:8080/api/group-members?groupID=8"
# Expected: "header failed" with 401 status

# Test with token (should succeed)
TOKEN="your_jwt_token_here"
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/group-members?groupID=8"
```

#### Test Groups Endpoints

```zsh
# Get user's groups
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/groups" | jq .

# Get group members (basic data)
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/group-members?groupID=8" | jq .

# Get group member contributions (with activity data)
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/group-members-contribution?groupID=8&startDate=2025-01-01&endDate=2025-12-31" | jq .

# Get group goals
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/group-goals?groupID=8" | jq .
```

## Common Issues and Solutions

### Database Issues

**Issue**: "relation does not exist" errors

```zsh
# Solution: Check database and run initialization
docker exec run-goals-db psql -U postgres -d run_goals -c "\dt"
# If empty, run initialization scripts
```

**Issue**: Database connection failed

```zsh
# Solution: Check database container status and logs
docker logs run-goals-db --tail 20
docker restart run-goals-db
```

### Backend Issues

**Issue**: "header failed" on API calls

```zsh
# Solution: Generate fresh JWT token (tokens expire after 1 hour)
# Use the Python script above to generate new token
```

**Issue**: Backend container not starting

```zsh
# Solution: Check logs and rebuild if needed
docker logs run-goals-backend --tail 20
docker-compose build backend
docker-compose up -d backend
```

### Frontend Issues

**Issue**: Frontend not accessible

```zsh
# Solution: Check container status and restart if needed
curl -s "http://localhost:4200" | head -5
docker logs run-goals-frontend --tail 20
docker restart run-goals-frontend
```

## Expected API Responses

### Group Members (Basic)

```json
{
  "members": [
    {
      "id": 2,
      "group_id": 8,
      "user_id": 1,
      "role": "admin",
      "joined_at": "2025-06-25T20:09:20.054643Z"
    }
  ]
}
```

### Group Members (With Contributions)

```json
{
  "members": [
    {
      "group_member_id": 2,
      "group_id": 8,
      "user_id": 1,
      "role": "admin",
      "joined_at": "2025-06-25T20:09:20.054643Z",
      "total_activities": 116,
      "total_distance": 608826.9,
      "total_unique_summits": 2,
      "total_summits": 3
    }
  ]
}
```

## Configuration Details

### Database Configuration

- **Host**: localhost:5432 (from host machine)
- **Database**: run_goals
- **User**: postgres
- **Password**: postgres

### Backend Configuration

- **Port**: 8080
- **JWT Secret**: "secret" (from backend/.env)
- **Database Connection**: Configured via docker-compose environment variables

### Frontend Configuration

- **Port**: 4200
- **API Proxy**: Configured to proxy /api/\* requests to backend:8080

## Testing Checklist

- [ ] All Docker containers running
- [ ] Database initialized with tables and sample data
- [ ] Backend API accessible with JWT authentication
- [ ] Frontend accessible and loading
- [ ] API endpoints returning expected data structure
- [ ] Frontend components receiving and displaying data correctly

## Useful Commands

```zsh
# Quick health check
docker ps
curl -s "http://localhost:4200" | head -1
curl -s "http://localhost:8080/api/groups" | head -1

# Generate fresh token and test API
TOKEN=$(python3 -c "import json,base64,hmac,hashlib,time; h={'alg':'HS256','typ':'JWT'}; p={'sub':1.0,'exp':int(time.time())+3600,'iat':int(time.time())}; he=base64.urlsafe_b64encode(json.dumps(h,separators=(',',':')).encode()).decode().rstrip('='); pe=base64.urlsafe_b64encode(json.dumps(p,separators=(',',':')).encode()).decode().rstrip('='); m=f'{he}.{pe}'; s=base64.urlsafe_b64encode(hmac.new('secret'.encode(),m.encode(),hashlib.sha256).digest()).decode().rstrip('='); print(f'{he}.{pe}.{s}')")

curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/group-members?groupID=8" | jq .

# Clean restart all services
docker-compose down
docker-compose up -d

# View logs for debugging
docker logs run-goals-backend --tail 20
docker logs run-goals-frontend --tail 20
docker logs run-goals-db --tail 20
```

---

This testing methodology was developed while fixing the members table functionality and provides a reliable way to test and validate the Run Goals application.
