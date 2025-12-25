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
		"SELECT id FROM brands WHERE name = ? LIMIT 1",
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
		SELECT id FROM brands 
		WHERE name = ? AND id != ?
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
		if err := rows.Scan(&brand.Id, &brand.Name, &brand.CreatedAt, &brand.UpdatedAt); err != nil {
			panic(err)
		}

		brands = append(brands, brand)
	}

	return brands
}

func Create(brand entities.Brand) error {
	exists, err := IsBrandExists(brand.Name)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateBrand
	}

	_, err = config.DB.Exec(`
		INSERT INTO brands (name, created_at, updated_at)
		VALUES (?, ?, ?)`,
		brand.Name,
		brand.CreatedAt,
		brand.UpdatedAt,
	)

	return err
}

func Detail(id int) entities.Brand {
	row := config.DB.QueryRow(`SELECT id, name FROM brands WHERE id = ?`, id)

	var brand entities.Brand
	if err := row.Scan(&brand.Id, &brand.Name); err != nil {
		panic(err.Error())
	}

	return brand
}

func Update(id int, brand entities.Brand) error {
	exists, err := IsBrandExistsExceptID(brand.Name, id)
	if err != nil {
		return err
	}

	if exists {
		return ErrDuplicateBrand
	}

	_, err = config.DB.Exec(`
		UPDATE brands 
		SET name = ?, updated_at = ?
		WHERE id = ?`,
		brand.Name,
		brand.UpdatedAt,
		id,
	)

	return err
}

func Delete(id int) error {
	_, err := config.DB.Exec(`DELETE FROM brands WHERE id = ?`, id)
	return err
}
