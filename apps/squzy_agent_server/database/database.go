package database

import (
	"context"
	"github.com/squzy/mongo_helper"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Database interface {
	Add(ctx context.Context, agent *apiPb.RegisterRequest) (string, error)
	UpdateStatus(ctx context.Context, agentId primitive.ObjectID, status apiPb.AgentStatus) error
	GetAll(ctx context.Context, filter bson.M) ([]*apiPb.AgentItem, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*apiPb.AgentItem, error)
}

type db struct {
	connector mongo_helper.Connector
}

type AgentDao struct {
	Id        primitive.ObjectID `bson:"_id"`
	AgentName string             `bson:"agentName,omitempty"`
	Status    apiPb.AgentStatus  `bson:"status"`
	HostInfo  *HostInfo          `bson:"hostInfo,omitempty"`
	History   []*HistoryItem     `bson:"history"`
}

type HostInfo struct {
	HostName     string        `bson:"hostName,omitempty"`
	Os           string        `bson:"os,omitempty"`
	PlatFormInfo *PlatFormInfo `bson:"platFormInfo,omitempty"`
}

type PlatFormInfo struct {
	Name    string `bson:"name,omitempty"`
	Family  string `bson:"family,omitempty"`
	Version string `bson:"version,omitempty"`
}

type HistoryItem struct {
	Status    apiPb.AgentStatus `bson:"status"`
	Timestamp time.Time         `bson:"time"`
}

func dbToPb(agent *AgentDao) *apiPb.AgentItem {
	a := &apiPb.AgentItem{
		Id:        agent.Id.Hex(),
		AgentName: agent.AgentName,
		Status:    agent.Status,
	}
	if agent.HostInfo != nil {
		a.HostInfo = &apiPb.HostInfo{
			HostName: agent.HostInfo.HostName,
			Os:       agent.HostInfo.Os,
		}
		if agent.HostInfo.PlatFormInfo != nil {
			a.HostInfo.PlatformInfo = &apiPb.PlatformInfo{
				Name:    agent.HostInfo.PlatFormInfo.Name,
				Family:  agent.HostInfo.PlatFormInfo.Family,
				Version: agent.HostInfo.PlatFormInfo.Version,
			}
		}
	}
	return a
}

func (d *db) Add(ctx context.Context, agent *apiPb.RegisterRequest) (string, error) {
	id := primitive.NewObjectID()

	agentData := &AgentDao{
		Id:        id,
		AgentName: agent.AgentName,
		Status:    apiPb.AgentStatus_REGISTRED,
		History: []*HistoryItem{
			{
				Status:    apiPb.AgentStatus_REGISTRED,
				Timestamp: time.Now(),
			},
		},
	}

	if agent.HostInfo != nil {
		agentData.HostInfo = &HostInfo{
			HostName: agent.HostInfo.HostName,
			Os:       agent.HostInfo.Os,
		}
		if agent.HostInfo.PlatformInfo != nil {
			agentData.HostInfo.PlatFormInfo = &PlatFormInfo{
				Name:    agent.HostInfo.PlatformInfo.Name,
				Family:  agent.HostInfo.PlatformInfo.Family,
				Version: agent.HostInfo.PlatformInfo.Version,
			}
		}
	}
	_, err := d.connector.InsertOne(ctx, agentData)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (d *db) GetAll(ctx context.Context, filter bson.M) ([]*apiPb.AgentItem, error) {
	agents := []*AgentDao{}
	err := d.connector.FindAll(ctx, filter, &agents)
	if err != nil {
		return nil, err
	}
	res := []*apiPb.AgentItem{}
	for _, v := range agents {
		res = append(res, dbToPb(v))
	}
	return res, nil
}

func (d *db) GetById(ctx context.Context, id primitive.ObjectID) (*apiPb.AgentItem, error) {
	agentDao := &AgentDao{}
	err := d.connector.FindOne(ctx, bson.M{
		"id": bson.M{
			"$eq": id,
		},
	}, agentDao)
	if err != nil {
		return nil, err
	}
	return dbToPb(agentDao), nil
}

func (d *db) UpdateStatus(ctx context.Context, agentId primitive.ObjectID, status apiPb.AgentStatus) error {
	historyItem := &HistoryItem{
		Status:    status,
		Timestamp: time.Now(),
	}
	_, err := d.connector.UpdateOne(ctx, bson.M{
		"_id": agentId,
	}, bson.M{
		"$set": bson.M{
			"status": status,
		},
		"$push": bson.M{
			"history": historyItem,
		},
	})
	return err
}

func New(connector mongo_helper.Connector) Database {
	return &db{
		connector: connector,
	}
}
