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
	assertNotZero(t, id)
	assertEqual(t, len(id.String()), DefaultNanoIdLength)
}

func TestNanoIdWithLength(t *testing.T) {
	id := NewNanoIdWithLength(10)
	assertEqual(t, len(id.String()), 10)
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

	assertEqual(t, original.String(), parsed.String())
}

func TestNanoIdJSONEmpty(t *testing.T) {
	var empty NanoId
	data, err := json.Marshal(empty)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	// Empty NanoId marshals to empty string (via MarshalText returning nil, nil)
	assertJSONEquals(t, data, `""`)
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

	assertNotZero(t, dp.Id())
	assertEqual(t, dp.Payload().Value, "test")
	assertEqual(t, dp.Payload().Count, 42)
	assertEqual(t, dp.Actor().Name, "Alice")
	assertEqual(t, dp.Reason(), "test reason")
}

func TestDataPointNow(t *testing.T) {
	payload := TestPayload{Value: "now", Count: 1}
	actor := BotActor(NewID[struct{}, string]("bot-1"), "TestBot")

	before := time.Now().UTC()
	dp := NewDataPointNow(payload, actor)
	after := time.Now().UTC()

	// Recorded should be between before and after
	assertTrue(t, !dp.Recorded().Before(before) && !dp.Recorded().After(after), "recorded time not in expected range")

	// Occurred equals recorded for NewDataPointNow
	assertTrue(t, dp.Occurred().Equal(dp.Recorded().Time), "occurred should equal recorded for NewDataPointNow")
}

func TestDataPointWithoutReason(t *testing.T) {
	payload := TestPayload{Value: "no-reason", Count: 0}
	actor := SystemActor[string]()
	occurred := Now()
	recorded := Now()

	dp := NewDataPoint(payload, actor, occurred, recorded)

	assertEqual(t, dp.Reason(), "")
}

func TestDataPointWithReason(t *testing.T) {
	payload := TestPayload{Value: "test", Count: 1}
	actor := SystemActor[string]()
	occurred := Now()
	recorded := Now()

	dp := NewDataPoint(payload, actor, occurred, recorded).WithReason("updated reason")

	assertEqual(t, dp.Reason(), "updated reason")
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

	assertEqual(t, original.Id().String(), parsed.Id().String())
	assertEqual(t, original.Payload().Value, parsed.Payload().Value)
	assertEqual(t, original.Payload().Count, parsed.Payload().Count)
	assertEqual(t, original.Actor().Name, parsed.Actor().Name)
	assertEqual(t, original.Reason(), parsed.Reason())
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
	assertJSONEquals(t, data, string(data2))
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
	assertFalse(t, strings.Contains(string(data), `"reason":""`), "empty reason should be omitted from JSON")
}

func TestDataPointGenericTypes(t *testing.T) {
	t.Run("string payload", func(t *testing.T) {
		dp := NewDataPointNow("hello", SystemActor[string]())
		assertEqual(t, dp.Payload(), "hello")
	})

	t.Run("int payload", func(t *testing.T) {
		dp := NewDataPointNow(42, SystemActor[string]())
		assertEqual(t, dp.Payload(), 42)
	})

	t.Run("struct payload", func(t *testing.T) {
		type SimpleStruct struct {
			X int
			Y int
		}
		dp := NewDataPointNow(SimpleStruct{X: 1, Y: 2}, SystemActor[string]())
		assertTrue(t, dp.Payload().X == 1 && dp.Payload().Y == 2, "expected {1,2}")
	})
}

// === Bitemporal Tests ===

func TestBitemporal(t *testing.T) {
	now := Now()
	b := NewBitemporal(now)

	assertNotZero(t, b.ValidFrom())
	assertZero(t, b.ValidUntil())
	assertNotZero(t, b.Recorded())
	assertFalse(t, b.IsCorrection(), "should not be a correction")
}

