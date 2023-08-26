// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameProject = "projects"

// Project mapped from table <projects>
type Project struct {
	ID         int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"Id"`
	Name       *string    `gorm:"column:name" json:"Name"`
	Status     *int32     `gorm:"column:status" json:"Status"`
	Branch     *string    `gorm:"column:branch" json:"Branch"`
	ProjectURL *string    `gorm:"column:project_url" json:"ProjectUrl"`
	WorkDir    *string    `gorm:"column:work_dir" json:"WorkDir"`
	PrivateKey *string    `gorm:"column:private_key" json:"PrivateKey"`
	Bin        *string    `gorm:"column:bin" json:"Bin"`
	Parameters *string    `gorm:"column:parameters" json:"Parameters"`
	CreatedAt  *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"CreatedAt"`
}

// TableName Project's table name
func (*Project) TableName() string {
	return TableNameProject
}
