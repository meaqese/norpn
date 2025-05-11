package services

type TestCase struct {
	Name           string
	Expression     string
	Expected       float64
	ShouldFail     bool
	TimeoutSeconds int
}

func GetTestCases() []TestCase {
	cases := []TestCase{
		{
			Name:           "simple",
			Expression:     "2+2",
			Expected:       4,
			ShouldFail:     false,
			TimeoutSeconds: 2,
		},
		{
			Name:           "priority",
			Expression:     "(2+2)*2",
			Expected:       8,
			ShouldFail:     false,
			TimeoutSeconds: 4,
		},
		{
			Name:           "priority 2",
			Expression:     "2+2*2",
			Expected:       6,
			ShouldFail:     false,
			TimeoutSeconds: 3,
		},
		{
			Name:           "/",
			Expression:     "1/2",
			Expected:       0.5,
			ShouldFail:     false,
			TimeoutSeconds: 2,
		},
		{
			Name:           "priority 3",
			Expression:     "2+2*(5-(2 + 1))",
			Expected:       6,
			ShouldFail:     false,
			TimeoutSeconds: 5,
		},
		{
			Name:           "division by zero",
			Expression:     "1/0",
			Expected:       0,
			ShouldFail:     true,
			TimeoutSeconds: 1,
		},
		{
			Name:           "invalid Expression",
			Expression:     "1+2**",
			Expected:       0,
			ShouldFail:     true,
			TimeoutSeconds: 2,
		},
		{
			Name:           "invalid Expression 2",
			Expression:     "1*-2",
			Expected:       0,
			ShouldFail:     true,
			TimeoutSeconds: 2,
		},
		{
			Name:           "empty",
			Expression:     "",
			Expected:       0,
			ShouldFail:     true,
			TimeoutSeconds: 1,
		},
		{
			Name:           "priority fail",
			Expression:     "((2+2-*)",
			Expected:       0,
			ShouldFail:     true,
			TimeoutSeconds: 2,
		},
	}

	return cases
}