func TestBitemporalWithRange(t *testing.T) {
	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	until := NewTimestamp(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC))
	recorded := NewTimestamp(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC))

	b := NewBitemporalWithRange(from, until, recorded)

	assertTrue(t, b.ValidFrom().Equal(from.Time), "valid from mismatch")
	assertTrue(t, b.ValidUntil().Equal(until.Time), "valid until mismatch")
}

func TestBitemporalIsValidAt(t *testing.T) {
	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	until := NewTimestamp(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC))
	recorded := Now()

	b := NewBitemporalWithRange(from, until, recorded)

	// Before valid range
	before := NewTimestamp(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC))
	assertFalse(t, b.IsValidAt(before), "should not be valid before validFrom")

	// During valid range
	during := NewTimestamp(time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC))
	assertTrue(t, b.IsValidAt(during), "should be valid during range")

	// After valid range
	after := NewTimestamp(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	assertFalse(t, b.IsValidAt(after), "should not be valid after validUntil")
}

func TestBitemporalIndefinite(t *testing.T) {
	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	b := NewBitemporalWithRange(from, Timestamp{}, Now())

	// Should be valid for any time after from
	future := NewTimestamp(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))
	assertTrue(t, b.IsValidAt(future), "indefinite validity should extend to future")
}

func TestBitemporalCorrection(t *testing.T) {
	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	b := NewCorrection(from, Timestamp{}, Now())
	assertTrue(t, b.IsCorrection(), "should be marked as correction")
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

	assertTrue(t, parsed.ValidFrom().Equal(original.ValidFrom().Time), "validFrom mismatch")
	assertTrue(t, parsed.ValidUntil().Equal(original.ValidUntil().Time), "validUntil mismatch")
}

// === Context Tests ===

func TestContext(t *testing.T) {
	ctx := NewContext("order-service")
	assertEqual(t, ctx.Source(), "order-service")
}

func TestContextWithFields(t *testing.T) {
	ctx := NewContext("payment-svc").
		WithEnvironment("production").
		WithSession("sess-123").
		WithRequest("req-456").
		WithTag("version", "1.0.0").
		WithTag("region", "us-east-1")

	assertEqual(t, ctx.Environment(), "production")
	assertEqual(t, ctx.Session(), "sess-123")
	assertEqual(t, ctx.Request(), "req-456")

	v, ok := ctx.Tag("version")
	assertTrue(t, ok && v == "1.0.0", "tag version mismatch")

	tags := ctx.Tags()
	assertEqual(t, len(tags), 2)
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

	assertEqual(t, parsed.Source(), "api-gateway")
	assertEqual(t, parsed.Environment(), "staging")
}

// === Trigger Tests ===

func TestTrigger(t *testing.T) {
	// Test enum generation works
	assertEqual(t, TriggerManual.String(), "Manual")

	// Test parsing
	trigger, err := ParseTrigger("Webhook")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	assertEqual(t, trigger, TriggerWebhook)

	// Test invalid parse
	_, err = ParseTrigger("Invalid")
	assertError(t, err, "invalid trigger")
}

// === DataPoint with Phase 2 fields ===

func TestDataPointWithTrigger(t *testing.T) {
	payload := TestPayload{Value: "webhook-trigger", Count: 1}
	actor := SystemActor[string]()

	dp := NewDataPointNow(payload, actor).WithTrigger(TriggerWebhook)

	assertEqual(t, dp.Trigger(), TriggerWebhook)
}

func TestDataPointWithContext(t *testing.T) {
	payload := TestPayload{Value: "context-test", Count: 1}
	actor := SystemActor[string]()
	ctx := NewContext("test-service").WithEnvironment("development")

	dp := NewDataPointNow(payload, actor).WithContext(ctx)

	assertEqual(t, dp.Context().Source(), "test-service")
}

