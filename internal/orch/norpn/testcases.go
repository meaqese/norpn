package norpn

type TestCase struct {
	Name       string
	Expression string
	Expected   float64
	ShouldFail bool
}

func GetTestCases() []TestCase {
	cases := []TestCase{
		{
			Name:       "simple",
			Expression: "2+2",
			Expected:   4,
			ShouldFail: false,
		},
		{
			Name:       "priority",
			Expression: "(2+2)*2",
			Expected:   8,
			ShouldFail: false,
		},
		{
			Name:       "priority 2",
			Expression: "2+2*2",
			Expected:   6,
			ShouldFail: false,
		},
		{
			Name:       "/",
			Expression: "1/2",
			Expected:   0.5,
			ShouldFail: false,
		},
		{
			Name:       "priority 3",
			Expression: "2+2*(5-(2 + 1))",
			Expected:   6,
			ShouldFail: false,
		},
		{
			Name:       "division by zero",
			Expression: "1/0",
			Expected:   0,
			ShouldFail: true,
		},
		{
			Name:       "invalid Expression",
			Expression: "1+2**",
			Expected:   0,
			ShouldFail: true,
		},
		{
			Name:       "invalid Expression 2",
			Expression: "1*-2",
			Expected:   0,
			ShouldFail: true,
		},
		{
			Name:       "empty",
			Expression: "",
			Expected:   0,
			ShouldFail: true,
		},
		{
			Name:       "priority fail",
			Expression: "((2+2-*)",
			Expected:   0,
			ShouldFail: true,
		},
	}

	return cases
}
