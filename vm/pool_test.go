package vm

import (
	"sync"
	"testing"

	"github.com/mredencom/expr/types"
)

// TestNewValuePool tests pool creation
func TestNewValuePool(t *testing.T) {
	pool := NewValuePool()

	if pool == nil {
		t.Fatal("Expected pool to not be nil")
	}

	// Test that all pools are initialized
	if pool.intPool.New == nil {
		t.Error("Expected intPool to be initialized")
	}
	if pool.floatPool.New == nil {
		t.Error("Expected floatPool to be initialized")
	}
	if pool.stringPool.New == nil {
		t.Error("Expected stringPool to be initialized")
	}
	if pool.boolPool.New == nil {
		t.Error("Expected boolPool to be initialized")
	}
}

// TestValuePoolGetInt tests int value operations
func TestValuePoolGetInt(t *testing.T) {
	pool := NewValuePool()

	// Test getting int values
	val1 := pool.GetInt(42)
	val2 := pool.GetInt(-123)
	val3 := pool.GetInt(0)

	if val1 == nil {
		t.Fatal("Expected non-nil int value")
	}
	if val1.Value() != 42 {
		t.Errorf("Expected 42, got %d", val1.Value())
	}

	if val2.Value() != -123 {
		t.Errorf("Expected -123, got %d", val2.Value())
	}

	if val3.Value() != 0 {
		t.Errorf("Expected 0, got %d", val3.Value())
	}

	// Test that values have correct type
	typeInfo := val1.Type()
	if typeInfo.Kind != types.KindInt64 {
		t.Errorf("Expected KindInt64, got %v", typeInfo.Kind)
	}
	if typeInfo.Name != "int" {
		t.Errorf("Expected 'int', got %s", typeInfo.Name)
	}
}

// TestValuePoolPutInt tests returning int values to pool
func TestValuePoolPutInt(t *testing.T) {
	pool := NewValuePool()

	val := pool.GetInt(42)

	// Put should not panic (even though it's a no-op for immutable values)
	pool.PutInt(val)

	// Should be able to get values after put
	val2 := pool.GetInt(123)
	if val2 == nil {
		t.Error("Expected to be able to get int after put")
	}
	if val2.Value() != 123 {
		t.Errorf("Expected 123, got %d", val2.Value())
	}
}

// TestValuePoolGetFloat tests float value operations
func TestValuePoolGetFloat(t *testing.T) {
	pool := NewValuePool()

	// Test getting float values
	val1 := pool.GetFloat(3.14)
	val2 := pool.GetFloat(-2.71)
	val3 := pool.GetFloat(0.0)

	if val1 == nil {
		t.Fatal("Expected non-nil float value")
	}
	if val1.Value() != 3.14 {
		t.Errorf("Expected 3.14, got %f", val1.Value())
	}

	if val2.Value() != -2.71 {
		t.Errorf("Expected -2.71, got %f", val2.Value())
	}

	if val3.Value() != 0.0 {
		t.Errorf("Expected 0.0, got %f", val3.Value())
	}

	// Test that values have correct type
	typeInfo := val1.Type()
	if typeInfo.Kind != types.KindFloat64 {
		t.Errorf("Expected KindFloat64, got %v", typeInfo.Kind)
	}
	if typeInfo.Name != "float" {
		t.Errorf("Expected 'float', got %s", typeInfo.Name)
	}
}

// TestValuePoolPutFloat tests returning float values to pool
func TestValuePoolPutFloat(t *testing.T) {
	pool := NewValuePool()

	val := pool.GetFloat(3.14)

	// Put should not panic (even though it's a no-op for immutable values)
	pool.PutFloat(val)

	// Should be able to get values after put
	val2 := pool.GetFloat(2.71)
	if val2 == nil {
		t.Error("Expected to be able to get float after put")
	}
	if val2.Value() != 2.71 {
		t.Errorf("Expected 2.71, got %f", val2.Value())
	}
}

// TestValuePoolGetString tests string value operations
func TestValuePoolGetString(t *testing.T) {
	pool := NewValuePool()

	// Test getting string values
	val1 := pool.GetString("hello")
	val2 := pool.GetString("")
	val3 := pool.GetString("world with spaces")

	if val1 == nil {
		t.Fatal("Expected non-nil string value")
	}
	if val1.Value() != "hello" {
		t.Errorf("Expected 'hello', got %s", val1.Value())
	}

	if val2.Value() != "" {
		t.Errorf("Expected '', got %s", val2.Value())
	}

	if val3.Value() != "world with spaces" {
		t.Errorf("Expected 'world with spaces', got %s", val3.Value())
	}

	// Test that values have correct type
	typeInfo := val1.Type()
	if typeInfo.Kind != types.KindString {
		t.Errorf("Expected KindString, got %v", typeInfo.Kind)
	}
	if typeInfo.Name != "string" {
		t.Errorf("Expected 'string', got %s", typeInfo.Name)
	}
}

