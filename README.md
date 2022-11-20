# FinalCD
Hashed content delivery. FinalCD delivers content by its SHA256 hash. 
This way the same file is never stored more than once. 


# Endpoints
Replace URI scheme, localhost, port if needed. Localhost is used to simplify developing by copying the curl commands.

### [POST] /u
Create a new file if it doesn't exist yet
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
Download a file by its hash
#### Request
```shell
curl --location --request GET 'http://localhost:8080/d/<hash of the file>
```
#### Response
_If the file is found_  
Status 200 OK
| Header  | Values |
| ------------- | ------------- |
| Content-Length  | Length of the file in bytes  |
| Content-Type  | Media type of the file, e.g. application/pdf  |
| X-Served-From | _disk on server_ or _cache on server_ |

_If the file is not found_  
Status 404 Not Found

### [GET] /l
Retrieve a list of available files
#### Response
Status 200 OK
```
[
    {
        "hash": "<hash of a file>",
        "size": <length of the file in bytes>
    },
    {
        "hash": "<hash of a file>",
        "size": <length of the file in bytes>
    },    
]
```


## Development
Backend: `go run .`  
Frontend: `cd frontend && yarn start`

## Deployment
1. Build frontend `cd frontend && yarn build`
2. Run `./deploy.sh`

