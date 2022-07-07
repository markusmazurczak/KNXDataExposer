package handler

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DataHandler struct {
	DB *gorm.DB
}

type Dataset struct {
	gorm.Model
	Group_Address string `gorm:"unique;not null"`
	Value         string `gorm:"not null"`
	Unit          string
}

//Updates or creates a new row in the dataset table.
//	- ga: The group address
//	- val: The value to save
//	- dh: DataHandler to communicate with the database
func Insert_dataset(ga string, val string, dh *DataHandler) error {

	s := strings.Split(val, " ")
	dataset := Dataset{
		Group_Address: ga,
		Value:         s[0],
		Unit:          s[1],
	}

	dh.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "group_address"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "unit"}),
	}).Create(&dataset)

	return nil
}

//Looks up a DataSet in the database
//	- ga: The group address to search a value for
//	- dh: The DataHandler for
func Get_dataset(ga string, dh *DataHandler) (Dataset, error) {
	var r Dataset

	if err := dh.DB.Where("group_address = ?", ga).First(&r).Error; err != nil {
		return r, err
	}
	return r, nil
}
