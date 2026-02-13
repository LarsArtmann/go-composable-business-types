package cbt

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// === NanoId Tests ===

func TestNanoId(t *testing.T) {
	id := NewNanoId()
	if id.IsEmpty() {
		t.Error("expected non-empty NanoId")
	}
	if len(id.String()) != DefaultNanoIdLength {
		t.Errorf("expected length %d, got %d", DefaultNanoIdLength, len(id.String()))
	}
}

func TestNanoIdWithLength(t *testing.T) {
	id := NewNanoIdWithLength(10)
	if len(id.String()) != 10 {
		t.Errorf("expected length 10, got %d", len(id.String()))
	}
}

func TestNanoIdUniqueness(t *testing.T) {
	ids := make(map[string]bool)
	for range 1000 {
		id := NewNanoId()
		if ids[id.String()] {
			t.Error("generated duplicate NanoId")
		}
		ids[id.String()] = true
	}
}

func TestNanoIdAlphabet(t *testing.T) {
	id := NewNanoId()
	for _, r := range id.String() {
		if !isNanoIdChar(r) {
			t.Errorf("invalid character in NanoId: %c", r)
		}
	}
}

func TestParseNanoId(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid", "abcdefgh12345678", nil},
		{"valid with special", "abc-def_ghi12345", nil},
		{"empty", "", ErrNanoIdEmpty},
		{"too short", "short", ErrNanoIdTooShort},
		{"too long", strings.Repeat("a", 257), ErrNanoIdTooLong},
		{"invalid char", "abcdefgh@2345678", ErrNanoIdInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseNanoId(tt.input)
			if err != tt.wantErr {
				t.Errorf("ParseNanoId(%q) = %v, want %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestNanoIdJSON(t *testing.T) {
	original := NewNanoId()
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed NanoId
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if original.String() != parsed.String() {
		t.Errorf("expected %s, got %s", original.String(), parsed.String())
	}
}

func TestNanoIdJSONEmpty(t *testing.T) {
	var empty NanoId
	data, err := json.Marshal(empty)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	// Empty NanoId marshals to empty string (via MarshalText returning nil, nil)
	if string(data) != `""` {
		t.Errorf("expected empty string for empty NanoId, got %s", string(data))
	}
}

// === DataPoint Tests ===

type TestPayload struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

func TestDataPoint(t *testing.T) {
	payload := TestPayload{Value: "test", Count: 42}
	actor := UserActor(NewID[struct{}, string]("user-1"), "Alice")
	occurred := NewTimestamp(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))
	recorded := NewTimestamp(time.Date(2024, 1, 15, 10, 31, 0, 0, time.UTC))

	dp := NewDataPoint(payload, actor, occurred, recorded, "test reason")

	if dp.Id().IsEmpty() {
		t.Error("expected non-empty DataPoint id")
	}
	if dp.Payload().Value != "test" {
		t.Errorf("expected payload value 'test', got %s", dp.Payload().Value)
	}
	if dp.Payload().Count != 42 {
		t.Errorf("expected payload count 42, got %d", dp.Payload().Count)
	}
	if dp.Actor().Name != "Alice" {
		t.Errorf("expected actor name 'Alice', got %s", dp.Actor().Name)
	}
	if dp.Reason() != "test reason" {
		t.Errorf("expected reason 'test reason', got %s", dp.Reason())
	}
}

func TestDataPointNow(t *testing.T) {
	payload := TestPayload{Value: "now", Count: 1}
	actor := BotActor(NewID[struct{}, string]("bot-1"), "TestBot")

	before := time.Now().UTC()
	dp := NewDataPointNow(payload, actor)
	after := time.Now().UTC()

	// Recorded should be between before and after
	if dp.Recorded().Before(before) || dp.Recorded().After(after) {
		t.Error("recorded time not in expected range")
	}

	// Occurred equals recorded for NewDataPointNow
	if !dp.Occurred().Equal(dp.Recorded().Time) {
		t.Error("occurred should equal recorded for NewDataPointNow")
	}
}

func TestDataPointWithoutReason(t *testing.T) {
	payload := TestPayload{Value: "no-reason", Count: 0}
	actor := SystemActor[string]()
	occurred := Now()
	recorded := Now()

	dp := NewDataPoint(payload, actor, occurred, recorded)

	if dp.Reason() != "" {
		t.Errorf("expected empty reason, got %s", dp.Reason())
	}
}

func TestDataPointWithReason(t *testing.T) {
	payload := TestPayload{Value: "test", Count: 1}
	actor := SystemActor[string]()
	occurred := Now()
	recorded := Now()

	dp := NewDataPoint(payload, actor, occurred, recorded).WithReason("updated reason")

	if dp.Reason() != "updated reason" {
		t.Errorf("expected 'updated reason', got %s", dp.Reason())
	}
}

func TestDataPointJSON(t *testing.T) {
	payload := TestPayload{Value: "json-test", Count: 99}
	actor := UserActor(NewID[struct{}, string]("user-42"), "Bob")
	occurred := NewTimestamp(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))
	recorded := NewTimestamp(time.Date(2024, 6, 1, 12, 1, 0, 0, time.UTC))

	original := NewDataPoint(payload, actor, occurred, recorded, "json serialization test")

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed DataPoint[TestPayload]
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if original.Id().String() != parsed.Id().String() {
		t.Errorf("id mismatch: %s vs %s", original.Id().String(), parsed.Id().String())
	}
	if original.Payload().Value != parsed.Payload().Value {
		t.Errorf("payload value mismatch")
	}
	if original.Payload().Count != parsed.Payload().Count {
		t.Errorf("payload count mismatch")
	}
	if original.Actor().Name != parsed.Actor().Name {
		t.Errorf("actor name mismatch")
	}
	if original.Reason() != parsed.Reason() {
		t.Errorf("reason mismatch")
	}
}

func TestDataPointJSONRoundTrip(t *testing.T) {
	payload := TestPayload{Value: "round-trip", Count: 123}
	actor := ServiceActor(NewID[struct{}, string]("svc-1"), "OrderService")

	dp1 := NewDataPointNow(payload, actor, "first version")
	data, _ := json.Marshal(dp1)

	var dp2 DataPoint[TestPayload]
	if err := json.Unmarshal(data, &dp2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	data2, _ := json.Marshal(dp2)
	if string(data) != string(data2) {
		t.Errorf("round trip failed:\n%s\nvs\n%s", string(data), string(data2))
	}
}

func TestDataPointJSONWithoutReason(t *testing.T) {
	payload := TestPayload{Value: "no-reason", Count: 0}
	actor := SystemActor[string]()

	dp := NewDataPointNow(payload, actor)
	data, err := json.Marshal(dp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	// Reason should be omitted when empty
	if strings.Contains(string(data), `"reason":""`) {
		t.Error("empty reason should be omitted from JSON")
	}
}

func TestDataPointGenericTypes(t *testing.T) {
	t.Run("string payload", func(t *testing.T) {
		dp := NewDataPointNow("hello", SystemActor[string]())
		if dp.Payload() != "hello" {
			t.Errorf("expected 'hello', got %s", dp.Payload())
		}
	})

	t.Run("int payload", func(t *testing.T) {
		dp := NewDataPointNow(42, SystemActor[string]())
		if dp.Payload() != 42 {
			t.Errorf("expected 42, got %d", dp.Payload())
		}
	})

	t.Run("struct payload", func(t *testing.T) {
		type SimpleStruct struct {
			X int
			Y int
		}
		dp := NewDataPointNow(SimpleStruct{X: 1, Y: 2}, SystemActor[string]())
		if dp.Payload().X != 1 || dp.Payload().Y != 2 {
			t.Errorf("expected {1,2}, got %+v", dp.Payload())
		}
	})
}

// === Bitemporal Tests ===

func TestBitemporal(t *testing.T) {
	now := Now()
	b := NewBitemporal(now)

	if b.ValidFrom().IsZero() {
		t.Error("valid from should not be zero")
	}
	if !b.ValidUntil().IsZero() {
		t.Error("valid until should be zero (indefinite)")
	}
	if b.Recorded().IsZero() {
		t.Error("recorded should not be zero")
	}
	if b.IsCorrection() {
		t.Error("should not be a correction")
	}
}

func TestBitemporalWithRange(t *testing.T) {
	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	until := NewTimestamp(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC))
	recorded := NewTimestamp(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))

	b := NewBitemporalWithRange(from, until, recorded)

	if !b.ValidFrom().Equal(from.Time) {
		t.Error("valid from mismatch")
	}
	if !b.ValidUntil().Equal(until.Time) {
		t.Error("valid until mismatch")
	}
}

func TestBitemporalIsValidAt(t *testing.T) {
	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	until := NewTimestamp(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC))
	recorded := Now()

	b := NewBitemporalWithRange(from, until, recorded)

	// Before valid range
	before := NewTimestamp(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC))
	if b.IsValidAt(before) {
		t.Error("should not be valid before validFrom")
	}

	// During valid range
	during := NewTimestamp(time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC))
	if !b.IsValidAt(during) {
		t.Error("should be valid during range")
	}

	// After valid range
	after := NewTimestamp(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	if b.IsValidAt(after) {
		t.Error("should not be valid after validUntil")
	}
}

