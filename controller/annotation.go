package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/helpers"
	"main/manager"
	"main/model"
	"time"
)

type AnnotationController struct {
	annotationManager *manager.AnnotationManager
}

func NewAnnotationController(annotationManager *manager.AnnotationManager) *AnnotationController {
	return &AnnotationController{annotationManager: annotationManager}
}

type annotationGetResponse struct {
	ID          string `json:"id" binding:"required"`
	VideoURL    string `json:"videoURL" binding:"required"`
	Start       string `json:"start" binding:"required"`
	End         string `json:"end" binding:"required"`
	UserCreated string `json:"userCreated" binding:"required"`
	Name        string `json:"name" binding:"optional"`
	Comment     string `json:"comment" binding:"optional"`
}

type annotationAddRequest struct {
	VideoURL string `json:"videoURL" binding:"required"`
	Start    string `json:"start" binding:"required"`
	End      string `json:"end" binding:"required"`
	Name     string `json:"name" binding:"optional"`
	Comment  string `json:"comment" binding:"optional"`
}

type annotationUpdateRequest struct {
	ID       string `json:"id" binding:"required"`
	VideoURL string `json:"videoURL" binding:"required"`
	Start    string `json:"start" binding:"optional"`
	End      string `json:"end" binding:"optional"`
	Name     string `json:"name" binding:"optional"`
	Comment  string `json:"comment" binding:"optional"`
}

type annotationAddResponse struct {
	ID string `json:"id" binding:"required"`
}

func (a *AnnotationController) GetAnnotationsForVideo(c *gin.Context) {
	videoURLStr := c.Request.URL.Query().Get("video_url")
	if !helpers.IsValidURL(videoURLStr) {
		responseWithError(c, 406, "Please provide valid video URL")
		return
	}

	anns, err := a.annotationManager.GetAllAnnotationsForVideo(videoURLStr)
	if err != nil {
		responseWithError(c, 500, "Could not get annotations: "+err.Error())
		return
	}

	result := make([]*annotationGetResponse, 0, len(anns))
	for _, ann := range anns {
		result = append(result, &annotationGetResponse{
			ID:          ann.ID.String(),
			VideoURL:    ann.VideoURL,
			Start:       ann.Start.String(),
			End:         ann.End.String(),
			UserCreated: ann.UserCreated.String(),
			Name:        ann.Name,
			Comment:     ann.Comment,
		})
	}

	c.JSON(200, result)
}

func (a *AnnotationController) CreateAnnotation(c *gin.Context) {
	var data annotationAddRequest

	if err := c.BindJSON(&data); err != nil {
		responseWithError(c, 406, "Please provide valid videoAddRequest metadata: "+err.Error())
		return
	}

	if !helpers.IsValidURL(data.VideoURL) {
		responseWithError(c, 406, "Invalid video URL")
		return
	}

	usr, ok := c.Get("User")
	if !ok {
		responseWithError(c, 500, "No user in the context")
	}

	start, err := time.ParseDuration(data.Start)
	if err != nil {
		responseWithError(c, 406, "Invalid start")
		return
	}

	end, err := time.ParseDuration(data.End)
	if err != nil {
		responseWithError(c, 406, "Invalid end")
		return
	}

	ann := &model.Annotation{
		ID:          uuid.New(),
		VideoURL:    data.VideoURL,
		Start:       start,
		End:         end,
		UserCreated: usr.(*model.User).ID,
		Name:        data.Name,
		Comment:     data.Comment,
	}

	if err := a.annotationManager.CreateAnnotation(ann); err != nil {
		responseWithError(c, 500, "Problem saving your annotation: "+err.Error())
		return
	}

	c.JSON(200, &annotationAddResponse{ID: ann.ID.String()})
}

func (a *AnnotationController) UpdateAnnotation(c *gin.Context) {
	var data annotationUpdateRequest

	if err := c.BindJSON(&data); err != nil {
		responseWithError(c, 406, "Please provide valid videoAddRequest metadata: "+err.Error())
		return
	}

	if !helpers.IsValidURL(data.VideoURL) {
		responseWithError(c, 406, "Invalid video URL")
		return
	}

	annID, err := uuid.Parse(data.ID)
	if err != nil {
		responseWithError(c, 406, "Please provide valid annotion ID")
	}

	usr, ok := c.Get("User")
	if !ok {
		responseWithError(c, 500, "No user in the context")
	}

	start, err := time.ParseDuration(data.Start)
	if err != nil {
		responseWithError(c, 406, "Invalid start")
		return
	}

	end, err := time.ParseDuration(data.End)
	if err != nil {
		responseWithError(c, 406, "Invalid end")
		return
	}

	ann := &model.Annotation{
		ID:          annID,
		VideoURL:    data.VideoURL,
		Start:       start,
		End:         end,
		UserCreated: usr.(*model.User).ID,
		Name:        data.Name,
		Comment:     data.Comment,
	}

	if err := a.annotationManager.UpdateAnnotation(ann); err != nil {
		responseWithError(c, 500, "Problem updating your annotation: "+err.Error())
		return
	}

	c.JSON(200, gin.H{"message": "Annotation updated"})
}

func (a *AnnotationController) DeleteAnnotation(c *gin.Context) {
	annIDStr := c.Request.URL.Query().Get("id")
	annID, err := uuid.Parse(annIDStr)
	if err != nil {
		responseWithError(c, 406, "Please provide valid annotion ID")
	}

	videoURLStr := c.Request.URL.Query().Get("video_url")
	if !helpers.IsValidURL(videoURLStr) {
		responseWithError(c, 406, "Please provide valid video URL")
		return
	}

	usr, ok := c.Get("User")
	if !ok {
		responseWithError(c, 500, "No user in the context")
	}

	err = a.annotationManager.DeleteAnnotation(videoURLStr, annID, usr.(*model.User).ID)
	if err != nil {
		responseWithError(c, 500, "Could not delete annotation: "+err.Error())
		return
	}

	c.JSON(200, gin.H{"message": "annotation deleted"})

}
