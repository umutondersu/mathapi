package routes

type DivisionbyZeroError struct{}

func (d DivisionbyZeroError) Error() string {
	return "Cannot divide by zero"
}

type InvalidKeysError struct{}

func (i InvalidKeysError) Error() string {
	return "Bad Request"
}

type BOValuesError struct{}

func (i BOValuesError) Error() string {
	return "number1 and number2 must be numbers"
}

type SOValueError struct{}

func (s SOValueError) Error() string {
	return "numbers must be an array of numbers"
}
