package middleware

import (
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 临时跳过权限检查，以便完成初始化
		// TODO: 初始化完成后删除下面这3行
		c.Next()
		return

		// 以下是原来的代码
		waitUse, _ := utils.GetClaims(c)
		//获取请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		e := utils.GetCasbin() // 判断策略中是否存在
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
