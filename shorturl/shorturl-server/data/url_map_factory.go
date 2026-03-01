package data

import (
	"database/sql"
	"shorturl/pkg/constants"
	"shorturl/pkg/log"
)

type IUrlMapDataFactory interface {
	NewUrlMapData(isPublic bool) IUrlMapData
}
type urlMapDataFactory struct {
	log log.ILogger
	db  *sql.DB
}

func NewUrlMapDataFactory(log log.ILogger, db *sql.DB) IUrlMapDataFactory {
	return &urlMapDataFactory{
		log: log,
		db:  db,
	}
}
func (f *urlMapDataFactory) NewUrlMapData(isPublic bool) IUrlMapData {
	tableName := constants.TABLENAME_URL_MAP
	if !isPublic {
		tableName = constants.TABLENAME_URL_MAP_USER
	}
	return newUrlMapData(f.log, f.db, tableName)
}
