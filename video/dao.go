package video

import (
	"encoding/json"
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
	"main/model"
)

type DAO interface {
	FetchVideo(videoURL string) (*model.Video, error)
	VideoExists(videoURL string) (bool, error)
	RemoveVideo(videoURL string) error
	UpsertVideo(video *model.Video) error
}

type dao struct {
	db *leveldb.DB
}

func NewDAO(db *leveldb.DB) *dao {
	return &dao{db: db}
}

func (m *dao) FetchVideo(videoURL string) (*model.Video, error) {
	bytes, err := m.db.Get([]byte(videoURL), nil)
	if err != nil {
		return nil, err
	}
	video := &model.Video{}
	if err := json.Unmarshal(bytes, video); err != nil {
		return nil, err
	}

	return video, nil
}

func (m *dao) VideoExists(videoURL string) (bool, error) {
	_, err := m.db.Get([]byte(videoURL), nil)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, leveldb.ErrNotFound) {
		return false, nil
	}
	return false, err
}

func (m *dao) RemoveVideo(videoURL string) error {
	return m.db.Delete([]byte(videoURL), nil)
}

func (m *dao) UpsertVideo(video *model.Video) error {
	valueBytes, err := json.Marshal(video)
	if err != nil {
		return err
	}
	return m.db.Put([]byte(video.URL), valueBytes, nil)
}
