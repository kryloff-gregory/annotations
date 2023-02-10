package video_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"main/model"
	"main/video"
	"os"
	"testing"
)

func TestDAO(t *testing.T) {
	conn, err := leveldb.OpenFile("testdb/videos", nil)
	assert.NoError(t, err)

	dao := video.NewDAO(conn)
	vid := &model.Video{
		URL:         "url",
		Duration:    300,
		UserCreated: uuid.New(),
	}

	assert.NoError(t, dao.UpsertVideo(vid))
	exists, err := dao.VideoExists("url")
	assert.NoError(t, err)
	assert.True(t, exists)

	videoGot, err := dao.FetchVideo("url")
	assert.NoError(t, err)
	assert.Equal(t, videoGot, vid)

	exists, err = dao.VideoExists("url1")
	assert.NoError(t, err)
	assert.False(t, exists)

	_, err = dao.FetchVideo("url1")
	assert.True(t, errors.Is(err, leveldb.ErrNotFound))

	conn.Close()
	os.RemoveAll("testdb")
}
