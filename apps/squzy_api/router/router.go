package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"net/http"
	"squzy/apps/squzy_api/handlers"
	"time"
)

var (
	errMissingConfig      = errors.New("missing config of scheduler")
	errNotFoundConfigType = errors.New("not found config type")
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
	Data interface{} `json:"data,omitempty"`
}

type HistoryFilterRequest struct {
	DateFrom *time.Time          `form:"dateFrom" time_format:"2006-01-02T15:04:05Z07:00"`
	DateTo   *time.Time          `form:"dateTo" time_format:"2006-01-02T15:04:05Z07:00"`
	Page     int32               `form:"page"`
	Limit    int32               `form:"limit"`
	Type     apiPb.TypeAgentStat `form:"type"`
}

type Scheduler struct {
	Type            apiPb.SchedulerType        `json:"type"`
	Interval        int32                      `json:"interval" binding:"required"`
	Timeout         int32                      `json:"timeout"`
	Name            string                     `json:"name"`
	HTTPConfig      *apiPb.HttpConfig          `json:"httpConfig"`
	TCPConfig       *apiPb.TcpConfig           `json:"tcpConfig"`
	HTTPValueConfig *apiPb.HttpJsonValueConfig `json:"httpValueConfig"`
	GRPCConfig      *apiPb.GrpcConfig          `json:"grpcConfig"`
	SiteMapConfig   *apiPb.SiteMapConfig       `json:"siteMapConfig"`
}

type Application struct {
	Host string `json:"host"`
	Name string `json:"name" binding:"required"`
}

