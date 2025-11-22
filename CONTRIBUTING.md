# Contributing to uniwidth

Thank you for considering contributing to uniwidth! This document outlines the development workflow and guidelines.

## Git Workflow (Simple Main-Based)

This project uses a simple main-based workflow for development.

### Branch Structure

```
main                 # Production-ready code (tagged releases)
  ├─ feature/*       # New features
  └─ fix/*           # Bug fixes
```

### Branch Purposes

- **main**: Production-ready code. All releases are tagged here.
- **feature/\***: New features. Branch from `main`, merge back to `main`.
- **fix/\***: Bug fixes. Branch from `main`, merge back to `main`.

### Workflow Commands

#### Starting a New Feature

```bash
# Create feature branch from main
git checkout main
git pull origin main
git checkout -b feature/my-new-feature

# Work on your feature...
git add .
git commit -m "feat: add my new feature"

# When done, push and create PR
git push origin feature/my-new-feature
# Create PR on GitHub → main
```

#### Fixing a Bug

```bash
# Create fix branch from main
git checkout main
git pull origin main
git checkout -b fix/issue-123

# Fix the bug...
git add .
git commit -m "fix: resolve issue #123"

# Push and create PR
git push origin fix/issue-123
# Create PR on GitHub → main
```

#### Creating a Release

```bash
# Ensure main is clean and tests pass
bash scripts/pre-release-check.sh

# Create release tag
git tag -a v0.1.0 -m "Release v0.1.0: First Stable Release"

# Push tag
git push origin v0.1.0

# Create GitHub Release from tag
```

## Commit Message Guidelines

Follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, etc.)
- **refactor**: Code refactoring
- **test**: Adding or updating tests
- **chore**: Maintenance tasks (build, dependencies, etc.)
- **perf**: Performance improvements

### Examples

```bash
feat: add support for grapheme cluster width calculation
fix: correct width calculation for variation selectors
docs: update README with performance benchmarks
refactor: optimize ASCII fast path
test: add tests for regional indicator pairs
chore: update golangci-lint to v2.5
perf: improve CJK range check performance
```

## Code Quality Standards

### Before Committing

1. **Format code**:
   ```bash
   go fmt ./...
   ```

2. **Run linter**:
   ```bash
   golangci-lint run
   ```

3. **Run tests**:
   ```bash
   go test -v ./...
   ```

4. **Check coverage**:
   ```bash
   go test -cover ./...
   ```

5. **All-in-one** (recommended):
   ```bash
   bash scripts/pre-release-check.sh
   ```

### Pull Request Requirements

