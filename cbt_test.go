package cbt

import (
	"encoding/json"
	"testing"
	"time"
)

func TestId(t *testing.T) {
	id := NewID[struct{}, string]("user-123")
	if id.Value() != "user-123" {
		t.Errorf("expected user-123, got %s", id.Value())
	}
	if id.IsZero() {
		t.Error("expected non-zero id")
	}
}

func TestIdInt(t *testing.T) {
	id := NewID[struct{}, int](42)
	if id.Value() != 42 {
		t.Errorf("expected 42, got %d", id.Value())
	}
}

func TestActor(t *testing.T) {
	chain := NewActorChain(UserActor(NewID[struct{}, string]("u-1"), "Alice")).
		Append(ServiceActor(NewID[struct{}, string]("api-gateway"), "API Gateway")).
		Append(ServiceActor(NewID[struct{}, string]("order-svc"), "Order Service"))

	if chain.IsEmpty() {
		t.Error("expected non-empty chain")
	}
	if chain.Origin().Kind != ActorKindUser {
		t.Error("expected origin to be user")
	}
	if chain.Origin().Name != "Alice" {
		t.Error("expected origin name to be Alice")
	}
	if chain.Current().Kind != ActorKindService {
		t.Error("expected current to be service")
	}
	if !chain.HasKind(ActorKindUser) {
		t.Error("expected chain to have user")
	}
	if len(chain.ByKind(ActorKindService)) != 2 {
		t.Error("expected 2 services in chain")
	}
}

func TestPercentage(t *testing.T) {
	p := NewPercentage(50)
	if p.Float64() != 0.5 {
		t.Errorf("expected 0.5, got %f", p.Float64())
	}

	clamped := NewPercentage(150)
	if clamped != 100 {
		t.Errorf("expected 100, got %d", clamped)
	}
}

func TestCents(t *testing.T) {
	c := NewCents(1099)
	if c.Int64() != 1099 {
		t.Errorf("expected 1099, got %d", c.Int64())
	}
	if c.Float64() != 10.99 {
		t.Errorf("expected 10.99, got %f", c.Float64())
	}
}

func TestCents_Arithmetic(t *testing.T) {
	a := NewCents(1000)
	b := NewCents(250)

	// Add
	sum := a.Add(b)
	if sum.Int64() != 1250 {
		t.Errorf("Add: expected 1250, got %d", sum.Int64())
	}

	// Sub
	diff := a.Sub(b)
	if diff.Int64() != 750 {
		t.Errorf("Sub: expected 750, got %d", diff.Int64())
	}

	// Mul
	prod := a.Mul(3)
	if prod.Int64() != 3000 {
		t.Errorf("Mul: expected 3000, got %d", prod.Int64())
	}

	// Div
	quot := a.Div(4)
	if quot.Int64() != 250 {
		t.Errorf("Div: expected 250, got %d", quot.Int64())
	}
}

func TestCents_SignOperations(t *testing.T) {
	pos := NewCents(100)
	neg := NewCents(-100)
	zero := NewCents(0)

	// Sign
	if pos.Sign() != 1 {
		t.Errorf("positive Sign: expected 1, got %d", pos.Sign())
	}
	if neg.Sign() != -1 {
		t.Errorf("negative Sign: expected -1, got %d", neg.Sign())
	}
	if zero.Sign() != 0 {
		t.Errorf("zero Sign: expected 0, got %d", zero.Sign())
	}

	// Abs
	if neg.Abs().Int64() != 100 {
		t.Errorf("Abs: expected 100, got %d", neg.Abs().Int64())
	}

	// Predicates
	if !pos.IsPositive() || pos.IsNegative() || pos.IsZero() {
		t.Error("IsPositive failed for positive value")
	}
	if !neg.IsNegative() || neg.IsPositive() || neg.IsZero() {
		t.Error("IsNegative failed for negative value")
	}
	if !zero.IsZero() || zero.IsPositive() || zero.IsNegative() {
		t.Error("IsZero failed for zero value")
	}
}

func TestBoundedString(t *testing.T) {
	bs, err := NewBoundedString(1, 10, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bs.String() != "hello" {
		t.Errorf("expected hello, got %s", bs.String())
	}
	if bs.Len() != 5 {
		t.Errorf("expected length 5, got %d", bs.Len())
	}
	if bs.MinLen() != 1 || bs.MaxLen() != 10 {
		t.Errorf("expected bounds [1,10], got [%d,%d]", bs.MinLen(), bs.MaxLen())
	}
}

func TestBoundedStringValidation(t *testing.T) {
	_, err := NewBoundedString(1, 10, "")
	if err == nil {
		t.Error("expected error for empty string with min 1")
	}

	_, err = NewBoundedString(1, 5, "too long string")
	if err == nil {
		t.Error("expected error for string exceeding max")
	}

	_, err = NewBoundedString(5, 3, "test")
	if err == nil {
		t.Error("expected error for invalid bounds")
	}
}

func TestBoundedStringHelperMethods(t *testing.T) {
	bs := MustBoundedString(2, 10, "test")
	if bs.IsEmpty() {
		t.Error("expected non-empty")
	}

	short := MustBoundedString(2, 10, "ab")
	if !short.IsMinLength() {
		t.Error("expected to be exactly at min length")
	}

	maxed := MustBoundedString(2, 4, "abcd")
	if !maxed.IsMaxLength() {
		t.Error("expected to be at max length")
	}
}

func TestBoundedStringOf(t *testing.T) {
	NewName := BoundedStringOf(1, 50)
	name, err := NewName("John")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name.String() != "John" {
		t.Errorf("expected John, got %s", name.String())
	}
}

func TestTrimmedBoundedString(t *testing.T) {
	bs, err := TrimmedBoundedString(1, 10, "  hello  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bs.String() != "hello" {
		t.Errorf("expected trimmed 'hello', got '%s'", bs.String())
	}
}

func TestNonEmptyString(t *testing.T) {
	_, err := NonEmptyString(10, "")
	if err == nil {
		t.Error("expected error for empty string")
	}

	bs, err := NonEmptyString(10, "x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bs.String() != "x" {
		t.Errorf("expected x, got %s", bs.String())
	}
}

// Branded ID tests

type (
	UserBrand  struct{}
	OrderBrand struct{}
)

type (
	UserID  = ID[UserBrand, string]
	OrderID = ID[OrderBrand, int64]
)

func TestBrandedID_String(t *testing.T) {
	uid := NewID[UserBrand, string]("user-123")
	if uid.Value() != "user-123" {
		t.Errorf("expected user-123, got %s", uid.Value())
	}
	if uid.String() != "user-123" {
		t.Errorf("expected String() = user-123, got %s", uid.String())
	}
	if uid.GoString() != "user-123" {
		t.Errorf("expected GoString() = user-123, got %s", uid.GoString())
	}
}

func TestBrandedID_Int64(t *testing.T) {
	oid := NewID[OrderBrand, int64](42)
	if oid.Value() != 42 {
		t.Errorf("expected 42, got %d", oid.Value())
	}
	if oid.String() != "42" {
		t.Errorf("expected String() = 42, got %s", oid.String())
	}
}

func TestBrandedID_IsZero(t *testing.T) {
	uid := NewID[UserBrand, string]("user-123")
	if uid.IsZero() {
		t.Error("expected non-zero ID to return false")
	}

	var zeroUserID UserID
	if !zeroUserID.IsZero() {
		t.Error("expected zero ID to return true")
	}

	var zeroOrderID OrderID
	if !zeroOrderID.IsZero() {
		t.Error("expected zero int64 ID to return true")
	}
}

func TestBrandedID_JSON_String(t *testing.T) {
	uid := NewID[UserBrand, string]("user-abc")

	// Marshal
	data, err := json.Marshal(uid)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != `"user-abc"` {
		t.Errorf("expected JSON \"user-abc\", got %s", string(data))
	}

	// Unmarshal
	var parsed UserID
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if parsed.Value() != "user-abc" {
		t.Errorf("expected user-abc, got %s", parsed.Value())
	}
}

func TestBrandedID_JSON_Zero(t *testing.T) {
	var uid UserID

	// Zero value should marshal to null
	data, err := json.Marshal(uid)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected JSON null for zero value, got %s", string(data))
	}
}