func TestBitemporalIndefinite(t *testing.T) {
	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	b := NewBitemporalWithRange(from, Timestamp{}, Now())

	// Should be valid for any time after from
	future := NewTimestamp(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))
	if !b.IsValidAt(future) {
		t.Error("indefinite validity should extend to future")
	}
}

func TestBitemporalCorrection(t *testing.T) {
	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	b := NewCorrection(from, Timestamp{}, Now())

	if !b.IsCorrection() {
		t.Error("should be marked as correction")
	}
}

func TestBitemporalJSON(t *testing.T) {
	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	until := NewTimestamp(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC))
	recorded := NewTimestamp(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))

	original := NewBitemporalWithRange(from, until, recorded)

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed Bitemporal
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !parsed.ValidFrom().Equal(original.ValidFrom().Time) {
		t.Error("validFrom mismatch")
	}
	if !parsed.ValidUntil().Equal(original.ValidUntil().Time) {
		t.Error("validUntil mismatch")
	}
}

// === Context Tests ===

func TestContext(t *testing.T) {
	ctx := NewContext("order-service")

	if ctx.Source() != "order-service" {
		t.Errorf("expected 'order-service', got %s", ctx.Source())
	}
}

func TestContextWithFields(t *testing.T) {
	ctx := NewContext("payment-svc").
		WithEnvironment("production").
		WithSession("sess-123").
		WithRequest("req-456").
		WithTag("version", "1.0.0").
		WithTag("region", "us-east-1")

	if ctx.Environment() != "production" {
		t.Error("environment mismatch")
	}
	if ctx.Session() != "sess-123" {
		t.Error("session mismatch")
	}
	if ctx.Request() != "req-456" {
		t.Error("request mismatch")
	}

	v, ok := ctx.Tag("version")
	if !ok || v != "1.0.0" {
		t.Error("tag version mismatch")
	}

	tags := ctx.Tags()
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
}

