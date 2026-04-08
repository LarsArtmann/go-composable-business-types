package datapoint_test

import (
	"encoding/json"
	"iter"
	"maps"
	"testing"

	"github.com/larsartmann/go-composable-business-types/actor"
	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/id"
	"github.com/larsartmann/go-composable-business-types/nanoid"
)

func TestNewDataPoint(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID, "Test User")
	dp := NewDataPoint("test payload", actorEntry)

	if dp.IsZero() {
		t.Error("DataPoint should not be zero")
	}

	if dp.ID().IsZero() {
		t.Error("DataPoint should have an ID")
	}

	if dp.Payload() != "test payload" {
		t.Errorf("expected payload 'test payload', got %v", dp.Payload())
	}

	if dp.Actor().Name != "Test User" {
		t.Errorf("expected actor name 'Test User', got %s", dp.Actor().Name)
	}

	if dp.Version() != 1 {
		t.Errorf("expected version 1, got %d", dp.Version())
	}

	if dp.Trigger() != enums.TriggerManual {
		t.Errorf("expected trigger Manual, got %v", dp.Trigger())
	}
}

func TestDataPointWithMethods(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	dp := NewDataPoint("payload", actorEntry)

	// Test WithTrigger
	dp2 := dp.WithTrigger(enums.TriggerSystem)
	if dp2.Trigger() != enums.TriggerSystem {
		t.Error("WithTrigger failed")
	}
	// Original should be unchanged
	if dp.Trigger() != enums.TriggerManual {
		t.Error("original DataPoint should be unchanged")
	}

	// Test WithReason
	dp3 := dp.WithReason("test reason")
	if dp3.Reason() != "test reason" {
		t.Errorf("expected reason 'test reason', got %s", dp3.Reason())
	}

	// Test WithVersion
	dp4 := dp.WithVersion(5)
	if dp4.Version() != 5 {
		t.Errorf("expected version 5, got %d", dp4.Version())
	}

	// Test WithTag
	dp5 := dp.WithTag("key", "value")
	if dp5.Tag("key") != "value" {
		t.Errorf("expected tag 'value', got %s", dp5.Tag("key"))
	}

	// Test WithTags
	dp6 := dp.WithTags(map[string]string{"a": "1", "b": "2"})
	if dp6.Tag("a") != "1" || dp6.Tag("b") != "2" {
		t.Error("WithTags failed")
	}
}

func TestDataPointWithReference(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	dp := NewDataPoint("payload", actorEntry)

	ref := NewReference("ref-id", "parent")
	dp2 := dp.WithReference(ref)

	refs := dp2.References()
	if len(refs) != 1 {
		t.Errorf("expected 1 reference, got %d", len(refs))
	}

	if refs[0].Relation() != "parent" {
		t.Errorf("expected relation 'parent', got %s", refs[0].Relation())
	}
}

func TestDataPointWithCause(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	dp := NewDataPoint("payload", actorEntry)

	causeID := nanoid.New()
	cause := NewCauseDirect[string](causeID)
	dp2 := dp.WithCause(cause)

	causes := dp2.Causes()
	if len(causes) != 1 {
		t.Errorf("expected 1 cause, got %d", len(causes))
	}

	if causes[0].Kind() != enums.CauseKindDirect {
		t.Errorf("expected kind 'direct', got %s", causes[0].Kind())
	}
}

func TestDataPointWithContext(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	dp := NewDataPoint("payload", actorEntry)

	ctx := NewContext().
		WithEnvironment("production").
		WithSource("test-service").
		WithTag("region", "us-east-1")

	dp2 := dp.WithContext(ctx)
	if dp2.Context().Environment() != "production" {
		t.Errorf("expected environment 'production', got %s", dp2.Context().Environment())
	}

	if dp2.Context().Source() != "test-service" {
		t.Errorf("expected source 'test-service', got %s", dp2.Context().Source())
	}

	if dp2.Context().Tag("region") != "us-east-1" {
		t.Error("context tag not set correctly")
	}
}

func TestDataPointJSON(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID, "Test User")
	dp := NewDataPoint("test payload", actorEntry).
		WithReason("test reason").
		WithTag("key", "value")

	data, err := json.Marshal(dp)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Verify JSON contains expected fields
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	if _, ok := raw["id"]; !ok {
		t.Error("JSON should contain 'id'")
	}

	if _, ok := raw["payload"]; !ok {
		t.Error("JSON should contain 'payload'")
	}

	if _, ok := raw["actor"]; !ok {
		t.Error("JSON should contain 'actor'")
	}
}

func TestDataPointUnmarshalJSON(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID, "Test User")
	original := NewDataPoint("test payload", actorEntry).
		WithReason("test reason").
		WithVersion(42)

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var parsed DataPoint[string]
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if parsed.Payload() != "test payload" {
		t.Errorf("expected payload 'test payload', got %v", parsed.Payload())
	}

	if parsed.Reason() != "test reason" {
		t.Errorf("expected reason 'test reason', got %s", parsed.Reason())
	}

	if parsed.Version() != 42 {
		t.Errorf("expected version 42, got %d", parsed.Version())
	}

	if parsed.ID().IsZero() {
		t.Error("parsed DataPoint should have an ID")
	}
}

func TestDataPointIsZero(t *testing.T) {
	t.Parallel()

	var zero DataPoint[string]
	if !zero.IsZero() {
		t.Error("zero DataPoint should be zero")
	}

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)

	dp := NewDataPoint("payload", actorEntry)
	if dp.IsZero() {
		t.Error("non-zero DataPoint should not be zero")
	}
}