func TestBrandedID_JSON_EmptyString(t *testing.T) {
	data := []byte(`""`)
	var uid UserID
	if err := json.Unmarshal(data, &uid); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if uid.Value() != "" {
		t.Errorf("expected empty string, got %s", uid.Value())
	}
}

func TestBrandedID_UnbrandedID(t *testing.T) {
	// Unbranded ID using struct{} as brand
	id := NewID[struct{}, string]("test-id")
	if id.Value() != "test-id" {
		t.Errorf("expected test-id, got %s", id.Value())
	}

	// Should be comparable
	id2 := NewID[struct{}, string]("test-id")
	if id != id2 {
		t.Error("expected equal IDs to be equal")
	}
}

func TestBrandedID_JSON_Int64_Marshal(t *testing.T) {
	// int64 IDs serialize as strings (by design)
	oid := NewID[OrderBrand, int64](42)

	data, err := json.Marshal(oid)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	// int64 serializes as string "42" (not number 42)
	if string(data) != `"42"` {
		t.Errorf("expected JSON \"42\", got %s", string(data))
	}
}

func TestBrandedID_JSON_Int64_Unmarshal_Error(t *testing.T) {
	// Unmarshaling into int64 ID should return error (documented limitation)
	data := []byte(`"42"`)
	var oid OrderID
	err := json.Unmarshal(data, &oid)
	if err == nil {
		t.Error("expected error when unmarshaling string into int64 ID")
	}
}

func TestBrandedID_Comparability(t *testing.T) {
	// IDs with same value should be equal
	uid1 := NewID[UserBrand, string]("user-123")
	uid2 := NewID[UserBrand, string]("user-123")
	if uid1 != uid2 {
		t.Error("expected equal IDs to be equal")
	}

	// IDs with different values should not be equal
	uid3 := NewID[UserBrand, string]("user-456")
	if uid1 == uid3 {
		t.Error("expected different IDs to not be equal")
	}
}

func TestBrandedID_JSON_Null_Unmarshal(t *testing.T) {
	data := []byte("null")
	var uid UserID
	if err := json.Unmarshal(data, &uid); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if !uid.IsZero() {
		t.Error("expected null to unmarshal to zero value")
	}
}

func TestEmail_Validation(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
	}{
		{"test@example.com", false},
		{"user.name@domain.org", false},
		{"user+tag@example.co.uk", false},
		{"", true},
		{"not-an-email", true},
		{"@example.com", true},
		{"user@", true},
		{"user@.com", true},
		{"user@example", true},
		{"user @example.com", true},
		{"user@exam ple.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			_, err := NewEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestEmail_WithDisplayName(t *testing.T) {
	// mail.ParseAddress accepts display names - we extract just the address
	email, err := NewEmail("John Doe <john@example.com>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if email.String() != "john@example.com" {
		t.Errorf("expected john@example.com, got %s", email.String())
	}
}

func TestEmail_MustParseEmail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid email")
		}
	}()
	MustParseEmail("invalid")
}

func TestEmail_Helpers(t *testing.T) {
	email := MustParseEmail("user.name@example.org")

	if email.IsZero() {
		t.Error("expected IsZero to be false for valid email")
	}
	if email.LocalPart() != "user.name" {
		t.Errorf("expected LocalPart 'user.name', got %q", email.LocalPart())
	}
	if email.Domain() != "example.org" {
		t.Errorf("expected Domain 'example.org', got %q", email.Domain())
	}

	// Test empty email
	var empty Email
	if !empty.IsZero() {
		t.Error("expected IsZero to be true for empty email")
	}
	if empty.LocalPart() != "" {
		t.Errorf("expected empty LocalPart, got %q", empty.LocalPart())
	}
	if empty.Domain() != "" {
		t.Errorf("expected empty Domain, got %q", empty.Domain())
	}
}

func TestURL_Validation(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"https://example.com", false},
		{"http://example.com", false},
		{"https://example.com/path", false},
		{"https://example.com/path?query=1", false},
		{"https://sub.example.com:8080/path", false},
		{"", true},
		{"not-a-url", true},
		{"ftp://example.com", true}, // only http/https allowed
		{"//example.com", true},     // no scheme
		{"https://", true},          // no host
		{"http://", true},           // no host
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			_, err := NewURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestURL_MustParseURL(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid URL")
		}
	}()
	MustParseURL("invalid")
}