func TestContextJSON(t *testing.T) {
	original := NewContext("api-gateway").
		WithEnvironment("staging").
		WithSession("sess-abc").
		WithTag("team", "platform")

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed Context
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if parsed.Source() != "api-gateway" {
		t.Error("source mismatch")
	}
	if parsed.Environment() != "staging" {
		t.Error("environment mismatch")
	}
}

// === Trigger Tests ===

func TestTrigger(t *testing.T) {
	// Test enum generation works
	if TriggerManual.String() != "Manual" {
		t.Errorf("expected 'Manual', got %s", TriggerManual.String())
	}

	// Test parsing
	trigger, err := ParseTrigger("Webhook")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if trigger != TriggerWebhook {
		t.Errorf("expected TriggerWebhook, got %v", trigger)
	}

	// Test invalid parse
	_, err = ParseTrigger("Invalid")
	if err == nil {
		t.Error("expected error for invalid trigger")
	}
}

// === DataPoint with Phase 2 fields ===

func TestDataPointWithTrigger(t *testing.T) {
	payload := TestPayload{Value: "webhook-trigger", Count: 1}
	actor := SystemActor[string]()

	dp := NewDataPointNow(payload, actor).WithTrigger(TriggerWebhook)

	if dp.Trigger() != TriggerWebhook {
		t.Errorf("expected TriggerWebhook, got %v", dp.Trigger())
	}
}

