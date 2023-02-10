package annotation

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"main/model"
)

type DAO interface {
	UpsertAnnotation(annotation *model.Annotation) error
	AnnotationExists(videoURL string, id uuid.UUID) (bool, error)
	FetchAnnotation(videoURL string, id uuid.UUID) (*model.Annotation, error)
	FetchAnnotationsForVideo(videoURL string) ([]*model.Annotation, error)
	RemoveAnnotation(videoURL string, id uuid.UUID) error
}

type dao struct {
	db *leveldb.DB
}

func NewDAO(db *leveldb.DB) *dao {
	return &dao{db: db}
}

func (dao *dao) UpsertAnnotation(annotation *model.Annotation) error {
	valueBytes, err := json.Marshal(annotation)
	if err != nil {
		return err
	}
	return dao.db.Put([]byte(dao.getKey(annotation.VideoURL, annotation.ID)), valueBytes, nil)
}

func (m *dao) AnnotationExists(videoURL string, id uuid.UUID) (bool, error) {
	_, err := m.db.Get([]byte(m.getKey(videoURL, id)), nil)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, leveldb.ErrNotFound) {
		return false, nil
	}
	return false, err
}

func (m *dao) FetchAnnotation(videoURL string, id uuid.UUID) (*model.Annotation, error) {
	ann := &model.Annotation{}
	bytes, err := m.db.Get([]byte(m.getKey(videoURL, id)), nil)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, ann); err != nil {
		return nil, err
	}

	return ann, nil

}

func (dao *dao) FetchAnnotationsForVideo(videoURL string) ([]*model.Annotation, error) {
	result := make([]*model.Annotation, 0)
	iter := dao.db.NewIterator(util.BytesPrefix([]byte(videoURL+":")), nil)
	defer iter.Release()

	for iter.Next() {
		annotation := &model.Annotation{}
		if err := json.Unmarshal(iter.Value(), annotation); err != nil {
			return nil, err
		}
		result = append(result, annotation)
	}
	return result, iter.Error()
}

func (dao *dao) RemoveAnnotation(videoURL string, id uuid.UUID) error {
	return dao.db.Delete([]byte(dao.getKey(videoURL, id)), nil)
}

func (dao *dao) getKey(videoURL string, id uuid.UUID) string {
	return fmt.Sprintf("%v:%v", videoURL, id.String())
}