type Transaction struct {
	Id       string                  `json:"id" binding:"required"`
	ParentID string                  `json:"parentId"`
	Name     string                  `json:"name" binding:"required"`
	DateFrom time.Time               `json:"dateFrom" time_format:"2006-01-02T15:04:05Z07:00" binding:"required"`
	DateTo   time.Time               `json:"dateTo" time_format:"2006-01-02T15:04:05Z07:00" binding:"required"`
	Status   apiPb.TransactionStatus `json:"status"`
	Type     apiPb.TransactionType   `json:"type"`
	Meta     *struct {
		Host   string `json:"host"`
		Path   string `json:"path"`
		Method string `json:"method"`
	} `json:"meta,omitempty"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
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
	{
		applications := v1.Group("applications")
		{
			applications.POST("", func(context *gin.Context) {
				application := &Application{}
				err := context.ShouldBindJSON(application)
				app, err := r.handlers.RegisterApplication(context, &apiPb.ApplicationInfo{
					Name:     application.Name,
					HostName: application.Host,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}

				successWrap(context, http.StatusOK, app)
			})
			application := applications.Group(":applicationId")
			{
				application.POST("transactions", func(context *gin.Context) {
					applicationId := context.Param("applicationId")
					trx := &Transaction{}
					err := context.ShouldBindJSON(trx)
					if err != nil {
						errWrap(context, http.StatusUnprocessableEntity, err)
						return
					}
					timeFrom, err := ptypes.TimestampProto(trx.DateFrom)
					if err != nil {
						errWrap(context, http.StatusUnprocessableEntity, err)
						return
					}
					timeTo, err := ptypes.TimestampProto(trx.DateTo)
					if err != nil {
						errWrap(context, http.StatusUnprocessableEntity, err)
						return
					}
					var meta *apiPb.TransactionInfo_Meta
					if trx.Meta != nil {
						meta = &apiPb.TransactionInfo_Meta{
							Host:   trx.Meta.Host,
							Path:   trx.Meta.Path,
							Method: trx.Meta.Method,
						}
					}
					var trxError *apiPb.TransactionInfo_Error
					if trx.Error != nil {
						trxError = &apiPb.TransactionInfo_Error{
							Message: trx.Error.Message,
						}
					}
					_, err = r.handlers.SaveTransaction(context, &apiPb.TransactionInfo{
						Id:            trx.Id,
						ApplicationId: applicationId,
						ParentId:      trx.ParentID,
						Meta:          meta,
						Name:          trx.Name,
						StartTime:     timeFrom,
						EndTime:       timeTo,
						Status:        trx.Status,
						Type:          trx.Type,
						Error:         trxError,
					})
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusAccepted, gin.H{
						"id": trx.Id,
					})
				})
			}
		}
		agents := v1.Group("agents")
		{
			agents.GET("", func(context *gin.Context) {
				list, err := r.handlers.GetAgentList(context)
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, list)
			})
			agent := agents.Group(":agentId")
			{
				agent.GET("", func(context *gin.Context) {
					agentID := context.Param("agentId")
					agent, err := r.handlers.GetAgentByID(context, agentID)
					if err != nil {
						errWrap(context, http.StatusNotFound, err)
						return
					}
					successWrap(context, http.StatusOK, agent)
				})
				//History
				agent.GET("/history", func(context *gin.Context) {
					agentID := context.Param("agentId")
					rq := HistoryFilterRequest{
						DateFrom: nil,
						DateTo:   nil,
						Page:     0,
						Limit:    0,
					}

					err := context.ShouldBind(&rq)
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					pagination, timeRange, err := GetFilters(&rq)
					if err != nil {
						errWrap(context, http.StatusUnprocessableEntity, err)
						return
					}
					res, err := r.handlers.GetAgentHistoryByID(context, &apiPb.GetAgentInformationRequest{
						AgentId:    agentID,
						Pagination: pagination,
						TimeRange:  timeRange,
						Type:       rq.Type,
					})

					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}

					successWrap(context, http.StatusOK, res)
				})
			}
		}

		schedulers := v1.Group("schedulers")
		{
			schedulers.GET("", func(context *gin.Context) {
				list, err := r.handlers.GetSchedulerList(context)
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, list)
			})
			schedulers.POST("", func(context *gin.Context) {
				request := new(Scheduler)
				err := context.ShouldBindJSON(request)
				if err != nil {
					errWrap(context, http.StatusUnprocessableEntity, err)
					return
				}
				var addReq *apiPb.AddRequest

				switch request.Type {
				case apiPb.SchedulerType_TCP:
					if request.TCPConfig == nil {
						errWrap(context, http.StatusUnprocessableEntity, errMissingConfig)
						return
					}
					addReq = &apiPb.AddRequest{
						Config: &apiPb.AddRequest_Tcp{
							Tcp: request.TCPConfig,
						},
					}

				case apiPb.SchedulerType_GRPC:
					if request.GRPCConfig == nil {
						errWrap(context, http.StatusUnprocessableEntity, errMissingConfig)
						return
					}
					addReq = &apiPb.AddRequest{
						Config: &apiPb.AddRequest_Grpc{
							Grpc: request.GRPCConfig,
						},
					}

				case apiPb.SchedulerType_HTTP:
					if request.HTTPConfig == nil {
						errWrap(context, http.StatusUnprocessableEntity, errMissingConfig)
						return
					}
					addReq = &apiPb.AddRequest{
						Config: &apiPb.AddRequest_Http{
							Http: request.HTTPConfig,
						},
					}

				case apiPb.SchedulerType_SITE_MAP:
					if request.SiteMapConfig == nil {
						errWrap(context, http.StatusUnprocessableEntity, errMissingConfig)
						return
					}
					addReq = &apiPb.AddRequest{
						Config: &apiPb.AddRequest_Sitemap{
							Sitemap: request.SiteMapConfig,
						},
					}

				case apiPb.SchedulerType_HTTP_JSON_VALUE:
					if request.HTTPValueConfig == nil {
						errWrap(context, http.StatusUnprocessableEntity, errMissingConfig)
						return
					}
					addReq = &apiPb.AddRequest{
						Config: &apiPb.AddRequest_HttpValue{
							HttpValue: request.HTTPValueConfig,
						},
					}

				default:
					errWrap(context, http.StatusUnprocessableEntity, errNotFoundConfigType)
					return
				}

				addReq.Interval = request.Interval
				addReq.Timeout = request.Timeout
				addReq.Name = request.Name
				res, err := r.handlers.AddScheduler(context, addReq)
				if err != nil {
					errWrap(context, http.StatusUnprocessableEntity, err)
					return
				}
				successWrap(context, http.StatusCreated, res)
			})
			scheduler := schedulers.Group(":schedulerId")
			{
				// Get by ID
				scheduler.GET("", func(context *gin.Context) {
					schedulerID := context.Param("schedulerId")
					scheduler, err := r.handlers.GetSchedulerByID(context, schedulerID)
					if err != nil {
						errWrap(context, http.StatusNotFound, err)
						return
					}
					successWrap(context, http.StatusOK, scheduler)
				})
				// Run by ID
				scheduler.PUT("run", func(context *gin.Context) {
					schedulerID := context.Param("schedulerId")
					err := r.handlers.RunScheduler(context, schedulerID)
					if err != nil {
						errWrap(context, http.StatusNotFound, err)
						return
					}
					successWrap(context, http.StatusAccepted, nil)
				})
				// Remove by ID
				scheduler.DELETE("", func(context *gin.Context) {
					schedulerID := context.Param("schedulerId")
					err := r.handlers.RemoveScheduler(context, schedulerID)
					if err != nil {
						errWrap(context, http.StatusNotFound, err)
						return
					}
					successWrap(context, http.StatusAccepted, nil)
				})
				// Stop by ID
				scheduler.PUT("stop", func(context *gin.Context) {
					schedulerID := context.Param("schedulerId")
					err := r.handlers.StopScheduler(context, schedulerID)
					if err != nil {
						errWrap(context, http.StatusNotFound, err)
						return
					}
					successWrap(context, http.StatusAccepted, nil)
				})

				//History
				scheduler.GET("/history", func(context *gin.Context) {
					schedulerID := context.Param("schedulerId")
					rq := HistoryFilterRequest{
						DateFrom: nil,
						DateTo:   nil,
						Page:     0,
						Limit:    0,
					}
					err := context.ShouldBind(&rq)

					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					pagination, timeRange, err := GetFilters(&rq)

					if err != nil {
						errWrap(context, http.StatusUnprocessableEntity, err)
						return
					}

					res, err := r.handlers.GetSchedulerHistoryByID(context, &apiPb.GetSchedulerInformationRequest{
						SchedulerId: schedulerID,
						Pagination:  pagination,
						TimeRange:   timeRange,
					})

					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusOK, res)
				})
			}
		}
	}

	return engine
}

func GetFilters(rq *HistoryFilterRequest) (*apiPb.Pagination, *apiPb.TimeFilter, error) {
	pagination := &apiPb.Pagination{
		Page:  rq.Page,
		Limit: rq.Limit,
	}

	if rq.Page == 0 && rq.Limit == 0 {
		pagination = nil
	}

	if rq.DateFrom != nil || rq.DateTo != nil {
		timeFilter := &apiPb.TimeFilter{}
		if rq.DateFrom != nil {
			t, err := ptypes.TimestampProto(*rq.DateFrom)
			if err != nil {
				return nil, nil, err
			}
			timeFilter.From = t
		}
		if rq.DateTo != nil {
			t, err := ptypes.TimestampProto(*rq.DateTo)
			if err != nil {
				return nil, nil, err
			}
			timeFilter.To = t
		}
		return pagination, timeFilter, nil
	}

	return pagination, nil, nil
}

func New(handlers handlers.Handlers) Router {
	return &router{
		handlers: handlers,
	}
}
