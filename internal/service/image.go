package service

import (
	"fmt"
	"io"
	"net/http"
	"ogimg/internal/model"
	"ogimg/internal/repository"
	"ogimg/pkg/log"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

type ImageService interface {
	GetOgImageByUrl(ctx *gin.Context, userUrl string) error
	GetOgDescByUrl(ctx *gin.Context, userUrl string) error
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
	imageBytes, err := s.repository.GetWebsiteOgImgFromCache(ctx, userUrl)
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

func (s *imageService) GetOgDescByUrl(ctx *gin.Context, userUrl string) error {
	// 检查缓存
	descFromCache, err := s.repository.GetWebSiteDescToCache(ctx, userUrl)
	fmt.Println("descFromCache", descFromCache)
	if err == nil && descFromCache != (model.WebSiteDescType{}) {
		// 把 bytes 转为图片返回
		ctx.JSON(http.StatusOK, descFromCache)
		return nil
	}

	urlResp, err := http.Get(userUrl)
	if err != nil {
		return err
	}
	defer urlResp.Body.Close()

	doc, err := html.Parse(urlResp.Body)
	if err != nil {
		return err
	}

	desc := model.WebSiteDescType{}
	findWebSiteDesc(doc, &desc)

	// 如果 logo 以 / 开头，则认为是相对路径，需要拼接上域名
	if strings.HasPrefix(desc.Logo, "/") {
		// 确保 userUrl 和 desc.Logo 之间只有一个 /
		// 去除 userUrl 末尾的 /
		userUrl = strings.TrimRight(userUrl, "/")
		desc.Logo = fmt.Sprintf("%s%s", userUrl, desc.Logo)
	}

	s.service.logger.Info("desc", zap.Any("desc", desc))

	err = s.repository.SetWebSiteDescToCache(ctx, userUrl, desc)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, desc)

	return nil
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
	err = repo.SetWebsiteOgImgToCache(ctx, userUrl, body)
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

func findWebSiteDesc(n *html.Node, desc *model.WebSiteDescType) {
	// 找到 head 标签并遍历里面的内容
	var headNode *html.Node
	var findHead func(*html.Node) *html.Node
	findHead = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "head" {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if result := findHead(c); result != nil {
				return result
			}
		}
		return nil
	}
	headNode = findHead(n)
	if headNode != nil {
		for c := headNode.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				switch c.Data {
				case "meta":
					var content string
					var isDescription bool
					for _, attr := range c.Attr {
						if attr.Key == "name" && attr.Val == "description" {
							isDescription = true
						}
						if attr.Key == "content" {
							content = attr.Val
						}
					}
					if isDescription {
						desc.Description = content
					}
				case "link":
					var content string
					var isIcon bool
					for _, attr := range c.Attr {
						if attr.Key == "rel" && attr.Val == "icon" {
							isIcon = true
						}
						if attr.Key == "href" {
							content = attr.Val
						}
					}
					if isIcon && content != "" {
						desc.Logo = content
					}
				case "title":
					if c.FirstChild != nil {
						desc.Title = c.FirstChild.Data
					}
				}
			}
		}
	}
}
