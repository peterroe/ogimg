package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"ogpic/internal/handler"
	"ogpic/internal/middleware"
	"ogpic/internal/repository"
	"ogpic/pkg/helper/resp"
	"ogpic/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

func NewServerHTTP(
	logger *log.Logger,
	userHandler *handler.UserHandler,
	repo *repository.Repository,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(
		middleware.CORSMiddleware(),
	)
	r.GET("/", func(ctx *gin.Context) {
		userUrl := ctx.Query("url")
		if userUrl == "" {
			resp.HandleError(ctx, http.StatusBadRequest, 1, "Url is required", nil)
			return
		}

		// Try to get image from cache
		if err := handleImageRequest(ctx, userUrl, repo, logger); err != nil {
			resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		}

	})
	r.GET("/user", userHandler.GetUserById)

	return r
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

// 封装处理图像请求的逻辑
func handleImageRequest(ctx *gin.Context, userUrl string, repo *repository.Repository, logger *log.Logger) error {
	// 检查缓存
	cachedBase64Val, err := repo.GetFromCache(ctx, userUrl)
	if err == nil && cachedBase64Val != "" {
		// 把 base64 的字符串转为图片返回
		return sendImageFromBase64(ctx, cachedBase64Val)
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
	return fetchAndCacheImage(ctx, ogImageUrl, userUrl, repo, logger)
}

// 从 base64 字符串发送图像
func sendImageFromBase64(ctx *gin.Context, base64Val string) error {
	decoded, err := base64.StdEncoding.DecodeString(base64Val)
	if err != nil {
		return fmt.Errorf("failed to decode base64: %v", err)
	}
	ctx.Data(http.StatusOK, "image/jpeg", decoded)
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

	// body 转为 base64 的字符串
	base64Val := base64.StdEncoding.EncodeToString(body)
	err = repo.SetToCache(ctx, userUrl, base64Val)
	if err != nil {
		logger.Error("Set cache error", zap.Error(err))
	}

	ctx.Data(http.StatusOK, imageResp.Header.Get("Content-Type"), body)
	return nil
}
