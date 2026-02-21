package email

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	FromName string
}

type EmailService struct {
	config SMTPConfig
}

func NewEmailService(cfg SMTPConfig) *EmailService {
	return &EmailService{config: cfg}
}

func (s *EmailService) SendOTP(to, code, purpose string) error {
	subject := "Kode OTP BukaOutlet"
	purposeLabel := "Login Admin Panel"
	if purpose == "mitra_login" {
		purposeLabel = "Login Mitra Portal"
	} else if purpose == "reset_password" {
		purposeLabel = "Reset Password"
	}

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family:'Segoe UI',Arial,sans-serif;background:#f4f7fb;margin:0;padding:0;">
<div style="max-width:480px;margin:40px auto;background:#fff;border-radius:16px;overflow:hidden;box-shadow:0 4px 24px rgba(0,0,0,0.08);">
  <div style="background:linear-gradient(135deg,#2563eb,#7c3aed);padding:32px;text-align:center;">
    <h1 style="color:#fff;margin:0;font-size:24px;">🔐 Kode OTP</h1>
    <p style="color:rgba(255,255,255,0.8);margin:8px 0 0;font-size:14px;">%s</p>
  </div>
  <div style="padding:32px;text-align:center;">
    <p style="color:#64748b;margin:0 0 24px;font-size:15px;">Masukkan kode berikut untuk melanjutkan:</p>
    <div style="background:#f1f5f9;border-radius:12px;padding:20px;display:inline-block;min-width:200px;">
      <span style="font-size:36px;font-weight:800;letter-spacing:8px;color:#1e293b;">%s</span>
    </div>
    <p style="color:#94a3b8;margin:24px 0 0;font-size:13px;">Kode ini berlaku selama <strong>5 menit</strong>.</p>
    <p style="color:#94a3b8;margin:8px 0 0;font-size:13px;">Jika Anda tidak meminta kode ini, abaikan email ini.</p>
  </div>
  <div style="background:#f8fafc;padding:16px;text-align:center;border-top:1px solid #e2e8f0;">
    <p style="color:#94a3b8;margin:0;font-size:12px;">© BukaOutlet — Sistem Manajemen Kemitraan</p>
  </div>
</div>
</body>
</html>`, purposeLabel, code)

	return s.sendMail(to, subject, body)
}

func (s *EmailService) sendMail(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.From)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	err := smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(msg.String()))
	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("✅ OTP email sent to %s", to)
	return nil
}