func TestURL_Helpers(t *testing.T) {
	url := MustParseURL("https://example.com:8080/api/v1/users?limit=10")

	if url.IsZero() {
		t.Error("expected IsZero to be false for valid URL")
	}
	if url.Scheme() != "https" {
		t.Errorf("expected Scheme 'https', got %q", url.Scheme())
	}
	if url.Host() != "example.com:8080" {
		t.Errorf("expected Host 'example.com:8080', got %q", url.Host())
	}
	if url.Path() != "/api/v1/users" {
		t.Errorf("expected Path '/api/v1/users', got %q", url.Path())
	}

	parsed, err := url.Parse()
	if err != nil {
		t.Errorf("Parse returned unexpected error: %v", err)
	}
	if parsed.Query().Get("limit") != "10" {
		t.Errorf("expected query param limit=10, got %s", parsed.Query().Get("limit"))
	}

	// Test empty URL
	var empty URL
	if !empty.IsZero() {
		t.Error("expected IsZero to be true for empty URL")
	}
	if empty.Scheme() != "" {
		t.Errorf("expected empty Scheme, got %q", empty.Scheme())
	}
	if empty.Host() != "" {
		t.Errorf("expected empty Host, got %q", empty.Host())
	}
	if empty.Path() != "" {
		t.Errorf("expected empty Path, got %q", empty.Path())
	}
}

func TestLocale_Parse(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"en-US", "en-US", false},
		{"en_US", "en-US", false},
		{"de-DE", "de-DE", false},
		{"fr-FR", "fr-FR", false},
		{"invalid-locale", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			locale, err := ParseLocale(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLocale(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && locale.String() != tt.want {
				t.Errorf("ParseLocale(%q) = %q, want %q", tt.input, locale.String(), tt.want)
			}
		})
	}
}

func TestLocale_Constants(t *testing.T) {
	if LocaleEnUS.String() != "en-US" {
		t.Errorf("LocaleEnUS: expected en-US, got %s", LocaleEnUS.String())
	}
	if LocaleEnGB.String() != "en-GB" {
		t.Errorf("LocaleEnGB: expected en-GB, got %s", LocaleEnGB.String())
	}
	if LocaleDeDE.String() != "de-DE" {
		t.Errorf("LocaleDeDE: expected de-DE, got %s", LocaleDeDE.String())
	}
}

func TestLocale_Methods(t *testing.T) {
	locale := LocaleEnUS

	if locale.Base() != "en" {
		t.Errorf("Base: expected en, got %s", locale.Base())
	}
	if locale.Region() != "US" {
		t.Errorf("Region: expected US, got %s", locale.Region())
	}
	if locale.IsZero() {
		t.Error("IsZero: expected false for valid locale")
	}

	var zeroLocale Locale
	if !zeroLocale.IsZero() {
		t.Error("IsZero: expected true for zero locale")
	}
}

func TestLocale_MarshalText(t *testing.T) {
	data, err := LocaleEnUS.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText error: %v", err)
	}
	if string(data) != "en-US" {
		t.Errorf("MarshalText: expected en-US, got %s", string(data))
	}

	// Test JSON serialization
	type wrapper struct {
		Locale Locale `json:"locale"`
	}
	w := wrapper{Locale: LocaleEnUS}
	jsonData, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("JSON marshal error: %v", err)
	}
	expected := `{"locale":"en-US"}`
	if string(jsonData) != expected {
		t.Errorf("JSON: expected %s, got %s", expected, string(jsonData))
	}

	// Test JSON deserialization
	var parsed wrapper
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}
	if parsed.Locale.String() != "en-US" {
		t.Errorf("JSON unmarshal: expected en-US, got %s", parsed.Locale.String())
	}
}

func TestLocale_MustParseLocale(t *testing.T) {
	locale := MustParseLocale("de-DE")
	if locale.String() != "de-DE" {
		t.Errorf("MustParseLocale: expected de-DE, got %s", locale.String())
	}
}

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name       string
		amount     string
		currency   string
		wantErr    bool
		wantString string
	}{
		{"valid USD", "10.99", "USD", false, "10.99 USD"},
		{"valid EUR", "100.00", "EUR", false, "100.00 EUR"},
		{"valid JPY (no decimals)", "1000", "JPY", false, "1000 JPY"},
		{"zero amount", "0.00", "USD", false, "0.00 USD"},
		{"negative amount", "-50.25", "USD", false, "-50.25 USD"},
		{"invalid currency", "10.99", "INVALID", true, ""},
		{"invalid amount", "abc", "USD", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoney(tt.amount, tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && money.String() != tt.wantString {
				t.Errorf("NewMoney().String() = %s, want %s", money.String(), tt.wantString)
			}
		})
	}
}

func TestNewMoneyFromCents(t *testing.T) {
	tests := []struct {
		name       string
		cents      int64
		currency   string
		wantErr    bool
		wantString string
	}{
		{"valid USD cents", 1099, "USD", false, "10.99 USD"},
		{"valid EUR cents", 10000, "EUR", false, "100.00 EUR"},
		{"valid JPY (no decimals)", 1000, "JPY", false, "1000 JPY"},
		{"zero cents", 0, "USD", false, "0.00 USD"},
		{"negative cents", -5025, "USD", false, "-50.25 USD"},
		{"invalid currency", 1099, "INVALID", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoneyFromCents(tt.cents, tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMoneyFromCents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && money.String() != tt.wantString {
				t.Errorf("NewMoneyFromCents().String() = %s, want %s", money.String(), tt.wantString)
			}
		})
	}
}

