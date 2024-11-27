package domain

type Profile struct {
	ID             string  `json:"ID"`
	FirstName      string  `json:"firstName"`
	MiddleName     string  `json:"middle_name"`
	LastName       string  `json:"last_name"`
	FullName       string  `json:"full_name"`
	Balance        float64 `json:"balance"`
	ToPay          float64 `json:"to_pay"`
	Tariff         string  `json:"tariff"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	Password       string  `json:"password"`
	InternetStatus bool    `json:"internet_status"`
	NextPayDate    string  `json:"next_pay_date"`
}
