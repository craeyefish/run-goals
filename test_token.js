const jwt = require('jsonwebtoken');

// Use the same secret and algorithm as the backend
const secret = 'secret';
const userID = 1; // Using the user that exists in the database

// Generate access token (1 hour expiry like the backend)
const accessToken = jwt.sign(
  {
    sub: userID,
    exp: Math.floor(Date.now() / 1000) + 60 * 60, // 1 hour
    iat: Math.floor(Date.now() / 1000),
  },
  secret,
  { algorithm: 'HS256' }
);

console.log('Access Token:');
console.log(accessToken);
console.log('\nTest command:');
console.log(
  `curl -H "Authorization: Bearer ${accessToken}" "http://localhost:8080/api/group-members?groupID=8"`
);