func TestIsValidCurrency(t *testing.T) {
	tests := []struct {
		code string
		want bool
	}{
		{"USD", true},
		{"EUR", true},
		{"GBP", true},
		{"JPY", true},
		{"CHF", true},
		{"INVALID", false},
		{"usd", false}, // case sensitive
		{"US", false},  // not a currency code
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if got := IsValidCurrency(tt.code); got != tt.want {
				t.Errorf("IsValidCurrency(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestCurrencyDigits(t *testing.T) {
	tests := []struct {
		code      string
		want      uint8
		wantFound bool
	}{
		{"USD", 2, true},
		{"EUR", 2, true},
		{"JPY", 0, true},
		{"INVALID", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, found := CurrencyDigits(tt.code)
			if got != tt.want || found != tt.wantFound {
				t.Errorf("CurrencyDigits(%q) = (%d, %v), want (%d, %v)", tt.code, got, found, tt.want, tt.wantFound)
			}
		})
	}
}

func TestCurrencySymbol(t *testing.T) {
	tests := []struct {
		code      string
		locale    string
		wantFound bool
	}{
		{"USD", "en-US", true},
		{"EUR", "de-DE", true},
		{"JPY", "ja-JP", true},
		{"INVALID", "en-US", false},
	}

	for _, tt := range tests {
		t.Run(tt.code+"_"+tt.locale, func(t *testing.T) {
			_, found := CurrencySymbol(tt.code, tt.locale)
			if found != tt.wantFound {
				t.Errorf("CurrencySymbol(%q, %q) found = %v, want %v", tt.code, tt.locale, found, tt.wantFound)
			}
		})
	}
}

func TestCurrencySymbolForLocale(t *testing.T) {
	tests := []struct {
		code      string
		locale    Locale
		wantFound bool
	}{
		{"USD", LocaleEnUS, true},
		{"EUR", LocaleDeDE, true},
		{"JPY", LocaleJaJP, true},
	}

	for _, tt := range tests {
		t.Run(tt.code+"_"+tt.locale.String(), func(t *testing.T) {
			_, found := CurrencySymbolForLocale(tt.code, tt.locale)
			if found != tt.wantFound {
				t.Errorf("CurrencySymbolForLocale(%q, %v) found = %v, want %v", tt.code, tt.locale, found, tt.wantFound)
			}
		})
	}
}

func TestAllCurrencyCodes(t *testing.T) {
	codes := AllCurrencyCodes()
	if len(codes) == 0 {
		t.Error("AllCurrencyCodes() returned empty slice")
	}

	// Verify some common currencies are present
	foundUSD := false
	foundEUR := false
	foundJPY := false
	for _, code := range codes {
		switch code {
		case "USD":
			foundUSD = true
		case "EUR":
			foundEUR = true
		case "JPY":
			foundJPY = true
		}
	}
	if !foundUSD {
		t.Error("AllCurrencyCodes() missing USD")
	}
	if !foundEUR {
		t.Error("AllCurrencyCodes() missing EUR")
	}
	if !foundJPY {
		t.Error("AllCurrencyCodes() missing JPY")
	}
}

func TestFormatMoney(t *testing.T) {
	money, _ := NewMoney("10.99", "USD")

	tests := []struct {
		name     string
		locale   string
		contains string // locale-dependent formatting, just check currency code present
	}{
		{"en-US", "en-US", "USD"},
		{"de-DE", "de-DE", "USD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatMoney(money, tt.locale)
			if formatted == "" {
				t.Error("FormatMoney() returned empty string")
			}
			// Verify it contains the amount
			if formatted != "" && !containsAnyDigit(formatted) {
				t.Errorf("FormatMoney() = %q, expected to contain digits", formatted)
			}
		})
	}
}

func TestFormatMoneyForLocale(t *testing.T) {
	money, _ := NewMoney("10.99", "USD")

	tests := []struct {
		name   string
		locale Locale
	}{
		{"en-US", LocaleEnUS},
		{"de-DE", LocaleDeDE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatMoneyForLocale(money, tt.locale)
			if formatted == "" {
				t.Error("FormatMoneyForLocale() returned empty string")
			}
		})
	}
}

func TestNewMoneyFormatter(t *testing.T) {
	formatter := NewMoneyFormatter("en-US")
	if formatter == nil {
		t.Error("NewMoneyFormatter() returned nil")
	}

	money, _ := NewMoney("10.99", "USD")
	formatted := formatter.Format(money)
	if formatted == "" {
		t.Error("MoneyFormatter.Format() returned empty string")
	}
}

func TestNewMoneyFormatterForLocale(t *testing.T) {
	formatter := NewMoneyFormatterForLocale(LocaleEnUS)
	if formatter == nil {
		t.Error("NewMoneyFormatterForLocale() returned nil")
	}

	money, _ := NewMoney("10.99", "USD")
	formatted := formatter.Format(money)
	if formatted == "" {
		t.Error("MoneyFormatter.Format() returned empty string")
	}
}

// containsAnyDigit checks if string contains any digit
func containsAnyDigit(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Duration
		expected time.Duration
	}{
		{"seconds", 5 * time.Second, 5 * time.Second},
		{"minutes", 2 * time.Minute, 2 * time.Minute},
		{"hours", 3 * time.Hour, 3 * time.Hour},
		{"zero", 0, 0},
		{"negative", -10 * time.Second, -10 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDuration(tt.input)
			if d.Duration != tt.expected {
				t.Errorf("NewDuration(%v).Duration = %v, want %v", tt.input, d.Duration, tt.expected)
			}
		})
	}
}

func TestURLString(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid http", "http://example.com", false},
		{"valid https", "https://example.com/path", false},
		{"valid with query", "https://example.com/path?query=1", false},
		{"invalid scheme", "ftp://example.com", true},
		{"missing host", "http://", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if u.String() != tt.url {
					t.Errorf("URL.String() = %q, want %q", u.String(), tt.url)
				}
			}
		})
	}
}

func TestMustParseURL(t *testing.T) {
	u := MustParseURL("https://example.com")
	if u.String() != "https://example.com" {
		t.Errorf("MustParseURL: expected https://example.com, got %s", u.String())
	}
}

func TestMustParseURL_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParseURL with invalid URL should panic")
		}
	}()
	MustParseURL("not-a-valid-url")
}

func TestNanoIdMustParseNanoId(t *testing.T) {
	id := MustParseNanoId("VwSt1Xx5")
	if id.String() != "VwSt1Xx5" {
		t.Errorf("MustParseNanoId: expected VwSt1Xx5, got %s", id.String())
	}
}

