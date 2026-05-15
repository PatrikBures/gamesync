

### create user
```sh
curl -X POST localhost:8080/api/v1/users -H "Content-Type: application/json" -d '{"userName": "something", "roleID":50}' 
```

### create role
```sh
curl -X POST localhost:8080/api/v1/roles \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer 9SpgWqQRRkK-v4uxp7nwvJhOm5FPiYWsG-_u8-V17Onm" \
    -d '{"roleName": "something"}'
```