func TestDataPointWithContext(t *testing.T) {
	payload := TestPayload{Value: "context-test", Count: 1}
	actor := SystemActor[string]()
	ctx := NewContext("test-service").WithEnvironment("development")

	dp := NewDataPointNow(payload, actor).WithContext(ctx)

	if dp.Context().Source() != "test-service" {
		t.Error("context source mismatch")
	}
}

func TestDataPointWithTemporal(t *testing.T) {
	payload := TestPayload{Value: "temporal-test", Count: 1}
	actor := SystemActor[string]()

	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	until := NewTimestamp(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC))
	recorded := Now()
	temporal := NewBitemporalWithRange(from, until, recorded)

	dp := NewDataPointNow(payload, actor).WithTemporal(temporal)

	if !dp.Temporal().ValidFrom().Equal(from.Time) {
		t.Error("temporal validFrom mismatch")
	}
	if !dp.Temporal().ValidUntil().Equal(until.Time) {
		t.Error("temporal validUntil mismatch")
	}
}

func TestDataPointTemporalBackwardCompat(t *testing.T) {
	// Occurred() should delegate to Temporal().ValidFrom()
	payload := TestPayload{Value: "compat", Count: 1}
	actor := SystemActor[string]()
	occurred := NewTimestamp(time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC))
	recorded := NewTimestamp(time.Date(2024, 6, 15, 10, 1, 0, 0, time.UTC))

	dp := NewDataPoint(payload, actor, occurred, recorded)

	if !dp.Occurred().Equal(occurred.Time) {
		t.Error("Occurred() should equal passed occurred time")
	}
	if !dp.Recorded().Equal(recorded.Time) {
		t.Error("Recorded() should equal passed recorded time")
	}
}

// === Reference Tests ===

func TestReference(t *testing.T) {
	ref := NewReference("order-123", "parent")

	if ref.Id() != "order-123" {
		t.Errorf("expected 'order-123', got %s", ref.Id())
	}
	if ref.Relation() != "parent" {
		t.Errorf("expected 'parent', got %s", ref.Relation())
	}
	if ref.Version() != 0 {
		t.Errorf("expected version 0, got %d", ref.Version())
	}
}

func TestReferenceWithVersion(t *testing.T) {
	ref := NewReferenceWithVersion("doc-456", "source", 5)

	if ref.Version() != 5 {
		t.Errorf("expected version 5, got %d", ref.Version())
	}
}

func TestReferenceTags(t *testing.T) {
	ref := NewReference("user-789", "owner").
		WithTag("department", "engineering").
		WithTag("level", "senior")

	if ref.Tags()["department"] != "engineering" {
		t.Error("tag department mismatch")
	}

	v, ok := ref.Tag("level")
	if !ok || v != "senior" {
		t.Error("tag level mismatch")
	}

	_, ok = ref.Tag("nonexistent")
	if ok {
		t.Error("nonexistent tag should return false")
	}
}

func TestReferenceJSON(t *testing.T) {
	original := NewReferenceWithVersion("entity-1", "reference", 3).
		WithTag("key", "value")

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed Reference[string]
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if parsed.Id() != "entity-1" {
		t.Error("id mismatch")
	}
	if parsed.Relation() != "reference" {
		t.Error("relation mismatch")
	}
	if parsed.Version() != 3 {
		t.Error("version mismatch")
	}
	if parsed.Tags()["key"] != "value" {
		t.Error("tag mismatch")
	}
}

// === Cause Tests ===