func TestNanoIdMustParseNanoId_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParseNanoId with invalid ID should panic")
		}
	}()
	MustParseNanoId("invalid!@#")
}

func TestNanoIdGoString(t *testing.T) {
	id := MustParseNanoId("VwSt1Xx5")
	if id.GoString() != "VwSt1Xx5" {
		t.Errorf("GoString: expected VwSt1Xx5, got %s", id.GoString())
	}
}

func TestBitemporal_IsCurrentlyValid(t *testing.T) {
	now := Now()

	// Fact that started in the past, no end date - should be valid
	pastOngoing := NewBitemporal(now)
	if !pastOngoing.IsCurrentlyValid() {
		t.Error("past ongoing fact should be currently valid")
	}

	// Fact with explicit end date in the future - should be valid
	futureEnd := NewBitemporalWithRange(now, NewTimestamp(now.Time.Add(24*time.Hour)), now)
	if !futureEnd.IsCurrentlyValid() {
		t.Error("fact ending in future should be currently valid")
	}

	// Fact that ended in the past - should not be valid
	pastEnd := NewBitemporalWithRange(
		NewTimestamp(now.Time.Add(-2*time.Hour)),
		NewTimestamp(now.Time.Add(-1*time.Hour)),
		now,
	)
	if pastEnd.IsCurrentlyValid() {
		t.Error("fact that ended in past should not be currently valid")
	}
}

func TestBitemporal_WithValidUntil(t *testing.T) {
	now := Now()
	b := NewBitemporal(now)

	endTime := NewTimestamp(now.Time.Add(24 * time.Hour))
	updated := b.WithValidUntil(endTime)

	if updated.ValidUntil() != endTime {
		t.Errorf("WithValidUntil: expected %v, got %v", endTime, updated.ValidUntil())
	}

	// Original should be unchanged
	if !b.ValidUntil().IsZero() {
		t.Error("original Bitemporal should still have zero ValidUntil")
	}
}

func TestContext_WithSource(t *testing.T) {
	ctx := NewContext("original-source")
	updated := ctx.WithSource("new-source")

	if updated.Source() != "new-source" {
		t.Errorf("WithSource: expected new-source, got %s", updated.Source())
	}

	// Original should be unchanged
	if ctx.Source() != "original-source" {
		t.Error("original Context should still have original source")
	}
}

func TestContext_WithTags(t *testing.T) {
	ctx := NewContext("test-service").
		WithTag("key1", "value1")

	updated := ctx.WithTags(map[string]string{
		"key2": "value2",
		"key3": "value3",
	})

	tags := updated.Tags()
	if tags["key1"] != "value1" {
		t.Errorf("WithTags: expected key1=value1, got %s", tags["key1"])
	}
	if tags["key2"] != "value2" {
		t.Errorf("WithTags: expected key2=value2, got %s", tags["key2"])
	}
	if tags["key3"] != "value3" {
		t.Errorf("WithTags: expected key3=value3, got %s", tags["key3"])
	}

	// Original should be unchanged
	originalTags := ctx.Tags()
	if _, exists := originalTags["key2"]; exists {
		t.Error("original Context should not have key2")
	}
}

func TestReference_WithVersion(t *testing.T) {
	ref := NewReference("entity-123", "parent")
	updated := ref.WithVersion(5)

	if updated.Version() != 5 {
		t.Errorf("WithVersion: expected 5, got %d", updated.Version())
	}

	// Original should be unchanged
	if ref.Version() != 0 {
		t.Errorf("original Reference should still have version 0, got %d", ref.Version())
	}
}

func TestID_MarshalText(t *testing.T) {
	// Non-zero ID should marshal to its string representation
	uid := NewID[UserBrand, string]("user-123")
	data, err := uid.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText error: %v", err)
	}
	if string(data) != "user-123" {
		t.Errorf("MarshalText: expected user-123, got %s", string(data))
	}

	// Zero ID should return nil
	var zeroUserID UserID
	data, err = zeroUserID.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText zero error: %v", err)
	}
	if data != nil {
		t.Errorf("MarshalText zero: expected nil, got %s", string(data))
	}
}

func TestID_UnmarshalText(t *testing.T) {
	// Unmarshal valid text
	var uid UserID
	err := uid.UnmarshalText([]byte("user-abc"))
	if err != nil {
		t.Fatalf("UnmarshalText error: %v", err)
	}
	if uid.Value() != "user-abc" {
		t.Errorf("UnmarshalText: expected user-abc, got %s", uid.Value())
	}

	// Unmarshal empty text should give zero value
	var uid2 UserID
	err = uid2.UnmarshalText([]byte{})
	if err != nil {
		t.Fatalf("UnmarshalText empty error: %v", err)
	}
	if !uid2.IsZero() {
		t.Error("UnmarshalText empty: expected zero ID")
	}
}

func TestLocale_NewLocale(t *testing.T) {
	tag := MustParseLocale("fr-FR").Tag()
	locale := NewLocale(tag)
	if locale.String() != "fr-FR" {
		t.Errorf("NewLocale: expected fr-FR, got %s", locale.String())
	}
}

func TestLocale_Tag(t *testing.T) {
	locale := LocaleEnUS
	tag := locale.Tag()
	if tag.String() != "en-US" {
		t.Errorf("Tag: expected en-US, got %s", tag.String())
	}
}

