package persistence

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"log"
	"text/template"

	"github.com/zlowram/godan/model"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLPersistenceManager struct {
	db *sql.DB
}

const (
	mySQLQueryTemplate = `SELECT DISTINCT INET_NTOA(ip), port, service, content FROM banners
{{if or .Ip .Ports .Services .Regexp}} WHERE{{end}}
{{if .Ip}} ip = INET_ATON((?)){{end}}
{{if and .Ip .Ports}} AND{{end}}{{if .Ports}} port IN ({{range $i, $v := .Ports}}{{if $i}},{{end}}?{{end}}){{end}}
{{if and .Services (or .Ip .Ports)}} AND{{end}}{{if .Services}} service IN ({{range $i, $v := .Services}}{{if $i}},{{end}}?{{end}}){{end}}
{{if and .Regexp (or .Ip .Ports .Services)}} AND{{end}}{{if .Regexp}} content regexp (?){{end}};`
)

func NewMySQLPersistenceManager(dsn string) *MySQLPersistenceManager {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS banners (ip INT UNSIGNED, port INT UNSIGNED, service VARCHAR(50), content MEDIUMTEXT)")
	if err != nil {
		log.Fatal(err)
	}

	return &MySQLPersistenceManager{db: db}
}

func (mp *MySQLPersistenceManager) SaveBanner(b model.Banner) {
	_, err := mp.db.Exec("INSERT INTO banners (ip, port, service, content) VALUES (INET_ATON(?), ?, ?, ?)", b.Ip, b.Port, b.Service, b.Content)
	if err != nil {
		log.Fatal(err)
	}
}

func (mp *MySQLPersistenceManager) QueryBanners(f model.Filters) ([]model.Banner, error) {
	t := template.Must(template.New("query").Parse(mySQLQueryTemplate))
	query := &bytes.Buffer{}
	err := t.Execute(query, f)
	if err != nil {
		return nil, err
	}
	stmt, err := mp.db.Prepare(query.String())
	data := make([]interface{}, 0, len(f.Ports))
	if f.Ip != "" {
		data = append(data, f.Ip)
	}
	for _, v := range f.Ports {
		data = append(data, v)
	}
	for _, v := range f.Services {
		data = append(data, v)
	}
	if f.Regexp != "" {
		data = append(data, f.Regexp)
	}

	rows, err := stmt.Query(data...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Banner, 0)
	for rows.Next() {
		var curr model.Banner
		if err := rows.Scan(&curr.Ip, &curr.Port, &curr.Service, &curr.Content); err != nil {
			return nil, err
		}
		curr.Content = base64.StdEncoding.EncodeToString([]byte(curr.Content))
		result = append(result, curr)
	}
	return result, nil
}

func (mp *MySQLPersistenceManager) Close() {
	mp.db.Close()
}
