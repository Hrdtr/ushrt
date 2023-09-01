# UShrt - Dead Simple Headless URL Shortener Service

UShrt is a lightweight and headless URL shortener service designed for seamless integration into your applications. It's built using Golang v1.21, PostgreSQL, Fiber, sqlc, golang-migrate, and provides API documentation using Swagger (generated with Swag).

## Features

- **Authentication:** Protect your shortening service with API key authentication.
- **CORS Configuration:** Easily configure cross-origin resource sharing.
- **Custom IDs/Slugs:** Create short links with custom IDs or slugs.
- **Dockerized:** Deploy with ease using Docker and find the image on Docker Hub at [hrdtr/ushrt](https://hub.docker.com/r/hrdtr/ushrt).
- **MIT License:** Use it freely in your projects.

## Quick Deploy

1. Create a `docker-compose.yml` file with the following contents:

```yaml
version: "3"

services:
  app:
    image: hrdtr/ushrt:latest
    command: ["./app"]
    environment:
      - APP_ENV=development
      - APP_API_KEY=changeme
      - APP_BASE_URL=http://localhost:3000
      - APP_CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:8000

      - POSTGRES_HOST=localhost
      - POSTGRES_PORT=5432
      - POSTGRES_USER=pguser
      - POSTGRES_PASSWORD=pgpassword
      - POSTGRES_DB=ushrt
      - POSTGRES_SSL_MODE=disable
    ports:
      - "3000:3000"
```

2. Run the following command to start the service:

```bash
docker-compose up -d
```

3. Your UShrt instance is now up and running!

## API Documentation

Explore the API endpoints and test them using the Swagger documentation. Simply visit the `/swagger/index.html` route of your deployed UShrt instance.

## Contributing

We welcome and encourage contributions from the community! If you'd like to contribute to the development of UShrt, here's how you can get started:

- Check out the [GitHub repository](https://github.com/Hrdtr/ushrt) for the latest code and open issues.

- Fork the repository and create a new branch for your contributions.

- Work on your changes, whether it's bug fixes, new features, or improvements to the documentation.

- Submit a pull request (PR) with your changes, and our team will review it.

- Join our community discussions on [GitHub Issues](https://github.com/Hrdtr/ushrt/issues) to share your ideas, report bugs, or ask questions.

By contributing to UShrt, you become part of an open and collaborative project that aims to simplify URL shortening for everyone. We appreciate your support and look forward to your contributions!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Special thanks to the open-source community for their contributions.

## Issues and Feedback

Please report issues or provide feedback on [GitHub Issues](https://github.com/Hrdtr/ushrt/issues).

---

Thank you for choosing UShrt! We hope it simplifies your URL shortening needs.
