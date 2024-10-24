package service

import (
	"github.com/gin-gonic/gin"

	"chainmaker_web/src/chain"
)

// GetChainConfigHandler get
type GetChainConfigHandler struct {
}

// Handle GetChainConfigHandler是否展示订阅按钮
func (getChainConfigHandler *GetChainConfigHandler) Handle(ctx *gin.Context) {
	ConvergeDataResponse(ctx, chain.GetConfigShow(), nil)
}
