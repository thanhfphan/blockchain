package ids

import (
	"testing"
)

func Test_IDString(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		label  string
		id     ID
		expect string
	}{
		{
			label:  "ID{}",
			id:     ID{},
			expect: "11111111111111111111111111111111",
		},
		{
			label:  "ID{}",
			id:     ID{33},
			expect: "3DpTLLwnVbLZUpy9nenjYabYYKeSJMMvGdBJcSEqnYpP",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.label, func(t *testing.T) {
			result := tc.id.String()
			if result != tc.expect {
				t.Errorf("exptect %s but return %s\n", tc.expect, result)
			}
		})
	}
}