func TestDataPointWithTemporal(t *testing.T) {
	payload := TestPayload{Value: "temporal-test", Count: 1}
	actor := SystemActor[string]()

	from := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	until := NewTimestamp(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC))
	recorded := Now()
	temporal := NewBitemporalWithRange(from, until, recorded)

	dp := NewDataPointNow(payload, actor).WithTemporal(temporal)

	assertTrue(t, dp.Temporal().ValidFrom().Equal(from.Time), "temporal validFrom mismatch")
	assertTrue(t, dp.Temporal().ValidUntil().Equal(until.Time), "temporal validUntil mismatch")
}

func TestDataPointTemporalBackwardCompat(t *testing.T) {
	// Occurred() should delegate to Temporal().ValidFrom()
	payload := TestPayload{Value: "compat", Count: 1}
	actor := SystemActor[string]()
	occurred := NewTimestamp(time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC))
	recorded := NewTimestamp(time.Date(2024, 6, 15, 10, 1, 0, 0, time.UTC))

	dp := NewDataPoint(payload, actor, occurred, recorded)

	assertTrue(t, dp.Occurred().Equal(occurred.Time), "Occurred() should equal passed occurred time")
	assertTrue(t, dp.Recorded().Equal(recorded.Time), "Recorded() should equal passed recorded time")
}

// === Reference Tests ===

func TestReference(t *testing.T) {
	ref := NewReference("order-123", "parent")

	assertEqual(t, ref.Id(), "order-123")
	assertEqual(t, ref.Relation(), "parent")
	assertEqual(t, ref.Version(), 0)
}

func TestReferenceWithVersion(t *testing.T) {
	ref := NewReferenceWithVersion("doc-456", "source", 5)
	assertEqual(t, ref.Version(), 5)
}

func TestReferenceTags(t *testing.T) {
	ref := NewReference("user-789", "owner").
		WithTag("department", "engineering").
		WithTag("level", "senior")

	assertTag(t, ref.Tags(), "department", "engineering")

	v, ok := ref.Tag("level")
	assertTrue(t, ok && v == "senior", "tag level mismatch")

	_, ok = ref.Tag("nonexistent")
	assertFalse(t, ok, "nonexistent tag should return false")
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

	assertEqual(t, parsed.Id(), "entity-1")
	assertEqual(t, parsed.Relation(), "reference")
	assertEqual(t, parsed.Version(), 3)
	assertTag(t, parsed.Tags(), "key", "value")
}

// === Cause Tests ===

func TestCause(t *testing.T) {
	id := NewNanoId()
	cause := NewCause(id, "command", "created")

	assertEqual(t, cause.Id().String(), id.String())
	assertEqual(t, cause.Kind(), "command")
	assertEqual(t, cause.Effect(), "created")
	assertFalse(t, cause.HasTrace(), "should not have trace")
}

func TestCauseDirect(t *testing.T) {
	id := NewNanoId()
	cause := NewCauseDirect(id)

	assertEqual(t, cause.Kind(), "direct")
	assertEqual(t, cause.Effect(), "caused")
}

func TestCauseCommand(t *testing.T) {
	id := NewNanoId()
	cause := NewCauseCommand(id, "approved")

	assertEqual(t, cause.Kind(), "command")
	assertEqual(t, cause.Effect(), "approved")
}

func TestCauseEvent(t *testing.T) {
	id := NewNanoId()
	cause := NewCauseEvent(id, "triggered")

	assertEqual(t, cause.Kind(), "event")
	assertEqual(t, cause.Effect(), "triggered")
}

func TestCauseWithTrace(t *testing.T) {
	id1 := NewNanoId()
	id2 := NewNanoId()
	id3 := NewNanoId()

	cause := NewCauseDirect(id1).
		WithTrace([]NanoId{id2, id3})

	assertTrue(t, cause.HasTrace(), "should have trace")
	assertEqual(t, len(cause.Trace()), 2)
}

