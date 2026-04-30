package main

import (
	"crypto/tls"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"os"
	"strings"

	"github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
)

//go:embed website/*
var websiteFS embed.FS

const GmailIMAP = "imap.gmail.com:993"

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type VerifyRequest struct {
	TxnID string `json:"txn_id"`
}

type VerifyResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

func main() {
	_ = godotenv.Load()

	// Embed website files — strip "website/" prefix so / serves index.html
	port := getEnv("SERVER_PORT", ":8080")

	webFS, _ := fs.Sub(websiteFS, "website")
	http.Handle("/", http.FileServer(http.FS(webFS)))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/api/config", handleConfig)
	http.HandleFunc("/api/verify-payment", handleVerifyPayment)

	fmt.Println("🚀 Server running at http://localhost" + port)
	fmt.Println("📦 HTML embedded in binary")
	fmt.Println("📁 Assets served from ./assets/")
	fmt.Println("🔐 Payment verification at /api/verify-payment")
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]string{
		"upi_id": getEnv("UPI_ID", ""),
	})
}

func handleVerifyPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(VerifyResponse{false, "Method not allowed"})
		return
	}

	var req VerifyRequest
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &req); err != nil || req.TxnID == "" {
		json.NewEncoder(w).Encode(VerifyResponse{false, "Please enter a valid Transaction ID"})
		return
	}

	txnID := strings.TrimSpace(req.TxnID)
	fmt.Printf("🔍 Verifying TXN: %s\n", txnID)

	result := searchGmail(txnID)

	switch result {
	case "verified":
		fmt.Printf("✅ VERIFIED: %s\n", txnID)
		json.NewEncoder(w).Encode(VerifyResponse{true, "Payment verified! ✅"})
	case "wrong_amount":
		fmt.Printf("⚠️ WRONG AMOUNT: %s\n", txnID)
		json.NewEncoder(w).Encode(VerifyResponse{false, "❌ Payment amount is not ₹299. Please pay the correct amount of ₹299 to get access."})
	default:
		fmt.Printf("❌ NOT FOUND: %s\n", txnID)
		json.NewEncoder(w).Encode(VerifyResponse{false, "Transaction not found. Please wait 2-3 minutes after payment and try again. If issue persists, contact support on WhatsApp."})
	}
}

func searchGmail(txnID string) string {
	tlsConf := &tls.Config{ServerName: "imap.gmail.com"}
	c, err := imapclient.DialTLS(GmailIMAP, tlsConf)
	if err != nil {
		log.Printf("IMAP connect fail: %v", err)
		return "not_found"
	}
	defer c.Logout()

	email := getEnv("GMAIL_EMAIL", "")
	pass := getEnv("GMAIL_APP_PASSWORD", "")
	if err := c.Login(email, pass); err != nil {
		log.Printf("IMAP login fail: %v", err)
		return "not_found"
	}

	mbox, err := c.Select("INBOX", true)
	if err != nil {
		log.Printf("Select INBOX fail: %v", err)
		return "not_found"
	}
	if mbox.Messages == 0 {
		return "not_found"
	}

	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 100 {
		from = mbox.Messages - 99
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 100)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	result := "not_found"
	for msg := range messages {
		if result == "verified" {
			continue
		}
		for _, literal := range msg.Body {
			if literal == nil {
				continue
			}
			fullText := extractAllText(literal)
			upper := strings.ToUpper(fullText)

			hasTxn := strings.Contains(upper, strings.ToUpper(txnID))
			hasAmount := strings.Contains(fullText, "299")

			if hasTxn && hasAmount {
				result = "verified"
				fmt.Printf("   ✅ Match: TXN=%s + ₹299 verified\n", txnID)
			} else if hasTxn && !hasAmount {
				result = "wrong_amount"
				fmt.Printf("   ⚠️  TXN found but amount NOT ₹299\n")
			}
		}
	}
	<-done
	return result
}

func extractAllText(r io.Reader) string {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		raw, _ := io.ReadAll(r)
		return string(raw)
	}

	contentType := msg.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/plain"
	}

	return decodeBody(msg.Body, contentType, msg.Header.Get("Content-Transfer-Encoding"))
}

func decodeBody(r io.Reader, contentType string, encoding string) string {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		raw, _ := io.ReadAll(r)
		return decodeTransfer(string(raw), encoding)
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			raw, _ := io.ReadAll(r)
			return decodeTransfer(string(raw), encoding)
		}

		mr := multipart.NewReader(r, boundary)
		var allText strings.Builder

		for {
			part, err := mr.NextPart()
			if err != nil {
				break
			}
			partCT := part.Header.Get("Content-Type")
			partEnc := part.Header.Get("Content-Transfer-Encoding")
			if partCT == "" {
				partCT = "text/plain"
			}
			allText.WriteString(decodeBody(part, partCT, partEnc))
			allText.WriteString(" ")
		}
		return allText.String()
	}

	raw, _ := io.ReadAll(r)
	return decodeTransfer(string(raw), encoding)
}

func decodeTransfer(s string, encoding string) string {
	encoding = strings.ToLower(strings.TrimSpace(encoding))
	if encoding == "base64" {
		s = strings.ReplaceAll(s, "\r\n", "")
		s = strings.ReplaceAll(s, "\n", "")
		decoded, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			decoded, err = base64.RawStdEncoding.DecodeString(s)
			if err != nil {
				return s
			}
		}
		return string(decoded)
	}
	if encoding == "quoted-printable" {
		result := strings.ReplaceAll(s, "=\r\n", "")
		result = strings.ReplaceAll(result, "=\n", "")
		return result
	}
	return s
}
