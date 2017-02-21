package pundit

import (
	"log"
	"math/rand"
	"time"

	"github.com/Knetic/govaluate"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
}

// Rule is a decision rule
type Rule struct {
	Description string

	// Expression is a truth expression such as "name == 'brett'"
	// Expressions are evaluated using github.com/Knetic/govaluate
	Expression string

	// Negate flips the value of the expression
	Negate bool

	// Outcome is the outcome of a matched rule
	Outcome interface{}

	// Title is an arbitrary title for the given rule
	Title string
}

// DecisionTable is the primary object in pundit. It is a table of rules
type DecisionTable struct {
	// BreakOnMatch tells this table to stop when a match is met. This is good
	// for short-circuiting tables where the value is determined quickly.
	BreakOnMatch bool

	// DefaultOutcome is the value of the table if none of the rules match
	DefaultOutcome interface{}

	Description string

	// Rnd is a random number between 0 and 100 created with the decision table
	Rnd int32

	Rules []Rule

	Title string
}

// Evaluate evaluates a request against the rules
func (e *DecisionTable) Evaluate(input map[string]interface{}) (*ResultSet, error) {
	input["rnd"] = e.Rnd

	rs := &ResultSet{
		Rules: make([]Rule, len(e.Rules)),
	}

	// finished is set to true when a match was found and e.BreakOnMatch is set to true. This will
	// prevent other rules from executing, but still completes the loop setting each item in rs.Rules
	// to the same value of it's original rule.
	finished := false
	for i, r := range e.Rules {
		rs.Rules[i] = r
		// clear the outcome to be determined below
		rs.Rules[i].Outcome = nil

		// e.BreakOnMatch is true, and a previous rule already matched, skip evaluation
		if finished {
			continue
		}

		evaluator, err := govaluate.NewEvaluableExpression(r.Expression)
		if err != nil {
			log.Fatalf("could not evaluate expression (%s) from rule %s: %s", r.Expression, r.Title, err)
		}
		result, err := evaluator.Evaluate(input)
		if err != nil {
			return nil, err
		}

		// XOR result and 'negate'
		// This switch is activated when the rule applies to this input
		if result.(bool) != r.Negate {
			rs.FinalDecision = r.Outcome
			rs.Rules[i].Outcome = r.Outcome
			if e.BreakOnMatch {
				finished = true
			}
		}

	}

	return rs, nil
}

// ScoreTable is the primary object in pundit. It is a table of rules
type ScoreTable struct {
	Title       string
	Description string
	Rules       []Rule
	Score       float64
}

// ResultSet is a result of a table
type ResultSet struct {
	Rules         []Rule
	FinalDecision interface{}
}
