package endpoint

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestCondition_Validate(t *testing.T) {
	scenarios := []struct {
		condition   Condition
		expectedErr error
	}{
		{condition: "[STATUS] == 200", expectedErr: nil},
		{condition: "[STATUS] != 200", expectedErr: nil},
		{condition: "[STATUS] <= 200", expectedErr: nil},
		{condition: "[STATUS] >= 200", expectedErr: nil},
		{condition: "[STATUS] < 200", expectedErr: nil},
		{condition: "[STATUS] > 200", expectedErr: nil},
		{condition: "[STATUS] == any(200, 201, 202, 203)", expectedErr: nil},
		{condition: "[STATUS] == [BODY].status", expectedErr: nil},
		{condition: "[CONNECTED] == true", expectedErr: nil},
		{condition: "[RESPONSE_TIME] < 500", expectedErr: nil},
		{condition: "[IP] == 127.0.0.1", expectedErr: nil},
		{condition: "[BODY] == 1", expectedErr: nil},
		{condition: "[BODY].test == wat", expectedErr: nil},
		{condition: "[BODY].test.wat == wat", expectedErr: nil},
		{condition: "[BODY].age == [BODY].id", expectedErr: nil},
		{condition: "[BODY].users[0].id == 1", expectedErr: nil},
		{condition: "len([BODY].users) == 100", expectedErr: nil},
		{condition: "len([BODY].data) < 5", expectedErr: nil},
		{condition: "has([BODY].errors) == false", expectedErr: nil},
		{condition: "has([BODY].users[0].name) == true", expectedErr: nil},
		{condition: "[BODY].name == pat(john*)", expectedErr: nil},
		{condition: "[CERTIFICATE_EXPIRATION] > 48h", expectedErr: nil},
		{condition: "[DOMAIN_EXPIRATION] > 720h", expectedErr: nil},
		{condition: "raw == raw", expectedErr: nil},
		{condition: "[STATUS] ? 201", expectedErr: errors.New("invalid condition: [STATUS] ? 201")},
		{condition: "[STATUS]==201", expectedErr: errors.New("invalid condition: [STATUS]==201")},
		{condition: "[STATUS] = = 201", expectedErr: errors.New("invalid condition: [STATUS] = = 201")},
		{condition: "[STATUS] ==", expectedErr: errors.New("invalid condition: [STATUS] ==")},
		{condition: "[STATUS]", expectedErr: errors.New("invalid condition: [STATUS]")},
	}
	for _, scenario := range scenarios {
		t.Run(string(scenario.condition), func(t *testing.T) {
			if err := scenario.condition.Validate(); fmt.Sprint(err) != fmt.Sprint(scenario.expectedErr) {
				t.Errorf("expected err %v, got %v", scenario.expectedErr, err)
			}
		})
	}
}

func TestCondition_Evaluate(t *testing.T) {
	tests := []struct {
		name             string
		condition        Condition
		result           *Result
		expectedSuccess  bool
		expectedOutput   string
		expectedSeverity SeverityStatus
	}{
		{"IP matches", Condition("[IP] == 127.0.0.1"), &Result{IP: "127.0.0.1"}, true, "[IP] == 127.0.0.1", None},
		{"Status is 200", Condition("[STATUS] == 200"), &Result{HTTPStatus: 200}, true, "[STATUS] == 200", None},
		{"Status is not 200", Condition("[STATUS] == 200"), &Result{HTTPStatus: 500}, false, "[STATUS] (500) == 200", Critical},
		{"Custom severity low", Condition("Low :: [STATUS] == 200"), &Result{HTTPStatus: 500}, false, "[STATUS] (500) == 200", Low},
		{"Custom severity medium", Condition("Medium :: [STATUS] == 200"), &Result{HTTPStatus: 500}, false, "[STATUS] (500) == 200", Medium},
		{"Custom severity high", Condition("High :: [STATUS] == 200"), &Result{HTTPStatus: 500}, false, "[STATUS] (500) == 200", High},
		{"Custom severity critical", Condition("Critical :: [STATUS] == 200"), &Result{HTTPStatus: 500}, false, "[STATUS] (500) == 200", Critical},
		{"Response time under limit", Condition("[RESPONSE_TIME] < 500"), &Result{Duration: 50 * time.Millisecond}, true, "[RESPONSE_TIME] < 500", None},
		{"Response time over limit", Condition("[RESPONSE_TIME] > 500"), &Result{Duration: 750 * time.Millisecond}, true, "[RESPONSE_TIME] > 500", None},
		{"Body JSONPath value match", Condition("[BODY].status == UP"), &Result{Body: []byte(`{"status":"UP"}`)}, true, "[BODY].status == UP", None},
		{"Body JSONPath complex", Condition("[BODY].data.name == john"), &Result{Body: []byte(`{"data": {"id": 1, "name": "john"}}`)}, true, "[BODY].data.name == john", None},
		{"Body JSONPath array", Condition("[BODY][0].id == 1"), &Result{Body: []byte(`[{"id": 1}, {"id": 2}]`)}, true, "[BODY][0].id == 1", None},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.condition.evaluate(test.result, false)
			if test.result.ConditionResults[0].Success != test.expectedSuccess {
				t.Errorf("Expected success=%v, got %v", test.expectedSuccess, test.result.ConditionResults[0].Success)
			}
			if test.result.ConditionResults[0].Condition != test.expectedOutput {
				t.Errorf("Expected condition='%s', got '%s'", test.expectedOutput, test.result.ConditionResults[0].Condition)
			}
			if test.result.ConditionResults[0].SeverityStatus != test.expectedSeverity {
				t.Errorf("Expected severity='%v', got '%v'", test.expectedSeverity, test.result.ConditionResults[0].SeverityStatus)
			}
		})
	}
}

func TestCondition_SanitizeSeverityCondition(t *testing.T) {
	tests := []struct {
		condition         Condition
		expectedSeverity  SeverityStatus
		expectedCondition string
	}{
		{Condition("[STATUS] == 200"), Critical, "[STATUS] == 200"},
		{Condition("Low :: [STATUS] == 201"), Low, "[STATUS] == 201"},
		{Condition("Medium :: [STATUS] == 404"), Medium, "[STATUS] == 404"},
		{Condition("High :: [STATUS] == 500"), High, "[STATUS] == 500"},
		{Condition("Critical :: [STATUS] == 500"), Critical, "[STATUS] == 500"},
	}

	for _, test := range tests {
		t.Run(string(test.condition), func(t *testing.T) {
			severity, condition := test.condition.sanitizeSeverityCondition()
			if severity != test.expectedSeverity {
				t.Errorf("Expected severity=%v, got %v", test.expectedSeverity, severity)
			}
			if condition != test.expectedCondition {
				t.Errorf("Expected condition='%s', got '%s'", test.expectedCondition, condition)
			}
		})
	}
}

func TestCondition_EvaluateWithInvalidOperator(t *testing.T) {
	condition := Condition("[STATUS] ? 201")
	result := &Result{HTTPStatus: 201}
	condition.evaluate(result, false)
	if result.Success {
		t.Error("Expected failure due to invalid operator")
	}
	if len(result.Errors) != 1 {
		t.Error("Expected one error due to invalid condition syntax")
	}
}
