package responses

type GetUserAddressesResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Street   string `json:"street"`
	House    string `json:"house"`
	Flat     string `json:"flat"`
	PostCode string `json:"post_code"`
	City     string `json:"city"`
	Phone    string `json:"phone"`
	Country  string `json:"country"`
}