- [ ] Code is formatted (`go fmt ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] All tests pass (`go test -v ./...`)
- [ ] New code has tests (minimum 90% coverage)
- [ ] Documentation updated (if applicable)
- [ ] Commit messages follow conventions
- [ ] No sensitive data (credentials, tokens, etc.)
- [ ] Benchmarks run (for performance-critical changes)

## Development Setup

### Prerequisites

- Go 1.25 or later
- golangci-lint v2.5+

### Install Dependencies

```bash
# Clone repository
git clone https://github.com/unilibs/uniwidth.git
cd uniwidth

# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Running Tests

```bash
# Run all tests
go test -v

# Run with coverage
go test -cover

# Run with race detector (if GCC available)
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem

# Run fuzzing (Go 1.18+)
go test -fuzz=FuzzStringWidth -fuzztime=30s
```

### Running Linter

```bash
# Run linter
golangci-lint run

# Run with auto-fix
golangci-lint run --fix
```

## Project Structure

```
uniwidth/
├── .github/              # GitHub configuration and CI/CD
│   ├── workflows/        # GitHub Actions
│   └── CODEOWNERS        # Code ownership
├── .gitignore            # Git ignore rules
├── .golangci.yml         # Linter configuration
├── bench/                # Benchmark comparisons vs go-runewidth
├── cmd/                  # Command-line utilities
│   └── generate-tables/  # Unicode table generator
├── docs/                 # Documentation
│   └── ARCHITECTURE.md   # Technical deep dive
├── scripts/              # Automation scripts
│   └── pre-release-check.sh # Pre-release validation
├── uniwidth.go           # Core width calculation (PUBLIC API)
├── options.go            # Functional options API (PUBLIC)
├── tables_generated.go   # Generated Unicode width tables
├── *_test.go             # Test files
├── LICENSE               # MIT License
├── README.md             # Main documentation
├── CHANGELOG.md          # Version history
├── CONTRIBUTING.md       # This file
├── CODE_OF_CONDUCT.md    # Community guidelines
├── SECURITY.md           # Security policy
└── go.mod / go.sum       # Go module files
```

## Adding New Features

1. Check if issue exists, if not create one
2. Discuss approach in the issue
3. Create feature branch from `main`
4. Implement feature with tests
5. Update documentation
6. Run quality checks (`bash scripts/pre-release-check.sh`)
7. Create pull request to `main`
8. Wait for code review
9. Address feedback
10. Merge when approved

## Code Style Guidelines

### General Principles

- Follow Go conventions and idioms
- Write self-documenting code
- Add comments for complex logic (especially Unicode edge cases)
- Keep functions small and focused
- Use meaningful variable names
- Performance is critical - always benchmark changes

### Naming Conventions

- **Public types/functions**: `PascalCase` (e.g., `RuneWidth`, `StringWidth`)
- **Private types/functions**: `camelCase` (e.g., `isASCIIOnly`, `binarySearch`)
- **Constants**: `PascalCase` with context prefix (e.g., `EAWide`, `EANarrow`)
- **Test functions**: `Test*` (e.g., `TestRuneWidth_ASCII`)

### Error Handling

- Always check and handle errors
- Use descriptive error messages
- Return errors immediately, don't wrap unnecessarily
- Validate inputs before processing

### Testing

- Use table-driven tests when appropriate
- Test both success and error cases
- Test edge cases (variation selectors, regional indicators, combining marks)
- Test performance (benchmarks for critical paths)
- Target 90%+ coverage

## Performance Guidelines

### Critical Paths

uniwidth is a performance-critical library. Changes to these functions must be benchmarked:

1. **ASCII Fast Path** (`uniwidth.go:34-45`)
   - Target: O(1), zero allocations
   - Benchmark: `BenchmarkStringWidth_ASCII_*`

2. **CJK Fast Path** (`uniwidth.go:48-75`)
   - Target: O(1), zero allocations
   - Benchmark: `BenchmarkStringWidth_CJK_*`

3. **Emoji Fast Path** (`uniwidth.go:77-116`)
   - Target: O(1), zero allocations
   - Benchmark: `BenchmarkStringWidth_Mixed_*`

### Benchmarking

```bash
# Save baseline
go test -bench=. -benchmem -count=10 > old.txt

# Make changes...

# Compare
go test -bench=. -benchmem -count=10 > new.txt
benchstat old.txt new.txt
```

**Rule**: No performance regressions allowed. Any change that reduces performance must be discussed in the issue.

## Unicode Compliance

### Important Considerations

1. **Unicode 16.0 Standard**: Changes must comply with Unicode 16.0
2. **Variation Selectors**: U+FE0E (text) and U+FE0F (emoji) must be handled correctly
3. **Regional Indicators**: Flag emoji (U+1F1E6 - U+1F1FF) must count as width 2 total
4. **Combining Marks**: Zero-width combining characters must have width 0
5. **East Asian Width**: Ambiguous characters configurable via Options API

### Table Generation

```bash
# Only regenerate tables when updating Unicode version
go generate

# Or manually
go run cmd/generate-tables/main.go
```

**Important**: Only update tables for new Unicode versions. Current: Unicode 16.0

## Getting Help

- Check existing issues and discussions
- Review `docs/ARCHITECTURE.md` for technical details
- Ask questions in GitHub Issues
- Check performance benchmarks in `bench/`
- Read source code comments for implementation details

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to uniwidth!** 🎉

**Built by the Phoenix TUI Framework team** | **Powered by Go 1.25+**
