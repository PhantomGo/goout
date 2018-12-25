package http

import (
	"fmt"
	"net/http/httputil"
	"os"
	"runtime"
	"time"

	"goout/conf"
	"goout/errors"
	"goout/service"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"
)

const (
	contextErrCode = "context/err/code"
)

var (
	svc *service.Service
)

// Init init http
func Init(c *conf.Config) {
	svc = service.New(c)
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(loggerHandler, recoverHandler)
	innerRouter(engine)
	go engine.Run(c.HTTPServer.Addr)
}

// innerRouter init local router api path.
func innerRouter(e *gin.Engine) {
	group := e.Group("/goout")
	group.POST("/users", createUser)
	group.POST("/home", myHome)
	group.GET("/tasks/:index", taskInfo)
	group.POST("/tasks/join/:index", joinTask)
	group.POST("/task/submit", submit)
	group.POST("/tasks/import", fillTasks)
	group.POST("/coin/consume", consumeCoin)
	group.GET("/invite/join/:code", joinInvite)
	group.POST("/invite", invite)
	group.POST("/dailybonus", dailybonus)
	group.GET("/ranking", ranking)
	group.GET("/ranklist", rankList)
	group.GET("/ping", pong)
}

func loggerHandler(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	method := c.Request.Method

	// Process request
	c.Next()

	// Stop timer
	end := time.Now()
	latency := end.Sub(start)
	statusCode := c.Writer.Status()
	ecode := c.GetInt(contextErrCode)
	clientIP := c.ClientIP()
	if raw != "" {
		path = path + "?" + raw
	}
	log.Infof("METHOD:%s | PATH:%s | CODE:%d | IP:%s | TIME:%d | ECODE:%d", method, path, statusCode, clientIP, latency/time.Millisecond, ecode)
}

func recoverHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			httprequest, _ := httputil.DumpRequest(c.Request, false)
			pnc := fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s", time.Now().Format("2018-11-19 22:04:18"), string(httprequest), err, buf)
			fmt.Fprintf(os.Stderr, pnc)
			log.Error(pnc)
			c.AbortWithStatus(500)
		}
	}()
	c.Next()
}

type resp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func result(c *gin.Context, data interface{}, err error) {
	ee := errors.Code(err)
	c.Set(contextErrCode, ee.Code())
	c.JSON(200, resp{
		Code: ee.Code(),
		Data: data,
	})
}
