package storage

import (
	"SrvCat/util"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type AccessCategory uint

const (
	AccessCategoryCheck AccessCategory = iota
	AccessCategoryVerified
)

type methods interface {
	// 检查初始化
	GetInit() (bool, error)
	// 更新初始化
	UpdateInit() error
	// 添加IP记录
	AddAccess(ip string, category AccessCategory) error
	// 检查验证频率
	GetAccessCount(after int64) (int, error)
	// 检查IP白名单
	GetVerified(ip string, after int64) (bool, error)
}

type sqlite struct {
	db *sql.DB
}

var (
	_      methods = (*sqlite)(nil)
	Sqlite         = new(sqlite)
)

func init() {
	db, err := sql.Open("sqlite3", "storage.db")
	util.FailOnException("An error occurred while open sqlite", err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS init (
	  	"initialized" integer(1)
	);`)
	util.FailOnException("An error occurred while create table init", err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS access (
  		"ip" TEXT(16) NOT NULL,
  		"time" integer(16) NOT NULL,
    	"category" integer(1) NOT NULL,
    	"used" integer(0) NOT NULL DEFAULT 0
	);`)
	util.FailOnException("An error occurred while create table access", err)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ip_index ON access ("ip");`)
	util.FailOnException("An error occurred while add ip index", err)
	Sqlite.db = db
}

func (s sqlite) GetInit() (bool, error) {
	rows, err := s.db.Query(`SELECT 1 FROM init`)
	defer rows.Close()
	if err != nil {
		return true, err
	}
	if rows.Next() {
		return false, nil
	}
	return true, nil
}

func (s sqlite) UpdateInit() error {
	_, err := s.db.Exec(`INSERT INTO init("initialized") VALUES (1)`)
	if err != nil {
		return err
	}
	return nil
}

func (s sqlite) AddAccess(ip string, category AccessCategory) error {
	stmt, err := s.db.Prepare(`INSERT INTO access("ip","time","category") VALUES (?,?,?)`)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ip, time.Now().UnixMilli(), category)
	if err != nil {
		return err
	}
	return nil
}

func (s sqlite) GetAccessCount(after int64) (int, error) {
	stmt, err := s.db.Prepare(`SELECT COUNT(1) FROM access WHERE "category" = ? AND "time" > ?`)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	rows, err := stmt.Query(AccessCategoryCheck, after)
	defer rows.Close()
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0, nil
}

func (s sqlite) GetVerified(ip string, after int64) (bool, error) {
	stmt, err := s.db.Prepare(`SELECT 1 FROM access WHERE "ip" = ? AND "category" = ? AND "time" > ?`)
	defer stmt.Close()
	if err != nil {
		return false, err
	}
	rows, err := stmt.Query(ip, AccessCategoryVerified, after)
	defer rows.Close()
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, err
	}
	return false, nil
}

func (s sqlite) GetUnusedVerify(ip string, after int64) (bool, error) {
	stmt, err := s.db.Prepare(`SELECT 1 FROM access WHERE "ip" = ? AND "category" = ? AND "time" > ? AND "used" = 0`)
	defer stmt.Close()
	if err != nil {
		return false, err
	}
	rows, err := stmt.Query(ip, AccessCategoryVerified, after)
	defer rows.Close()
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, err
	}
	return false, nil
}

func (s sqlite) UpdateUsed(ip string, after int64) error {
	stmt, err := s.db.Prepare(`UPDATE access SET "used" = 1 WHERE "ip" = ? AND "category" = ? AND "time" > ? AND "used" = 0`)
	defer stmt.Close()
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(ip, AccessCategoryVerified, after); err != nil {
		return err
	}
	return nil
}
