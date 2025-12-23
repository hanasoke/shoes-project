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
			brands.name as brand_name,
			shoes.shoe_type,
			shoes.shoe_size,
			shoes.shoe_description,
			shoes.shoe_sku,
			shoes.shoe_price,
			shoes.shoe_stock,
			shoes.shoe_photo,
			shoes.shoe_status,
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
		err := rows.Scan(
			&shoe.Id,
			&shoe.Name,
			&shoe.Brand.Brand_Name,
			&shoe.Type,
			&shoe.Size,
			&shoe.Description,
			&shoe.Sku,
			&shoe.Price,
			&shoe.Stock,
			&shoe.Photo,
			&shoe.Status,
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
