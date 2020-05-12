package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"squzy/apps/squzy_api/handlers"
)

type Router interface {
	GetEngine() *gin.Engine
}

type router struct {
	handlers handlers.Handlers
}

type E struct {
	Error error `json:"error"`
}

type D struct {
	Data interface{} `json:"data"`
}

func errWrap(c *gin.Context, status int, err error) {
	c.AbortWithStatusJSON(status, E{
		Error: err,
	})
}

func successWrap(c *gin.Context, status int, data interface{}) {
	c.JSON(status, D{
		Data: data,
	})
}

func (r *router) GetEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	v1 := engine.Group("v1")
	agents := v1.Group("agents")
	agents.GET("", func(context *gin.Context) {
		list, err := r.handlers.GetAgentList(context)
		if err != nil {
			errWrap(context, http.StatusInternalServerError, err)
			return
		}
		successWrap(context, http.StatusOK, list)
	})
	agents.GET(":agentId", func(context *gin.Context) {
		agentId := context.Param("agentId")
		agent, err := r.handlers.GetAgentById(context, agentId)
		if err != nil {
			errWrap(context, http.StatusNotFound, err)
			return
		}
		successWrap(context, http.StatusOK, agent)
	})
	return engine
}

func New(handlers handlers.Handlers) Router {
	return &router{
		handlers,
	}
}
