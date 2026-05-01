# 🚀 UPIShield Gateway

### Self-hosted UPI Payment Verification Gateway for Digital Products

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)
![UPI](https://img.shields.io/badge/UPI-Payment_Verification-22c55e?style=for-the-badge)
![Self Hosted](https://img.shields.io/badge/Self--Hosted-Yes-purple?style=for-the-badge)

**UPIShield Gateway** is a simple, fast, and self-hosted payment verification system designed for selling digital products using UPI with **0% Transaction Charges**.

It automatically verifies a customer’s UPI transaction ID by securely checking your Gmail/IMAP inbox and unlocks the product download page instantly. Keep 100% of your earnings without any third-party gateway fees!

> ⭐ Star this repo and fork it if you want more self-hosted payment automation tools!

---

## ✨ Features

- ✅ **Automated UPI Payment Verification** using Transaction IDs
- ✅ **Gmail IMAP Integration** for real-time payment email scanning
- ✅ **Instant Digital Product Delivery** upon successful verification
- ✅ **High Performance** – Built with Go for speed and low memory usage
- ✅ **Docker Ready** – Simple deployment via container
- ✅ **Built-in Landing Page** – Clean and responsive UI included
- ✅ **0% Transaction Charges** – Keep 100% of your profits, absolutely zero hidden fees or API charges!

---

## 🛠️ How It Works

1. Customer scans your static UPI QR code on the landing page.
2. Customer makes the payment and enters the **UPI Transaction ID (UTR)**.
3. UPIShield Gateway connects to your Gmail via IMAP and searches for the bank payment receipt.
4. If the payment is verified, the user is redirected to a secure download page.

---

## 🚀 Quick Start

### Environment Setup

1. Rename `.env.example` to `.env`.
2. Configure the required environment variables:
   - `GMAIL_EMAIL`: Your Gmail address.
   - `GMAIL_APP_PASSWORD`: Your [Gmail App Password](https://myaccount.google.com/apppasswords).
   - `SERVER_PORT`: Server port (default is `:8080`).
   - `UPI_ID`: Your UPI ID (displayed on the landing page).
   - `REQUIRED_AMOUNT`: The minimum required payment amount (e.g., `299`).

### Run Locally

Make sure you have [Go](https://go.dev/dl/) installed.

```bash
# Install dependencies
go mod tidy

# Run the server
go run .
```

### Docker Run Command

```bash
# Build the Docker image
docker build -t upishield-gateway .

# Run the container
docker run -d -p 8080:8080 --env-file .env upishield-gateway
```

---

## 📁 Project Structure

```text
.
├── main.go             # Core Go backend server
├── Dockerfile          # Docker container configuration
├── .env.example        # Environment variable template
├── website/           # Embedded web templates and static files
└── assets/             # Images, QR code, and downloadable products
```

---

## 📸 Screenshots

| Landing Page | Payment Verify | Download Page |
| :---: | :---: | :---: |
| ![Landing Page](assets/screenshot-home.png) | ![Payment Verify](assets/screenshot-verify.png) | ![Download Page](assets/screenshot-download.png) |

---

## 🛡️ Security Notes

- **Never commit your `.env` file** to public repositories.
- Always use a **Gmail App Password**, not your primary Google account password.
- Run the server behind a reverse proxy with **HTTPS** in production.
- Do not upload your private product ZIP files publicly without proper access control.

---

## 🗺️ Roadmap

- [ ] Add support for webhook notifications
- [ ] Add support for multiple products
- [ ] Admin dashboard for transaction logs

---

## 🤝 Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to get started.

---

## ⚖️ License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## ⚠️ Disclaimer

This software is provided "as is", without warranty of any kind. The authors are not responsible for any financial loss or security breaches. Always test thoroughly before using in a production environment.
