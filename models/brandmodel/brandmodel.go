package brandmodel

import (
	"shoes-project/config"
	"shoes-project/entities"
)

func GetAll() []entities.Brand {
	rows, err := config.DB.Query(`SELECT * FROM brands`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var brands []entities.Brand

	for rows.Next() {
		var brand entities.Brand
		if err := rows.Scan(&brand.Brand_Id, &brand.Brand_Name, &brand.CreatedAt, &brand.UpdatedAt); err != nil {
			panic(err)
		}

		brands = append(brands, brand)
	}

	return brands
}

func Create(brand entities.Brand) bool {
	result, err := config.DB.Exec(`
		INSERT INTO brands (brand_name, created_at, updated_at) VALUE (?, ?, ?)`,
		brand.Brand_Name, brand.CreatedAt, brand.UpdatedAt,
	)

	if err != nil {
		panic(err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	return lastInsertId > 0
}
