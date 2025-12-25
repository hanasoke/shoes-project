package shoemodel

// func GetAll() []entities.Shoe {
// 	rows, err := config.DB.Query(`
// 		SELECT
// 			shoes.id,
// 			shoes.name,
// 			shoes.type,
// 			shoes.description,
// 			shoes.sku,
// 			shoes.price,
// 			shoes.stock,
// 			shoes.created_at,
// 			shoes.updated_at
// 		FROM shoes
// 		JOIN brands ON shoes.brand_id = brands.id
// 	`)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer rows.Close()

// 	var shoes []entities.Shoe

// 	for rows.Next() {
// 		var shoe entities.Shoe
// 		err := rows.Scan(
// 			&shoe.Id,
// 			&shoe.Name,
// 			&shoe.Type,
// 			&shoe.Description,
// 			&shoe.Sku,
// 			&shoe.Price,
// 			&shoe.Stock,
// 			&shoe.CreatedAt,
// 			&shoe.UpdatedAt,
// 		)

// 		if err != nil {
// 			panic(err)
// 		}

// 		shoes = append(shoes, shoe)
// 	}
// 	return shoes
// }

// func Create(shoe entities.Shoe) bool {
// 	result, err := config.DB.Exec(`
// 		INSERT INTO shoes(

// 		) VALUES ()
// 	`)
// }
