<h1 align="center">
    <br>
    <img src="https://www.google.com/url?sa=i&url=https%3A%2F%2Fmedium.com%2Fswlh%2Fbuilding-a-user-auth-system-with-jwt-using-golang-30892659cc0&psig=AOvVaw1zrDRs_07ozt2ZLCMe8NQd&ust=1722774963949000&source=images&cd=vfe&opi=89978449&ved=0CBEQjRxqFwoTCMDUwoDr2IcDFQAAAAAdAAAAABAE"
         alt="Markdownify" width="200">
    <br>
    Go Jwt Auth Boilerplate
    <br>
</h1>

<h4 align="center">Golang boilerplate project to have a quick setup of server with JWT Authentication</h4>

<p align="center">
    <a href="#key-features">Key Features</a> •
    <a href="#how-to-use">How To Use</a> •
    <a href="#download">Download</a> •
    <a href="#license">License</a>
</p>

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

# Install dependencies
$ docker-compose up -d

# Copy default envs
$ cp .env.dist .env
```
# Run the app
$ go run main.go
```

## License

MIT

---


