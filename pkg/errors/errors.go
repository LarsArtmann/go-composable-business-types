// Package errors provides centralized domain-specific errors for the library.
//
// This package consolidates all error definitions, making it easy to:
//
//   - Check error types with errors.Is()
//   - Extract structured data with errors.As()
//   - Find all error definitions in one place
//
// # Sentinel Errors
//
// Use errors.Is() to check for specific error conditions:
//
//	if errors.Is(err, errors.ErrInvalidEmail) {
//	    // handle invalid email
//	}
//
// # Structured Errors
//
// Some errors contain additional context. Use errors.As() to extract:
//
//	var parseErr *errors.UnmarshalError
//	if errors.As(err, &parseErr) {
//	    fmt.Println(parseErr.Type)    // The type that failed
//	    fmt.Println(parseErr.Input)    // The invalid input
//	}
//
// # Wrapping
//
// Use the Wrap functions to add context while preserving the error chain:
//
//	return errors.ErrBoundedStringMinLength
//	return errors.WrapMalformed(err, "id", input)
//	return errors.WrapInvalid(err, "email", value)
package errors

import (
	"errors"
	"fmt"
)

// =============================================================================
// Sentinel Errors - Use errors.Is() to check
// =============================================================================

// Email errors.
var (
	// ErrInvalidEmail is returned when an email address fails validation.
	ErrInvalidEmail = errors.New("invalid email address")

	// ErrEmailEmpty is returned when an email address is empty.
	ErrEmailEmpty = errors.New("email cannot be empty")
)

// URL errors.
var (
	// ErrInvalidURL is returned when a URL fails validation.
	ErrInvalidURL = errors.New("invalid URL")

	// ErrURLEmpty is returned when a URL is empty.
	ErrURLEmpty = errors.New("URL cannot be empty")

	// ErrURLScheme is returned when a URL has an invalid scheme.
	ErrURLScheme = errors.New("URL must use http or https scheme")

	// ErrURLHost is returned when a URL has no host.
	ErrURLHost = errors.New("URL must have a host")
)

// BoundedString errors.
var (
	// ErrBoundedStringMinLength is returned when a string is too short.
	ErrBoundedStringMinLength = errors.New("string below minimum length")

	// ErrBoundedStringMaxLength is returned when a string is too long.
	ErrBoundedStringMaxLength = errors.New("string exceeds maximum length")

	// ErrBoundedStringMinNegative is returned when minimum length is negative.
	ErrBoundedStringMinNegative = errors.New("minimum length cannot be negative")

	// ErrBoundedStringMaxLessThanMin is returned when max < min.
	ErrBoundedStringMaxLessThanMin = errors.New("maximum length cannot be less than minimum length")
)

// NanoID errors.
var (
	// ErrNanoIDEmpty is returned when a NanoID is empty.
	ErrNanoIDEmpty = errors.New("nanoid cannot be empty")

	// ErrNanoIDTooShort is returned when a NanoID is too short.
	ErrNanoIDTooShort = errors.New("nanoid below minimum length")

	// ErrNanoIDTooLong is returned when a NanoID is too long.
	ErrNanoIDTooLong = errors.New("nanoid exceeds maximum length")

	// ErrNanoIDInvalid is returned when a NanoID contains invalid characters.
	ErrNanoIDInvalid = errors.New("nanoid contains invalid characters")
)

// ID errors.
var (
	// ErrIDInvalid is returned when an ID value is invalid.
	ErrIDInvalid = errors.New("invalid ID value")

	// ErrIDTypeNotSupported is returned when an ID type is not supported.
	ErrIDTypeNotSupported = errors.New("unsupported ID type for this operation")

	// ErrIDInsufficientData is returned when there's not enough data to decode.
	ErrIDInsufficientData = errors.New("insufficient data for ID type")
)

// Importance errors.
var (
	// ErrImportanceExceedsMax is returned when an importance value exceeds 100.
	ErrImportanceExceedsMax = errors.New("importance exceeds maximum value of 100")

	// ErrImportanceEmpty is returned when an importance value is required but empty/zero.
	ErrImportanceEmpty = errors.New("importance cannot be empty")

	// ErrImportanceUnknownClassification is returned when a string cannot be parsed as importance.
	ErrImportanceUnknownClassification = errors.New("unknown importance classification")
)

