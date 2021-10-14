package taintsandtolerations

import (
	"github.com/microsoft/moc-sdk-for-go/services/common"
	wssdcommon "github.com/microsoft/moc/rpc/common"
)

func GetWssdTaints(taints *[]common.Taint) []*wssdcommon.Taint {
	result := []*wssdcommon.Taint{}
	if taints != nil {
		for _, taint := range *taints {
			result = append(result, GetWssdTaint(taint))
		}
	}
	return result
}

func GetTaints(taints []*wssdcommon.Taint) *[]common.Taint {
	result := []common.Taint{}
	for _, taint := range taints {
		if taint == nil {
			continue
		}
		result = append(result, GetTaint(*taint))
	}
	return &result
}

func GetWssdTaint(taint common.Taint) *wssdcommon.Taint {
	return &wssdcommon.Taint{
		Key:   taint.Key,
		Value: taint.Value,
	}
}

func GetTaint(taint wssdcommon.Taint) common.Taint {
	return common.Taint{
		Key:   taint.Key,
		Value: taint.Value,
	}
}

func GetWssdTolerations(tolerations *[]common.Toleration) []*wssdcommon.Toleration {
	result := []*wssdcommon.Toleration{}
	if tolerations != nil {
		for _, toleration := range *tolerations {
			result = append(result, GetWssdToleration(toleration))
		}
	}
	return result
}

func GetTolerations(tolerations []*wssdcommon.Toleration) *[]common.Toleration {
	result := []common.Toleration{}
	for _, toleration := range tolerations {
		if toleration == nil {
			continue
		}
		result = append(result, GetToleration(*toleration))
	}
	return &result
}

func GetWssdToleration(toleration common.Toleration) *wssdcommon.Toleration {
	return &wssdcommon.Toleration{
		Operator: GetWssdTolerationOperator(toleration.Operator),
		Key:      toleration.Key,
		Value:    toleration.Value,
		Required: toleration.Required,
	}
}

func GetToleration(toleration wssdcommon.Toleration) common.Toleration {
	return common.Toleration{
		Operator: GetTolerationOperator(toleration.Operator),
		Key:      toleration.Key,
		Value:    toleration.Value,
		Required: toleration.Required,
	}
}

func GetWssdTolerationOperator(operator common.TolerationOperator) wssdcommon.TolerationOperator {
	switch operator {
	case common.TolerationOperator_Exists:
		return wssdcommon.TolerationOperator_Exists

	case common.TolerationOperator_Equal:
		return wssdcommon.TolerationOperator_Equal

	default:
		return wssdcommon.TolerationOperator_NONE
	}
}

func GetTolerationOperator(operator wssdcommon.TolerationOperator) common.TolerationOperator {
	switch operator {
	case wssdcommon.TolerationOperator_Exists:
		return common.TolerationOperator_Exists

	case wssdcommon.TolerationOperator_Equal:
		return common.TolerationOperator_Equal

	default:
		return ""
	}
}
