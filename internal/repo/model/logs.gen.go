// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameLog = "logs"

// Log mapped from table <logs>
type Log struct {
	ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	BuildID   int32     `gorm:"column:build_id;not null" json:"buildId"`
	Content   string    `gorm:"column:content;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(3)" json:"createdAt"`
}

// TableName Log's table name
func (*Log) TableName() string {
	return TableNameLog
}
