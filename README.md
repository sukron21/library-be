<br />
<p align="center">
  <h3 align="center">Library API</h3>
</p>

<!-- TABLE OF CONTENTS -->

## Table of Contents

- [About the Project](#about-the-project)
  - [Project Structure](#project-structure)
  - [Package Modules](#package-modules)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Contributing](#contributing)
- [Related Project](#related-project)
- [Contact](#contact)

<!-- ABOUT THE PROJECT -->

## About The Project

This repository contains the backend source code for the **Library Management System** application. The backend is built using **Golang** with the **Fiber** framework and **PostgreSQL** database, and uses **Postman** for testing and API documentation..

### Project Structure

```
|── library-be
   |── config                       # Database configuration
   |── controller                   # Request controller
   |── helper                       # response
   |── middleware                   # Middleware configuration
   |── model                        # Database query model
   |── routes                       # API Endpoint routes
   |── .env                             # Environment variables
   |── .gitignore                       # Files that should be ignored
   |── Library.postman_collection.json   # Postman Documentation
   |── db.sql                           # SQL database creation
   |── main.go                         # Index file
   |── README.md                        # Readme
```

### Package Modules

Below are lists of modules used in this API:

- [Golang](https://golang.org/)
- [Fiber](https://gofiber.io/) — web framework yang cepat & ringan
- [PostgreSQL](https://www.postgresql.org/) — database relational
- [GORM](https://gorm.io/) — ORM untuk Go
- [JWT](https://jwt.io/) — authentication
- [Postman](https://www.postman.com/) — testing & dokumentasi API

<!-- GETTING STARTED -->

## Getting Started

### Prerequisites

This is an example of things you need to use the application and how to install them.

- [Go](https://go.dev/)

### Installation

1. Clone the repo

```sh
git clone https://github.com/sukron21/library-be.git
```

2. Install NPM packages

```sh
go mod tidy
```

3. Add .env file at your backend root folder project, and add the following

```sh
PORT=3000
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret_key
```

4. running the project

```sh
air

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b your/branch`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/yourbranch`)
5. Open a Pull Request

<!-- RELATED PROJECT -->

## Related Project

- [Frontend](https://github.com/sukron21/library-fe)

<!-- CONTACT -->

## Contact

Contributors name and contact info

- Rahmat Furqon [@sukron21](https://github.com/sukron21)
```