func TestPriorityEnum(t *testing.T) {
	// Test Names()
	names := PriorityNames()
	if len(names) != 4 {
		t.Fatalf("PriorityNames: expected 4, got %d", len(names))
	}
	expectedNames := []string{"Low", "Medium", "High", "Critical"}
	for i, name := range names {
		if name != expectedNames[i] {
			t.Errorf("PriorityNames[%d]: expected %s, got %s", i, expectedNames[i], name)
		}
	}

	// Test Values()
	values := PriorityValues()
	if len(values) != 4 {
		t.Fatalf("PriorityValues: expected 4, got %d", len(values))
	}
	if values[0] != PriorityLow || values[3] != PriorityCritical {
		t.Errorf("PriorityValues: unexpected values %v", values)
	}

	// Test String()
	if PriorityLow.String() != "Low" {
		t.Errorf("String: expected Low, got %s", PriorityLow.String())
	}

	// Test String() for invalid value
	invalid := Priority(99)
	if invalid.String() != "Priority(99)" {
		t.Errorf("String invalid: expected Priority(99), got %s", invalid.String())
	}

	// Test IsValid()
	if !PriorityMedium.IsValid() {
		t.Error("IsValid: expected Medium to be valid")
	}
	if invalid.IsValid() {
		t.Error("IsValid: expected 99 to be invalid")
	}

	// Test ParsePriority
	parsed, err := ParsePriority("High")
	if err != nil || parsed != PriorityHigh {
		t.Errorf("ParsePriority: expected High, got %v (err: %v)", parsed, err)
	}

	// Test ParsePriority with invalid
	_, err = ParsePriority("Invalid")
	if err == nil {
		t.Error("ParsePriority: expected error for invalid value")
	}

	// Test MustParsePriority
	mustParsed := MustParsePriority("Critical")
	if mustParsed != PriorityCritical {
		t.Errorf("MustParsePriority: expected Critical, got %v", mustParsed)
	}

	// Test MarshalText
	data, err := PriorityHigh.MarshalText()
	if err != nil || string(data) != "High" {
		t.Errorf("MarshalText: expected High, got %s (err: %v)", data, err)
	}

	// Test UnmarshalText
	var p Priority
	err = p.UnmarshalText([]byte("Low"))
	if err != nil || p != PriorityLow {
		t.Errorf("UnmarshalText: expected Low, got %v (err: %v)", p, err)
	}

	// Test UnmarshalText with invalid
	err = p.UnmarshalText([]byte("Invalid"))
	if err == nil {
		t.Error("UnmarshalText: expected error for invalid value")
	}

	// Test AppendText
	p2 := PriorityCritical
	appended, err2 := p2.AppendText([]byte("prefix-"))
	if err2 != nil || string(appended) != "prefix-Critical" {
		t.Errorf("AppendText: expected prefix-Critical, got %s (err: %v)", appended, err2)
	}

	// Test Scan (database → Go)
	var scanPriority Priority

	// Scan nil
	if err := scanPriority.Scan(nil); err != nil {
		t.Errorf("Scan nil: unexpected error %v", err)
	}
	if scanPriority != PriorityLow {
		t.Errorf("Scan nil: expected Low (0), got %v", scanPriority)
	}

	// Scan int64
	if err := scanPriority.Scan(int64(2)); err != nil {
		t.Errorf("Scan int64: unexpected error %v", err)
	}
	if scanPriority != PriorityHigh {
		t.Errorf("Scan int64: expected High, got %v", scanPriority)
	}

	// Scan string
	if err := scanPriority.Scan("Critical"); err != nil {
		t.Errorf("Scan string: unexpected error %v", err)
	}
	if scanPriority != PriorityCritical {
		t.Errorf("Scan string: expected Critical, got %v", scanPriority)
	}

	// Scan []byte
	if err := scanPriority.Scan([]byte("Low")); err != nil {
		t.Errorf("Scan []byte: unexpected error %v", err)
	}
	if scanPriority != PriorityLow {
		t.Errorf("Scan []byte: expected Low, got %v", scanPriority)
	}

	// Scan Priority directly
	if err := scanPriority.Scan(PriorityMedium); err != nil {
		t.Errorf("Scan Priority: unexpected error %v", err)
	}
	if scanPriority != PriorityMedium {
		t.Errorf("Scan Priority: expected Medium, got %v", scanPriority)
	}

	// Scan invalid string
	err = scanPriority.Scan("InvalidPriority")
	if err == nil {
		t.Error("Scan invalid: expected error for invalid string")
	}

	// Test Value (Go → database)
	val, err := PriorityHigh.Value()
	if err != nil {
		t.Errorf("Value: unexpected error %v", err)
	}
	if val != "High" {
		t.Errorf("Value: expected High, got %v", val)
	}

	// Test remaining Scan branches (uint, float64, pointers)
	intVal := 1
	int64Val := int64(1)
	uintVal := uint(1)
	uint64Val := uint64(1)
	float64Val := float64(1)
	strVal := "Medium"

	// Scan *int
	var pInt Priority
	if err := pInt.Scan(&intVal); err != nil {
		t.Errorf("Scan *int: unexpected error %v", err)
	}
	if pInt != PriorityMedium {
		t.Errorf("Scan *int: expected Medium, got %v", pInt)
	}

	// Scan *int64
	var pInt64 Priority
	if err := pInt64.Scan(&int64Val); err != nil {
		t.Errorf("Scan *int64: unexpected error %v", err)
	}
	if pInt64 != PriorityMedium {
		t.Errorf("Scan *int64: expected Medium, got %v", pInt64)
	}

	// Scan *uint
	var pUint Priority
	if err := pUint.Scan(&uintVal); err != nil {
		t.Errorf("Scan *uint: unexpected error %v", err)
	}
	if pUint != PriorityMedium {
		t.Errorf("Scan *uint: expected Medium, got %v", pUint)
	}

	// Scan *uint64
	var pUint64 Priority
	if err := pUint64.Scan(&uint64Val); err != nil {
		t.Errorf("Scan *uint64: unexpected error %v", err)
	}
	if pUint64 != PriorityMedium {
		t.Errorf("Scan *uint64: expected Medium, got %v", pUint64)
	}

	// Scan *float64
	var pFloat64 Priority
	if err := pFloat64.Scan(&float64Val); err != nil {
		t.Errorf("Scan *float64: unexpected error %v", err)
	}
	if pFloat64 != PriorityMedium {
		t.Errorf("Scan *float64: expected Medium, got %v", pFloat64)
	}

	// Scan *string
	var pStr Priority
	if err := pStr.Scan(&strVal); err != nil {
		t.Errorf("Scan *string: unexpected error %v", err)
	}
	if pStr != PriorityMedium {
		t.Errorf("Scan *string: expected Medium, got %v", pStr)
	}

	// Scan uint
	var pUintDirect Priority
	if err := pUintDirect.Scan(uint(2)); err != nil {
		t.Errorf("Scan uint: unexpected error %v", err)
	}
	if pUintDirect != PriorityHigh {
		t.Errorf("Scan uint: expected High, got %v", pUintDirect)
	}

	// Scan uint64
	var pUint64Direct Priority
	if err := pUint64Direct.Scan(uint64(3)); err != nil {
		t.Errorf("Scan uint64: unexpected error %v", err)
	}
	if pUint64Direct != PriorityCritical {
		t.Errorf("Scan uint64: expected Critical, got %v", pUint64Direct)
	}

	// Scan float64
	var pFloat64Direct Priority
	if err := pFloat64Direct.Scan(float64(0)); err != nil {
		t.Errorf("Scan float64: unexpected error %v", err)
	}
	if pFloat64Direct != PriorityLow {
		t.Errorf("Scan float64: expected Low, got %v", pFloat64Direct)
	}

	// Test Scan *Priority (pointer to same type)
	pOrig := PriorityCritical
	var pPtr Priority
	if err := pPtr.Scan(&pOrig); err != nil {
		t.Errorf("Scan *Priority: unexpected error %v", err)
	}
	if pPtr != PriorityCritical {
		t.Errorf("Scan *Priority: expected Critical, got %v", pPtr)
	}

	// Test Scan nil pointer (*Priority) - should error
	var pNil Priority
	err = pNil.Scan((*Priority)(nil))
	if err == nil {
		t.Error("Scan *Priority(nil): expected error for nil pointer")
	}

	// Test Scan nil pointer (*int) - should error
	var pNilInt Priority
	err = pNilInt.Scan((*int)(nil))
	if err == nil {
		t.Error("Scan *int(nil): expected error for nil pointer")
	}

	// Test Scan nil pointer (*string) - should error
	var pNilStr Priority
	err = pNilStr.Scan((*string)(nil))
	if err == nil {
		t.Error("Scan *string(nil): expected error for nil pointer")
	}

	// Test Scan int
	var pIntDirect Priority
	if err := pIntDirect.Scan(3); err != nil {
		t.Errorf("Scan int: unexpected error %v", err)
	}
	if pIntDirect != PriorityCritical {
		t.Errorf("Scan int: expected Critical, got %v", pIntDirect)
	}
}

