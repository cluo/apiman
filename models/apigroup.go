package models

import (
	//"fmt"
	//"errors"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type ApiGroup struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CreatorID   uint      `json:"creator" gorm:"default 0"`
	ProjectID   uint      `json:"project" gorm:"default 0"`
	//DeletedAt   *time.Time `json:"-"`
}

func (ag *ApiGroup) AfterSave() (err error) {
	d := new(ApiGroupIndex)
	copier.Copy(d, ag)
	d.SearchType = "api_group"
	d.ID = strconv.Itoa(int(ag.ID))
	err = BleveIndex.Index("user:"+d.ID, d)
	return
}

func (ag *ApiGroup) AfterDelete() (err error) {
	err = BleveIndex.Delete("api_group:" + strconv.Itoa(int(ag.ID)))
	return
}

func CreateApiGroup(apigroup *ApiGroup) error {
	err := db.Create(apigroup).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":       err.Error(),
			"apigroup": *apigroup,
		}).Error("create api group error")
		return err
	}

	log.WithFields(log.Fields{
		"apigroup": *apigroup,
	}).Info("create api group success")

	return nil
}

func GetApiGroupByID(id uint) (*ApiGroup, error) {
	apigroup := new(ApiGroup)
	err := db.First(apigroup, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("get api group error")
		return nil, err
	}

	return apigroup, nil
}

func UpdateApiGroup(apigroup *ApiGroup) error {
	err := db.Model(apigroup).Updates(apigroup).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":       err.Error(),
			"apigroup": *apigroup,
		}).Error("update api group error")
		return err
	}

	log.WithFields(log.Fields{
		"apigroup": *apigroup,
	}).Info("update api group success")

	return nil
}

func DeleteApiGroupByID(id uint) error {
	err := db.Where("id = ?", id).Delete(ApiGroup{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("delete api group error")
		return err
	}

	return nil
}

type Apis struct {
	Api
	Creator string `json:"creator"`
}

func (Apis) TableName() string {
	return "apis"
}

func GetApiGroupApis(ag_id uint) ([]*Apis, error) {
	//fmt.Println(ag_id)
	apis := make([]*Apis, 0)
	err := db.Where("group_id = ?", ag_id).Find(&apis).Error
	if err != nil {
		return nil, err
	}

	for _, api := range apis {
		u, _ := GetUserByID(api.CreatorID)
		api.Creator = u.Name
	}

	return apis, nil
}