// Tag errors.
var (
	// ErrTagEmpty is returned when a tag is empty.
	ErrTagEmpty = errors.New("tag cannot be empty")

	// ErrTagTooLong is returned when a tag exceeds the maximum length.
	ErrTagTooLong = errors.New("tag exceeds maximum length")

	// ErrTagInvalidChars is returned when a tag contains invalid characters.
	ErrTagInvalidChars = errors.New("tag contains invalid characters")

	// ErrProjectCoreNil is returned when a ProjectCore is nil.
	ErrProjectCoreNil = errors.New("projectcore is nil")

	// ErrProjectCoreNameRequired is returned when a ProjectCore has no name.
	ErrProjectCoreNameRequired = errors.New("projectcore name is required")

	// ErrProjectCorePathRequired is returned when a ProjectCore has no path.
	ErrProjectCorePathRequired = errors.New("projectcore path is required")
)

// Parse/Marshal errors - structured errors for detailed context.
var (
	// ErrMalformedInput is returned when input cannot be parsed.
	ErrMalformedInput = errors.New("malformed input")

	// ErrUnsupportedType is returned when a type is not supported.
	ErrUnsupportedType = errors.New("unsupported type")

	// ErrInvalidJSON is returned when JSON is malformed.
	ErrInvalidJSON = errors.New("invalid JSON")
)

// =============================================================================
// Structured Error Types - Use errors.As() to extract details
// =============================================================================

// ContextualError is an error with additional context about what caused it.
type ContextualError interface {
	error
	Wrapped() error
}

// wrappedError provides shared Unwrap/Wrapped implementation.
type wrappedError struct {
	Err error
}

func (w *wrappedError) Unwrap() error {
	return w.Err
}

func (w *wrappedError) Wrapped() error {
	return w.Err
}

// UnmarshalError represents a failure to parse/unmarshal data.
type UnmarshalError struct {
	Type  string // The type that failed (e.g., "JSON", "XML", "Text")
	Input string // The invalid input that caused the failure
	wrappedError
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("unmarshal %s: %s: %v", e.Type, e.Input, e.Err)
}