func TestStatusEnum(t *testing.T) {
	// Test basic operations
	if StatusActive.String() != "Active" {
		t.Errorf("String: expected Active, got %s", StatusActive.String())
	}

	if !StatusDraft.IsValid() {
		t.Error("IsValid: expected Draft to be valid")
	}

	parsed, err := ParseStatus("Paused")
	if err != nil || parsed != StatusPaused {
		t.Errorf("ParseStatus: expected Paused, got %v (err: %v)", parsed, err)
	}

	// Test MustParse panics on invalid
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParseStatus: expected panic for invalid value")
		}
	}()
	MustParseStatus("Invalid")

	// Test Names and Values
	if len(StatusNames()) != 5 {
		t.Errorf("StatusNames: expected 5, got %d", len(StatusNames()))
	}
	if len(StatusValues()) != 5 {
		t.Errorf("StatusValues: expected 5, got %d", len(StatusValues()))
	}

	// Test AppendText
	s2 := StatusArchived
	appended, err2 := s2.AppendText([]byte("prefix-"))
	if err2 != nil || string(appended) != "prefix-Archived" {
		t.Errorf("AppendText: expected prefix-Archived, got %s (err: %v)", appended, err2)
	}

	// Test Scan (database → Go)
	var scanStatus Status

	// Scan nil
	if err := scanStatus.Scan(nil); err != nil {
		t.Errorf("Scan nil: unexpected error %v", err)
	}
	if scanStatus != StatusDraft {
		t.Errorf("Scan nil: expected Draft (0), got %v", scanStatus)
	}

	// Scan int64
	if err := scanStatus.Scan(int64(3)); err != nil {
		t.Errorf("Scan int64: unexpected error %v", err)
	}
	if scanStatus != StatusArchived {
		t.Errorf("Scan int64: expected Archived, got %v", scanStatus)
	}

	// Scan string
	if err := scanStatus.Scan("Deleted"); err != nil {
		t.Errorf("Scan string: unexpected error %v", err)
	}
	if scanStatus != StatusDeleted {
		t.Errorf("Scan string: expected Deleted, got %v", scanStatus)
	}

	// Scan []byte
	if err := scanStatus.Scan([]byte("Active")); err != nil {
		t.Errorf("Scan []byte: unexpected error %v", err)
	}
	if scanStatus != StatusActive {
		t.Errorf("Scan []byte: expected Active, got %v", scanStatus)
	}

	// Scan Status directly
	if err := scanStatus.Scan(StatusPaused); err != nil {
		t.Errorf("Scan Status: unexpected error %v", err)
	}
	if scanStatus != StatusPaused {
		t.Errorf("Scan Status: expected Paused, got %v", scanStatus)
	}

	// Scan invalid string
	err = scanStatus.Scan("InvalidStatus")
	if err == nil {
		t.Error("Scan invalid: expected error for invalid string")
	}

	// Test Value (Go → database)
	val, err := StatusActive.Value()
	if err != nil {
		t.Errorf("Value: unexpected error %v", err)
	}
	if val != "Active" {
		t.Errorf("Value: expected Active, got %v", val)
	}
}

