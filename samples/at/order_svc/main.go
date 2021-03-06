package main

import (
	"github.com/dk-lockdown/seata-golang/client"
	"github.com/dk-lockdown/seata-golang/client/at/sql/struct/cache"
	"github.com/dk-lockdown/seata-golang/client/context"
	"net/http"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/dk-lockdown/seata-golang/client/at/exec"
	"github.com/dk-lockdown/seata-golang/client/config"
	"github.com/dk-lockdown/seata-golang/samples/at/order_svc/dao"
)

const configPath="/Users/scottlewis/dksl/git/1/seata-golang/samples/at/order_svc/conf/client.yml"

func main() {
	r := gin.Default()
	config.InitConf(configPath)
	client.NewRpcClient()
	cache.SetTableMetaCache(cache.NewMysqlTableMetaCache(config.GetClientConfig().ATConfig.DSN))
	exec.InitDataResourceManager()

	db,err := exec.NewDB(config.GetClientConfig().ATConfig)
	if err != nil {
		panic(err)
	}
	d := &dao.Dao{
		DB: db,
	}

	r.POST("/createSo", func(c *gin.Context) {
		type req struct {
			Req []*dao.SoMaster
		}
		var q req
		if err := c.ShouldBindJSON(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rootContext := &context.RootContext{Context:c}
		rootContext.Bind(c.Request.Header.Get("Xid"))

		d.CreateSO(rootContext,q.Req)

		c.JSON(200, gin.H{
			"success": true,
			"message": "success",
		})
	})
	r.Run(":8002")
}
