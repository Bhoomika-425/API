package package

import (
    "fmt"
    "net/smtp"
)

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	otp := rand.Intn(max-min+1) + min
	otpString := fmt.Sprintf("%06d", otp)

	return otpString
}


 
func main() {
    // Sender's email address and password
    from := "bhoomikasabalur@gmail.com"
    password := "hghf uxjo oqgf qmzx"
 
    // Recipient's email address
    to := "niki16809@gmail.com"
 
    // SMTP server details
    smtpServer := "smtp.gmail.com"
    smtpPort := 587
 
    // Message content
    message := []byte("Subject: Test Email\n\nThis is a test email body.")
 
    // Authentication information
    auth := smtp.PlainAuth("", from, password, smtpServer)
 
    // SMTP connection
    smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
    err := smtp.SendMail(smtpAddr, auth, from, []string{to}, message)
    if err != nil {
        fmt.Println("Error sending email:", err)
        return
    }
 
    fmt.Println("Email sent successfully!")
}