func TestTriggerEnum(t *testing.T) {
	// Test Names and Values
	if len(TriggerNames()) != 7 {
		t.Errorf("TriggerNames: expected 7, got %d", len(TriggerNames()))
	}
	if len(TriggerValues()) != 7 {
		t.Errorf("TriggerValues: expected 7, got %d", len(TriggerValues()))
	}

	// Test String
	if TriggerWebhook.String() != "Webhook" {
		t.Errorf("String: expected Webhook, got %s", TriggerWebhook.String())
	}

	// Test text marshaling
	data, _ := TriggerManual.MarshalText()
	var t2 Trigger
	_ = t2.UnmarshalText(data)
	if t2 != TriggerManual {
		t.Errorf("Marshal/Unmarshal: expected Manual, got %v", t2)
	}

	// Test IsValid
	if !TriggerSystem.IsValid() {
		t.Error("IsValid: expected System to be valid")
	}
	invalid := Trigger(99)
	if invalid.IsValid() {
		t.Error("IsValid: expected 99 to be invalid")
	}

	// Test AppendText
	trig := TriggerMigration
	appended, err := trig.AppendText([]byte("prefix-"))
	if err != nil || string(appended) != "prefix-Migration" {
		t.Errorf("AppendText: expected prefix-Migration, got %s (err: %v)", appended, err)
	}

	// Test Scan (database → Go)
	var scanTrigger Trigger

	// Scan nil
	if err := scanTrigger.Scan(nil); err != nil {
		t.Errorf("Scan nil: unexpected error %v", err)
	}
	if scanTrigger != TriggerManual {
		t.Errorf("Scan nil: expected Manual (0), got %v", scanTrigger)
	}

	// Scan int64
	if err := scanTrigger.Scan(int64(5)); err != nil {
		t.Errorf("Scan int64: unexpected error %v", err)
	}
	if scanTrigger != TriggerSystem {
		t.Errorf("Scan int64: expected System, got %v", scanTrigger)
	}

	// Scan string
	if err := scanTrigger.Scan("Correction"); err != nil {
		t.Errorf("Scan string: unexpected error %v", err)
	}
	if scanTrigger != TriggerCorrection {
		t.Errorf("Scan string: expected Correction, got %v", scanTrigger)
	}

	// Scan []byte
	if err := scanTrigger.Scan([]byte("Scheduled")); err != nil {
		t.Errorf("Scan []byte: unexpected error %v", err)
	}
	if scanTrigger != TriggerScheduled {
		t.Errorf("Scan []byte: expected Scheduled, got %v", scanTrigger)
	}

	// Scan Trigger directly
	if err := scanTrigger.Scan(TriggerImport); err != nil {
		t.Errorf("Scan Trigger: unexpected error %v", err)
	}
	if scanTrigger != TriggerImport {
		t.Errorf("Scan Trigger: expected Import, got %v", scanTrigger)
	}

	// Scan invalid string
	err = scanTrigger.Scan("InvalidTrigger")
	if err == nil {
		t.Error("Scan invalid: expected error for invalid string")
	}

	// Test Value (Go → database)
	val, err := TriggerWebhook.Value()
	if err != nil {
		t.Errorf("Value: unexpected error %v", err)
	}
	if val != "Webhook" {
		t.Errorf("Value: expected Webhook, got %v", val)
	}

	// Test Parse error path
	_, err = ParseTrigger("Invalid")
	if err == nil {
		t.Error("ParseTrigger: expected error for invalid value")
	}
}

func TestActorKindEnum(t *testing.T) {
	// Test basic operations
	if ActorKindUser.String() != "User" {
		t.Errorf("String: expected User, got %s", ActorKindUser.String())
	}

	if !ActorKindBot.IsValid() {
		t.Error("IsValid: expected Bot to be valid")
	}

	parsed, err := ParseActorKind("System")
	if err != nil || parsed != ActorKindSystem {
		t.Errorf("ParseActorKind: expected System, got %v (err: %v)", parsed, err)
	}

	// Test invalid kind
	invalid := ActorKind(99)
	if invalid.IsValid() {
		t.Error("IsValid: expected 99 to be invalid")
	}
	if invalid.String() != "ActorKind(99)" {
		t.Errorf("String invalid: expected ActorKind(99), got %s", invalid.String())
	}

	// Test Names and Values
	if len(ActorKindNames()) != 4 {
		t.Errorf("ActorKindNames: expected 4, got %d", len(ActorKindNames()))
	}
	if len(ActorKindValues()) != 4 {
		t.Errorf("ActorKindValues: expected 4, got %d", len(ActorKindValues()))
	}

	// Test AppendText
	ak := ActorKindService
	appended, err := ak.AppendText([]byte("prefix-"))
	if err != nil || string(appended) != "prefix-Service" {
		t.Errorf("AppendText: expected prefix-Service, got %s (err: %v)", appended, err)
	}

	// Test Scan (database → Go)
	var scanActorKind ActorKind

	// Scan nil
	if err := scanActorKind.Scan(nil); err != nil {
		t.Errorf("Scan nil: unexpected error %v", err)
	}
	if scanActorKind != ActorKindUser {
		t.Errorf("Scan nil: expected User (0), got %v", scanActorKind)
	}

	// Scan int64
	if err := scanActorKind.Scan(int64(2)); err != nil {
		t.Errorf("Scan int64: unexpected error %v", err)
	}
	if scanActorKind != ActorKindSystem {
		t.Errorf("Scan int64: expected System, got %v", scanActorKind)
	}

	// Scan string
	if err := scanActorKind.Scan("Service"); err != nil {
		t.Errorf("Scan string: unexpected error %v", err)
	}
	if scanActorKind != ActorKindService {
		t.Errorf("Scan string: expected Service, got %v", scanActorKind)
	}

	// Scan []byte
	if err := scanActorKind.Scan([]byte("Bot")); err != nil {
		t.Errorf("Scan []byte: unexpected error %v", err)
	}
	if scanActorKind != ActorKindBot {
		t.Errorf("Scan []byte: expected Bot, got %v", scanActorKind)
	}

	// Scan ActorKind directly
	if err := scanActorKind.Scan(ActorKindBot); err != nil {
		t.Errorf("Scan ActorKind: unexpected error %v", err)
	}
	if scanActorKind != ActorKindBot {
		t.Errorf("Scan ActorKind: expected Bot, got %v", scanActorKind)
	}

	// Scan invalid string
	err = scanActorKind.Scan("InvalidActorKind")
	if err == nil {
		t.Error("Scan invalid: expected error for invalid string")
	}

	// Test Value (Go → database)
	val, err := ActorKindService.Value()
	if err != nil {
		t.Errorf("Value: unexpected error %v", err)
	}
	if val != "Service" {
		t.Errorf("Value: expected Service, got %v", val)
	}

	// Test MarshalText/UnmarshalText error path
	var ak2 ActorKind
	err = ak2.UnmarshalText([]byte("Invalid"))
	if err == nil {
		t.Error("UnmarshalText: expected error for invalid value")
	}

	// Test MustParse panics on invalid
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParseActorKind: expected panic for invalid value")
		}
	}()
	MustParseActorKind("Invalid")
}