func TestDataPointIntPayload(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	dp := NewDataPoint(42, actorEntry)

	if dp.Payload() != 42 {
		t.Errorf("expected payload 42, got %d", dp.Payload())
	}

	// Test JSON round-trip
	data, err := json.Marshal(dp)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var parsed DataPoint[int]
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if parsed.Payload() != 42 {
		t.Errorf("expected payload 42 after unmarshal, got %d", parsed.Payload())
	}
}

func TestDataPointComplexChain(t *testing.T) {
	t.Parallel()
	// Create a complex DataPoint with all fields set
	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID, "John Doe")

	ctx := NewContext().
		WithEnvironment("production").
		WithSession("session-abc").
		WithRequest("req-xyz").
		WithSource("order-service")

	ref := NewReference("order-456", "parent-order").
		WithVersion(3).
		WithTag("type", "subscription")

	causeID := nanoid.New()
	cause := NewCauseCommand[string](causeID, "create-order")

	dp := NewDataPoint("widget-order", actorEntry).
		WithTrigger(enums.TriggerWebhook).
		WithReason("Customer checkout").
		WithContext(ctx).
		WithVersion(2).
		WithTag("priority", "high").
		WithTag("channel", "web").
		WithReference(ref).
		WithCause(cause)

	// Verify all fields
	if dp.Trigger() != enums.TriggerWebhook {
		t.Error("trigger mismatch")
	}

	if dp.Reason() != "Customer checkout" {
		t.Error("reason mismatch")
	}

	if dp.Context().Environment() != "production" {
		t.Error("environment mismatch")
	}

	if len(dp.References()) != 1 {
		t.Errorf("expected 1 reference, got %d", len(dp.References()))
	}

	if len(dp.Causes()) != 1 {
		t.Errorf("expected 1 cause, got %d", len(dp.Causes()))
	}

	if dp.Tag("priority") != "high" {
		t.Error("tag mismatch")
	}
}

func TestDataPointAllReferences(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	dp := NewDataPoint("payload", actorEntry).
		WithReference(NewReference("ref-1", "parent")).
		WithReference(NewReference("ref-2", "child"))

	ids := make([]string, 0, 2)
	for ref := range dp.AllReferences() {
		ids = append(ids, ref.ID())
	}

	if len(ids) != 2 {
		t.Fatalf("expected 2 references, got %d", len(ids))
	}

	if ids[0] != "ref-1" || ids[1] != "ref-2" {
		t.Errorf("expected [ref-1, ref-2], got %v", ids)
	}
}

func TestDataPointAllReferencesBreak(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	dp := NewDataPoint("payload", actorEntry).
		WithReference(NewReference("ref-1", "parent")).
		WithReference(NewReference("ref-2", "child"))

	var count int
	for range dp.AllReferences() {
		count++

		break
	}

	if count != 1 {
		t.Errorf("expected break after 1, got %d", count)
	}
}

func TestDataPointAllCauses(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	cause1 := NewCauseDirect[string](nanoid.New())
	cause2 := NewCauseCommand[string](nanoid.New(), "create")
	dp := NewDataPoint("payload", actorEntry).
		WithCause(cause1).
		WithCause(cause2)

	kinds := make([]string, 0, 2)
	for c := range dp.AllCauses() {
		kinds = append(kinds, c.Kind().String())
	}

	if len(kinds) != 2 {
		t.Fatalf("expected 2 causes, got %d", len(kinds))
	}
}

func TestDataPointAllTags(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	dp := NewDataPoint("payload", actorEntry).
		WithTag("env", "prod").
		WithTag("region", "us-east-1")

	tags := maps.Collect(dp.AllTags())

	if tags["env"] != "prod" {
		t.Errorf("expected env=prod, got %s", tags["env"])
	}

	if tags["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %s", tags["region"])
	}
}

func TestDataPointAllTagsBreak(t *testing.T) {
	t.Parallel()

	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	dp := NewDataPoint("payload", actorEntry).
		WithTag("a", "1").
		WithTag("b", "2")

	var count int
	for range dp.AllTags() {
		count++

		break
	}

	if count != 1 {
		t.Errorf("expected break after 1, got %d", count)
	}
}

func countIterator[T any](seq iter.Seq[T]) int {
	var count int
	for range seq {
		count++
	}

	return count
}

func countSeq2Iterator[K, V any](seq iter.Seq2[K, V]) int {
	var count int
	for range seq {
		count++
	}

	return count
}

func testDataPointIteratorEmpty(t *testing.T, name string, count int) {
	t.Helper()
	if count != 0 {
		t.Errorf("expected 0 iterations for %s, got %d", name, count)
	}
}

func newTestDataPointEmpty(t *testing.T) *DataPoint {
	t.Helper()
	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID)
	return NewDataPoint("payload", actorEntry)
}

func testDataPointSeqIteratorEmpty[T any](
	t *testing.T,
	name string,
	iter func(*DataPoint) iter.Seq[T],
) {
	t.Parallel()
	dp := newTestDataPointEmpty(t)
	testDataPointIteratorEmpty(t, name, countIterator(iter(dp)))
}

func TestDataPointAllReferencesEmpty(t *testing.T) {
	testDataPointSeqIteratorEmpty(
		t,
		"AllReferences",
		func(dp *DataPoint) iter.Seq[Reference[string]] { return dp.AllReferences() },
	)
}

func TestDataPointAllCausesEmpty(t *testing.T) {
	testDataPointSeqIteratorEmpty(
		t,
		"AllCauses",
		func(dp *DataPoint) iter.Seq[Cause[string]] { return dp.AllCauses() },
	)
}

func TestDataPointAllTagsEmpty(t *testing.T) {
	t.Parallel()
	dp := newTestDataPointEmpty(t)
	testDataPointIteratorEmpty(t, "AllTags", countSeq2Iterator(dp.AllTags()))
}
