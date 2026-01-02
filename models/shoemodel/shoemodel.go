package shoemodel

import (
	"database/sql"
	"errors"
	"shoes-project/config"
	"shoes-project/entities"
)

var ErrDuplicateShoe = errors.New("shoe already exists")

func IsShoeExists(name string) (bool, error) {
	var id uint
	err := config.DB.QueryRow(
		"SELECT id FROM shoes WHERE name = ? LIMIT 1",
		name,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func IsShoeExistsExceptID(name string, id int) (bool, error) {
	var Id int
	err := config.DB.QueryRow(`
		SELECT id FROM shoes 
		WHERE name = ? AND id != ?
		LIMIT 1`,
		name, id,
	).Scan(&Id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func GetAll() []entities.Shoe {
	rows, err := config.DB.Query(`
		SELECT 
			shoes.id,		
			shoes.name,
			brands.name,
			shoes.type,
			shoes.description,
			shoes.sku,
			shoes.price,
			shoes.stock,
			shoes.created_at,
			brands.id as id_brand
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
		var brand entities.Brand
		var idBrand uint

		err := rows.Scan(
			&shoe.Id,
			&shoe.Name,
			&brand.Name, // Brand name
			&shoe.Type,
			&shoe.Description,
			&shoe.SKU,
			&shoe.Price,
			&shoe.Stock,
			&shoe.CreatedAt,
			&brand.Id, // Brand ID
		)

		if err != nil {
			panic(err)
		}

		brand.Id = idBrand
		shoe.Brand = brand
		shoes = append(shoes, shoe)
	}

	return shoes
}

func Create(shoe entities.ShoeCreate) error {
	exists, err := IsShoeExists(shoe.Name)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateShoe
	}

	_, err = config.DB.Exec(`
		INSERT INTO shoes (
			name, 
			id_brand, 
			type, 
			description, 
			sku, 
			price, 
			stock, 
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		shoe.Name,
		shoe.IdBrand,
		shoe.Type,
		shoe.Description,
		shoe.SKU,
		shoe.Price,
		shoe.Stock,
		shoe.CreatedAt,
	)

	return err
}

func Detail(id int) (entities.Shoe, error) {
	row := config.DB.QueryRow(`
		SELECT 
			shoes.id,		
			shoes.name,
			brands.id as brand_id,
			brands.name as brand_name,
			shoes.type,
			shoes.description,
			shoes.sku,
			shoes.price,
			shoes.stock,
			shoes.created_at
		FROM shoes
		JOIN brands ON shoes.id_brand = brands.id
		WHERE shoes.id = ? 
	`, id)

	var shoe entities.Shoe
	var brand entities.Brand

	err := row.Scan(
		&shoe.Id,
		&shoe.Name,
		&brand.Id,   // Brand ID
		&brand.Name, // Brand Name
		&shoe.Type,
		&shoe.Description,
		&shoe.SKU,
		&shoe.Price,
		&shoe.Stock,
		&shoe.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Shoe{}, errors.New("shoe not found")
		}
		return entities.Shoe{}, err
	}

	// Assign brand to shoe
	shoe.Brand = brand

	return shoe, nil
}

func Update(id int, shoe entities.ShoeUpdate) error {
	// Cek apakah nama sudah digunakan oleh shoe lain
	exists, err := IsShoeExistsExceptID(shoe.Name, id)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateShoe
	}

	// Update data shoe
	result, err := config.DB.Exec(`
			UPDATE shoes 
			SET 
				name = ?, 
				id_brand = ?,
				type = ?, 
				description = ?,
				sku = ?,
				price = ?, 
				stock = ?, 
				updated_at = ?
			WHERE id = ?`,
		shoe.Name,
		shoe.IdBrand,
		shoe.Type,
		shoe.Description,
		shoe.SKU,
		shoe.Price,
		shoe.Stock,
		shoe.UpdatedAt,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no shoe found to update")
	}

	return nil
}

func Delete(id int) error {
	_, err := config.DB.Exec(`DELETE FROM shoes WHERE id = ?`, id)
	return err
}
