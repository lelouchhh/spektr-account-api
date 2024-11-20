package domain

type Profile struct {
	ID             int    `json:"ID"`
	Name           string `json:"name"`
	FirstName      string `json:"firstName"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	FullName       string `json:"full_name"`
	Balance        string `json:"balance"`
	Tariff         string `json:"tariff"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"`
	InternetStatus bool   `json:"internet_status"`
}
