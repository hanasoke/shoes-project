package shoemodel

import (
	"shoes-project/config"
	"shoes-project/entities"
)

func GetAll() []entities.Shoe {
	rows, err := config.DB.Query(`
		SELECT 
			shoes.shoe_id, 
			shoes.shoe_name, 
			brands.brand_name,
			shoes.shoe_type,
			shoes.shoe_description,
			shoes.shoe_sku,
			shoes.shoe_price,
			shoes.shoe_stock,
			shoes.created_at,
			shoes.updated_at 
		FROM shoes 
		JOIN brands ON shoes.brand_id = brands.brand_id
		ORDER BY shoes.created_at DESC 
	`)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var shoes []entities.Shoe

	for rows.Next() {
		var shoe entities.Shoe
		err := rows.Scan(
			&shoe.Id,
			&shoe.Name,
			&shoe.Brand.Brand_Name,
			&shoe.Type,
			&shoe.Description,
			&shoe.Sku,
			&shoe.Price,
			&shoe.Stock,
			&shoe.CreatedAt,
			&shoe.UpdatedAt,
		)

		if err != nil {
			panic(err)
		}

		shoes = append(shoes, shoe)
	}
	return shoes
}

func IsSkuExists(sku string) bool {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM shoes WHERE shoe_sku = ?", sku).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func Create(shoe entities.Shoe) bool {
	result, err := config.DB.Exec(`
		INSERT INTO shoes (shoe_name, brand_id, shoe_type, shoe_description, shoe_sku, shoe_price, shoe_stock, created_at)
	VALUES (?,?,?,?,?,?,,?,?)`,
		shoe.Name,
		shoe.Brand.Brand_Id,
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
	result.LastInsertId()
	if err != nil {
		panic(err)
	}

	return LastInsertId > 0
}
