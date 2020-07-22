package tag_service

import "github.com/miaozhang/webservice/models"

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}

	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}
	return tags, nil
}
