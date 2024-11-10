package config

type Permission string

/*
Product Permissions
*/
const (
	ReadProduct   Permission = "product_read"
	CreateProduct Permission = "product_create"
	EditProduct   Permission = "product_edit"
	DeleteProduct Permission = "product_delete"
)

/*
User Permissions
*/
const (
	ReadUser   Permission = "user_read"
	EditUser   Permission = "user_edit"
	DeleteUser Permission = "user_delete"
)
