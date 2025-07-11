$token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3NTIyNjA2NTUsImlhdCI6MTc1MjE3NDI1NSwidXNlcl9pZCI6MX0.MNE4JObbvcp-Yqtkt0UWv99u9BeLednfsGkekz_yAzU'
$headers = @{Authorization="Bearer $token"}
Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/auth/me' -Headers $headers
