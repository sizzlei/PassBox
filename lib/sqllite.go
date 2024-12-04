package lib 

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"fmt"
	"time"
)

type DBO struct {
	Client 		*sql.DB
}

type Configure struct {
	Upper		string 
	Lower 		string 
	Digits 		string 
	Special 	string
	UseSpecial 	int
	UseDigits	int 
	PassLength  int
}

func GetTime() string {
	now := time.Now().UTC()
	curKst := now.Add(time.Hour * 9)

	return curKst.Format("2006-01-02 15:04:05")
}

func InitDB(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3",fmt.Sprintf("%s",file))
	if err != nil {
		return nil, err
	}

	createBox := `
		CREATE TABLE IF NOT EXISTS passbox (
		box_id integer primary key  autoincrement,
		box_name text,
		box_password text,
		created_at text,
		updated_at text,
		unique (box_name)
		)
	`

	createBoxHist := `
		CREATE TABLE IF NOT EXISTS passbox_hist (
		box_name text,
		box_password text,
		created_at text
		)
	`

	createHistIndex := `
		CREATE INDEX idx_passbox_hist_01 ON passbox_hist(box_name)
	`

	createBoxSetup := `
		CREATE TABLE IF NOT EXISTS passbox_config (
			conf_id integer primary key autoincrement,
			upper_char text,
			lower_char text,
			digits_char text,
			special_char text,
			use_sepcial integer,
			use_digits integer,
			pass_length integer
		)
	`

	_, err = db.Exec(createBox)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createBoxHist)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createHistIndex)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createBoxSetup)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectStorage(file string) (DBO, error){
	db, err := sql.Open("sqlite3",fmt.Sprintf("%s",file))
	if err != nil {
		return DBO{}, err
	}

	return DBO{
		Client: db,
	}, nil
}

func (d DBO) CheckConfigure() (Configure, error) {
	query := `
		SELECT upper_char,lower_char,digits_char,special_char,use_sepcial,use_digits,pass_length FROM passbox_config WHERE conf_id = 1
	`

	var cnf Configure
	err := d.Client.QueryRow(query).Scan(
		&cnf.Upper,
		&cnf.Lower,
		&cnf.Digits,
		&cnf.Special,
		&cnf.UseSpecial,
		&cnf.UseDigits,
		&cnf.PassLength,
	)
	if err != nil {
		return Configure{}, err
	}

	return cnf, nil
}

func (d DBO) UpdateConfigure(cnf Configure) error {
	query := `
		UPDATE passbox_config
		SET 
			upper_char = ?,
			lower_char = ?,
			digits_char = ?,
			special_char = ?,
			use_sepcial = ?,
			use_digits = ?,
			pass_length = ?
		WHERE conf_id = 1

	`

	_, err := d.Client.Exec(query,
		cnf.Upper,
		cnf.Lower,
		cnf.Digits,
		cnf.Special,
		cnf.UseSpecial,
		cnf.UseDigits,
		cnf.PassLength,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d DBO) SetConfigure(cnf Configure) error {
	query := `
		INSERT INTO passbox_config(upper_char,lower_char,digits_char,special_char,use_sepcial,use_digits,pass_length)
		values (?,?,?,?,?,?,?)
	`

	_, err := d.Client.Exec(query,
		cnf.Upper,
		cnf.Lower,
		cnf.Digits,
		cnf.Special,
		cnf.UseSpecial,
		cnf.UseDigits,
		cnf.PassLength,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d DBO) AddPassBox(box Passbox) error {
	query := `
		INSERT INTO passbox(box_name,box_password,created_at,updated_at)
		values (?,?,?,?)
	`

	_, err := d.Client.Exec(query,box.Name,box.Pass,GetTime(),GetTime())
	if err != nil {
		return err
	}

	return nil
}

func (d DBO) GetPassList() ([]Passbox, error) {
	query := `
		SELECT
			box_id, box_name, box_password,updated_at
		FROM passbox
	`
	var box []Passbox
	data, err := d.Client.Query(query)
	if err != nil {
		return box, err
	}
	defer data.Close() 

	for data.Next() {
		var b Passbox 
		err := data.Scan(
			&b.Id,
			&b.Name,
			&b.Pass,
			&b.Updated,
		)
		if err != nil {
			return box, err
		}
		box = append(box,b)
	}

	return box, nil
}

func (d DBO) UpdatePass(boxId int,pass string) error {
	query := `
		UPDATE passbox 
		SET 
			box_password = ?,
			updated_at = ?
		WHERE box_id = ?
	`

	_, err := d.Client.Exec(query,pass,GetTime(),boxId)
	if err != nil {
		return err
	}
	return nil
}

func (d DBO) DeletePass(boxId int) error {
	query := `
		DELETE FROM passbox WHERE box_id = ?
	`

	_, err := d.Client.Exec(query,boxId)
	if err != nil {
		return err
	}
	return nil
}

func (d DBO) WriteHist(box Passbox) error {
	query := `
		INSERT INTO passbox_hist (box_name,box_password,created_at)
		VALUES (?,?,?)
	`

	_, err := d.Client.Exec(query, box.Name,box.Pass,GetTime())
	if err != nil {
		return err
	}
	return nil
}

func (d DBO) GetPassHist(boxName string) ([]Passbox, error) {
	query := `
		SELECT
			box_name, box_password, created_at
		FROM passbox_hist 
		WHERE box_name = ?
		ORDER BY created_at desc
	`
	var box []Passbox
	data, err := d.Client.Query(query, boxName)
	if err != nil {
		return box, err
	}
	defer data.Close() 

	for data.Next() {
		var b Passbox 
		err := data.Scan(
			&b.Name,
			&b.Pass,
			&b.Created,
		)
		if err != nil {
			return box, err
		}
		box = append(box,b)
	}

	return box, nil
}