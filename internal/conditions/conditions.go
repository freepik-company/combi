package conditions

import (
	"fmt"

	"gcmerge/api/v1alpha1"
	"gcmerge/internal/globals"
	"gcmerge/internal/template"
)

const (
	mandatoryConditionErrorMessage = "mandatory condition '%s' fail with value { %s } and result { %s }"
	optionalConditionErrorMessage  = "optional condition '%s' fail with value { %s } and result { %s }"
)

func EvalConditions(conditions *v1alpha1.ConditionsT, config *map[string]interface{}) (success bool, err error) {
	for _, condition := range conditions.Mandatory {
		result, err := template.EvaluateTemplate(condition.Template, *config)
		if err != nil {
			return success, err
		}

		if condition.Value != result {
			err = fmt.Errorf(mandatoryConditionErrorMessage, condition.Name, condition.Value, result)
			return success, err
		}
	}

	for _, condition := range conditions.Optional {
		result, err := template.EvaluateTemplate(condition.Template, *config)
		if err != nil {
			return success, err
		}

		if condition.Value != result {
			globals.ExecContext.Logger.Warnf(optionalConditionErrorMessage, condition.Name, condition.Value, result)
		}
	}

	success = true

	return success, err
}