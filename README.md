Client for Pushbullet

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
