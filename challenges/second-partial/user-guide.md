User Guide
==========

## Installations

- Download the folder *second-partial*
 
- Open a terminal and type the next command:
 
```
$ go get github.com/dgrijalva/jwt-go
```

This library is needed to generate a unique token for each user.
# Using the API
It is required to open two terminals, the first terminal will run the following command:

```
$ go run api.go
```
This will receive the requests from the second terminal.

The existing commands are the following.

## Login

The *Login* command will receive an username and a password, and will generate a unique token required for the other commands.

```
$ curl -u username:password http://localhost:8080/login
```

Once its done the terminal will display the following message:

```
{
	"message": "Hi username, welcome to the DPIP System",
	"token" "OjIE89GzFw"
}
```

**Note: Instead of username it will display the name introduced in the Login command**

**Note 2: token provided above is just an example, in reality the token will be generated once the command executes**

## Logout

The *Logout* command will receive the token of the user and will revoke it.

```
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/logout
```

Once this is done the following message will appear:

```
{
	"message": "Bye username, your token has been revoked"
}
```

**Note: Because the token is revoked it will no longer be valid to use, please use this command at the very end**

## Upload

The *Upload* command will receive the **absolute** path of the image that you want to upload and the token of the user.

```
$ curl -F 'data=@/path/to/local/image.png' -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/upload
```

If the user introduced correctly the absolute path and the token, the following message will be displayed:

```
{
	"message": "An image has been successfully uploaded",
	"filename": "image.png",
	"size": "500kb"
}
```

## Status

the *Status* command will receive the token of the user.

```
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status
```

If the token exists it will the display the user name and the current time.

```
{
	"message": "Hi username, the DPIP System is Up and Running"
	"time": "2015-03-07 11:06:39"
}
```
