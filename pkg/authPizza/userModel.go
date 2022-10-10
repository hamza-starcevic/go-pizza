package authPizza

type User struct {
	ID         string `json:"id,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Role       string `json:"role,omitempty"`
	Date_added string `json:"date_added,omitempty"`
}
