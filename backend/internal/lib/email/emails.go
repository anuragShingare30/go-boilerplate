package email

// @dev Utility functions to send different types of emails
func (c *Client) SendWelcomeEmail(to string, firstName string) error {
	data := map[string]string{
		"UserFirstName": firstName,
	}

	return c.SendEmail(
		to,
		"Welcome to our application",
		TemplateWelcome,
		data,
	)
}
