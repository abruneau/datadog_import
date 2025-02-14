package synthetic

import (
	"fmt"
)

func UnknownSyntheticsAssertionOperatorError(operator, def string) error {
	message := fmt.Sprintf("unknown SyntheticsAssertionOperator: %s", operator)
	if def != "" {
		message += fmt.Sprintf(", the operator has been set to the default: %s", def)

	}
	return fmt.Errorf(message)
}

func InvalidSyntheticsAssertionOperatorError(operator, def string) error {
	message := fmt.Sprintf("invalid SyntheticsAssertionOperator: %s", operator)
	if def != "" {
		message += fmt.Sprintf(", the operator has been set to the default: %s", def)

	}
	return fmt.Errorf(message)
}
