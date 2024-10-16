package routes

type DivisionbyZeroError struct{}

func (d DivisionbyZeroError) Error() string {
	return "Cannot divide by zero"
}

type InvalidKeysError struct{}

func (i InvalidKeysError) Error() string {
	return "Bad Request"
}

type InvalidKeyValuesError struct{}

func (i InvalidKeyValuesError) Error() string {
	return "number1 and number2 must be numbers"
}

type SumOperationKeyValueError struct{}

func (s SumOperationKeyValueError) Error() string {
	return "numbers must be an array of numbers"
}