func TestCause(t *testing.T) {
	id := NewNanoId()
	cause := NewCause(id, "command", "created")

	if cause.Id().String() != id.String() {
		t.Error("id mismatch")
	}
	if cause.Kind() != "command" {
		t.Errorf("expected 'command', got %s", cause.Kind())
	}
	if cause.Effect() != "created" {
		t.Errorf("expected 'created', got %s", cause.Effect())
	}
	if cause.HasTrace() {
		t.Error("should not have trace")
	}
}

func TestCauseDirect(t *testing.T) {
	id := NewNanoId()
	cause := NewCauseDirect(id)

	if cause.Kind() != "direct" {
		t.Errorf("expected 'direct', got %s", cause.Kind())
	}
	if cause.Effect() != "caused" {
		t.Errorf("expected 'caused', got %s", cause.Effect())
	}
}

func TestCauseCommand(t *testing.T) {
	id := NewNanoId()
	cause := NewCauseCommand(id, "approved")

	if cause.Kind() != "command" {
		t.Error("kind should be command")
	}
	if cause.Effect() != "approved" {
		t.Error("effect mismatch")
	}
}

func TestCauseEvent(t *testing.T) {
	id := NewNanoId()
	cause := NewCauseEvent(id, "triggered")

	if cause.Kind() != "event" {
		t.Error("kind should be event")
	}
	if cause.Effect() != "triggered" {
		t.Error("effect mismatch")
	}
}

func TestCauseWithTrace(t *testing.T) {
	id1 := NewNanoId()
	id2 := NewNanoId()
	id3 := NewNanoId()

	cause := NewCauseDirect(id1).
		WithTrace([]NanoId{id2, id3})

	if !cause.HasTrace() {
		t.Error("should have trace")
	}
	if len(cause.Trace()) != 2 {
		t.Errorf("expected 2 trace items, got %d", len(cause.Trace()))
	}
}

func TestCauseAppendTrace(t *testing.T) {
	id1 := NewNanoId()
	id2 := NewNanoId()

	cause := NewCauseDirect(id1).AppendTrace(id2)

	if !cause.HasTrace() {
		t.Error("should have trace")
	}
	if cause.Trace()[0].String() != id2.String() {
		t.Error("trace item mismatch")
	}
}

func TestCauseJSON(t *testing.T) {
	id := NewNanoId()
	traceId := NewNanoId()

	original := NewCauseCommand(id, "executed").
		WithTrace([]NanoId{traceId})

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed Cause[string]
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if parsed.Id().String() != id.String() {
		t.Error("id mismatch")
	}
	if parsed.Kind() != "command" {
		t.Error("kind mismatch")
	}
	if parsed.Effect() != "executed" {
		t.Error("effect mismatch")
	}
	if !parsed.HasTrace() {
		t.Error("should have trace")
	}
}

// === DataPoint Phase 3 Tests ===

func TestDataPointWithVersion(t *testing.T) {
	dp := NewDataPointNow(TestPayload{Value: "versioned", Count: 1}, SystemActor[string]()).
		WithVersion(5)

	if dp.Version() != 5 {
		t.Errorf("expected version 5, got %d", dp.Version())
	}
}

func TestDataPointWithTags(t *testing.T) {
	tags := map[string]string{"env": "prod", "region": "us-east-1"}
	dp := NewDataPointNow(TestPayload{Value: "tagged", Count: 1}, SystemActor[string]()).
		WithTags(tags)

	if dp.Tags()["env"] != "prod" {
		t.Error("tag env mismatch")
	}
	if dp.Tags()["region"] != "us-east-1" {
		t.Error("tag region mismatch")
	}
}

func TestDataPointAddTag(t *testing.T) {
	dp := NewDataPointNow(TestPayload{Value: "add-tag", Count: 1}, SystemActor[string]()).
		AddTag("key1", "value1").
		AddTag("key2", "value2")

	if dp.Tags()["key1"] != "value1" {
		t.Error("tag key1 mismatch")
	}
	if dp.Tags()["key2"] != "value2" {
		t.Error("tag key2 mismatch")
	}

	v, ok := dp.Tag("key1")
	if !ok || v != "value1" {
		t.Error("Tag() method failed")
	}
}

