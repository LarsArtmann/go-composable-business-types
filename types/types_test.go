package types

import (
	"errors"
	"testing"

	pkgerrors "github.com/larsartmann/go-composable-business-types/pkg/errors"
	"github.com/larsartmann/go-composable-business-types/testutil"
)

func TestEmailParts(t *testing.T) {
	t.Parallel()

	email, _ := NewEmail("user@example.com")
	if email.LocalPart() != "user" {
		t.Errorf("expected local part 'user', got %s", email.LocalPart())
	}

	if email.Domain() != "example.com" {
		t.Errorf("expected domain 'example.com', got %s", email.Domain())
	}
}

func runParseValidationTest[T testutil.ParseTester](
	t *testing.T,
	typeName string,
	constructor func(string) (T, error),
	cases []struct {
		name    string
		input   string
		wantErr bool
	},
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

	runParseValidationTest(t, "Email", NewEmail, []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "test@example.com", false},
		{"valid with dots", "first.last@example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"empty", "", true},
		{"no at", "testexample.com", true},
		{"no domain", "test@", true},
		{"no local", "@example.com", true},
	})

	runParseValidationTest(t, "URL", NewURL, []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid https", "https://example.com", false},
		{"valid http", "http://example.com", false},
		{"valid with path", "https://example.com/path", false},
		{"empty", "", true},
		{"no scheme", "example.com", true},
		{"ftp not allowed", "ftp://example.com", true},
		{"no host", "https:///path", true},
	})
}

// testParseInvalidAndValid tests both invalid input (with specific error type) and valid input (with expected output).
func testParseInvalidAndValid[T any](
	t *testing.T,
	name string,
	constructor func(string) (T, error),
	cases []struct {
		input          string
		expectedErr    error
		expectedOutput string
	},
) {
	t.Helper()

	for _, c := range cases {
		t.Run(name+"/"+c.input, func(t *testing.T) {
			t.Parallel()

			got, err := constructor(c.input)
			if c.expectedErr != nil {
				if err == nil {
					t.Error("expected error")
				} else if !errors.Is(err, c.expectedErr) {
					t.Errorf("expected %v, got %v", c.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotStr := any(got).(interface{ String() string }).String(); gotStr != c.expectedOutput {
				t.Errorf("expected %s, got %s", c.expectedOutput, gotStr)
			}
		})
	}
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
	t.Parallel()

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
	t.Parallel()

	testutil.RunIsZeroTest(t, func() (URL, error) {
		return NewURL("https://example.com")
	})
}

func TestParseURL(t *testing.T) {
	t.Parallel()

	testParseInvalidAndValid(t, "Email", NewEmail, []struct {
		input          string
		expectedErr    error
		expectedOutput string
	}{
		{"invalid-email", pkgerrors.ErrInvalidEmail, ""},
		{"test@example.com", nil, "test@example.com"},
	})

	testParseInvalidAndValid(t, "URL", NewURL, []struct {
		input          string
		expectedErr    error
		expectedOutput string
	}{
		{"not-a-valid-url", pkgerrors.ErrInvalidURL, ""},
		{"https://example.com", nil, "https://example.com"},
	})
}
