package conditions

import (
	"fmt"
	"gcmerge/internal/config"
	"gcmerge/internal/globals"
	"gcmerge/internal/template"
)

const (
	mandatoryConditionErrorMessage = "mandatory condition '%s' fail with value { %s } and result { %s }"
	optionalConditionErrorMessage  = "optional condition '%s' fail with value { %s } and result { %s }"
)

func RunConditions(conditions *config.ConditionsT, config map[string]interface{}) (err error) {
	for _, condition := range conditions.Mandatory {
		result, err := template.EvaluateTemplate(condition.Template, config)
		if err != nil {
			return err
		}

		if condition.Value != result {
			err = fmt.Errorf(mandatoryConditionErrorMessage, condition.Name, condition.Value, result)
			return err
		}
	}

	for _, condition := range conditions.Optional {
		result, err := template.EvaluateTemplate(condition.Template, config)
		if err != nil {
			return err
		}

		if condition.Value != result {
			globals.ExecContext.Logger.Warnf(optionalConditionErrorMessage, condition.Name, condition.Value, result)
		}
	}

	return err
}
