package utils

import (
	"errors"

	"crypto/rand"

	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FetchUserIDFromCtx(ctx *gin.Context) (uuid.UUID, error) {
	userIdInterface, _ := ctx.Get("userId")
	userID, ok := userIdInterface.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("failed to fetch user id from context")
	}
	return userID, nil
}

func CreateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // Handle the error properly in production code
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func Contains(slice []uuid.UUID, item uuid.UUID) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}
	return false
}

func ExtractDomainAndSubdomain(urlString string) (bool, string) {

	urlPattern := regexp.MustCompile(`(?i)\b((?:[a-z][\w-]+:(?:/{1,3}|[a-z0-9%])|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s` + "`" + `!()[\]{};:'".,<>?«»“”‘’]))`)

	if urlPattern.MatchString(urlString) {
		parsedURL, err := url.Parse(urlString)
		if err != nil {
			return false, ""
		}

		if parsedURL.Host == "" {
			// The regex matches, but URL parsing failed to identify a host. Try with "https://".
			parsedURL, err = url.Parse("https://" + urlString)
			if err != nil {
				return false, ""
			}
		}

		host := parsedURL.Hostname()
		return true, host
	}
	return false, ""
}
