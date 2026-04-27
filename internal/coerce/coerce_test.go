package coerce

import (
	"strings"
	"testing"
)

func TestRunBoolNormalises(t *testing.T) {
	vars := map[string]string{
		"ENABLED": "TRUE",
		"DEBUG":   "1",
		"VERBOSE": "false",
	}
	res := Run(vars, TypeBool, nil)
	if len(res.Failed) != 0 {
		t.Fatalf("expected no failures, got %v", res.Failed)
	}
	if res.Coerced["ENABLED"] != "true" {
		t.Errorf("ENABLED: want true, got %s", res.Coerced["ENABLED"])
	}
	if res.Coerced["DEBUG"] != "true" {
		t.Errorf("DEBUG: want true, got %s", res.Coerced["DEBUG"])
	}
	if res.Changed < 2 {
		t.Errorf("expected at least 2 changed, got %d", res.Changed)
	}
}

func TestRunIntNormalises(t *testing.T) {
	vars := map[string]string{
		"PORT":    " 8080 ",
		"TIMEOUT": "30",
	}
	res := Run(vars, TypeInt, nil)
	if len(res.Failed) != 0 {
		t.Fatalf("unexpected failures: %v", res.Failed)
	}
	if res.Coerced["PORT"] != "8080" {
		t.Errorf("PORT: want 8080, got %s", res.Coerced["PORT"])
	}
}

func TestRunFloatNormalises(t *testing.T) {
	vars := map[string]string{"RATIO": " 3.14000 "}
	res := Run(vars, TypeFloat, nil)
	if len(res.Failed) != 0 {
		t.Fatalf("unexpected failures: %v", res.Failed)
	}
	if res.Coerced["RATIO"] != "3.14" {
		t.Errorf("RATIO: want 3.14, got %s", res.Coerced["RATIO"])
	}
}

func TestRunFailsInvalidInt(t *testing.T) {
	vars := map[string]string{"PORT": "not-a-number"}
	res := Run(vars, TypeInt, nil)
	if _, ok := res.Failed["PORT"]; !ok {
		t.Error("expected PORT to be in Failed")
	}
	// original value preserved in Coerced
	if res.Coerced["PORT"] != "not-a-number" {
		t.Errorf("expected original value preserved, got %s", res.Coerced["PORT"])
	}
}

func TestRunFilterKeys(t *testing.T) {
	vars := map[string]string{
		"A": "TRUE",
		"B": "TRUE",
	}
	res := Run(vars, TypeBool, []string{"A"})
	if res.Coerced["A"] != "true" {
		t.Errorf("A should be coerced, got %s", res.Coerced["A"])
	}
	if res.Coerced["B"] != "TRUE" {
		t.Errorf("B should be unchanged, got %s", res.Coerced["B"])
	}
}

func TestFormatSummary(t *testing.T) {
	vars := map[string]string{"X": "bad", "Y": "1"}
	res := Run(vars, TypeInt, nil)
	out := Format(res)
	if !strings.Contains(out, "FAIL") {
		t.Errorf("expected FAIL in output, got: %s", out)
	}
	if !strings.Contains(out, "bad") {
		t.Errorf("expected original value in output, got: %s", out)
	}
}
