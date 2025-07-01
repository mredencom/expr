package tests

import (
	"testing"

	expr "github.com/mredencom/expr"
)

// TestSimplePipelineTypeMethod tests basic type method calls in pipeline
func TestSimplePipelineTypeMethod(t *testing.T) {
	t.Run("simple pipeline test", func(t *testing.T) {
		// First test: can we call a type method directly?
		t.Log("Testing direct type method call...")
		env := map[string]interface{}{
			"text": "hello",
		}

		// Test direct usage of upper function
		result1, err1 := expr.Eval(`upper(text)`, env)
		if err1 != nil {
			t.Logf("⚠️  Direct upper() call failed: %v", err1)
		} else {
			t.Logf("✅ Direct upper() call: %v", result1)
		}

		// Now test simple pipeline without type methods
		t.Log("Testing simple pipeline...")
		result2, err2 := expr.Eval(`["hello", "world"] | map(upper(#))`, env)
		if err2 != nil {
			t.Logf("⚠️  Simple pipeline failed: %v", err2)
		} else {
			t.Logf("✅ Simple pipeline: %v", result2)
		}

		// Now test type method in pipeline
		t.Log("Testing type method in pipeline...")
		result3, err3 := expr.Eval(`["hello", "world"] | map(#.upper())`, env)
		if err3 != nil {
			t.Logf("⚠️  Type method in pipeline failed: %v", err3)
		} else {
			t.Logf("✅ Type method in pipeline: %v", result3)
		}
	})
}