// ValidationError represents a validation failure.
type ValidationError struct {
	Field string // The field that failed validation
	Value any    // The invalid value
	Err   error  // The underlying error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s=%v: %v", e.Field, e.Value, e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// RangeError represents a value outside valid range.
type RangeError struct {
	Value      any  // The invalid value
	Min        any  // Minimum allowed value
	Max        any  // Maximum allowed value
	OutOfRange bool // true if above max, false if below min
}

func (e *RangeError) Error() string {
	if e.OutOfRange {
		return fmt.Sprintf("value %v exceeds maximum %v", e.Value, e.Max)
	}

	return fmt.Sprintf("value %v below minimum %v", e.Value, e.Min)
}

// ScanError represents a database scan failure.
type ScanError struct {
	SourceType string // The type we're scanning from
	TargetType string // The type we're scanning into
	wrappedError
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("cannot scan %s into %s: %v", e.SourceType, e.TargetType, e.Err)
}

// =============================================================================
// Error Wrapping Helpers
// =============================================================================

// WrapMalformed wraps an error with malformed input context.
func WrapMalformed(err error, typeName, input string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: malformed %s %q", ErrMalformedInput, typeName, input)
}

// WrapInvalid wraps an error with invalid value context.
func WrapInvalid(err error, field string, value any) error {
	if err == nil {
		return nil
	}

	return &ValidationError{
		Field: field,
		Value: value,
		Err:   err,
	}
}

// WrapRange wraps an error as a range error.
func WrapRange(value, minVal, maxVal any, outOfRange bool) error {
	return &RangeError{
		Value:      value,
		Min:        minVal,
		Max:        maxVal,
		OutOfRange: outOfRange,
	}
}

// WrapContextual wraps err with context using constructor.
func WrapContextual[T ContextualError](err error, constructor func(err error) T) error {
	if err == nil {
		return nil
	}

	return constructor(err)
}

// WrapScan wraps an error as a scan error.
func WrapScan(err error, sourceType, targetType string) error {
	return WrapContextual(err, func(e error) *ScanError {
		return &ScanError{
			SourceType:   sourceType,
			TargetType:   targetType,
			wrappedError: wrappedError{Err: e},
		}
	})
}

// WrapUnmarshal wraps an error as an unmarshal error.
func WrapUnmarshal(err error, typeName, input string) error {
	return WrapContextual(err, func(e error) *UnmarshalError {
		return &UnmarshalError{Type: typeName, Input: input, wrappedError: wrappedError{Err: e}}
	})
}

// =============================================================================
// Generic Error Extraction - Go 1.26+ errors.AsType helpers
// =============================================================================

// AsUnmarshalError extracts an UnmarshalError from the error chain.
func AsUnmarshalError(err error) (*UnmarshalError, bool) {
	return errors.AsType[*UnmarshalError](err)
}

// AsValidationError extracts a ValidationError from the error chain.
func AsValidationError(err error) (*ValidationError, bool) {
	return errors.AsType[*ValidationError](err)
}

// AsRangeError extracts a RangeError from the error chain.
func AsRangeError(err error) (*RangeError, bool) {
	return errors.AsType[*RangeError](err)
}

// AsScanError extracts a ScanError from the error chain.
func AsScanError(err error) (*ScanError, bool) {
	return errors.AsType[*ScanError](err)
}

// =============================================================================
// Error Checking Helpers
// =============================================================================

// IsInvalidEmail checks if the error is related to invalid email.
func IsInvalidEmail(err error) bool {
	return errors.Is(err, ErrInvalidEmail) || errors.Is(err, ErrEmailEmpty)
}

// IsInvalidURL checks if the error is related to invalid URL.
func IsInvalidURL(err error) bool {
	return IsOneOf(err, ErrInvalidURL, ErrURLEmpty, ErrURLScheme, ErrURLHost)
}

// IsBoundedStringError checks if the error is related to bounded string validation.
func IsBoundedStringError(err error) bool {
	return IsOneOf(err, ErrBoundedStringMinLength, ErrBoundedStringMaxLength,
		ErrBoundedStringMinNegative, ErrBoundedStringMaxLessThanMin)
}

// IsNanoIDError checks if the error is related to NanoID validation.
func IsNanoIDError(err error) bool {
	return IsOneOf(err, ErrNanoIDEmpty, ErrNanoIDTooShort, ErrNanoIDTooLong, ErrNanoIDInvalid)
}

// IsOneOf checks if the error matches any of the provided sentinel errors.
func IsOneOf(err error, sentinels ...error) bool {
	for _, sentinel := range sentinels {
		if errors.Is(err, sentinel) {
			return true
		}
	}

	return false
}

// IsIDError checks if the error is related to ID validation.
func IsIDError(err error) bool {
	return IsOneOf(err, ErrIDInvalid, ErrIDTypeNotSupported, ErrIDInsufficientData)
}

// IsImportanceError checks if the error is related to importance validation.
func IsImportanceError(err error) bool {
	return IsOneOf(
		err,
		ErrImportanceExceedsMax,
		ErrImportanceEmpty,
		ErrImportanceUnknownClassification,
	)
}

// IsTagError checks if the error is related to tag validation.
func IsTagError(err error) bool {
	return IsOneOf(err, ErrTagEmpty, ErrTagTooLong, ErrTagInvalidChars)
}

// IsProjectCoreError checks if the error is related to ProjectCore validation.
func IsProjectCoreError(err error) bool {
	return IsOneOf(err, ErrProjectCoreNil, ErrProjectCoreNameRequired, ErrProjectCorePathRequired)
}

// IsParseError checks if the error is a parse/marshal error.
func IsParseError(err error) bool {
	return IsOneOf(err, ErrMalformedInput, ErrUnsupportedType, ErrInvalidJSON)
}
