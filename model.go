package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type participant struct {
	ID     int    `json:"id"`
	AID    string `json:"aid"`
	Name   string `json:"name"`
	CID    int    `json:"cid"`
	Year   int    `json:"year"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
}

type college struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type dbresult interface {
	Scan(...interface{}) error
}

func (p *participant) getParticipant(db *sql.DB, by string) error {
	switch by {
	case "aid":
		statement := fmt.Sprintf("select id, aid, name, cid, year, mobile, email from participants where aid='%s'", p.AID)
		return p.scanParticipant(db.QueryRow(statement))
	case "id":
		statement := fmt.Sprintf("select id, aid, name, cid, year, mobile, email from participants where id='%d'", p.ID)
		return p.scanParticipant(db.QueryRow(statement))
	}
	return errors.New("Invalid option")
}

func (p *participant) updateParticipant(db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE participants SET aid=?, name=?, cid=?, year=?, mobile=? WHERE id=?")
	if err != nil {
		return err
	}
	p.AID = fmt.Sprintf("ab%05d", p.ID+925)
	//statement := fmt.Sprintf("INSERT INTO participants(aid, name, cid, year, mobile, email) VALUES('%s', '%s', %d, %d, '%s', '%s')", p.AID, p.Name, p.CID, p.Year, p.Mobile, p.Email)
	_, err = stmt.Exec(p.AID, p.Name, p.CID, p.Year, p.Mobile, p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *participant) deleteParticipant(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *participant) createParticipant(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO participants(email) VALUES(?)")
	if err != nil {
		return err
	}
	//statement := fmt.Sprintf("INSERT INTO participants(aid, name, cid, year, mobile, email) VALUES('%s', '%s', %d, %d, '%s', '%s')", p.AID, p.Name, p.CID, p.Year, p.Mobile, p.Email)
	res, err := stmt.Exec(p.Email)
	if err != nil {
		return err
	}
	lid, _ := res.LastInsertId()
	p.ID = int(lid)
	return nil
}

func getParticipants(db *sql.DB) ([]participant, error) {
	statement := "SELECT id, aid, name, cid, year, mobile, email FROM participants"
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	participants := []participant{}
	for rows.Next() {
		var p participant
		if err := p.scanParticipant(rows); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}

func getColleges(db *sql.DB) ([]college, error) {
	statement := "SELECT id, name FROM college"
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	colleges := []college{}
	for rows.Next() {
		var c college
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		colleges = append(colleges, c)
	}
	return colleges, nil
}

func (c *college) getCollege(db *sql.DB) error {
	statement := fmt.Sprintf("select id, name from college where id='%d'", c.ID)
	return db.QueryRow(statement).Scan(&c.ID, &c.Name)
}

func (p *participant) scanParticipant(res dbresult) error {
	var aid sql.NullString
	var name sql.NullString
	var cid sql.NullInt64
	var year sql.NullInt64
	var mobile sql.NullString
	if err := res.Scan(&p.ID, &aid, &name, &cid, &year, &mobile, &p.Email); err != nil {
		return err
	}
	p.AID = ""
	if aid.Valid {
		p.AID = aid.String
	}
	p.Name = ""
	if name.Valid {
		p.Name = name.String
	}
	p.CID = 0
	if cid.Valid {
		p.CID = int(cid.Int64)
	}
	p.Year = 0
	if year.Valid {
		p.Year = int(year.Int64)
	}
	p.Mobile = ""
	if mobile.Valid {
		p.Mobile = mobile.String
	}
	return nil
}
