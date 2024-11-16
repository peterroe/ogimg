package server

import (
	"io"
	"net/http"
	"ogpic/internal/handler"
	"ogpic/internal/middleware"
	"ogpic/pkg/helper/resp"
	"ogpic/pkg/log"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

func NewServerHTTP(
	logger *log.Logger,
	userHandler *handler.UserHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(
		middleware.CORSMiddleware(),
	)
	r.GET("/", func(ctx *gin.Context) {
		userUrl := ctx.Query("url")

		// Get HTML Content
		urlResp, err := http.Get(userUrl)
		if err != nil {
			resp.HandleError(ctx, http.StatusBadRequest, 1, "Get html of url error", nil)
			return
		}
		defer urlResp.Body.Close()

		doc, err := html.Parse(urlResp.Body)
		if err != nil {
			resp.HandleError(ctx, http.StatusBadRequest, 1, "Parse html error", nil)
			return
		}

		ogImageUrl := findOGImage(doc)

		if ogImageUrl == "" {
			resp.HandleError(ctx, http.StatusBadRequest, 1, "No og:image found", nil)
			return
		}
		// 重定向到图片的地址
		// ctx.Redirect(http.StatusTemporaryRedirect, ogImageUrl)
		// 或者：返回 ogImageUrl 图片的资源
		imageResp, err := http.Get(ogImageUrl)
		if err != nil {
			resp.HandleError(ctx, http.StatusBadRequest, 1, "Get image error", nil)
			return
		}
		defer imageResp.Body.Close()
		body, err := io.ReadAll(imageResp.Body)
		if err != nil {
			resp.HandleError(ctx, http.StatusInternalServerError, 1, "Failed to read image", nil)
			return
		}
		ctx.Data(http.StatusOK, imageResp.Header.Get("Content-Type"), body)
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
