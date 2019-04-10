# server
LabelTong server with Golang and Node.js

# Usage
```$xslt
go run *.go
```
- port : 19432

# API
- /dataset/list [GET] Login function, Oauth from client required
    - return : list of dataset in json (see model.go dataset struct)
- /dataset/list/{dsid}/get [GET] Logout function, Oauth from client required
    - dsid = name of dataset,result of /dataset/list
    - return: json 
        - fileid = (fileid in model.go datatolabel}
         - Base64data = {image data with base64 encoded}
- /dataset/list/{dsid}/{id}/ans [POST] Check if user is authenticated not yet implemented
- /dataset/list/{dsid}/{id}/info [GET] Check if user is authenticated not yet implemented

# References
- https://mingrammer.com/getting-started-with-oauth2-in-go/ 
- https://blog.kowalczyk.info/article/f/accessing-github-api-from-go.html
