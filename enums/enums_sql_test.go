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

func TestCauseKindSQL(t *testing.T) {
	t.Parallel()
	val, err := CauseKindDirect.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "Direct" {
		t.Errorf("expected Direct, got %v", val)
	}

	var ck CauseKind
	if err := ck.Scan("Command"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ck != CauseKindCommand {
		t.Errorf("expected Command, got %v", ck)
	}

	// Test Scan with []byte
	var ck2 CauseKind
	if err := ck2.Scan([]byte("Event")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ck2 != CauseKindEvent {
		t.Errorf("expected Event, got %v", ck2)
	}

	// Test Scan with nil
	ck3 := CauseKindDirect
	if err := ck3.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ck3 != 0 {
		t.Errorf("expected zero value, got %v", ck3)
	}

	// Test Scan with int
	var ck4 CauseKind
	if err := ck4.Scan(1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ck4 != CauseKindCommand {
		t.Errorf("expected CauseKindCommand(1), got %v", ck4)
	}
}

func TestCauseKindMarshal(t *testing.T) {
	t.Parallel()
	data, err := CauseKindEvent.MarshalText()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(data) != "Event" {
		t.Errorf("expected Event, got %s", string(data))
	}

	var ck CauseKind
	if err := ck.UnmarshalText([]byte("Direct")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ck != CauseKindDirect {
		t.Errorf("expected Direct, got %v", ck)
	}
}

// Test AppendText methods
func TestActorKindAppendText(t *testing.T) {
	t.Parallel()
	var buf []byte
	ak := ActorKindUser
	n, err := ak.AppendText(buf)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(n) != "User" {
		t.Errorf("expected User, got %s", string(n))
	}

	// Test with existing buffer
	buf = []byte("prefix:")
	ak2 := ActorKindBot
	n, err = ak2.AppendText(buf)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(n) != "prefix:Bot" {
		t.Errorf("expected prefix:Bot, got %s", string(n))
	}
}

func TestPriorityAppendText(t *testing.T) {
	t.Parallel()
	p := PriorityHigh
	n, err := p.AppendText(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(n) != "High" {
		t.Errorf("expected High, got %s", string(n))
	}
}

func TestStatusAppendText(t *testing.T) {
	t.Parallel()
	s := StatusActive
	n, err := s.AppendText(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(n) != "Active" {
		t.Errorf("expected Active, got %s", string(n))
	}
}

func TestTriggerAppendText(t *testing.T) {
	t.Parallel()
	tr := TriggerWebhook
	n, err := tr.AppendText(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(n) != "Webhook" {
		t.Errorf("expected Webhook, got %s", string(n))
	}
}

func TestCauseKindAppendText(t *testing.T) {
	t.Parallel()
	ck := CauseKindDirect
	n, err := ck.AppendText(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(n) != "Direct" {
		t.Errorf("expected Direct, got %s", string(n))
	}
}

// Test invalid enum String() formatting
func TestInvalidEnumStrings(t *testing.T) {
	t.Parallel()
	t.Run("ActorKind", func(t *testing.T) {
		t.Parallel()
		invalid := ActorKind(99)
		expected := "ActorKind(99)"
		if invalid.String() != expected {
			t.Errorf("expected %s, got %s", expected, invalid.String())
		}
	})

	t.Run("Priority", func(t *testing.T) {
		t.Parallel()
		invalid := Priority(99)
		expected := "Priority(99)"
		if invalid.String() != expected {
			t.Errorf("expected %s, got %s", expected, invalid.String())
		}
	})

	t.Run("Status", func(t *testing.T) {
		t.Parallel()
		invalid := Status(99)
		expected := "Status(99)"
		if invalid.String() != expected {
			t.Errorf("expected %s, got %s", expected, invalid.String())
		}
	})

	t.Run("Trigger", func(t *testing.T) {
		t.Parallel()
		invalid := Trigger(99)
		expected := "Trigger(99)"
		if invalid.String() != expected {
			t.Errorf("expected %s, got %s", expected, invalid.String())
		}
	})

	t.Run("CauseKind", func(t *testing.T) {
		t.Parallel()
		invalid := CauseKind(99)
		expected := "CauseKind(99)"
		if invalid.String() != expected {
			t.Errorf("expected %s, got %s", expected, invalid.String())
		}
	})
}

// Test UnmarshalText error cases
func TestUnmarshalTextErrors(t *testing.T) {
	t.Parallel()
	t.Run("ActorKind", func(t *testing.T) {
		t.Parallel()
		var ak ActorKind
		err := ak.UnmarshalText([]byte("Invalid"))
		if err == nil {
			t.Error("expected error for invalid ActorKind")
		}
	})

	t.Run("Priority", func(t *testing.T) {
		t.Parallel()
		var p Priority
		err := p.UnmarshalText([]byte("Invalid"))
		if err == nil {
			t.Error("expected error for invalid Priority")
		}
	})

	t.Run("Status", func(t *testing.T) {
		t.Parallel()
		var s Status
		err := s.UnmarshalText([]byte("Invalid"))
		if err == nil {
			t.Error("expected error for invalid Status")
		}
	})

	t.Run("Trigger", func(t *testing.T) {
		t.Parallel()
		var tr Trigger
		err := tr.UnmarshalText([]byte("Invalid"))
		if err == nil {
			t.Error("expected error for invalid Trigger")
		}
	})

	t.Run("CauseKind", func(t *testing.T) {
		t.Parallel()
		var ck CauseKind
		err := ck.UnmarshalText([]byte("Invalid"))
		if err == nil {
			t.Error("expected error for invalid CauseKind")
		}
	})
}

// Test all Value methods
func TestValueMethods(t *testing.T) {
	t.Parallel()
	t.Run("ActorKind", func(t *testing.T) {
		t.Parallel()
		val, err := ActorKindSystem.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "System" {
			t.Errorf("expected System, got %v", val)
		}
	})

	t.Run("Priority", func(t *testing.T) {
		t.Parallel()
		val, err := PriorityMedium.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "Medium" {
			t.Errorf("expected Medium, got %v", val)
		}
	})

	t.Run("Status", func(t *testing.T) {
		t.Parallel()
		val, err := StatusPaused.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "Paused" {
			t.Errorf("expected Paused, got %v", val)
		}
	})

	t.Run("Trigger", func(t *testing.T) {
		t.Parallel()
		val, err := TriggerImport.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "Import" {
			t.Errorf("expected Import, got %v", val)
		}
	})

	t.Run("CauseKind", func(t *testing.T) {
		t.Parallel()
		val, err := CauseKindEvent.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "Event" {
			t.Errorf("expected Event, got %v", val)
		}
	})
}

// Test comprehensive Scan types for all enums
func TestActorKindScanAllTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		src   any
		want  ActorKind
		valid bool
	}{
		{"int64", int64(1), ActorKindBot, true},
		{"string", "System", ActorKindSystem, true},
		{"bytes", []byte("Service"), ActorKindService, true},
		{"int", int(0), ActorKindUser, true},
		{"uint", uint(2), ActorKindSystem, true},
		{"uint64", uint64(3), ActorKindService, true},
		{"float64", float64(1), ActorKindBot, true},
		{"nil", nil, ActorKind(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var ak ActorKind
			err := ak.Scan(tt.src)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if ak != tt.want {
				t.Errorf("expected %v, got %v", tt.want, ak)
			}
		})
	}
}

func TestPriorityScanAllTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		src   any
		want  Priority
		valid bool
	}{
		{"int64", int64(2), PriorityHigh, true},
		{"string", "Critical", PriorityCritical, true},
		{"bytes", []byte("Low"), PriorityLow, true},
		{"int", int(1), PriorityMedium, true},
		{"uint", uint(0), PriorityLow, true},
		{"uint64", uint64(3), PriorityCritical, true},
		{"float64", float64(2), PriorityHigh, true},
		{"nil", nil, Priority(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var p Priority
			err := p.Scan(tt.src)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if p != tt.want {
				t.Errorf("expected %v, got %v", tt.want, p)
			}
		})
	}
}

func TestStatusScanAllTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		src   any
		want  Status
		valid bool
	}{
		{"int64", int64(1), StatusActive, true},
		{"string", "Archived", StatusArchived, true},
		{"bytes", []byte("Deleted"), StatusDeleted, true},
		{"int", int(0), StatusDraft, true},
		{"uint", uint(2), StatusPaused, true},
		{"uint64", uint64(4), StatusDeleted, true},
		{"float64", float64(1), StatusActive, true},
		{"nil", nil, Status(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var s Status
			err := s.Scan(tt.src)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if s != tt.want {
				t.Errorf("expected %v, got %v", tt.want, s)
			}
		})
	}
}

func TestTriggerScanAllTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		src   any
		want  Trigger
		valid bool
	}{
		{"int64", int64(2), TriggerWebhook, true},
		{"string", "Correction", TriggerCorrection, true},
		{"bytes", []byte("Import"), TriggerImport, true},
		{"int", int(0), TriggerManual, true},
		{"uint", uint(5), TriggerSystem, true},
		{"uint64", uint64(6), TriggerCorrection, true},
		{"float64", float64(2), TriggerWebhook, true},
		{"nil", nil, Trigger(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var tr Trigger
			err := tr.Scan(tt.src)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tr != tt.want {
				t.Errorf("expected %v, got %v", tt.want, tr)
			}
		})
	}
}

func TestCauseKindScanAllTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		src   any
		want  CauseKind
		valid bool
	}{
		{"int64", int64(1), CauseKindCommand, true},
		{"string", "Event", CauseKindEvent, true},
		{"bytes", []byte("Direct"), CauseKindDirect, true},
		{"int", int(0), CauseKindDirect, true},
		{"uint", uint(2), CauseKindEvent, true},
		{"uint64", uint64(1), CauseKindCommand, true},
		{"float64", float64(1), CauseKindCommand, true},
		{"nil", nil, CauseKind(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var ck CauseKind
			err := ck.Scan(tt.src)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if ck != tt.want {
				t.Errorf("expected %v, got %v", tt.want, ck)
			}
		})
	}
}
