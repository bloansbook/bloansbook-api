package email

import (
	"context"
	"fmt"
	"html"

	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/resend/resend-go/v3"
)

// Service handles email sending via Resend API
type Service struct {
	apiKey string
	from   string
}

// NewService creates a new email service
func NewService(fromEmail string) *Service {
	return &Service{
		apiKey: config.ApplicationConfig.Resend.APIKey,
		from:   fromEmail,
	}
}

// SendEmail sends an email using the Resend API
func (s *Service) SendEmail(ctx context.Context, to []string, subject, htmlContent string) error {
	if s.apiKey == "" {
		return fmt.Errorf("Resend API key is not configured")
	}

	client := resend.NewClient(s.apiKey)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      to,
		Subject: subject,
		Html:    htmlContent,
	}
	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// SendStaffWelcomeEmail sends a welcome email with credentials to a new staff member
func (s *Service) SendStaffWelcomeEmail(ctx context.Context, toEmail, staffID, password, firstName, lastName string) error {
	subject := "Welcome to BloansBook - Your Login Credentials"

	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Welcome to BloansBook</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background: #3498db; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background: #f9f9f9; }
				.credentials { background: white; padding: 15px; border-radius: 5px; margin: 20px 0; }
				.credential-row { margin: 10px 0; }
				.label { font-weight: bold; color: #3498db; }
				.footer { text-align: center; margin-top: 20px; padding: 20px; color: #777; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Welcome to BloansBook</h1>
				</div>

				<div class="content">
					<p>Dear %s %s,</p>

					<p>Welcome to BloansBook! Your staff account has been successfully created. Here are your login credentials:</p>

					<div class="credentials">
						<div class="credential-row">
							<span class="label">Staff ID:</span> %s
						</div>
						<div class="credential-row">
							<span class="label">Password:</span> %s
						</div>
					</div>

					<p>For security reasons, we recommend you change your password upon first login.</p>
					<p>If you have any questions, please contact your administrator.</p>

					<p>Best regards,<br>The BloansBook Team</p>
				</div>

				<div class="footer">
					<p>This is an automated message. Please do not reply to this email.</p>
				</div>
			</div>
		</body>
		</html>
	`, html.EscapeString(firstName), html.EscapeString(lastName), html.EscapeString(staffID), html.EscapeString(password))

	return s.SendEmail(ctx, []string{toEmail}, subject, htmlContent)
}
