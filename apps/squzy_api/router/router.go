package router

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/squzy/squzy/apps/squzy_api/handlers"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"net/http"
	"strconv"
	"time"
)

var (
	errMissingConfig         = errors.New("missing config of scheduler")
	errNotFoundConfigType    = errors.New("not found config type")
	errWrongNotificationType = errors.New("wrong notification type")
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

type SchedulerUptimeRequest struct {
	TimeRange *TimeFilterRequest
}

type SchedulerHistory struct {
	Pagination    *PaginationRequest
	TimeFilters   *TimeFilterRequest
	Status        apiPb.SchedulerCode     `form:"status"`
	SortDirection apiPb.SortDirection     `form:"sort_direction"`
	SortBy        apiPb.SortSchedulerList `form:"sort_by"`
}

type AgentHistory struct {
	Pagination  *PaginationRequest
	TimeFilters *TimeFilterRequest
	Type        apiPb.TypeAgentStat `form:"type"`
}

type GetIncidentListRequest struct {
	Pagination    *PaginationRequest
	TimeFilters   *TimeFilterRequest
	Status        apiPb.IncidentStatus   `form:"status"`
	RuleId        string                 `form:"ruleId"`
	SortBy        apiPb.SortIncidentList `form:"sort_by"`
	SortDirection apiPb.SortDirection    `form:"sort_direction"`
}

type GetTransactionListRequest struct {
	Pagination        *PaginationRequest
	TimeFilters       *TimeFilterRequest
	SortBy            apiPb.SortTransactionList `form:"sort_by"`
	SortDirection     apiPb.SortDirection       `form:"sort_direction"`
	TransactionType   apiPb.TransactionType     `form:"transaction_type"`
	TransactionStatus apiPb.TransactionStatus   `form:"transaction_status"`
	HostFilter        string                    `form:"host"`
	NameFilter        string                    `form:"name"`
	PathFilter        string                    `form:"path"`
	MethodFilter      string                    `form:"method"`
}

type GetTransactionGroupRequest struct {
	TimeFilters       *TimeFilterRequest
	GroupType         apiPb.GroupTransaction  `form:"group_by"`
	TransactionType   apiPb.TransactionType   `form:"transaction_type"`
	TransactionStatus apiPb.TransactionStatus `form:"transaction_status"`
}

type TimeFilterRequest struct {
	DateFrom *time.Time `form:"dateFrom" time_format:"2006-01-02T15:04:05Z07:00"`
	DateTo   *time.Time `form:"dateTo" time_format:"2006-01-02T15:04:05Z07:00"`
}

type PaginationRequest struct {
	Page  int32 `form:"page"`
	Limit int32 `form:"limit"`
}

type ValidateRuleRequest struct {
	OwnerType apiPb.ComponentOwnerType `json:"ownerType"`
	Rule      string                   `json:"rule" binding:"required"`
}

type CreateRuleRequest struct {
	Rule      string                   `json:"rule" binding:"required"`
	Name      string                   `json:"name"`
	AutoClose bool                     `json:"autoClose"`
	OwnerType apiPb.ComponentOwnerType `json:"ownerType"`
	OwnerId   string                   `json:"ownerId" binding:"required"`
}

type CreateNotificationMethod struct {
	Type          apiPb.NotificationMethodType `json:"type"`
	Name          string                       `json:"name" binding:"required"`
	SlackConfig   *apiPb.SlackMethod           `json:"slackConfig,omitempty"`
	WebHookConfig *apiPb.WebHookMethod         `json:"webhookConfig,omitempty"`
}

type RuleIdRequest struct {
	RuleId string `json:"ruleId"`
}

type ListOfNotificationMethods struct {
	OwnerId   string                   `form:"ownerId"`
	OwnerType apiPb.ComponentOwnerType `form:"ownerType"`
}

type LinkNotificationMethod struct {
	OwnerId   string                   `json:"ownerId" binding:"required"`
	OwnerType apiPb.ComponentOwnerType `json:"ownerType" binding:"required"`
}

type ListRulesByOwnerIdRequest struct {
	OwnerType apiPb.ComponentOwnerType `form:"ownerType"`
	OwnerId   string                   `form:"ownerId"  binding:"required"`
}

type Scheduler struct {
	Type                apiPb.SchedulerType        `json:"type"`
	Interval            int32                      `json:"interval" binding:"required"`
	Timeout             int32                      `json:"timeout"`
	Name                string                     `json:"name"`
	HTTPConfig          *apiPb.HttpConfig          `json:"httpConfig,omitempty"`
	TCPConfig           *apiPb.TcpConfig           `json:"tcpConfig,omitempty"`
	HTTPValueConfig     *apiPb.HttpJsonValueConfig `json:"httpValueConfig,omitempty"`
	GRPCConfig          *apiPb.GrpcConfig          `json:"grpcConfig,omitempty"`
	SiteMapConfig       *apiPb.SiteMapConfig       `json:"siteMapConfig,omitempty"`
	SSLExpirationConfig *apiPb.SslExpirationConfig `json:"sslExpirationConfig,omitempty"`
}

type Application struct {
	Host    string `json:"host"`
	Name    string `json:"name" binding:"required"`
	AgentId string `json:"agentId"`
}

type transactionTime timestamp.Timestamp

var _ json.Unmarshaler = &transactionTime{}

func (mt *transactionTime) UnmarshalJSON(bs []byte) error {
	var stringTime string
	err := json.Unmarshal(bs, &stringTime)
	if err != nil {
		return err
	}

	intTime, err := strconv.ParseInt(stringTime, 10, 64)
	if err != nil {
		return err
	}
	*mt = transactionTime{
		Seconds: intTime / 1e9,
		Nanos:   int32(intTime % 1e9),
	}
	return nil
}

func (mt *transactionTime) ToTimeStamp() *timestamp.Timestamp {
	t := timestamp.Timestamp(*mt)
	return &t
}

type Transaction struct {
	Id       string                  `json:"id" binding:"required"`
	ParentID string                  `json:"parentId"`
	Name     string                  `json:"name" binding:"required"`
	DateFrom transactionTime         `json:"dateFrom" time_format:"unixNano" binding:"required"`
	DateTo   transactionTime         `json:"dateTo" time_format:"unixNano" binding:"required"`
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
		notifications := v1.Group("notifications")
		{
			notifications.POST("", func(context *gin.Context) {
				createNotificationReq := &CreateNotificationMethod{}

				err := context.ShouldBindJSON(createNotificationReq)
				if err != nil {
					errWrap(context, http.StatusUnprocessableEntity, err)
					return
				}
				var req *apiPb.CreateNotificationMethodRequest
				switch createNotificationReq.Type {
				case apiPb.NotificationMethodType_NOTIFICATION_METHOD_SLACK:
					req = &apiPb.CreateNotificationMethodRequest{
						Type: apiPb.NotificationMethodType_NOTIFICATION_METHOD_SLACK,
						Name: createNotificationReq.Name,
						Method: &apiPb.CreateNotificationMethodRequest_Slack{
							Slack: createNotificationReq.SlackConfig,
						},
					}
				case apiPb.NotificationMethodType_NOTIFICATION_METHOD_WEBHOOK:
					req = &apiPb.CreateNotificationMethodRequest{
						Type: apiPb.NotificationMethodType_NOTIFICATION_METHOD_WEBHOOK,
						Name: createNotificationReq.Name,
						Method: &apiPb.CreateNotificationMethodRequest_Webhook{
							Webhook: createNotificationReq.WebHookConfig,
						},
					}
				default:
					errWrap(context, http.StatusBadRequest, errWrongNotificationType)
				}
				method, err := r.handlers.CreateNotificationMethod(context, req)
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusCreated, method)
			})
			notifications.GET("", func(context *gin.Context) {
				rq := &ListOfNotificationMethods{}
				err := context.ShouldBind(rq)

				if err != nil {
					errWrap(context, http.StatusBadRequest, err)
					return
				}

				list, err := r.handlers.GetNotificationMethods(context, GetNotificationList(rq.OwnerType, rq.OwnerId))

				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}

				successWrap(context, http.StatusOK, list)
			})
			notification := notifications.Group(":method_id")
			notification.GET("", func(context *gin.Context) {
				methodId := context.Param("method_id")
				method, err := r.handlers.GetMethodById(context, &apiPb.NotificationMethodIdRequest{
					Id: methodId,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, method)
			})
			notification.PUT("activate", func(context *gin.Context) {
				methodId := context.Param("method_id")
				method, err := r.handlers.ActivateById(context, &apiPb.NotificationMethodIdRequest{
					Id: methodId,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, method)
			})
			notification.POST("link", func(context *gin.Context) {
				methodId := context.Param("method_id")
				linkReq := &LinkNotificationMethod{}
				err := context.ShouldBindJSON(linkReq)
				if err != nil {
					errWrap(context, http.StatusUnprocessableEntity, err)
					return
				}
				method, err := r.handlers.LinkById(context, &apiPb.NotificationMethodRequest{
					OwnerType:            linkReq.OwnerType,
					OwnerId:              linkReq.OwnerId,
					NotificationMethodId: methodId,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, method)
			})
			notification.POST("unlink", func(context *gin.Context) {
				methodId := context.Param("method_id")
				linkReq := &LinkNotificationMethod{}
				err := context.ShouldBindJSON(linkReq)
				if err != nil {
					errWrap(context, http.StatusUnprocessableEntity, err)
					return
				}
				method, err := r.handlers.UnLinkById(context, &apiPb.NotificationMethodRequest{
					OwnerType:            linkReq.OwnerType,
					OwnerId:              linkReq.OwnerId,
					NotificationMethodId: methodId,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, method)
			})
			notification.PUT("deactivate", func(context *gin.Context) {
				methodId := context.Param("method_id")
				method, err := r.handlers.DeactivateById(context, &apiPb.NotificationMethodIdRequest{
					Id: methodId,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, method)
			})
			notification.DELETE("", func(context *gin.Context) {
				methodId := context.Param("method_id")
				method, err := r.handlers.DeleteById(context, &apiPb.NotificationMethodIdRequest{
					Id: methodId,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, method)
			})
		}
		ruleMethod := v1.Group("rule")
		{
			ruleMethod.POST("validate", func(context *gin.Context) {
				validateRuleReq := &ValidateRuleRequest{}

				err := context.ShouldBindJSON(validateRuleReq)
				if err != nil {
					errWrap(context, http.StatusUnprocessableEntity, err)
					return
				}

				rule, err := r.handlers.ValidateRule(context, &apiPb.ValidateRuleRequest{
					Rule:      validateRuleReq.Rule,
					OwnerType: validateRuleReq.OwnerType,
				})

				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}

				successWrap(context, http.StatusOK, rule)
			})
		}
		incidents := v1.Group("incidents")
		{
			incidents.GET("", func(context *gin.Context) {
				rq := &GetIncidentListRequest{}
				err := context.ShouldBind(rq)

				if err != nil {
					errWrap(context, http.StatusBadRequest, err)
					return
				}

				pagination, timeRange, err := GetFilters(rq.Pagination, rq.TimeFilters)
				if err != nil {
					errWrap(context, http.StatusUnprocessableEntity, err)
					return
				}

				res, err := r.handlers.GetIncidentList(context, &apiPb.GetIncidentsListRequest{
					Pagination: pagination,
					TimeRange:  timeRange,
					Status:     rq.Status,
					RuleId:     GetStringValueFromString(rq.RuleId),
					Sort:       GetIncidentListSorting(rq.SortDirection, rq.SortBy),
				})

				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}

				successWrap(context, http.StatusOK, res)
			})

			incident := incidents.Group(":incident_id")

			incident.GET("", func(context *gin.Context) {
				id := context.Param("incident_id")
				inc, err := r.handlers.GetIncidentById(context, &apiPb.IncidentIdRequest{
					IncidentId: id,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}

				successWrap(context, http.StatusOK, inc)
			})

			incident.PUT("close", func(context *gin.Context) {
				id := context.Param("incident_id")
				inc, err := r.handlers.CloseIncident(context, &apiPb.IncidentIdRequest{
					IncidentId: id,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}

				successWrap(context, http.StatusOK, inc)
			})

			incident.PUT("study", func(context *gin.Context) {
				id := context.Param("incident_id")
				inc, err := r.handlers.StudyIncident(context, &apiPb.IncidentIdRequest{
					IncidentId: id,
				})
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}

				successWrap(context, http.StatusOK, inc)
			})
		}
		rules := v1.Group("rules")
		{
			rules.GET("", func(context *gin.Context) {
				rq := &ListRulesByOwnerIdRequest{}
				err := context.ShouldBind(rq)

				if err != nil {
					errWrap(context, http.StatusBadRequest, err)
					return
				}

				rules, err := r.handlers.GetRulesByOwnerId(context, &apiPb.GetRulesByOwnerIdRequest{
					OwnerType: rq.OwnerType,
					OwnerId:   rq.OwnerId,
				})

				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}

				successWrap(context, http.StatusOK, rules)
			})

			rules.POST("", func(context *gin.Context) {
				ruleReq := &CreateRuleRequest{}
				err := context.ShouldBindJSON(ruleReq)
				if err != nil {
					errWrap(context, http.StatusBadRequest, err)
					return
				}

				rule, err := r.handlers.CreateRule(context, &apiPb.CreateRuleRequest{
					Rule:      ruleReq.Rule,
					Name:      ruleReq.Name,
					AutoClose: ruleReq.AutoClose,
					OwnerType: ruleReq.OwnerType,
					OwnerId:   ruleReq.OwnerId,
				})

				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}

				successWrap(context, http.StatusCreated, rule)
			})

			singleRule := rules.Group(":rule_id")
			{
				singleRule.GET("", func(context *gin.Context) {
					ruleId := context.Param("rule_id")
					rule, err := r.handlers.GetRuleById(context, &apiPb.RuleIdRequest{
						RuleId: ruleId,
					})
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusOK, rule)
				})
				singleRule.DELETE("", func(context *gin.Context) {
					ruleId := context.Param("rule_id")
					rule, err := r.handlers.RemoveRuleById(context, &apiPb.RuleIdRequest{
						RuleId: ruleId,
					})
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusOK, rule)
				})

				singleRule.PUT("activate", func(context *gin.Context) {
					ruleId := context.Param("rule_id")
					rule, err := r.handlers.ActivateRuleById(context, &apiPb.RuleIdRequest{
						RuleId: ruleId,
					})
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusOK, rule)
				})
				singleRule.PUT("deactivate", func(context *gin.Context) {
					ruleId := context.Param("rule_id")
					rule, err := r.handlers.DeactivateRuleById(context, &apiPb.RuleIdRequest{
						RuleId: ruleId,
					})
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusOK, rule)
				})
			}
		}

		transaction := v1.Group("transaction")
		{
			transaction.GET(":transaction_id", func(context *gin.Context) {
				trxId := context.Param("transaction_id")
				res, err := r.handlers.GetTransactionById(context, trxId)
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, res)
			})
		}
		applications := v1.Group("applications")
		{
			applications.GET("", func(context *gin.Context) {
				res, err := r.handlers.GetApplicationList(context)
				if err != nil {
					errWrap(context, http.StatusInternalServerError, err)
					return
				}
				successWrap(context, http.StatusOK, res)
			})
			applications.POST("", func(context *gin.Context) {
				application := &Application{}
				err := context.ShouldBindJSON(application)
				if err != nil {
					errWrap(context, http.StatusUnprocessableEntity, err)
					return
				}
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
				application.GET("", func(context *gin.Context) {
					applicationId := context.Param("applicationId")
					res, err := r.handlers.GetApplicationById(context, applicationId)

					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusOK, res)
				})

				application.PUT("enabled", func(context *gin.Context) {
					applicationId := context.Param("applicationId")
					res, err := r.handlers.EnabledApplicationById(context, applicationId)
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusAccepted, res)
				})

				application.DELETE("archived", func(context *gin.Context) {
					applicationId := context.Param("applicationId")
					res, err := r.handlers.ArchivedApplicationById(context, applicationId)
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusAccepted, res)
				})

				application.PUT("disabled", func(context *gin.Context) {
					applicationId := context.Param("applicationId")
					res, err := r.handlers.DisabledApplicationById(context, applicationId)
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusAccepted, res)
				})

				transactions := application.Group("transactions")
				{
					transactions.GET("list", func(context *gin.Context) {
						applicationId := context.Param("applicationId")
						rq := &GetTransactionListRequest{}
						err := context.ShouldBind(rq)
						if err != nil {
							errWrap(context, http.StatusBadRequest, err)
							return
						}
						pagination, timeRange, err := GetFilters(rq.Pagination, rq.TimeFilters)
						if err != nil {
							errWrap(context, http.StatusUnprocessableEntity, err)
							return
						}
						res, err := r.handlers.GetTransactionsList(context, &apiPb.GetTransactionsRequest{
							ApplicationId: applicationId,
							Pagination:    pagination,
							TimeRange:     timeRange,
							Type:          rq.TransactionType,
							Status:        rq.TransactionStatus,
							Host:          GetStringValueFromString(rq.HostFilter),
							Name:          GetStringValueFromString(rq.NameFilter),
							Path:          GetStringValueFromString(rq.PathFilter),
							Method:        GetStringValueFromString(rq.MethodFilter),
							Sort:          GetTransactionListSorting(rq.SortDirection, rq.SortBy),
						})
						if err != nil {
							errWrap(context, http.StatusInternalServerError, err)
							return
						}
						successWrap(context, http.StatusOK, res)
					})
					transactions.GET("group", func(context *gin.Context) {
						applicationId := context.Param("applicationId")
						rq := &GetTransactionGroupRequest{}
						err := context.ShouldBind(rq)
						if err != nil {
							errWrap(context, http.StatusBadRequest, err)
							return
						}
						_, timeRange, err := GetFilters(nil, rq.TimeFilters)
						if err != nil {
							errWrap(context, http.StatusUnprocessableEntity, err)
							return
						}

						res, err := r.handlers.GetTransactionGroups(context, &apiPb.GetTransactionGroupRequest{
							ApplicationId: applicationId,
							TimeRange:     timeRange,
							GroupType:     rq.GroupType,
							Type:          rq.TransactionType,
							Status:        rq.TransactionStatus,
						})
						if err != nil {
							errWrap(context, http.StatusInternalServerError, err)
							return
						}
						successWrap(context, http.StatusOK, res)
					})
					transactions.POST("", func(context *gin.Context) {
						applicationId := context.Param("applicationId")
						trx := &Transaction{}
						err := context.ShouldBindJSON(trx)
						if err != nil {
							successWrap(context, http.StatusAccepted, nil)
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
							StartTime:     trx.DateFrom.ToTimeStamp(),
							EndTime:       trx.DateTo.ToTimeStamp(),
							Status:        trx.Status,
							Type:          trx.Type,
							Error:         trxError,
						})
						if err != nil {
							// we will skip error here
							successWrap(context, http.StatusAccepted, nil)
							return
						}
						successWrap(context, http.StatusAccepted, nil)
					})
				}
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
					rq := &AgentHistory{}

					err := context.ShouldBind(rq)
					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					pagination, timeRange, err := GetFilters(rq.Pagination, rq.TimeFilters)
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
				case apiPb.SchedulerType_SSL_EXPIRATION:
					if request.SSLExpirationConfig == nil {
						errWrap(context, http.StatusUnprocessableEntity, errMissingConfig)
						return
					}
					addReq = &apiPb.AddRequest{
						Config: &apiPb.AddRequest_SslExpiration{
							SslExpiration: request.SSLExpirationConfig,
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

				scheduler.GET("uptime", func(context *gin.Context) {
					schedulerID := context.Param("schedulerId")
					req := &SchedulerUptimeRequest{}
					err := context.ShouldBind(req)

					if err != nil {
						errWrap(context, http.StatusUnprocessableEntity, err)
						return
					}

					_, timeRange, err := GetFilters(nil, req.TimeRange)

					if err != nil {
						errWrap(context, http.StatusUnprocessableEntity, err)
						return
					}
					res, err := r.handlers.GetSchedulerUptime(context, &apiPb.GetSchedulerUptimeRequest{
						SchedulerId: schedulerID,
						TimeRange:   timeRange,
					})

					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					successWrap(context, http.StatusOK, res)
				})

				//History
				scheduler.GET("/history", func(context *gin.Context) {
					schedulerID := context.Param("schedulerId")
					rq := &SchedulerHistory{}
					err := context.ShouldBind(rq)

					if err != nil {
						errWrap(context, http.StatusInternalServerError, err)
						return
					}
					pagination, timeRange, err := GetFilters(rq.Pagination, rq.TimeFilters)

					if err != nil {
						errWrap(context, http.StatusUnprocessableEntity, err)
						return
					}

					res, err := r.handlers.GetSchedulerHistoryByID(context, &apiPb.GetSchedulerInformationRequest{
						SchedulerId: schedulerID,
						Pagination:  pagination,
						TimeRange:   timeRange,
						Sort:        GetSchedulerListSorting(rq.SortDirection, rq.SortBy),
						Status:      rq.Status,
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

func GetSchedulerListSorting(direction apiPb.SortDirection, sortBy apiPb.SortSchedulerList) *apiPb.SortingSchedulerList {
	if sortBy == apiPb.SortSchedulerList_SORT_SCHEDULER_LIST_UNSPECIFIED {
		return nil
	}
	return &apiPb.SortingSchedulerList{
		Direction: direction,
		SortBy:    sortBy,
	}
}

func GetTransactionListSorting(direction apiPb.SortDirection, sortBy apiPb.SortTransactionList) *apiPb.SortingTransactionList {
	if sortBy == apiPb.SortTransactionList_SORT_TRANSACTION_LIST_UNSPECIFIED {
		return nil
	}
	return &apiPb.SortingTransactionList{
		Direction: direction,
		SortBy:    sortBy,
	}
}

func GetIncidentListSorting(direction apiPb.SortDirection, sortBy apiPb.SortIncidentList) *apiPb.SortingIncidentList {
	if sortBy == apiPb.SortIncidentList_SORT_INCIDENT_LIST_UNSPECIFIED {
		return nil
	}
	return &apiPb.SortingIncidentList{
		Direction: direction,
		SortBy:    sortBy,
	}
}

func GetNotificationList(ownerType apiPb.ComponentOwnerType, ownerId string) *apiPb.GetListRequest {
	if ownerType == apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_UNSPECIFIED {
		return nil
	}
	return &apiPb.GetListRequest{
		OwnerId:   ownerId,
		OwnerType: ownerType,
	}
}

func GetStringValueFromString(str string) *wrappers.StringValue {
	if str == "" {
		return nil
	}
	return &wrappers.StringValue{
		Value: str,
	}
}

func GetFilters(paginationFilter *PaginationRequest, timeFilter *TimeFilterRequest) (*apiPb.Pagination, *apiPb.TimeFilter, error) {
	var pagination *apiPb.Pagination
	if paginationFilter == nil {
		pagination = nil
	} else {
		pagination = &apiPb.Pagination{
			Page:  paginationFilter.Page,
			Limit: paginationFilter.Limit,
		}
	}

	if timeFilter != nil {
		timeFilters := &apiPb.TimeFilter{}
		if timeFilter.DateFrom != nil {
			t := timestamp.New(*timeFilter.DateFrom)
			err := t.CheckValid()
			if err != nil {
				return nil, nil, err
			}
			timeFilters.From = t
		}
		if timeFilter.DateTo != nil {
			t := timestamp.New(*timeFilter.DateTo)
			err := t.CheckValid()
			if err != nil {
				return nil, nil, err
			}
			timeFilters.To = t
		}
		return pagination, timeFilters, nil
	}

	return pagination, nil, nil
}

func New(handlers handlers.Handlers) Router {
	return &router{
		handlers: handlers,
	}
}
