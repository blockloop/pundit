package pundit

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecision(t *testing.T) {
	dt := &DecisionTable{
		BreakOnMatch:   true,
		DefaultOutcome: "default",
		Title:          "test one",
		Rules: []Rule{
			{
				Title:      "Under 21",
				Expression: "age < 21",
				Outcome:    "Under 21",
			},
			{
				Title:      "Very Low Income",
				Expression: "income < 1",
				Outcome:    "Very Low Income",
			},
		},
	}

	res, err := dt.Evaluate(map[string]interface{}{
		"age":    20,
		"income": 10,
	})
	assert.NoError(t, err, "expected Evaluate() to not err")
	assert.Equal(t, res.Rules[0].Outcome, dt.Rules[0].Outcome)
	assert.Nil(t, res.Rules[1].Outcome)
}

func TestDecisionBreakOnMatch(t *testing.T) {
	dt := &DecisionTable{
		BreakOnMatch:   true,
		DefaultOutcome: "default",
		Title:          "test one",
		Rules: []Rule{
			{
				Title:      "Under 21",
				Expression: "age < 21",
				Outcome:    "Under 21",
			},
			{
				Title:      "Very Low Income",
				Expression: "income < 10000",
				Outcome:    "Very Low Income",
			},
		},
	}

	res, err := dt.Evaluate(map[string]interface{}{
		"age":    20,
		"income": 10,
	})
	assert.NoError(t, err, "expected Evaluate() to not err")
	assert.Equal(t, res.FinalDecision, dt.Rules[0].Outcome)
}

func TestDecisionNotBreakOnMatch(t *testing.T) {
	dt := &DecisionTable{
		BreakOnMatch:   false,
		DefaultOutcome: "default",
		Title:          "test two",
		Rules: []Rule{
			{
				Title:      "Under 21",
				Expression: "age < 21",
				Outcome:    "Under 21",
			},
			{
				Title:      "Very Low Income",
				Expression: "income < 10000",
				Outcome:    "Very Low Income",
			},
		},
	}

	res, err := dt.Evaluate(map[string]interface{}{
		"age":    20,
		"income": 10,
	})
	assert.NoError(t, err, "expected Evaluate() to not err")
	assert.Equal(t, res.FinalDecision, dt.Rules[1].Outcome)
}

func TestDecisionRndEvaluates(t *testing.T) {
	dt := &DecisionTable{
		BreakOnMatch:   true,
		DefaultOutcome: "default",
		Title:          "test nine",
		Rnd:            int32(100.0 * rand.Float64()),
		Rules: []Rule{
			{
				Expression: "rnd < 200",
				Outcome:    "left",
			},
			{
				Expression: "rnd >= 200",
				Outcome:    "right",
			},
		},
	}

	res, err := dt.Evaluate(map[string]interface{}{})
	assert.NoError(t, err, "expected Evaluate() to not err")
	assert.Equal(t, res.FinalDecision, dt.Rules[0].Outcome)
}

func TestDecisionSetsRulesNilWhenNotRun(t *testing.T) {
	dt := &DecisionTable{
		BreakOnMatch:   true,
		DefaultOutcome: "default",
		Title:          "test nine",
		Rules: []Rule{
			{
				Expression: "num == 10",
				Outcome:    "10",
			},
			{
				Expression: "num == 20",
				Outcome:    "20",
			},
			{
				Expression: "num == 30",
				Outcome:    "30",
			},
		},
	}

	res, err := dt.Evaluate(map[string]interface{}{
		"num": 10,
	})
	assert.NoError(t, err, "expected Evaluate() to not err")
	assert.Nil(t, res.Rules[1].Outcome)
	assert.Nil(t, res.Rules[2].Outcome)
}
