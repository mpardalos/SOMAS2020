package iigointernal

import (
	"testing"

	"reflect" // Used to compare two maps

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type allocInput struct {
	resourceRequests   map[int]int
	commonPoolResource int
}

func TestAllocationRequests(t *testing.T) {
	cases := []struct {
		name  string
		input allocInput
		reply map[int]int
		want  error
	}{
		{
			name:  "Limitied resources",
			input: allocInput{resourceRequests: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30}, commonPoolResource: 100},
			reply: map[int]int{1: 3, 2: 7, 3: 10, 4: 14, 5: 17, 6: 21},
			want:  nil,
		},
		{
			name:  "Enough resources",
			input: allocInput{resourceRequests: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30}, commonPoolResource: 150},
			reply: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
			want:  nil,
		},
		{
			name:  "No resources",
			input: allocInput{resourceRequests: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30}, commonPoolResource: 0},
			reply: map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0},
			want:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			president := &basePresident{
				Id:     int(shared.Team1),
				budget: 50,
			}

			val, err := president.EvaluateAllocationRequests(tc.input.resourceRequests, tc.input.commonPoolResource)
			if err == nil && tc.want == nil {
				if !reflect.DeepEqual(tc.reply, val) {
					t.Errorf("%v - Failed. Expected %v, got %v", tc.name, tc.reply, val)
				}
			} else {
				if err != tc.want {
					t.Errorf("%v - Failed. Returned error '%v' which was not equal to expected '%v'", tc.name, err, tc.want)
				}
			}
		})
	}
}

func checkIfInList(s string, lst []string) bool {
	for _, i := range lst {
		if s == i {
			return true
		}
	}
	return false
}
func TestPickRuleToVote(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		reply string
		want  error
	}{
		{
			name:  "Basic rule",
			input: []string{"rule"},
			reply: "rule",
			want:  nil,
		},
		{
			name:  "Empty string",
			input: []string{""},
			reply: "",
			want:  nil,
		},
		{
			name:  "Longer list",
			input: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
			reply: "",
			want:  nil,
		},
		{
			name:  "Empty list",
			input: []string{},
			reply: "",
			want:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			president := &basePresident{}

			val, err := president.PickRuleToVote(tc.input)
			if err == nil && tc.want == nil {
				if len(tc.input) == 0 {
					if val != "" {
						t.Errorf("%v - Failed. Returned '%v', but expectd an empty string", tc.name, val)
					}
				} else if !checkIfInList(val, tc.input) {
					t.Errorf("%v - Failed. Returned '%v', expected '%v'", tc.name, val, tc.input)
				}
			} else {
				if err != tc.want {
					t.Errorf("%v - Failed. Returned error '%v' which was not equal to expected '%v'", tc.name, err, tc.want)
				}
			}
		})
	}
}

func TestSetRuleProposals(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Single Rule",
			input:    []string{"rule"},
			expected: []string{"rule"},
		},
		{
			name:     "No Rule",
			input:    []string{""},
			expected: []string{""},
		},
		{
			name:     "Multiple Rules",
			input:    []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
			expected: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			president := &basePresident{
				Id:             int(shared.Team1),
				RulesProposals: []string{},
			}

			president.setRuleProposals(tc.input)

			if !reflect.DeepEqual(president.RulesProposals, tc.expected) {
				t.Errorf("%v - Failed. RulesProposals set to '%v', expected '%v'", tc.name, president.RulesProposals, tc.expected)
			}

		})
	}
}

