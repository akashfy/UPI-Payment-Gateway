# Security Policy

Security is a top priority for UPIShield Gateway. Please follow these guidelines to keep your deployment and users safe.

## Best Practices

- **Never commit your `.env` file**: Ensure it is included in your `.gitignore`.
- **Use a Gmail App Password**: Never use your primary Google account password in the `.env` file. Create a dedicated [App Password](https://myaccount.google.com/apppasswords).
- **Use HTTPS in production**: Always serve the application over HTTPS using a reverse proxy (e.g., Nginx, Caddy) or a similar service.
- **Restrict CORS in production**: Ensure proper CORS headers are set if your frontend and backend are hosted separately.
- **Do not upload private product ZIP files publicly**: Keep your digital products safe by placing them outside the public web root or restricting access properly.

## Reporting a Vulnerability

If you discover a security vulnerability, please **report it through GitHub issues** or contact the repository owner directly. Do not disclose it publicly until it has been patched.
