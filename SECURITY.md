# Security Policy

## Supported Versions

uniwidth is currently in stable release. We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1.0 | :x:                |

Future stable releases (v1.0+) will follow semantic versioning with LTS support.

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability in uniwidth, please report it responsibly.

### How to Report

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, please report security issues by:

1. **Private Security Advisory** (preferred):
   https://github.com/unilibs/uniwidth/security/advisories/new

2. **Email** to maintainers:
   Create a private GitHub issue or contact via discussions

### What to Include

Please include the following information in your report:

- **Description** of the vulnerability
- **Steps to reproduce** the issue
- **Affected versions** (which versions are impacted)
- **Potential impact** (DoS, information disclosure, etc.)
- **Suggested fix** (if you have one)
- **Your contact information** (for follow-up questions)

### Response Timeline

- **Initial Response**: Within 48-72 hours
- **Triage & Assessment**: Within 1 week
- **Fix & Disclosure**: Coordinated with reporter

We aim to:
1. Acknowledge receipt within 72 hours
2. Provide an initial assessment within 1 week
3. Work with you on a coordinated disclosure timeline
4. Credit you in the security advisory (unless you prefer to remain anonymous)

## Security Considerations for Unicode Width Calculation

uniwidth processes Unicode strings, which introduces potential security considerations.

### 1. Input Validation

**Risk**: Malformed Unicode strings could cause unexpected behavior.

**Attack Vectors**:
- Invalid UTF-8 sequences
- Extremely long strings (resource exhaustion)
- Strings with unexpected Unicode characters

**Mitigation in Library**:
- ✅ All string processing uses Go's native UTF-8 handling
- ✅ Invalid runes are handled gracefully
- ✅ No buffer overflows (memory-safe Go code)
- ✅ Performance limits on string processing

**User Recommendations**:
```go
// ✅ GOOD - Validate input before processing
if !utf8.ValidString(userInput) {
    return errors.New("invalid UTF-8")
}

if len(userInput) > maxAllowedLength {
    return errors.New("string too long")
}

width := uniwidth.StringWidth(userInput)
```

### 2. Resource Exhaustion

**Risk**: Extremely long strings could consume excessive resources.

**Example Attack**:
```
Input: 1GB string of Unicode characters
Result: High CPU usage calculating width
```

**Mitigation**:
- Library processes strings efficiently (O(n) complexity)
- No memory allocations for ASCII-only strings
- Minimal allocations for Unicode strings (1 allocation for `[]rune` conversion)

**User Best Practices**:
```go
// Set reasonable limits on input strings
const maxInputLength = 1024 * 1024 // 1MB

if len(input) > maxInputLength {
    return errors.New("input too large")
}

width := uniwidth.StringWidth(input)
```

### 3. Unicode Normalization

**Risk**: Different Unicode representations of the same visual character.

**Example**:
```
"é" can be represented as:
- Single character: U+00E9 (LATIN SMALL LETTER E WITH ACUTE)
- Two characters: U+0065 + U+0301 (e + COMBINING ACUTE ACCENT)
```

**Status**: uniwidth calculates width based on codepoints, not normalized forms.

**User Responsibility**: If your application requires normalization, apply it before width calculation:

```go
import "golang.org/x/text/unicode/norm"

// Normalize to NFC form before width calculation
normalized := norm.NFC.String(input)
width := uniwidth.StringWidth(normalized)
```

### 4. Variation Selectors

**Risk**: Variation selectors (U+FE0E, U+FE0F) change character presentation.

**Example**:
```
"❤" (U+2764) = width 2 (emoji presentation)
"❤︎" (U+2764 + U+FE0E) = width 1 (text presentation)
```

**Mitigation**: uniwidth correctly handles variation selectors as of v0.1.0-beta.

### 5. Regional Indicator Pairs (Flags)

**Risk**: Flag emoji consist of two regional indicator characters.

**Example**:
```
🇺🇸 = U+1F1FA + U+1F1F8 = width 2 (not 4!)
```

**Mitigation**: uniwidth correctly handles regional indicator pairs as of v0.1.0-beta.

## Known Security Considerations

### 1. Performance with Unicode Strings

**Status**: Mitigated through tiered lookup strategy.

**Risk Level**: Low

**Description**: Width calculation for Unicode strings requires `[]rune` conversion (1 allocation).

**Mitigation**:
- ASCII-only strings use zero-allocation fast path
- Unicode strings: 1 allocation only
- Tiered lookup strategy (O(1) for 90-95% of cases)

### 2. Memory Safety

**Status**: Guaranteed by Go runtime.

**Risk Level**: Very Low

**Description**: Go is memory-safe. Buffer overflows and use-after-free bugs are not possible.

**Mitigation**:
- All memory managed by Go runtime
- No unsafe pointer operations
- No C dependencies

### 3. Dependency Security

uniwidth has **ZERO external dependencies** (pure standard library).

**Monitoring**:
- ✅ No third-party dependencies
- ✅ Only Go standard library
- ✅ No C code or CGO
- ✅ Fully auditable codebase

## Security Best Practices for Users

### Input Validation

Always validate user input before processing:

```go
import (
    "unicode/utf8"
    "github.com/unilibs/uniwidth"
)

func SafeStringWidth(input string) (int, error) {
    // Validate UTF-8
    if !utf8.ValidString(input) {
        return 0, errors.New("invalid UTF-8")
    }

    // Validate length
    if len(input) > 1024*1024 { // 1MB
        return 0, errors.New("string too long")
    }

    // Safe to process
    return uniwidth.StringWidth(input), nil
}
```

### Resource Limits

Set limits when processing untrusted input:

```go
// Check string length before processing
const maxDisplayWidth = 1000 // Maximum terminal width

width := uniwidth.StringWidth(userInput)
if width > maxDisplayWidth {
    return errors.New("display width exceeds limit")
}
```

### Error Handling

Always check for edge cases:

```go
// Handle empty strings
if len(input) == 0 {
    return 0
}

// Calculate width
width := uniwidth.StringWidth(input)

// Validate result
if width < 0 {
    return errors.New("unexpected negative width")
}
```

## Security Testing

### Current Testing

- ✅ Unit tests with edge cases (variation selectors, regional indicators)
- ✅ Fuzzing tests (`FuzzStringWidth`, `FuzzRuneWidth`)
- ✅ Conformance tests (Unicode 16.0 compliance)
- ✅ Race detector (0 data races)
- ✅ Linting with golangci-lint (34+ linters)

### Planned for v1.0

- 🔄 Extended fuzzing with go-fuzz
- 🔄 Static analysis with gosec
- 🔄 SAST/DAST scanning in CI

## Security Contact

- **GitHub Security Advisory**: https://github.com/unilibs/uniwidth/security/advisories/new
- **Public Issues** (for non-sensitive bugs): https://github.com/unilibs/uniwidth/issues
- **Discussions**: https://github.com/unilibs/uniwidth/discussions

## Bug Bounty Program

uniwidth does not currently have a bug bounty program. We rely on responsible disclosure from the security community.

If you report a valid security vulnerability:
- ✅ Public credit in security advisory (if desired)
- ✅ Acknowledgment in CHANGELOG
- ✅ Our gratitude and recognition in README
- ✅ Priority review and quick fix

---

**Thank you for helping keep uniwidth secure!** 🔒

*Security is a journey, not a destination. We continuously improve our security posture with each release.*
