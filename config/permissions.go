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

/*
Variant Permissions
*/
const (
	ReadVariant   Permission = "variant_read"
	CreateVariant Permission = "variant_create"
	EditVariant   Permission = "variant_edit"
	DeleteVariant Permission = "variant_delete"
)

/*
Category Permissions
*/
const (
	CreateCategory Permission = "category_create"
	EditCategory   Permission = "category_edit"
	DeleteCategory Permission = "category_delete"
)

/*
Subcategory Permissions
*/
const (
	CreateSubcategory Permission = "subcategory_create"
	EditSubcategory   Permission = "subcategory_edit"
	DeleteSubcategory Permission = "subcategory_delete"
)

var AllPermissions = []Permission{
	ReadProduct, CreateProduct, EditProduct, DeleteProduct,
	ReadUser, EditUser, DeleteUser,
	ReadVariant, CreateVariant, EditVariant, DeleteVariant,
	CreateCategory, EditCategory, DeleteCategory,
	CreateSubcategory, EditSubcategory, DeleteSubcategory,
}
