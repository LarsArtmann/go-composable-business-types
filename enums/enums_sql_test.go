package enums

import (
	"testing"
)

// Test SQL Scanner/Valuer interfaces
func TestActorKindSQL(t *testing.T) {
	t.Parallel()
	// Test Value
	val, err := ActorKindUser.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "User" {
		t.Errorf("expected User, got %v", val)
	}

	// Test Scan with string
	var ak ActorKind
	if err := ak.Scan("Bot"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ak != ActorKindBot {
		t.Errorf("expected Bot, got %v", ak)
	}

	// Test Scan with []byte
	var ak2 ActorKind
	if err := ak2.Scan([]byte("System")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ak2 != ActorKindSystem {
		t.Errorf("expected System, got %v", ak2)
	}

	// Test Scan with nil
	ak3 := ActorKindUser
	if err := ak3.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ak3 != 0 {
		t.Errorf("expected zero value, got %v", ak3)
	}

	// Test Scan with int (supported type - sets raw value)
	var ak4 ActorKind
	if err := ak4.Scan(2); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ak4 != ActorKindSystem {
		t.Errorf("expected ActorKindSystem(2), got %v", ak4)
	}

	// Test Scan with int64 (supported type)
	var ak5 ActorKind
	if err := ak5.Scan(int64(3)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ak5 != ActorKindService {
		t.Errorf("expected ActorKindService(3), got %v", ak5)
	}
}

func TestPrioritySQL(t *testing.T) {
	t.Parallel()
	val, err := PriorityHigh.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "High" {
		t.Errorf("expected High, got %v", val)
	}

	var p Priority
	if err := p.Scan("Critical"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p != PriorityCritical {
		t.Errorf("expected Critical, got %v", p)
	}
}

func TestStatusSQL(t *testing.T) {
	t.Parallel()
	val, err := StatusActive.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "Active" {
		t.Errorf("expected Active, got %v", val)
	}

	var s Status
	if err := s.Scan("Archived"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if s != StatusArchived {
		t.Errorf("expected Archived, got %v", s)
	}
}

func TestTriggerSQL(t *testing.T) {
	t.Parallel()
	val, err := TriggerWebhook.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "Webhook" {
		t.Errorf("expected Webhook, got %v", val)
	}

	var tr Trigger
	if err := tr.Scan("Migration"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if tr != TriggerMigration {
		t.Errorf("expected Migration, got %v", tr)
	}
}

// Test MarshalText/UnmarshalText
func TestActorKindMarshal(t *testing.T) {
	t.Parallel()
	data, err := ActorKindService.MarshalText()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(data) != "Service" {
		t.Errorf("expected Service, got %s", string(data))
	}

	var ak ActorKind
	if err := ak.UnmarshalText([]byte("User")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ak != ActorKindUser {
		t.Errorf("expected User, got %v", ak)
	}
}

func TestPriorityMarshal(t *testing.T) {
	t.Parallel()
	data, err := PriorityLow.MarshalText()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(data) != "Low" {
		t.Errorf("expected Low, got %s", string(data))
	}

	var p Priority
	if err := p.UnmarshalText([]byte("Medium")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p != PriorityMedium {
		t.Errorf("expected Medium, got %v", p)
	}
}

func TestStatusMarshal(t *testing.T) {
	t.Parallel()
	data, err := StatusDeleted.MarshalText()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(data) != "Deleted" {
		t.Errorf("expected Deleted, got %s", string(data))
	}

	var s Status
	if err := s.UnmarshalText([]byte("Draft")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if s != StatusDraft {
		t.Errorf("expected Draft, got %v", s)
	}
}

func TestTriggerMarshal(t *testing.T) {
	t.Parallel()
	data, err := TriggerCorrection.MarshalText()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(data) != "Correction" {
		t.Errorf("expected Correction, got %s", string(data))
	}

	var tr Trigger
	if err := tr.UnmarshalText([]byte("Scheduled")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if tr != TriggerScheduled {
		t.Errorf("expected Scheduled, got %v", tr)
	}
}
