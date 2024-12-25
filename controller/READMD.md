## Create

>   `request`
>
>   ```bash
>   curl --silent --location --request POST 'http://localhost:8080/api/user' \
>   --header 'Content-Type: application/json' \
>   --header 'Authorization: Bearer -' \
>   --data-raw '{
>       "id": "user01",
>       "name": "user01",
>       "email": "user01@gmail.com"
>   }'
>   ```
>
>   `response`
>
>   ```json
>   {
>       "code": 0,
>       "data": {
>           "id": "user01",
>           "created_at": "2024-12-25T12:43:21.250766+08:00",
>           "updated_at": "2024-12-25T12:43:21.241+08:00",
>           "name": "user01",
>           "email": "user01@gmail.com"
>       },
>       "msg": "success"
>   }
>   ```



## Delete

### delete one resource (by route parameter)

>   `request`
>
>   ```bash
>   # Delete user whose id is 'user01'
>   curl --silent --location --request DELETE 'http://localhost:8080/api/user/user01' \
>   --header 'Authorization: Bearer -'
>   ```
>
>   `response`
>
>   ```json
>   {
>       "code": 0,
>       "data": "",
>       "msg": "success"
>   }
>   ```

### delete one resource (by http body)

>   `request`
>
>   ```bash
>   # Delete user whose id is 'user01'.
>   curl --silent --location --request DELETE 'http://localhost:8080/api/user' \
>   --header 'Content-Type: application/json' \
>   --header 'Authorization: Bearer -' \
>   --data '[
>       "user01"
>   ]'
>   ```
>
>   `response`
>
>   ```json
>   {
>       "code": 0,
>       "data": "",
>       "msg": "success"
>   }
>   ```

### delete multiple resources (by http body)

>   `request`
>
>   ```bash
>   # Delete user whose id are 'user01', 'user02'.
>   curl --silent --location --request DELETE 'http://localhost:8080/api/user' \
>   --header 'Content-Type: application/json' \
>   --header 'Authorization: Bearer -' \
>   --data '[
>       "user01", "user02"
>   ]'
>   ```
>
>   `response`
>
>   ```json
>   {
>       "code": 0,
>       "data": "",
>       "msg": "success"
>   }
>   ```

### delete multiple resource (by route parameter and http body)

>`request`
>
>```bash
># Delete user whose id are 'user01', 'user02', "user03".
>curl --silent --location --request DELETE 'http://localhost:8080/api/user/user01' \
>--header 'Content-Type: application/json' \
>--header 'Authorization: Bearer -' \
>--data '[
>    "user02", "user03"
>]'
>```
>
>`response`
>
>```json
>{
>    "code": 0,
>    "data": "",
>    "msg": "success"
>}
>```



## Update

>   `request`
>
>   ```bash
>   curl --silent --location --request PUT 'http://localhost:8080/api/user' \
>   --header 'Content-Type: application/json' \
>   --header 'Authorization: Bearer -' \
>   --data-raw '{
>       "id": "user01",
>       "name": "user01_modifed",
>       "email": "user01_modifed@gmail.com"
>   }'
>   ```
>
>   `response`
>
>   ```json
>   {
>       "code": 0,
>       "data": {
>           "id": "user01",
>           "created_at": "2024-12-25T13:01:01.634+08:00",
>           "updated_at": "2024-12-25T13:26:18.307+08:00",
>           "name": "user01_modifed",
>           "email": "user01_modifed@gmail.com"
>       },
>       "msg": "success"
>   }
>   ```



## UpdatePartial

>   `request`
>
>   ```bash
>   curl --silent --location --request PATCH 'http://localhost:8080/api/user' \
>   --header 'Content-Type: application/json' \
>   --header 'Authorization: Bearer -' \
>   --data '{
>       "id": "user01",
>       "name": "user01_fake"
>   }'
>   ```
>
>   `response`
>
>   ```json
>   {
>       "code": 0,
>       "data": {
>           "id": "user01",
>           "created_at": "2024-12-25T13:01:01.634+08:00",
>           "updated_at": "2024-12-25T13:28:55.709+08:00",
>           "name": "user01_fake",
>           "email": "user01_modifed@gmail.com"
>       },
>       "msg": "success"
>   }
>   ```



## List

>   `request`
>
>   ```bash
>   curl --silent --location --request GET 'http://localhost:8080/api/user' \
>   --header 'Authorization: Bearer -'
>   ```
>
>   `response`
>
>   ```json
>   {
>       "code": 0,
>       "data": {
>           "items": [
>               {
>                   "id": "user01",
>                   "created_at": "2024-12-25T11:21:21.134+08:00",
>                   "updated_at": "2024-12-25T11:21:21.135+08:00",
>                   "name": "user01",
>                   "email": "user01@gmail.com"
>               },
>               {
>                   "id": "user02",
>                   "created_at": "2024-12-25T11:23:18.017+08:00",
>                   "updated_at": "2024-12-25T11:23:18.017+08:00",
>                   "name": "user02",
>                   "email": "user02@gmail.com"
>               }
>           ],
>           "total": 2
>       },
>       "msg": "success"
>   }
>   ```



## Get

>   `request`
>
>   ```bash
>   curl --silent --location --request GET 'http://localhost:8080/api/user/user01' \
>   --header 'Authorization: Bearer -'
>   ```
>
>   `response`
>
>   ```json
>   {
>       "code": 0,
>       "data": {
>           "id": "user01",
>           "created_at": "2024-12-25T11:21:21.134+08:00",
>           "updated_at": "2024-12-25T11:21:21.135+08:00",
>           "name": "user01",
>           "email": "user01@gmail.com"
>       },
>       "msg": "success"
>   }
>   ```