// TestValuePoolPutString tests returning string values to pool
func TestValuePoolPutString(t *testing.T) {
	pool := NewValuePool()

	val := pool.GetString("test")

	// Put should not panic (even though it's a no-op for immutable values)
	pool.PutString(val)

	// Should be able to get values after put
	val2 := pool.GetString("another")
	if val2 == nil {
		t.Error("Expected to be able to get string after put")
	}
	if val2.Value() != "another" {
		t.Errorf("Expected 'another', got %s", val2.Value())
	}
}

// TestValuePoolGetBool tests bool value operations
func TestValuePoolGetBool(t *testing.T) {
	pool := NewValuePool()

	// Test getting bool values
	val1 := pool.GetBool(true)
	val2 := pool.GetBool(false)

	if val1 == nil {
		t.Fatal("Expected non-nil bool value")
	}
	if val1.Value() != true {
		t.Errorf("Expected true, got %v", val1.Value())
	}

	if val2.Value() != false {
		t.Errorf("Expected false, got %v", val2.Value())
	}

	// Test that values have correct type
	typeInfo := val1.Type()
	if typeInfo.Kind != types.KindBool {
		t.Errorf("Expected KindBool, got %v", typeInfo.Kind)
	}
	if typeInfo.Name != "bool" {
		t.Errorf("Expected 'bool', got %s", typeInfo.Name)
	}
}

// TestValuePoolPutBool tests returning bool values to pool
func TestValuePoolPutBool(t *testing.T) {
	pool := NewValuePool()

	val := pool.GetBool(true)

	// Put should not panic (even though it's a no-op for immutable values)
	pool.PutBool(val)

	// Should be able to get values after put
	val2 := pool.GetBool(false)
	if val2 == nil {
		t.Error("Expected to be able to get bool after put")
	}
	if val2.Value() != false {
		t.Errorf("Expected false, got %v", val2.Value())
	}
}

// TestValuePoolGetStats tests statistics functionality
func TestValuePoolGetStats(t *testing.T) {
	pool := NewValuePool()

	// Get initial stats
	stats := pool.GetStats()

	// Check cache sizes (these should be pre-populated)
	if stats.IntCacheSize <= 0 {
		t.Errorf("Expected IntCacheSize > 0, got %d", stats.IntCacheSize)
	}
	if stats.FloatCacheSize <= 0 {
		t.Errorf("Expected FloatCacheSize > 0, got %d", stats.FloatCacheSize)
	}
	if stats.StringCacheSize <= 0 {
		t.Errorf("Expected StringCacheSize > 0, got %d", stats.StringCacheSize)
	}

	// Get some values and put them back
	intVal := pool.GetInt(42)
	floatVal := pool.GetFloat(3.14)
	stringVal := pool.GetString("test")
	boolVal := pool.GetBool(true)

	pool.PutInt(intVal)
	pool.PutFloat(floatVal)
	pool.PutString(stringVal)
	pool.PutBool(boolVal)

	// Cache sizes should still be populated after operations
	stats2 := pool.GetStats()
	if stats2.IntCacheSize <= 0 {
		t.Errorf("Expected IntCacheSize > 0 after operations, got %d", stats2.IntCacheSize)
	}
}

// TestValuePoolConcurrency tests concurrent access to the pool
func TestValuePoolConcurrency(t *testing.T) {
	pool := NewValuePool()

	var wg sync.WaitGroup
	numGoroutines := 100
	operationsPerGoroutine := 100

	// Test concurrent access to all pool types
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				// Test int pool
				intVal := pool.GetInt(int64(id*1000 + j))
				if intVal.Value() != int64(id*1000+j) {
					t.Errorf("Expected %d, got %d", id*1000+j, intVal.Value())
				}
				pool.PutInt(intVal)

				// Test float pool
				floatVal := pool.GetFloat(float64(id) + 0.5)
				if floatVal.Value() != float64(id)+0.5 {
					t.Errorf("Expected %f, got %f", float64(id)+0.5, floatVal.Value())
				}
				pool.PutFloat(floatVal)

				// Test string pool
				stringVal := pool.GetString("goroutine_test")
				if stringVal.Value() != "goroutine_test" {
					t.Errorf("Expected 'goroutine_test', got %s", stringVal.Value())
				}
				pool.PutString(stringVal)

				// Test bool pool
				boolVal := pool.GetBool(j%2 == 0)
				expected := j%2 == 0
				if boolVal.Value() != expected {
					t.Errorf("Expected %v, got %v", expected, boolVal.Value())
				}
				pool.PutBool(boolVal)
			}
		}(i)
	}

	wg.Wait()

	// Verify pool is still functional after concurrent operations
	testVal := pool.GetInt(999)
	if testVal == nil {
		t.Error("Expected pool to be functional after concurrent operations")
	}
	if testVal.Value() != 999 {
		t.Errorf("Expected 999, got %d", testVal.Value())
	}
}

