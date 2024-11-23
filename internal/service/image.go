package service

import (
	"fmt"
	"io"
	"net/http"
	"ogimg/internal/repository"
	"ogimg/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

type ImageService interface {
	GetOgImageByUrl(ctx *gin.Context, userUrl string) error
}

type imageService struct {
	service    *Service
	repository *repository.Repository
}

func NewImageService(service *Service, repository *repository.Repository) ImageService {
	return &imageService{
		service:    service,
		repository: repository,
	}
}

func (s *imageService) GetOgImageByUrl(ctx *gin.Context, userUrl string) error {
	// 检查缓存
	imageBytes, err := s.repository.GetWebsiteDescFromCache(ctx, userUrl)
	if err == nil && imageBytes != nil {
		// 把 bytes 转为图片返回
		contentType := http.DetectContentType(imageBytes)
		ctx.Data(http.StatusOK, contentType, imageBytes)
		return nil
	}

	// 获取 HTML 内容
	urlResp, err := http.Get(userUrl)
	if err != nil {
		return err
	}
	defer urlResp.Body.Close()

	doc, err := html.Parse(urlResp.Body)
	if err != nil {
		return err
	}

	ogImageUrl := findOGImage(doc)
	if ogImageUrl == "" {
		return fmt.Errorf("no og:image found")
	}

	// 获取图像
	return fetchAndCacheImage(ctx, ogImageUrl, userUrl, s.repository, s.service.logger)
}

// 获取图像并缓存
func fetchAndCacheImage(ctx *gin.Context, ogImageUrl, userUrl string, repo *repository.Repository, logger *log.Logger) error {
	imageResp, err := http.Get(ogImageUrl)
	if err != nil {
		return err
	}
	defer imageResp.Body.Close()

	body, err := io.ReadAll(imageResp.Body)
	if err != nil {
		return err
	}

	// 缓存 bytes
	err = repo.SetWebsiteDescToCache(ctx, userUrl, body)
	if err != nil {
		logger.Error("Set cache error", zap.Error(err))
	}

	ctx.Data(http.StatusOK, imageResp.Header.Get("Content-Type"), body)
	return nil
}

func findOGImage(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var property, content string
		for _, attr := range n.Attr {
			if attr.Key == "property" {
				property = attr.Val
			}
			if attr.Key == "content" {
				content = attr.Val
			}
			if property == "og:image" && content != "" {
				return content
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := findOGImage(c)
		if result != "" {
			return result
		}
	}
	return ""
}
