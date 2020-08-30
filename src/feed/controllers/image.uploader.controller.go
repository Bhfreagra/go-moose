package controllers

import (
	"go-moose/database/models"
	"go-moose/src/feed/services"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadImage saves an image in the filesystem
func UploadImage(ctx *gin.Context) {

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := filepath.Base(file.Filename)

	user, _ := ctx.Get("user_token")
	assertedUserToken := user.(models.UserToken)

	imageID, fullpath := services.UploadImage(filename, &assertedUserToken.User)

	if err := ctx.SaveUploadedFile(file, fullpath); err != nil {

		services.RemoveUploadedImageOnSacingError(imageID)

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"image_id": imageID})
}
