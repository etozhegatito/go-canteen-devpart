# University Canteen Management System

In this project, we use the Go programming language to create an order menu for a university canteen. Using Docker and Docker Compose, we achieved containerization and made the project much easier to run.

## Project Architecture

- **Backend:** Go
- **Frontend:** HTML + JavaScript
- **Database:** PostgreSQL (Web interface: pgAdmin)

## System Requirements

- **OS:** Linux, Windows
- **Docker** and **docker-compose**

## Installation and Setup

To install our project, clone the repository to your computer using the command:

```bash
git clone https://github.com/eldos02/Meloman.git
```
Then open a terminal or linux-like terminial on Windows and navigate to the root folder of our project. Now you can run the project with the instructions below: 

## Dockerized Web Application - University Canteen Management System

**Install Docker and Docker Compose (if needed):**

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
```

```bash
curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
```
**Download required Go packages:**

```bash
go mod tidy
```

**Start the container:**

```bash
docker-compose build
docker-compose up
```

## Servers

- **The application will be available at:**

```bash
http://localhost:8070
```

- ***Database Management:* pgAdmin will be available at:**

```bash
http://localhost:5070
```

- *Login: prof.aka777@gmail.com*
- *Password: kaskelen*

After authorizations in pgAdmin create a server with the following parameters

- name: db
- host name / address: db
- maintenance database: sdu
- username: sdudent
- password: kek

## Usage Instructions
- *Authentication:*
  Create a user. By default, the first registered user becomes an admin.
- *Manage Dishes:*
  After logging in, you'll be prompted to add a dish. Once added, it will be visible at /orders. You can delete/update the menu if desired.
- *Statistics:*
  At /analytics, view order analytics, including the number of orders and average checkout.

## Recommendations
If you don't know how to use Docker , go to: Docker Overview and check out the Docker documentation. That's it, follow these instructions and you will successfully launch our project on your computer!

If app-1 fails to start on the first try, just do docker-compose up again. 

## Contributors
**The project was developed by:**
- **Sanabek Yeldos:** [GitHub](https://github.com/eldos02)
- **Yerdali Akarys:** [GitHub](https://github.com/profaka)
- **Bakirbayev Sanzhar:** [GitHub](https://github.com/etozhegatito)
