package requests

type InsertUserAddressRequest struct {
	Name     string `json:"name" binding:"required,max=60"`
	Surname  string `json:"surname" binding:"required,max=80"`
	Street   string `json:"street" binding:"required,max=60"`
	House    string `json:"house" binding:"required,max=20"`
	Flat     string `json:"flat,omitempty" binding:"max=20"`
	PostCode string `json:"post_code" binding:"required,max=10"`
	City     string `json:"city" binding:"required,max=60"`
	Phone    string `json:"phone" binding:"required,max=20"`
	Country  string `json:"country" binding:"required,max=60"`
}

type DeleteUserAddressRequest struct {
	AddressID string `json:"address_id" binding:"required"`
}
