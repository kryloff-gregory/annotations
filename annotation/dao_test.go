package annotation_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"main/annotation"
	"main/model"
	"os"
	"testing"
)

func TestDAO(t *testing.T) {
	conn, err := leveldb.OpenFile("testdb/annotations", nil)
	if err != nil {
		log.Fatal(err)
	}

	dao := annotation.NewDAO(conn)
	ann := &model.Annotation{
		ID:          uuid.New(),
		VideoURL:    "myurl",
		Start:       0,
		End:         10,
		UserCreated: uuid.New(),
		Name:        "name",
		Comment:     "comment",
	}

	assert.NoError(t, dao.UpsertAnnotation(ann))
	exists, err := dao.AnnotationExists("myurl", ann.ID)
	assert.NoError(t, err)
	assert.True(t, exists)

	ann1 := &model.Annotation{
		ID:          uuid.New(),
		VideoURL:    "myurl",
		Start:       0,
		End:         11,
		UserCreated: uuid.New(),
		Name:        "name1",
		Comment:     "comment1",
	}
	assert.NoError(t, dao.UpsertAnnotation(ann1))

	annGot, err := dao.FetchAnnotation("myurl", ann.ID)
	assert.NoError(t, err)
	assert.Equal(t, ann, annGot)

	ann2 := &model.Annotation{
		ID:          uuid.New(),
		VideoURL:    "myurl1",
		Start:       0,
		End:         10,
		UserCreated: uuid.New(),
		Name:        "name",
		Comment:     "comment",
	}
	assert.NoError(t, dao.UpsertAnnotation(ann2))

	gotAnns, err := dao.FetchAnnotationsForVideo("myurl")
	assert.NoError(t, err)

	assert.ElementsMatch(t, []*model.Annotation{ann, ann1}, gotAnns)

	err = dao.RemoveAnnotation("myurl", ann.ID)
	assert.NoError(t, err)

	exists, err = dao.AnnotationExists("myurl", ann.ID)
	assert.NoError(t, err)
	assert.False(t, exists)

	annGot, err = dao.FetchAnnotation("myurl", ann.ID)
	assert.True(t, errors.Is(err, leveldb.ErrNotFound))

	conn.Close()
	os.RemoveAll("testdb")
}
