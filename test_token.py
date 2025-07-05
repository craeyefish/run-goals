#!/usr/bin/env python3
import jwt
import time

# Use the same secret and algorithm as the backend
secret = 'secret'
user_id = 1  # Using the user that exists in the database

# Generate access token (1 hour expiry like the backend)
payload = {
    'sub': float(user_id),  # Backend expects float64
    'exp': int(time.time()) + (60 * 60),  # 1 hour
    'iat': int(time.time())
}

access_token = jwt.encode(payload, secret, algorithm='HS256')

print('Access Token:')
print(access_token)
print('\nTest command:')
print(f'curl -H "Authorization: Bearer {access_token}" "http://localhost:8080/api/group-members?groupID=8"')
