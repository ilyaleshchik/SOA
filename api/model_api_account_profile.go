package api

type AccountProfile struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date"`
	Bio       string `json:"bio"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Name      string `json:"name"`
}

type UpdateAccountProfileRequest struct {
	Name     *string `json:"name"`
	LastName *string `json:"last_name"`
	Birthday *string `json:"birtday"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}
