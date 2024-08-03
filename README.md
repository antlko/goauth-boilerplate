<h1 align="center">
    <br>
    <img src="https://miro.medium.com/v2/resize:fit:4800/format:webp/1*ovJrUZn9l-SXfEAWDpt2qQ.png"
         alt="Go JWT" width="200">
    <br>
    Go Jwt Auth Boilerplate
    <br>
</h1>

<h4 align="center">Golang boilerplate project to have a quick setup of server with JWT Authentication</h4>

<p align="center">
    <a href="#key-features">Key Features</a> •
    <a href="#how-to-use">How To Use</a> •
    <a href="#api">Api</a> •
    <a href="#license">License</a>
</p>

```
+--------+                                           +---------------+
|        |--(A)------- Authorization Grant --------->|               |
|        |                                           |               |
|        |<-(B)----------- Access Token -------------|               |
|        |               & Refresh Token             |               |
|        |                                           |               |
|        |                                           |               |
|        |--(C)----------- Access Token ------------>|               |
|        |                                           |               |
|        |<-(D)--------- Protected Resource ---------| Authorization |
| Client |                                           |     Server    |
|        |--(E)----------- Access Token ------------>|               |
|        |                                           |               |
|        |<-(F)-------- Invalid Token Error ---------|               |
|        |                                           |               |
|        |                                           |               |
|        |--(G)----------- Refresh Token ----------->|               |
|        |                                           |               |
|        |<-(H)----------- Access Token -------------|               |
+--------+           & Optional Refresh Token        +---------------+
```

## Key Features

* DB instance prepared - easy to create needed structures with goose library.
* Fiber framework - fast golang web library.
* SignUp - prepared endpoint to register the user.
* SignIn - authenticate user and get access & refresh tokens.
* Refresh - refresh tokens.
* Easy-to-test - project structured in a way to make it simple and easy to mock everything and test.
* Docker-Compose for DB

## How To Use

```bash
# Clone this repository
$ git clone https://github.com/antlko/goauth-boilerplate.git

# Go into the repository
$ cd goauth-boilerplate

# Deploy postgres container on demand
$ docker-compose up -d

# Copy default envs
$ cp .env.dist .env

# Run the app
$ go run main.go
```

## Api
```http
POST /api/v1/users/auth/signup
{
"login":"test",
"email":"test@gmail.com",
"password":"test"
}

POST /api/v1/users/auth/signin
{
"login":"test",
"password":"test"
}

POST /api/v1/users/auth/token/refresh
{
"refresh_token":"your_refresh_token"
}
```

Example of usage the protected endpoint
```http
POST /api/v1/protected/user
Authorization: Bearer your_access_token
```

## License

MIT