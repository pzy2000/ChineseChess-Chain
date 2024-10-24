package router

import (
	"chainmaker_web/src/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCors(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "测试跨域中间件",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.Use(Cors())
			router.GET("/test-cors", func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test-cors", nil)
			req.Header.Set("Origin", "http://example.com")
			router.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("Cors() failed, status code = %v, want %v", w.Code, http.StatusOK)
			}
		})
	}
}

func TestInitRouter(t *testing.T) {
	type args struct {
		router  *gin.Engine
		webConf *config.WebConf
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "测试初始化路由",
			args: args{
				router: gin.Default(),
				webConf: &config.WebConf{
					CrossDomain: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitRouter(tt.args.router, tt.args.webConf)
			// 添加断言以验证路由是否已正确初始化
		})
	}
}

func Test_initControllers(t *testing.T) {
	type args struct {
		routeGroup *gin.RouterGroup
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "测试初始化控制器",
			args: args{
				routeGroup: gin.Default().Group("/"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initControllers(tt.args.routeGroup)
			// 添加断言以验证控制器是否已正确初始化
		})
	}
}
