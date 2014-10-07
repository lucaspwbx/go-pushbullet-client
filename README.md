## Missing stuff

### Request authorization to upload file

#### POST

##### Endpoint: https://api.pushbullet.com/v2/upload-request
##### Params: file_name, file_type
##### Request: 
```curl -u <your_access_token_here>: -X POST https://api.pushbullet.com/v2/upload-req   uest -d file_name=image.png -d file_type=image/png```

##### Response: 
```{
  "file_type": "image/png",
  "file_name": "image.png",
  "file_url": "https://s3.amazonaws.com/pushbullet-uploads/ubd-VWb1dP5XrZzvHReWHCycIwPyuAMp2R9I/image.png",
  "upload_url": "https://s3.amazonaws.com/pushbullet-uploads",
  "data": {
    "awsaccesskeyid": "AKIAJJIUQPUDGPM4GD3W",
    "acl": "public-read",
    "key": "ubd-CWb1dP5XrZzvHReWHCycIwPyuAMp2R9I/image.png",
    "signature": "UX5s1uIy1ov6+xlj58JY7rGFKcs=",
    "policy": "eyKjb25kaXRpb25zIjTE6MzcuMjM0MTMwWiJ9",
    "content-type": "image/png"
  }
}
```

### Upload file

#### POST
#### Endpoint: "https://s3.amazonaws.com/pushbullet-uploads"
#### Params: 
Copy of all the parameters from the data object in the response to the upload request. In addition to that, the file should be uploaded as the parameter file. This request is more complicated than most of the other API requests and requires multipart/form-data encoding.
After the request completes, the file will be available at file_url from the upload request response.
#### Request:
```
curl -i -X POST https://s3.amazonaws.com/pushbullet-uploads \
  -F awsaccesskeyid=AKIAJJIUQPUDGPM4GD3W \
  -F acl=public-read \
  -F key=ubd-CWb1dP5XrZzvHReWHCycIwPyuAMp2R9I/image.png \
  -F signature=UX5s1uIy1ov6+xlj58JY7rGFKcs= \
  -F policy=eyKjb25kaXRpb25z6MzcuMjM0MTMwWiJ9 \
  -F content-type=image/png
  -F file=@test.txt
```
#### Response:
```
HTTP/1.1 204 No Content
```

### Push

#### Checklist
##### Params:

###### type - Set to list
###### title - The list's title.
###### items - The list items, a list of strings e.g. ["one", "two", "three"].

##### Request:
```
curl -u <your_access_token_here>: -X POST https://api.pushbullet.com/v2/pushes --header 'Content-Type: application/json' --data-binary '{"type": "list", "title": "List Title", "items": ["Item One", "Item Two"]}'

```

##### Response:
```
{
  "iden": "ubdpjAkaGXvUl2",
  "type": "list",
  "title": "List Title",
  "items": [{"checked": false, "text": "Item One"}, {"checked": false, "text": "Item Two"}],
  "created": 1411595195.1267679,
  "modified": 1411595195.1268303,
  "active": true,
  "dismissed": false,
  "sender_iden": "ubd",
  "sender_email": "ryan@pushbullet.com",
  "sender_email_normalized": "ryan@pushbullet.com",
  "receiver_iden": "ubd",
  "receiver_email": "ryan@pushbullet.com",
  "receiver_email_normalized": "ryan@pushbullet.com"
}
```

#### File
##### Params:

###### type - Set to file
###### file_name - The name of the file.
###### file_type - The MIME type of the file.
###### file_url - The url for the file. See pushing files for how to get a file_url
###### body - A message to go with the file.

##### Request:
```
Pushing files is a two-part process: first the file needs to be uploaded, then a push needs to be sent for that file.

To upload a new file, use the upload request endpoint.

Once the file has been uploaded, set the file_name, file_url, and file_type returned in the response to the upload request as the parameters for a new push with type=file.
```

### Request push history (GET)

#### Allow parameter passing:

##### modified_after - Request pushes modified after this timestamp.

### Update Push

#### Params:

##### dismissed
Set to true to mark the push as dismissed. All devices displaying this push should hide it from view.

##### items 
Update the items of a list push. The format should be the same as the items property of the push object, e.g. [{"checked": true, "text": "Item One"},
"checked": true, "text": "Item Two"}]

#### Request:
```
curl -u <your_access_token_here>: -X POST https://api.pushbullet.com/v2/pushes/ubdpjAkaGXvUl2 --header 'Content-Type: application/json' --data-binary '{"items": [{"checked": true, "text": "one"}, {"checked": true, "text": "two"}]}'
