package types

import (
	"errors"
	"testing"

	pkgerrors "github.com/larsartmann/go-composable-business-types/pkg/errors"
	"github.com/larsartmann/go-composable-business-types/testutil"
)

type parseTestCase struct {
	name    string
	input   string
	wantErr bool
}

var emailValidationCases []parseTestCase

func init() {
	emailValidationCases = append(
		emailValidationCases,
		parseTestCase{name: "valid simple", input: "test@example.com"},
	)
	emailValidationCases = append(
		emailValidationCases,
		parseTestCase{name: "valid with dots", input: "first.last@example.com"},
	)
	emailValidationCases = append(
		emailValidationCases,
		parseTestCase{name: "valid with plus", input: "user+tag@example.com"},
	)
	emailValidationCases = append(
		emailValidationCases,
		parseTestCase{name: "empty", input: "", wantErr: true},
	)
	emailValidationCases = append(
		emailValidationCases,
		parseTestCase{name: "no at", input: "testexample.com", wantErr: true},
	)
	emailValidationCases = append(
		emailValidationCases,
		parseTestCase{name: "no domain", input: "test@", wantErr: true},
	)
	emailValidationCases = append(
		emailValidationCases,
		parseTestCase{name: "no local", input: "@example.com", wantErr: true},
	)
}

var urlValidationCases []parseTestCase

func init() {
	urlValidationCases = append(
		urlValidationCases,
		parseTestCase{name: "valid https", input: "https://example.com"},
	)
	urlValidationCases = append(
		urlValidationCases,
		parseTestCase{name: "valid http", input: "http://example.com"},
	)
	urlValidationCases = append(
		urlValidationCases,
		parseTestCase{name: "valid with path", input: "https://example.com/path"},
	)
	urlValidationCases = append(
		urlValidationCases,
		parseTestCase{name: "empty", input: "", wantErr: true},
	)
	urlValidationCases = append(
		urlValidationCases,
		parseTestCase{name: "no scheme", input: "example.com", wantErr: true},
	)
	urlValidationCases = append(
		urlValidationCases,
		parseTestCase{name: "ftp not allowed", input: "ftp://example.com", wantErr: true},
	)
	urlValidationCases = append(
		urlValidationCases,
		parseTestCase{name: "no host", input: "https:///path", wantErr: true},
	)
}

var validationTestCases = map[string][]parseTestCase{
	"Email": emailValidationCases,
	"URL":   urlValidationCases,
}

func TestEmailParts(t *testing.T) {
	t.Parallel()

	email, _ := NewEmail("user@example.com")
	testutil.RunPartsTest(t, email, []testutil.PartAccessor[Email]{
		{"local part", func(e Email) string { return e.LocalPart() }, "user"},
		{"domain", func(e Email) string { return e.Domain() }, "example.com"},
	})
}

func runParseValidationTest[T testutil.ParseTester](
	t *testing.T,
	typeName string,
	constructor func(string) (T, error),
	cases []parseTestCase,
) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testutil.RunParseTest(t, typeName, tc.input, constructor, tc.wantErr)
		})
	}
}

func TestValidation(t *testing.T) {
	t.Parallel()

	runParseValidationTest(t, "Email", NewEmail, validationTestCases["Email"])
	runParseValidationTest(t, "URL", NewURL, validationTestCases["URL"])
}

func TestEmailNormalize(t *testing.T) {
	t.Parallel()

	email, _ := NewEmail("User@Example.COM")

	normalized := email.Normalize()
	if normalized.String() != "User@example.com" {
		t.Errorf("expected User@example.com, got %s", normalized.String())
	}
}

func TestEmailIsZero(t *testing.T) {
	testutil.RunIsZeroTest(t, func() (Email, error) {
		return NewEmail("test@example.com")
	})
}

func TestURLParts(t *testing.T) {
	t.Parallel()

	url, _ := NewURL("https://example.com:8080/path/to/resource")
	if url.Scheme() != "https" {
		t.Errorf("expected scheme https, got %s", url.Scheme())
	}

	if url.Host() != "example.com:8080" {
		t.Errorf("expected host example.com:8080, got %s", url.Host())
	}

	if url.Path() != "/path/to/resource" {
		t.Errorf("expected path /path/to/resource, got %s", url.Path())
	}
}

func TestURLIsZero(t *testing.T) {
	testutil.RunIsZeroTest(t, func() (URL, error) {
		return NewURL("https://example.com")
	})
}

func TestParseURL(t *testing.T) {
	t.Parallel()

	testParseURLCases := []struct {
		name        string
		constructor func(string) (any, error)
		input       string
		expectedErr error
		expectedOut string
	}{
		{
			"Email/invalid",
			func(s string) (any, error) { return NewEmail(s) },
			"invalid-email",
			pkgerrors.ErrInvalidEmail,
			"",
		},
		{
			"Email/valid",
			func(s string) (any, error) { return NewEmail(s) },
			"test@example.com",
			nil,
			"test@example.com",
		},
		{
			"URL/invalid",
			func(s string) (any, error) { return NewURL(s) },
			"not-a-valid-url",
			pkgerrors.ErrInvalidURL,
			"",
		},
		{
			"URL/valid",
			func(s string) (any, error) { return NewURL(s) },
			"https://example.com",
			nil,
			"https://example.com",
		},
	}

	for _, tc := range testParseURLCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := tc.constructor(tc.input)
			if tc.expectedErr != nil {
				if err == nil {
					t.Error("expected error")
				} else if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected %v, got %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotStr := any(got).(interface{ String() string }).String()
			if gotStr != tc.expectedOut {
				t.Errorf("expected %s, got %s", tc.expectedOut, gotStr)
			}
		})
	}
}
