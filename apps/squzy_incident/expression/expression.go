package expression

import (
	"errors"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/araddon/dateparse"
	apiPb "github.com/squzy/squzy_generated/generated/github.com/squzy/squzy_proto"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

type Expression interface {
	ProcessRule(ruleType apiPb.ComponentOwnerType, id string, rule string) (bool, error)
	IsValid(ruleType apiPb.ComponentOwnerType, rule string) error
}

type expressionStruct struct {
	storageClient apiPb.StorageClient
}

var (
	errRuleTypeNotProvided = errors.New("rule type not provided")
)

func NewExpression(storage apiPb.StorageClient) Expression {
	return &expressionStruct{
		storageClient: storage,
	}
}

func (e *expressionStruct) ProcessRule(ruleType apiPb.ComponentOwnerType, id string, rule string) (bool, error) {
	env, err := e.getEnv(ruleType, id)

	if err != nil {
		return false, err
	}

	program, err := expr.Compile(rule, expr.Env(env))
	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}
	value, err := strconv.ParseBool(fmt.Sprintf("%v", output))
	if err != nil {
		return false, err
	}
	return value, nil
}

func (e *expressionStruct) IsValid(ruleType apiPb.ComponentOwnerType, rule string) error {
	env, err := e.getEnv(ruleType, "id")
	if err != nil {
		return err
	}
	_, err = expr.Compile(rule, expr.Env(env))
	return err
}

func (e *expressionStruct) getEnv(owner apiPb.ComponentOwnerType, id string) (map[string]interface{}, error) {
	switch owner {
	case apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_SCHEDULER:
		return e.getSnapshotEnv(id), nil
	case apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_AGENT:
		return e.getAgentEnv(id), nil
	case apiPb.ComponentOwnerType_COMPONENT_OWNER_TYPE_APPLICATION:
		return e.getTransactionEnv(id), nil
	}
	return nil, errRuleTypeNotProvided
}

func convertToTimestamp(strTime string) *timestamp.Timestamp {
	t, err := dateparse.ParseAny(strTime)
	if err != nil {
		panic(err)
	}
	res := timestamp.New(t)
	err = res.CheckValid()
	if err != nil {
		panic(err)
	}
	return res
}

func getTimeRange(start, end *timestamp.Timestamp) int64 {
	err := start.CheckValid()
	if err != nil {
		panic("No start time")
	}
	startTime := start.AsTime()
	err = end.CheckValid()
	if err != nil {
		panic("No end time")
	}
	endTime := end.AsTime()

	return (endTime.UnixNano() - startTime.UnixNano()) / int64(time.Millisecond)
}