// TestValuePoolMemoryBehavior tests that the pool behaves correctly with memory
func TestValuePoolMemoryBehavior(t *testing.T) {
	pool := NewValuePool()

	// Test that getting the same value multiple times works
	val1 := pool.GetInt(42)
	val2 := pool.GetInt(42)

	// Since we create new values each time, they should be different objects
	// but have the same value
	if val1.Value() != val2.Value() {
		t.Error("Expected same values")
	}

	// Test with different types
	intVal := pool.GetInt(123)
	floatVal := pool.GetFloat(123.0)
	stringVal := pool.GetString("123")
	boolVal := pool.GetBool(true)

	// Verify all values are distinct types
	if intVal.Type().Kind == floatVal.Type().Kind {
		t.Error("Expected different types for int and float")
	}
	if floatVal.Type().Kind == stringVal.Type().Kind {
		t.Error("Expected different types for float and string")
	}
	if stringVal.Type().Kind == boolVal.Type().Kind {
		t.Error("Expected different types for string and bool")
	}
}

// TestValuePoolEdgeCases tests edge cases and boundary conditions
func TestValuePoolEdgeCases(t *testing.T) {
	pool := NewValuePool()

	t.Run("LargeValues", func(t *testing.T) {
		// Test with large int
		largeInt := pool.GetInt(9223372036854775807) // Max int64
		if largeInt.Value() != 9223372036854775807 {
			t.Errorf("Expected max int64, got %d", largeInt.Value())
		}

		// Test with very small int
		smallInt := pool.GetInt(-9223372036854775808) // Min int64
		if smallInt.Value() != -9223372036854775808 {
			t.Errorf("Expected min int64, got %d", smallInt.Value())
		}

		// Test with large float
		largeFloat := pool.GetFloat(1.7976931348623157e+308) // Close to max float64
		if largeFloat.Value() != 1.7976931348623157e+308 {
			t.Errorf("Expected large float, got %f", largeFloat.Value())
		}
	})

	t.Run("EmptyAndSpecialStrings", func(t *testing.T) {
		// Test empty string
		emptyStr := pool.GetString("")
		if emptyStr.Value() != "" {
			t.Errorf("Expected empty string, got '%s'", emptyStr.Value())
		}

		// Test string with special characters
		specialStr := pool.GetString("Hello\nWorld\t!")
		if specialStr.Value() != "Hello\nWorld\t!" {
			t.Errorf("Expected special string, got '%s'", specialStr.Value())
		}

		// Test very long string
		longStr := make([]byte, 10000)
		for i := range longStr {
			longStr[i] = byte('A' + (i % 26))
		}
		longString := pool.GetString(string(longStr))
		if len(longString.Value()) != 10000 {
			t.Errorf("Expected length 10000, got %d", len(longString.Value()))
		}
	})

	t.Run("SpecialFloatValues", func(t *testing.T) {
		// Test zero
		zeroFloat := pool.GetFloat(0.0)
		if zeroFloat.Value() != 0.0 {
			t.Errorf("Expected 0.0, got %f", zeroFloat.Value())
		}

		// Test negative zero (should be treated as zero)
		negZeroFloat := pool.GetFloat(-0.0)
		if negZeroFloat.Value() != 0.0 {
			t.Errorf("Expected 0.0 for -0.0, got %f", negZeroFloat.Value())
		}

		// Test very small positive number
		smallFloat := pool.GetFloat(4.9406564584124654e-324) // Smallest positive float64
		if smallFloat.Value() != 4.9406564584124654e-324 {
			t.Errorf("Expected smallest float, got %f", smallFloat.Value())
		}
	})
}

// TestValuePoolTypeConsistency tests that pool returns consistent types
func TestValuePoolTypeConsistency(t *testing.T) {
	pool := NewValuePool()

	// Get multiple values of each type
	intVal1 := pool.GetInt(1)
	intVal2 := pool.GetInt(2)
	floatVal1 := pool.GetFloat(1.0)
	floatVal2 := pool.GetFloat(2.0)
	stringVal1 := pool.GetString("a")
	stringVal2 := pool.GetString("b")
	boolVal1 := pool.GetBool(true)
	boolVal2 := pool.GetBool(false)

	// Check that type info is consistent
	if intVal1.Type().Kind != intVal2.Type().Kind {
		t.Error("Expected consistent int type")
	}
	if floatVal1.Type().Kind != floatVal2.Type().Kind {
		t.Error("Expected consistent float type")
	}
	if stringVal1.Type().Kind != stringVal2.Type().Kind {
		t.Error("Expected consistent string type")
	}
	if boolVal1.Type().Kind != boolVal2.Type().Kind {
		t.Error("Expected consistent bool type")
	}

	// Check that type names are consistent
	if intVal1.Type().Name != intVal2.Type().Name {
		t.Error("Expected consistent int type name")
	}
	if floatVal1.Type().Name != floatVal2.Type().Name {
		t.Error("Expected consistent float type name")
	}
	if stringVal1.Type().Name != stringVal2.Type().Name {
		t.Error("Expected consistent string type name")
	}
	if boolVal1.Type().Name != boolVal2.Type().Name {
		t.Error("Expected consistent bool type name")
	}
}
