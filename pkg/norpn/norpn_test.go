package norpn

import "testing"

func TestSuccess(t *testing.T) {
	cases := GetTestCases()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := Calc(tc.Expression)
			if err != nil && !tc.ShouldFail {
				t.Fatalf("successful case %s returns error %s", tc.Expression, err)
			}

			if !tc.ShouldFail && result != tc.Expected {
				t.Fatalf("expected \"%f\" got \"%f\" in expression - %s", tc.Expected, result, tc.Expression)
			}
		})
	}
}
