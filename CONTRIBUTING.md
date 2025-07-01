# Contributing Guide

Thank you for your interest in the Expr project! We welcome all forms of contributions, including but not limited to code, documentation, test cases, issue reports, and feature suggestions.

## 🤝 How to Contribute

### Reporting Issues

If you find a bug or have a feature suggestion, please:

1. Search existing [Issues](https://github.com/mredencom/expr/issues) to avoid duplicates
2. Create a new Issue with:
   - Detailed problem description
   - Steps to reproduce
   - Expected and actual behavior
   - Environment information (Go version, operating system, etc.)
   - If possible, provide minimal reproduction code

### Submitting Code

1. **Fork the project**
   ```bash
   # Click the Fork button on GitHub page
   git clone https://github.com/YOUR-USERNAME/expr.git
   cd expr
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Write code**
   - Follow project coding standards
   - Add necessary test cases
   - Ensure all tests pass
   - Update relevant documentation

4. **Commit changes**
   ```bash
   git add .
   git commit -m "Add: concise commit message"
   git push origin feature/your-feature-name
   ```

5. **Create Pull Request**
   - Describe your changes in detail
   - Link related Issues
   - Ensure CI checks pass

## 📝 Coding Standards

### Go Code Style

Follow standard Go coding conventions:

```go
// ✅ Correct function naming and comments
// CompileExpression compiles an expression into an executable program
func CompileExpression(expr string, options ...Option) (*Program, error) {
    if expr == "" {
        return nil, fmt.Errorf("expression cannot be empty")
    }
    
    // Implementation logic...
    return program, nil
}

// ✅ Correct struct definition
type CompileOptions struct {
    EnableOptimization bool   // Enable compilation optimization
    MaxIterations      int    // Maximum iterations
    Timeout           time.Duration // Timeout duration
}
```

### Error Handling

```go
// ✅ Correct error handling
func processData(data interface{}) (result interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic occurred while processing data: %v", r)
        }
    }()
    
    if data == nil {
        return nil, fmt.Errorf("data cannot be nil")
    }
    
    // Processing logic...
    return result, nil
}
```

### Comment Requirements

```go
// Package expr provides a high-performance expression evaluation engine
//
// This package supports modern syntax features including Lambda expressions,
// pipeline operations, and null safety.
// It uses a zero-reflection architecture for excellent execution performance.
//
// Basic usage:
//   result, err := expr.Eval("2 + 3 * 4", nil)
//
// Lambda expressions:
//   result, err := expr.Eval("users | filter(u => u.age > 18)", env)
package expr

