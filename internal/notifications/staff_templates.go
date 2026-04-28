package notifications

import (
	"context"
	"html"

	"github.com/bloansbook/bloansbook-api/pkg/email"
)

// StaffEmailNotifier handles staff-related email notifications
type StaffEmailNotifier struct {
	emailService *email.Service
}

// NewStaffEmailNotifier creates a new staff email notifier
func NewStaffEmailNotifier(emailService *email.Service) *StaffEmailNotifier {
	return &StaffEmailNotifier{
		emailService: emailService,
	}
}

// SendWelcomeEmail sends a welcome email to a new staff member
func (n *StaffEmailNotifier) SendWelcomeEmail(ctx context.Context, toEmail, staffID, password, firstName, lastName string) error {
	return n.emailService.SendStaffWelcomeEmail(ctx, toEmail, staffID, password, firstName, lastName)
}

// GetStaffWelcomeTemplate returns the HTML template for staff welcome email
func GetStaffWelcomeTemplate(firstName, lastName, staffID, password string) string {
	return `
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
					<p>Dear ` + html.EscapeString(firstName) + ` ` + html.EscapeString(lastName) + `,</p>

					<p>Welcome to BloansBook! Your staff account has been successfully created. Here are your login credentials:</p>

					<div class="credentials">
						<div class="credential-row">
							<span class="label">Staff ID:</span> ` + html.EscapeString(staffID) + `
						</div>
						<div class="credential-row">
							<span class="label">Password:</span> ` + html.EscapeString(password) + `
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
	`
}
