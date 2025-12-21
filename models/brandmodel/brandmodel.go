package brandmodel

import (
	"database/sql"
	"errors"
	"shoes-project/config"
	"shoes-project/entities"
)

var ErrDuplicateBrand = errors.New("brand already exists")

func IsBrandExists(name string) (bool, error) {
	var id uint
	err := config.DB.QueryRow(
		"SELECT brand_id FROM brands WHERE brand_name = ? LIMIT 1",
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

func IsBrandExistsExceptID(name string, id int) (bool, error) {
	var brandID int
	err := config.DB.QueryRow(`
		SELECT brand_id FROM brands 
		WHERE brand_name = ? AND brand_id != ?
		LIMIT 1`,
		name, id,
	).Scan(&brandID)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

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

func Create(brand entities.Brand) error {
	exists, err := IsBrandExists(brand.Brand_Name)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateBrand
	}

	_, err = config.DB.Exec(`
		INSERT INTO brands (brand_name, created_at, updated_at)
		VALUES (?, ?, ?)`,
		brand.Brand_Name,
		brand.CreatedAt,
		brand.UpdatedAt,
	)

	return err
}

func Detail(brand_id int) entities.Brand {
	row := config.DB.QueryRow(`SELECT brand_id, brand_name FROM brands WHERE brand_id = ?`, brand_id)

	var brand entities.Brand
	if err := row.Scan(&brand.Brand_Id, &brand.Brand_Name); err != nil {
		panic(err.Error())
	}

	return brand
}

func Update(brand_id int, brand entities.Brand) error {
	exists, err := IsBrandExistsExceptID(brand.Brand_Name, brand_id)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateBrand
	}

	_, err = config.DB.Exec(`
		UPDATE brands 
		SET brand_name = ?, updated_at = ?
		WHERE brand_id = ?`,
		brand.Brand_Name,
		brand.UpdatedAt,
		brand_id,
	)

	return err
}