func TestCauseAppendTrace(t *testing.T) {
	id1 := NewNanoId()
	id2 := NewNanoId()

	cause := NewCauseDirect(id1).AppendTrace(id2)

	assertTrue(t, cause.HasTrace(), "should have trace")
	assertEqual(t, cause.Trace()[0].String(), id2.String())
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

	assertEqual(t, parsed.Id().String(), id.String())
	assertEqual(t, parsed.Kind(), "command")
	assertEqual(t, parsed.Effect(), "executed")
	assertTrue(t, parsed.HasTrace(), "should have trace")
}

// === DataPoint Phase 3 Tests ===

func TestDataPointWithVersion(t *testing.T) {
	dp := NewDataPointNow(TestPayload{Value: "versioned", Count: 1}, SystemActor[string]()).
		WithVersion(5)

	assertEqual(t, dp.Version(), 5)
}

func TestDataPointWithTags(t *testing.T) {
	tags := map[string]string{"env": "prod", "region": "us-east-1"}
	dp := NewDataPointNow(TestPayload{Value: "tagged", Count: 1}, SystemActor[string]()).
		WithTags(tags)

	assertTag(t, dp.Tags(), "env", "prod")
	assertTag(t, dp.Tags(), "region", "us-east-1")
}

func TestDataPointAddTag(t *testing.T) {
	dp := NewDataPointNow(TestPayload{Value: "add-tag", Count: 1}, SystemActor[string]()).
		AddTag("key1", "value1").
		AddTag("key2", "value2")

	assertTag(t, dp.Tags(), "key1", "value1")
	assertTag(t, dp.Tags(), "key2", "value2")

	v, ok := dp.Tag("key1")
	assertTrue(t, ok && v == "value1", "Tag() method failed")
}

func TestDataPointWithReferences(t *testing.T) {
	ref1 := NewReference("order-1", "parent")
	ref2 := NewReference("customer-2", "owner")

	dp := NewDataPointNow(TestPayload{Value: "refs", Count: 1}, SystemActor[string]()).
		WithReferences([]Reference[string]{ref1, ref2})

	assertEqual(t, len(dp.References()), 2)
	assertEqual(t, dp.References()[0].Id(), "order-1")
}

func TestDataPointAddReference(t *testing.T) {
	ref := NewReference("product-123", "item")

	dp := NewDataPointNow(TestPayload{Value: "add-ref", Count: 1}, SystemActor[string]()).
		AddReference(ref)

	assertEqual(t, len(dp.References()), 1)
	assertEqual(t, dp.References()[0].Id(), "product-123")
}

func TestDataPointWithCauses(t *testing.T) {
	cause1 := NewCauseDirect(NewNanoId())
	cause2 := NewCauseCommand(NewNanoId(), "approved")

	dp := NewDataPointNow(TestPayload{Value: "causes", Count: 1}, SystemActor[string]()).
		WithCauses([]Cause[string]{cause1, cause2})

	assertEqual(t, len(dp.Causes()), 2)
}

func TestDataPointAddCause(t *testing.T) {
	cause := NewCauseEvent(NewNanoId(), "triggered")

	dp := NewDataPointNow(TestPayload{Value: "add-cause", Count: 1}, SystemActor[string]()).
		AddCause(cause)

	assertEqual(t, len(dp.Causes()), 1)
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

	assertEqual(t, parsed.Version(), 3)
	assertTag(t, parsed.Tags(), "env", "test")
	assertEqual(t, len(parsed.References()), 1)
	assertEqual(t, parsed.References()[0].Id(), "doc-1")
	assertEqual(t, parsed.References()[0].Version(), 2)
	assertEqual(t, len(parsed.Causes()), 1)
	assertEqual(t, parsed.Causes()[0].Kind(), "command")
	assertTrue(t, parsed.Causes()[0].HasTrace(), "cause should have trace")
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

	assertJSONEquals(t, data, string(data2))
}
