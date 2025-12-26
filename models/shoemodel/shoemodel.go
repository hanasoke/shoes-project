package shoemodel

import (
	"shoes-project/config"
	"shoes-project/entities"
)

func GetAll() []entities.Shoe {
	rows, err := config.DB.Query(`
		SELECT 		
			shoes.name,
			brands.name as brand_name,
			shoes.type,
			shoes.description,
			shoes.sku,
			shoes.price,
			shoes.stock,
			shoes.created_at
		FROM shoes
		JOIN brands ON shoes.id_brand = brands.id
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
			&shoe.Brand.Name,
			&shoe.Type,
			&shoe.Description,
			&shoe.SKU,
			&shoe.Price,
			&shoe.Stock,
			&shoe.CreatedAt,
		)

		if err != nil {
			panic(err)
		}

		shoes = append(shoes, shoe)
	}

	return shoes
}

func Create(shoe entities.Shoe) bool {
	result, err := config.DB.Exec(`
		INSERT INTO shoes(
			name, id_brand, type, description, sku, price, stock, created_at 
		) VALUES (?,?,?,?,?,?,?,?)`,
		shoe.Name,
		shoe.Brand.Id,
		shoe.Type,
		shoe.Description,
		shoe.SKU,
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
