package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
)

func ProxyUploadImage(file *multipart.FileHeader, response *CheveretoUploadResponse) error {
	token := GetToken()

	fileContent, err := file.Open()
	if err != nil {
		return err
	}

	content, err := io.ReadAll(fileContent)
	if err != nil {
		return err
	}

	fileData := &fiber.FormFile{
		Fieldname: "file",
		Name:      file.Filename,
		Content:   content,
	}

	agent := LskyBaseAgent(fiber.AcquireAgent(), fiber.MethodPost, "/upload").
		Set("Authorization", "Bearer "+token).
		ContentType(fiber.MIMEMultipartForm).
		FileData(fileData).MultipartForm(nil)
	defer fiber.ReleaseAgent(agent)

	for {
		agent.Set("Authorization", "Bearer "+token)
		if err = agent.Parse(); err != nil {
			return err
		}

		code, body, errs := agent.Bytes()
		if len(errs) != 0 {
			return errs[0]
		}
		if code == 200 {
			var lskyUploadResponse LskyUploadResponse
			err = json.Unmarshal(body, &lskyUploadResponse)
			if err != nil {
				return err
			}

			// url transform to direct url
			urlRaw := lskyUploadResponse.Data.Links.Url
			urlData, err := url.ParseRequestURI(urlRaw)
			if err != nil {
				return err
			}
			if Config.HostRewrite != "" {
				urlData.Host = Config.HostRewrite
			} else {
				urlData.Host = ProxyUrlData.Host
			}
			urlData.Scheme = ProxyUrlData.Scheme
			directUrl := urlData.String()
			log.Printf("image upload: %v\n", directUrl)

			response.StatusCode = 200
			response.StatusTxt = "Upload Success"
			response.Image = CheveretoImageInfo{
				Name:       strings.TrimSuffix(lskyUploadResponse.Data.Name, filepath.Ext(lskyUploadResponse.Data.Name)),
				Extension:  lskyUploadResponse.Data.Extension,
				Md5:        lskyUploadResponse.Data.Md5,
				Filename:   lskyUploadResponse.Data.Name,
				Mime:       lskyUploadResponse.Data.Mimetype,
				Url:        directUrl,
				DisplayUrl: directUrl,
			}

			return nil
		} else if code == 401 {
			// refresh token
			newToken := GetToken() // maybe another coroutine refresh the token
			if token != newToken {
				// another coroutine refresh the token
				token = newToken
			} else {
				// this coroutine refresh the token
				token, err = LskyRefreshToken()
				if err != nil {
					return err
				}
			}
		} else {
			message := fmt.Sprintf(`{"code": %v}`, code)
			return fiber.NewError(fiber.StatusInternalServerError, message)
		}
	}
}