func TestDataPointWithReferences(t *testing.T) {
	ref1 := NewReference("order-1", "parent")
	ref2 := NewReference("customer-2", "owner")

	dp := NewDataPointNow(TestPayload{Value: "refs", Count: 1}, SystemActor[string]()).
		WithReferences([]Reference[string]{ref1, ref2})

	if len(dp.References()) != 2 {
		t.Errorf("expected 2 references, got %d", len(dp.References()))
	}
	if dp.References()[0].Id() != "order-1" {
		t.Error("first reference id mismatch")
	}
}

func TestDataPointAddReference(t *testing.T) {
	ref := NewReference("product-123", "item")

	dp := NewDataPointNow(TestPayload{Value: "add-ref", Count: 1}, SystemActor[string]()).
		AddReference(ref)

	if len(dp.References()) != 1 {
		t.Errorf("expected 1 reference, got %d", len(dp.References()))
	}
	if dp.References()[0].Id() != "product-123" {
		t.Error("reference id mismatch")
	}
}

func TestDataPointWithCauses(t *testing.T) {
	cause1 := NewCauseDirect(NewNanoId())
	cause2 := NewCauseCommand(NewNanoId(), "approved")

	dp := NewDataPointNow(TestPayload{Value: "causes", Count: 1}, SystemActor[string]()).
		WithCauses([]Cause[string]{cause1, cause2})

	if len(dp.Causes()) != 2 {
		t.Errorf("expected 2 causes, got %d", len(dp.Causes()))
	}
}

func TestDataPointAddCause(t *testing.T) {
	cause := NewCauseEvent(NewNanoId(), "triggered")

	dp := NewDataPointNow(TestPayload{Value: "add-cause", Count: 1}, SystemActor[string]()).
		AddCause(cause)

	if len(dp.Causes()) != 1 {
		t.Errorf("expected 1 cause, got %d", len(dp.Causes()))
	}
}

func TestDataPointFullJSON(t *testing.T) {
	causeId := NewNanoId()
	traceId := NewNanoId()

	dp := NewDataPointNow(TestPayload{Value: "full", Count: 42}, SystemActor[string]()).
		WithTrigger(TriggerWebhook).
		WithReason("full integration test").
		WithVersion(3).
		AddTag("env", "test").
		AddReference(NewReferenceWithVersion("doc-1", "source", 2)).
		AddCause(NewCauseCommand(causeId, "created").WithTrace([]NanoId{traceId}))

	data, err := json.Marshal(dp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed DataPoint[TestPayload]
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if parsed.Version() != 3 {
		t.Errorf("version mismatch: %d", parsed.Version())
	}
	if parsed.Tags()["env"] != "test" {
		t.Error("tag env mismatch")
	}
	if len(parsed.References()) != 1 {
		t.Error("references count mismatch")
	}
	if parsed.References()[0].Id() != "doc-1" {
		t.Error("reference id mismatch")
	}
	if parsed.References()[0].Version() != 2 {
		t.Error("reference version mismatch")
	}
	if len(parsed.Causes()) != 1 {
		t.Error("causes count mismatch")
	}
	if parsed.Causes()[0].Kind() != "command" {
		t.Error("cause kind mismatch")
	}
	if !parsed.Causes()[0].HasTrace() {
		t.Error("cause should have trace")
	}
}

func TestDataPointFullRoundTrip(t *testing.T) {
	original := NewDataPointNow(TestPayload{Value: "roundtrip", Count: 99}, SystemActor[string]()).
		WithTrigger(TriggerScheduled).
		WithReason("round trip test").
		WithVersion(7).
		AddTag("team", "platform").
		AddTag("version", "2.0").
		AddReference(NewReference("parent-1", "parent")).
		AddReference(NewReference("child-2", "child")).
		AddCause(NewCauseDirect(NewNanoId())).
		AddCause(NewCauseEvent(NewNanoId(), "triggered"))

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed DataPoint[TestPayload]
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	data2, err := json.Marshal(parsed)
	if err != nil {
		t.Fatalf("second marshal error: %v", err)
	}

	if string(data) != string(data2) {
		t.Errorf("round trip failed:\n%s\nvs\n%s", string(data), string(data2))
	}
}
