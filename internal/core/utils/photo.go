package utils

import (
	"context"

	"github.com/0azis/bot"
)

// GetTelegramPhoto return URL for user avatar
func GetTelegramPhoto(b *bot.Bot, ctx context.Context, userID int64) string {
	var url string

	// get list of photos
	photos, _ := b.GetUserProfilePhotos(ctx, &bot.GetUserProfilePhotosParams{
		UserID: userID,
	})

	if photos.TotalCount == 0 {
		url = "" // if user don't have an avatar
	} else {
		file, _ := b.GetFile(ctx, &bot.GetFileParams{
			FileID: photos.Photos[0][0].FileID, // take the first avatar
		})
		url = b.FileDownloadLink(file) // generate link
	}

	return url
}
