# FinalCD
Hashed content delivery. FinalCD delivers content by its SHA256 hash. 
This way the same file is never stored more than once. 


# Endpoints

### [POST] /
Create a new file if it doesn't exist yet.
#### Request
```shell
curl --location --request POST 'http://localhost:8080/' --form 'f=@"/tmp/my_file"
```
#### Response
_If the file didn't exist yet_  
Status 201 Created
```json
  {"hash": "<hash of the created file>"}
```
_If the file already existed_  
Status 200 OK
```json
  {"hash": "<hash of the file>"}
```

### [GET] /d/:hash



## Development
Backend: `go run .`
Frontend: `cd frontend && yarn start`


