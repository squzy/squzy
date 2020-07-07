package expression

import (
	"errors"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/araddon/dateparse"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"strconv"
	"time"
)

type Expression interface {
	ProcessRule(ruleType apiPb.RuleOwnerType, id string, rule string) (bool, error)
	IsValid(ruleType apiPb.RuleOwnerType, rule string) error
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

func (e *expressionStruct) ProcessRule(ruleType apiPb.RuleOwnerType, id string, rule string) (bool, error) {
	env, err := e.getEnv(ruleType, id)

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

func (e *expressionStruct) IsValid(ruleType apiPb.RuleOwnerType, rule string) error {
	env, err := e.getEnv(ruleType, "id")
	if err != nil {
		return err
	}
	_, err = expr.Compile(rule, expr.Env(env))
	return err
}

func (e *expressionStruct) getEnv(owner apiPb.RuleOwnerType, id string) (map[string]interface{}, error) {
	switch owner {
	case apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_SCHEDULER:
		return e.getSnapshotEnv(id), nil
	case apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT:
		return e.getAgentEnv(id), nil
	case apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_APPLICATION:
		return e.getTransactionEnv(id), nil
	}
	return nil, errRuleTypeNotProvided
}

func convertToTimestamp(strTime string) *timestamp.Timestamp {
	t, err := dateparse.ParseAny(strTime)
	if err != nil {
		panic(err)
	}
	res, err := ptypes.TimestampProto(t)
	if err != nil {
		panic(err)
	}
	return res
}

func getTimeRange(start, end *timestamp.Timestamp) int64 {
	startTime, err := ptypes.Timestamp(start)
	if err != nil {
		panic("No start time")
	}
	endTime, err := ptypes.Timestamp(end)
	if err != nil {
		panic("No end time")
	}
	return (endTime.UnixNano() - startTime.UnixNano()) / int64(time.Millisecond)
}
