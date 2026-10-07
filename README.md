# Ikou API 🌐

This repository contains the API for the [Ikou](https://ikou-web.netlify.app/) project, a community-driven travel app, designed to make trip organizing with friends and peers a breeze, providing recommendations and allowing users to explore and refer back to places, activities, and trips created by others in the community. It's structured to provide support for the Ikou website frontend, serving as its backend counterpart. 

> 🚨 This is an ongoing project and subject to significant changes. Detailed documentation will be provided as the project matures.

> 🤖 Starting from October 2026, this project is built primarily with AI coding agents, following the workflow in [CLAUDE.md](CLAUDE.md).

## Table of Contents
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Contribution](#contribution)
- [Contact](#contact)

## Technology Stack 💻
- **Language:** Go
- **Database:** PostgreSQL ([Neon](https://neon.tech/))
- **Containerization:** Docker

## Getting Started
1. Clone the repository.
2. Navigate to the project directory and update the app.env file with your appropriate database details for local development.
4. Run the project using Makefile:
   ```sh
   make start

## Docker Support 🐳
For those who prefer Docker, a `docker-compose.yaml` file is included in the project. Feel free to utilize it if you prefer to run the application in containers.

## Usage 🛠️
This API is primarily structured to support the Ikou frontend, serving as its backend counterpart. However, it can also run independently as a standalone API server. For more interaction details with the frontend, please refer to the [Ikou Frontend Repository](https://github.com/ngfenglong/ikou-website).

## Project Structure 🌳
```plaintext
api
├─ config
├─ controllers
├─ dto
├─ mapper
├─ middleware
├─ models
├─ repositories
├─ routes
├─ store
└─ server.go
cmd
└─ main.go
internal
├─ helper
└─ util
```

> **Note:** The `data-seeding` and `dist` directories are not included in this repository as they are gitignored.


## Contribution 🤝
This project is currently in its infancy. Contributions, ideas, and bug reports are very welcome!

## Contact 📬
For any inquiries or clarifications related to this project, please contact [zell_dev@hotmail.com](mailto:zell_dev@hotmail.com).
