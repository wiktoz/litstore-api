package enums

type SelectType string

const (
	Button SelectType = "button"
	Select SelectType = "select"
)

func IsValidSelectType(t string) bool {
	switch SelectType(t) {
	case Button, Select:
		return true
	default:
		return false
	}
}
