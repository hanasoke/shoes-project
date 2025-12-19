package brand

import (
	"shoes-project/config"
	"shoes-project/entities"
)

func GetAll() []entities.Brand {
	rows, err := config.DB.Query(`SELECT * FROM brand`)
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
