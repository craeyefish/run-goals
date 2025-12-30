# Testing Guide

## Quick Start

### Start Services
```bash
# With production DB (recommended for dev)
docker compose -f docker-compose.prod-db.yaml up --build

# Full local stack (local DB)
docker compose up --build
```

### Generate JWT Token
```bash
TOKEN=$(python3 -c "import json,base64,hmac,hashlib,time; h={'alg':'HS256','typ':'JWT'}; p={'sub':1.0,'exp':int(time.time())+3600,'iat':int(time.time())}; he=base64.urlsafe_b64encode(json.dumps(h,separators=(',',':')).encode()).decode().rstrip('='); pe=base64.urlsafe_b64encode(json.dumps(p,separators=(',',':')).encode()).decode().rstrip('='); m=f'{he}.{pe}'; s=base64.urlsafe_b64encode(hmac.new('secret'.encode(),m.encode(),hashlib.sha256).digest()).decode().rstrip('='); print(f'{he}.{pe}.{s}')")
```

### Test API
```bash
# Without token (401)
curl http://localhost:8080/api/groups

# With token
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/groups | jq .

# Admin endpoint (no JWT needed)
curl -X POST "http://localhost:8080/admin/refresh-peaks?admin_key=dev-admin-key"
```

## Key Endpoints

| Endpoint | Auth | Description |
|----------|------|-------------|
| `GET /api/activities` | JWT | User's activities |
| `GET /api/peaks` | JWT | All peaks with is_summited |
| `GET /api/groups` | JWT | User's groups |
| `POST /hikegang/sync` | None | Trigger activity sync |
| `POST /admin/refresh-peaks` | admin_key | Refresh peak data from OSM |

## Database Access

```bash
# Production DB
get access string online.

# Local DB
docker exec -it run-goals-db psql -U postgres -d run_goals
```

## Common Issues

| Issue | Solution |
|-------|----------|
| "header failed" | Missing/invalid JWT token |
| "token validation failed" | Token expired, regenerate |
| DB connection failed | Check DATABASE_SSLMODE (require for prod, disable for local) |