func TestGetTaxMap(t *testing.T) {
	cases := []struct {
		name           string
		input          map[int]int
		bPresident     basePresident // base
		expectedLength int
	}{
		{
			name:  "Empty tax map base",
			input: map[int]int{},
			bPresident: basePresident{
				Id:              3,
				clientPresident: nil,
			},
			expectedLength: 0,
		},
		{
			name:  "Short tax map base",
			input: map[int]int{1: 5},
			bPresident: basePresident{
				Id:              3,
				clientPresident: nil,
			},
			expectedLength: 1,
		},
		{
			name:  "Long tax map base",
			input: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
			bPresident: basePresident{
				Id:              4,
				clientPresident: nil,
			},
			expectedLength: 6,
		},
		{
			name:  "Client empty tax map base",
			input: map[int]int{},
			bPresident: basePresident{
				Id: 5,
				clientPresident: &basePresident{
					Id:               3,
					ResourceRequests: nil,
				},
			},
			expectedLength: 0,
		},
		{
			name:  "Client short tax map base",
			input: map[int]int{1: 5},
			bPresident: basePresident{
				Id: 5,
				clientPresident: &basePresident{
					Id:               3,
					ResourceRequests: nil,
				},
			},
			expectedLength: 1,
		},
		{
			name:  "Client long tax map base",
			input: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
			bPresident: basePresident{
				Id: 5,
				clientPresident: &basePresident{
					Id:               3,
					ResourceRequests: nil,
				},
			},
			expectedLength: 6,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val := tc.bPresident.getTaxMap(tc.input)
			if len(val) != tc.expectedLength {
				t.Errorf("%v - Failed. RulesProposals set to '%v', expected '%v'", tc.name, val, tc.input)
			}

		})
	}
}

func TestGetRuleForSpeaker(t *testing.T) {
	cases := []struct {
		name       string
		bPresident basePresident // base
		expected   []string
	}{
		{
			name: "Empty tax map base",
			bPresident: basePresident{
				Id:              3,
				RulesProposals:  []string{},
				clientPresident: nil,
			},
			expected: []string{""},
		},
		{
			name: "Short tax map base",
			bPresident: basePresident{
				Id:              3,
				RulesProposals:  []string{"test"},
				clientPresident: nil,
			},
			expected: []string{"test"},
		},
		{
			name: "Long tax map base",
			bPresident: basePresident{
				Id:              3,
				RulesProposals:  []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
				clientPresident: nil,
			},
			expected: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
		},
		{
			name: "Client empty tax map base",
			bPresident: basePresident{
				Id:             5,
				RulesProposals: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
				clientPresident: &basePresident{
					Id:             3,
					RulesProposals: []string{},
				},
			},
			expected: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
		},
		{
			name: "Client tax map base override",
			bPresident: basePresident{
				Id:             5,
				RulesProposals: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
				clientPresident: &basePresident{
					Id:             3,
					RulesProposals: []string{"test1", "test2", "test3"},
				},
			},
			expected: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			val := tc.bPresident.getRuleForSpeaker()
			if len(tc.expected) == 0 {
				if val != "" {
					t.Errorf("%v - Failed. Returned '%v', but expectd an empty string", tc.name, val)
				}
			} else if !checkIfInList(val, tc.expected) {
				t.Errorf("%v - Failed. Returned '%v', expected '%v'", tc.name, val, tc.expected)
			}
		})
	}
}

func TestGetAllocationRequests(t *testing.T) {
	cases := []struct {
		name       string
		bPresident basePresident // base
		input      int
		expected   map[int]int
	}{
		{
			name: "Limited Resources",
			bPresident: basePresident{
				Id:               3,
				ResourceRequests: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
				clientPresident:  nil,
			},
			input:    100,
			expected: map[int]int{1: 3, 2: 7, 3: 10, 4: 14, 5: 17, 6: 21},
		},
		{
			name: "Excess Resources",
			bPresident: basePresident{
				Id:               3,
				ResourceRequests: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
				clientPresident:  nil,
			},
			input:    150,
			expected: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
		},
		{
			name: "Client override limited resourceRequests",
			bPresident: basePresident{
				Id:               5,
				ResourceRequests: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
				clientPresident: &basePresident{
					Id:               3,
					ResourceRequests: map[int]int{1: 15, 2: 20, 25: 3, 100: 4, 35: 5, 6: 30},
				},
			},
			input:    100,
			expected: map[int]int{1: 3, 2: 7, 3: 10, 4: 14, 5: 17, 6: 21},
		},
		{
			name: "Client override excess resourceRequests",
			bPresident: basePresident{
				Id:               5,
				ResourceRequests: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
				clientPresident: &basePresident{
					Id:               3,
					ResourceRequests: map[int]int{1: 15, 2: 20, 25: 3, 100: 4, 35: 5, 6: 30},
				},
			},
			input:    150,
			expected: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			val := tc.bPresident.getAllocationRequests(tc.input)
			if !reflect.DeepEqual(val, tc.expected) {
				t.Errorf("%v - Failed. Got '%v', but expected '%v'", tc.name, val, tc.expected)
			}
		})
	}
}
