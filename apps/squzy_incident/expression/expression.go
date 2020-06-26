package expression

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/araddon/dateparse"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"strconv"
)

type Expression interface {
	ProcessRule(ruleType apiPb.RuleOwnerType, id string, rule string) bool
	IsValid(ruleType apiPb.RuleOwnerType, rule string) error
}

type expressionStruct struct {
	storageClient apiPb.StorageClient
}

func NewExpression(storage apiPb.StorageClient) Expression {
	return &expressionStruct{
		storageClient: storage,
	}
}

func (e *expressionStruct) ProcessRule(ruleType apiPb.RuleOwnerType, applicationId string, rule string) bool {
	env := e.getTransactionEnv(applicationId)

	program, err := expr.Compile(rule, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	value, err := strconv.ParseBool(fmt.Sprintf("%v", output))
	if err == nil {
		return value
	}
	return false
}

func (e *expressionStruct) IsValid(ruleType apiPb.RuleOwnerType, rule string) error {
	env := e.getTransactionEnv("id")

	_, err := expr.Compile(rule, expr.Env(env))
	return err
}

func (e *expressionStruct) getEnv(owner apiPb.RuleOwnerType, id string) map[string]interface{} {
	switch owner {
	case apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_SCHEDULER:
		return e.getSnapshotEnv(id)
	case apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_AGENT:
		return e.getAgentEnv(id)
	case apiPb.RuleOwnerType_INCIDENT_OWNER_TYPE_APPLICATION:
		return e.getTransactionEnv(id)
	}
	return nil
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
