// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameBuildRecord = "build_records"

// BuildRecord mapped from table <build_records>
type BuildRecord struct {
	ID        int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"Id"`
	Pid       *int32     `gorm:"column:pid" json:"Pid"`
	Status    *int32     `gorm:"column:status" json:"Status"`
	CreatedAt *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"CreatedAt"`
}

// TableName BuildRecord's table name
func (*BuildRecord) TableName() string {
	return TableNameBuildRecord
}
