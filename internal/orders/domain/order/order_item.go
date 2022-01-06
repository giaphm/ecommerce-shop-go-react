package order

type OrderItem struct {
	uuid        string
	productUuid string
	quantity    int
}

func (ot OrderItem) GetUuid() string {
	return ot.uuid
}

func (ot OrderItem) GetProductUuid() string {
	return ot.productUuid
}

func (ot OrderItem) GetQuantity() int {
	return ot.quantity
}

func NewOrderItem(
	uuid string,
	productUuid string,
	quantity int,
) (*OrderItem, error) {

	if uuid == "" {
		return nil, ErrInvalidOrderItemUuid
	}
	if productUuid == "" {
		return nil, ErrInvalidOrderItemProductUuid
	}
	if quantity <= 0 {
		return nil, ErrInvalidOrderItemQuantity
	}

	return &OrderItem{
		uuid:        uuid,
		productUuid: productUuid,
		quantity:    quantity,
	}, nil
}
