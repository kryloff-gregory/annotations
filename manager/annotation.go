package manager

import (
	"errors"
	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"main/annotation"
	"main/model"
	"main/services"
	"main/video"
)

type AnnotationManager struct {
	videoDAO      video.DAO
	annotationDAO annotation.DAO
	lockService   *services.LockService
}

func NewAnnotationManager(videoDAO video.DAO, annotationDAO annotation.DAO, lockService *services.LockService) *AnnotationManager {
	return &AnnotationManager{videoDAO: videoDAO, annotationDAO: annotationDAO, lockService: lockService}
}

func (m *AnnotationManager) CreateAnnotation(ann *model.Annotation) error {
	m.lockService.LockItem(ann.VideoURL)
	defer m.lockService.UnlockItem(ann.VideoURL)

	vid, err := m.videoDAO.FetchVideo(ann.VideoURL)
	if err != nil {
		return err
	}

	exists, err := m.annotationDAO.AnnotationExists(ann.VideoURL, ann.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("annotation already exists")
	}

	if err := m.validateAnnotationVsVideo(ann, vid); err != nil {
		return err
	}

	return m.annotationDAO.UpsertAnnotation(ann)
}

func (m *AnnotationManager) UpdateAnnotation(ann *model.Annotation) error {
	m.lockService.LockItem(ann.VideoURL)
	defer m.lockService.UnlockItem(ann.VideoURL)

	vid, err := m.videoDAO.FetchVideo(ann.VideoURL)
	if err != nil {
		return err
	}

	annStored, err := m.annotationDAO.FetchAnnotation(ann.VideoURL, ann.ID)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return err
	}

	if annStored == nil {
		return errors.New("annotation does not exist")
	}

	if annStored.UserCreated != ann.UserCreated {
		return errors.New("user is not authorized for the action")
	}

	if err := m.validateAnnotationVsVideo(ann, vid); err != nil {
		return err
	}

	return m.annotationDAO.UpsertAnnotation(ann)
}

func (m *AnnotationManager) DeleteAnnotation(videoURL string, annID uuid.UUID, userID uuid.UUID) error {
	m.lockService.LockItem(videoURL)
	defer m.lockService.UnlockItem(videoURL)

	annStored, err := m.annotationDAO.FetchAnnotation(videoURL, annID)
	if err != nil {
		return err
	}

	if annStored.UserCreated != userID {
		return errors.New("user is not authorized for this action")
	}

	return m.annotationDAO.RemoveAnnotation(videoURL, annID)
}

func (m *AnnotationManager) GetAllAnnotationsForVideo(videoURL string) ([]*model.Annotation, error) {
	m.lockService.RLockItem(videoURL)
	defer m.lockService.RUnlockItem(videoURL)

	return m.annotationDAO.FetchAnnotationsForVideo(videoURL)
}

func (m *AnnotationManager) validateAnnotationVsVideo(ann *model.Annotation, vid *model.Video) error {
	if !(ann.Start.Nanoseconds() > 0 && ann.Start.Nanoseconds() < ann.End.Nanoseconds() &&
		ann.End.Nanoseconds() < vid.Duration.Nanoseconds()) {
		return errors.New("invalid annotation boundaries")
	}
	return nil
}