// Eval compiles and executes an expression
//
// This function is suitable for simple one-time execution scenarios.
// For expressions that need to be executed multiple times,
// it's recommended to use Compile and Run functions for better performance.
//
// Parameters:
//   expression - The expression string to execute
//   env - Execution environment, can be map[string]interface{} or struct
//
// Returns:
//   interface{} - Expression execution result
//   error - Compilation or execution error
//
// Example:
//   result, err := Eval("name + ' is ' + toString(age)", map[string]interface{}{
//       "name": "Alice",
//       "age": 30,
//   })
func Eval(expression string, env interface{}) (interface{}, error) {
    // Implementation...
}
```

## 🧪 Testing Requirements

### Unit Tests

Each new feature must include corresponding tests:

```go
func TestEval(t *testing.T) {
    tests := []struct {
        name       string
        expression string
        env        interface{}
        expected   interface{}
        wantErr    bool
    }{
        {
            name:       "simple arithmetic",
            expression: "2 + 3 * 4",
            env:        nil,
            expected:   14,
            wantErr:    false,
        },
        {
            name:       "variable access",
            expression: "name + ' is ' + toString(age)",
            env:        map[string]interface{}{"name": "Alice", "age": 30},
            expected:   "Alice is 30",
            wantErr:    false,
        },
        {
            name:       "invalid expression",
            expression: "invalid expression +",
            env:        nil,
            expected:   nil,
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Eval(tt.expression, tt.env)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Benchmark Tests

Performance-sensitive code needs benchmark tests:

```go
func BenchmarkEval(b *testing.B) {
    expression := "users | filter(u => u.active) | map(u => u.name)"
    env := map[string]interface{}{
        "users": []map[string]interface{}{
            {"name": "Alice", "active": true},
            {"name": "Bob", "active": false},
        },
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = Eval(expression, env)
    }
}

func BenchmarkCompileAndRun(b *testing.B) {
    expression := "users | filter(u => u.active) | map(u => u.name)"
    program, _ := Compile(expression)
    env := map[string]interface{}{
        "users": []map[string]interface{}{
            {"name": "Alice", "active": true},
            {"name": "Bob", "active": false},
        },
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = Run(program, env)
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./pkg/vm

# Run tests with coverage
go test -cover ./...

# Run benchmark tests
go test -bench=. ./...

# Run race detector
go test -race ./...
```

## 📋 Commit Message Standards

Use clear commit message format:

```
Type: Brief description

Detailed description (if needed)

Related Issue: #123
```

### Commit Types

- `Add`: Add new feature
- `Fix`: Fix bug
- `Update`: Update existing feature
- `Remove`: Remove feature or code
- `Refactor`: Refactor code
- `Test`: Add or modify tests
- `Doc`: Documentation updates
- `Style`: Code style adjustments
- `Perf`: Performance optimization

### Example

```bash
git commit -m "Add: Lambda expression support

- Implement Lambda expression parsing and compilation
- Support single and multi-parameter Lambdas
- Add related test cases
- Update API documentation

Closes #45"
```

## 🏗️ Development Environment Setup

### Requirements

- Go 1.19+
- Git
- Recommended IDE: VS Code with Go extension

### Project Setup

```bash
# Clone project
git clone https://github.com/mredencom/expr.git
cd expr

# Install dependencies
go mod tidy

# Run tests to ensure environment is working
go test ./...

# Run example
go run example/simple_example/main.go
```

### Development Tools

Recommended development tools and configuration:

```json
// .vscode/settings.json
{
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "go.formatTool": "goimports",
    "go.useLanguageServer": true,
    "go.testFlags": ["-v"],
    "go.coverOnSave": true
}
```

## 📖 Documentation Contributions

### Documentation Types

- **API Documentation**: Detailed function and method descriptions
- **Tutorials**: Usage guides and examples
- **Best Practices**: Recommended usage patterns
- **Troubleshooting**: Common issues and solutions

### Documentation Format

Use Markdown format, following this structure:

```markdown
# Title

## Overview

Brief description...

## Usage

### Basic Usage

```go
// Code example
```

### Advanced Usage

```go
// Advanced example
```

## Notes

- Important note 1
- Important note 2

## References

- [Related link](url)
```

## 🎯 Contribution Areas

We especially welcome contributions in these areas:

### Core Features
- New language features
- Performance optimizations
- Error handling improvements
- Memory usage optimizations

### Built-in Functions
- New math functions
- String processing functions
- Date/time functions
- JSON processing functions

### Module System
- New built-in modules
- Module management features
- Module documentation

### Development Tools
- Debugger features
- Performance profiling tools
- Testing tools
- Code generation tools

### Documentation and Examples
- Usage tutorials
- Best practices guides
- Real project examples
- Performance optimization guides

## ⚡ Quick Start Contributing

### Simple Tasks (Suitable for Beginners)

1. **Improve Documentation**
   - Fix spelling errors
   - Add code examples
   - Translate documentation

2. **Add Test Cases**
   - Boundary condition tests
   - Error scenario tests
   - Performance tests

3. **Fix Small Bugs**
   - Look for issues marked "good first issue"
   - Fix typos
   - Improve error messages

### Medium Tasks

1. **Add Built-in Functions**
   - Implement new math functions
   - Add string processing functions
   - Create date/time functions

2. **Performance Optimization**
   - Optimize existing algorithms
   - Reduce memory allocations
   - Improve caching mechanisms

### Advanced Tasks

1. **New Language Features**
   - Implement new syntax
   - Add type system features
   - Extend module system

2. **Architecture Improvements**
   - Optimize virtual machine
   - Improve compiler
   - Refactor core components

## 📞 Contact Us

- **GitHub Issues**: Technical issues and feature requests
- **GitHub Discussions**: General discussion and Q&A
- **Email**: support@mredencom.com

## 🎉 Contributor Recognition

All contributors will be recognized in:

- Contributor list in README.md
- Special thanks in release notes
- Acknowledgments section in project documentation

## 📄 License

By contributing code, you agree that your contributions will be licensed under the same MIT license as the project.

---

**Thank you for supporting the Expr project!** 🙏

Every contribution, big or small, makes this project better. We look forward to building this excellent expression engine with you! 