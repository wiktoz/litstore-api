package requests

type GetAvailableDeliveriesRequest struct {
	ItemIDs []string `gorm:"foreignKey:ItemID" json:"items"`
}
