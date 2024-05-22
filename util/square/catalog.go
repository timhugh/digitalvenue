package square

const (
	CatalogItemVariationType = "ITEM_VARIATION"
	CatalogItemType          = "ITEM"
)

type CatalogObject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type CatalogItemVariation struct {
	CatalogObject
	ItemVariationData struct {
		ItemID string `json:"item_id"`
	} `json:"item_variation_data"`
}

type CatalogItem struct {
	CatalogObject
	ItemData struct {
		Categories []struct {
			ID string `json:"id"`
		} `json:"categories"`
	} `json:"item_data"`
}
