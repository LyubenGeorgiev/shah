# Shah Chess Web App

Shah is a Golang web application for playing chess with complete server-side rendering. It offers features for playing against an opponent or a computer, incorporating cookie-based authentication with options for Login, Register, and Logout functionalities. The application includes a match history section within user accounts, as well as chat for live game communication using WebSockets and another chat for the chat section.

## Features

- Play chess against opponents or the computer.
- Complete server-side rendering.
- User authentication with Login, Register, and Logout functionalities.
- Match history tracking in the account section.
- Real-time chat and games using WebSockets.
- Admin panel for users with ADMIN role, allowing editing of user profiles.
- Picture uploading feature in the account section.
- News section with the latest chess updates, editable by admins.

## Technologies Used

- **Backend**: Golang
- **Frontend**: HTMX, Templ
- **Authentication**: Cookie-based
- **Game + Chat**: WebSockets
- **CI/CD Pipeline**: GitHub Actions

## Getting Started

To start the project, run the following command:

```bash
docker compose up --build
```

This will build and run the Docker containers required for the application.

## Chess Engine
The project includes a custom chess engine written in Golang, boasting a strength of around 2400 Elo. It is inspired by the [Chess Programming](https://www.youtube.com/@chessprogramming591) bitboard chess engine [series](https://www.youtube.com/playlist?list=PLmN0neTso3Jxh8ZIylk74JpwfiWNI76Cs)

## DevOps Pipeline

The project features a comprehensive CI/CD pipeline managed through GitHub Actions. The pipeline includes:

### Dev Pipeline

- gosec: Runs the Gosec Security Scanner on the codebase.
- gitleaks: Executes Gitleaks to check for sensitive information leakage.
- test: Runs tests after security checks.
- build: Builds Docker images, scans for vulnerabilities using Trivy.

### Main Pipeline

- gosec, gitleaks, test: Similar jobs as in the Dev Pipeline, executing security checks and testing.
- build: Builds Docker images, scans for vulnerabilities using Trivy, and uploads to Docker Hub.
- trigger-infra: Triggers deployment workflow in the [shah-infra](https://github.com/LyubenGeorgiev/shah-infra) repository after successful build and push to Docker Hub.

## Contributing

Contributions are welcome! Feel free to open issues for feature requests or bug reports. Pull requests are also appreciated

## License

This project is licensed under the [MIT License](LICENSE)