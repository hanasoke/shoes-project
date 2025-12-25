package shoemodel

import (
	"shoes-project/config"
	"shoes-project/entities"
)

func IsSkuExists(sku string) bool {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM shoes WHERE shoe_sku = ?", sku).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func GetAll() []entities.Shoe {
	rows, err := config.DB.Query(`
		SELECT 
			shoes.shoe_id, 
			shoes.shoe_name, 
			brands.brand_name as brand_name,
			shoes.shoe_type,
			shoes.shoe_description,
			shoes.shoe_sku,
			shoes.shoe_price,
			shoes.shoe_stock,
			shoes.created_at,
			shoes.updated_at 
		FROM shoes 
		JOIN brands ON shoes.brand_id = brands.brand_id 
	`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var shoes []entities.Shoe

	for rows.Next() {
		var shoe entities.Shoe
		if err := rows.Scan(&shoe.Id, &shoe.Name, &shoe.Brand.Name, &shoe.Type, &shoe.Description, &shoe.Sku, &shoe.Price, &shoe.Stock, &shoe.CreatedAt, &shoe.UpdatedAt); err != nil {
			panic(err)
		}

		shoes = append(shoes, shoe)
	}

	return shoes
}

func Create(shoe entities.Shoe) bool {
	// Check for duplicate SKU before insert (additional safety)
	if IsSkuExists(shoe.Sku) {
		return false
	}

	result, err := config.DB.Exec(`
		INSERT INTO shoes (name, brand_id, shoe_type, shoe_description, shoe_sku, shoe_price, shoe_stock, created_at)
	VALUES (?,?,?,?,?,?,?,?)`,
		shoe.Name,
		shoe.Brand.Id,
		shoe.Type,
		shoe.Description,
		shoe.Sku,
		shoe.Price,
		shoe.Stock,
		shoe.CreatedAt,
	)

	if err != nil {
		panic(err)
	}

	LastInsertId, err := result.LastInsertId()
	if err != nil {
		return false
	}

	return LastInsertId > 0
}

// Additional helper function for duplicate checking
func FindBySku(sku string) (entities.Shoe, error) {
	var shoe entities.Shoe
	err := config.DB.QueryRow(`
		SELECT id, name, sku
		FROM shoes WHERE sku = ?`, sku).Scan(
		&shoe.Id,
		&shoe.Name,
		&shoe.Sku,
	)
	return shoe, err
}
