package enums

import (
	"slices"
	"testing"
)

func TestActorKind(t *testing.T) {
	t.Parallel()
	// Test constants exist
	_ = ActorKindUser
	_ = ActorKindBot
	_ = ActorKindSystem
	_ = ActorKindService

	// Test String()
	if ActorKindUser.String() != "User" {
		t.Errorf("expected User, got %s", ActorKindUser.String())
	}
	if ActorKindSystem.String() != "System" {
		t.Errorf("expected System, got %s", ActorKindSystem.String())
	}

	// Test all ActorKind values
	tests := []struct {
		kind     ActorKind
		expected string
	}{
		{ActorKindUser, "User"},
		{ActorKindBot, "Bot"},
		{ActorKindSystem, "System"},
		{ActorKindService, "Service"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			if tt.kind.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.kind.String())
			}
		})
	}
}

func TestParseActorKind(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		want    ActorKind
		wantErr bool
	}{
		{"User", ActorKindUser, false},
		{"Bot", ActorKindBot, false},
		{"System", ActorKindSystem, false},
		{"Service", ActorKindService, false},
		{"Invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParseActorKind(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestParseActorKindError(t *testing.T) {
	t.Parallel()
	_, err := ParseActorKind("Invalid")
	if err == nil {
		t.Error("expected error for invalid ActorKind")
	}
}

func TestActorKindIsValid(t *testing.T) {
	t.Parallel()
	if !ActorKindUser.IsValid() {
		t.Error("ActorKindUser should be valid")
	}
	if !ActorKindService.IsValid() {
		t.Error("ActorKindService should be valid")
	}
	// Invalid value
	invalid := ActorKind(99)
	if invalid.IsValid() {
		t.Error("ActorKind(99) should not be valid")
	}
}

func TestPriority(t *testing.T) {
	t.Parallel()
	// Test constants
	_ = PriorityLow
	_ = PriorityMedium
	_ = PriorityHigh
	_ = PriorityCritical

	// Test ordering
	if PriorityLow >= PriorityHigh {
		t.Error("Low should be less than High")
	}

	// Test String() for all values
	tests := []struct {
		priority Priority
		expected string
	}{
		{PriorityLow, "Low"},
		{PriorityMedium, "Medium"},
		{PriorityHigh, "High"},
		{PriorityCritical, "Critical"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			if tt.priority.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.priority.String())
			}
		})
	}
}

func TestParsePriority(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		want    Priority
		wantErr bool
	}{
		{"Low", PriorityLow, false},
		{"Medium", PriorityMedium, false},
		{"High", PriorityHigh, false},
		{"Critical", PriorityCritical, false},
		{"Invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParsePriority(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestPriorityIsValid(t *testing.T) {
	t.Parallel()
	if !PriorityLow.IsValid() {
		t.Error("PriorityLow should be valid")
	}
	if !PriorityCritical.IsValid() {
		t.Error("PriorityCritical should be valid")
	}
	invalid := Priority(99)
	if invalid.IsValid() {
		t.Error("Priority(99) should not be valid")
	}
}

func TestStatus(t *testing.T) {
	t.Parallel()
	// Test all statuses exist
	statuses := []Status{StatusDraft, StatusActive, StatusPaused, StatusArchived, StatusDeleted}
	expected := []string{"Draft", "Active", "Paused", "Archived", "Deleted"}

	for i, s := range statuses {
		if s.String() != expected[i] {
			t.Errorf("Status %v should have string %s, got %s", s, expected[i], s.String())
		}
	}
}

func TestParseStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		want    Status
		wantErr bool
	}{
		{"Draft", StatusDraft, false},
		{"Active", StatusActive, false},
		{"Paused", StatusPaused, false},
		{"Archived", StatusArchived, false},
		{"Deleted", StatusDeleted, false},
		{"Invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParseStatus(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestStatusIsValid(t *testing.T) {
	t.Parallel()
	validStatuses := []Status{
		StatusDraft,
		StatusActive,
		StatusPaused,
		StatusArchived,
		StatusDeleted,
	}
	for _, s := range validStatuses {
		if !s.IsValid() {
			t.Errorf("Status %v should be valid", s)
		}
	}
	invalid := Status(99)
	if invalid.IsValid() {
		t.Error("Status(99) should not be valid")
	}
}

func TestTrigger(t *testing.T) {
	t.Parallel()
	// Test all triggers
	tests := []struct {
		trigger  Trigger
		expected string
	}{
		{TriggerManual, "Manual"},
		{TriggerScheduled, "Scheduled"},
		{TriggerWebhook, "Webhook"},
		{TriggerImport, "Import"},
		{TriggerMigration, "Migration"},
		{TriggerSystem, "System"},
		{TriggerCorrection, "Correction"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			if tt.trigger.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.trigger.String())
			}
		})
	}
}

func TestParseTrigger(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		want    Trigger
		wantErr bool
	}{
		{"Manual", TriggerManual, false},
		{"Scheduled", TriggerScheduled, false},
		{"Webhook", TriggerWebhook, false},
		{"Import", TriggerImport, false},
		{"Migration", TriggerMigration, false},
		{"System", TriggerSystem, false},
		{"Correction", TriggerCorrection, false},
		{"Invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParseTrigger(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestTriggerIsValid(t *testing.T) {
	t.Parallel()
	validTriggers := []Trigger{
		TriggerManual, TriggerScheduled, TriggerWebhook, TriggerImport,
		TriggerMigration, TriggerSystem, TriggerCorrection,
	}
	for _, tr := range validTriggers {
		if !tr.IsValid() {
			t.Errorf("Trigger %v should be valid", tr)
		}
	}
	invalid := Trigger(99)
	if invalid.IsValid() {
		t.Error("Trigger(99) should not be valid")
	}
}

func TestEnumValues(t *testing.T) {
	t.Parallel()
	// Test ActorKindValues
	akValues := ActorKindValues()
	if len(akValues) != 4 {
		t.Errorf("expected 4 ActorKind values, got %d", len(akValues))
	}

	// Test PriorityValues
	pValues := PriorityValues()
	if len(pValues) != 4 {
		t.Errorf("expected 4 Priority values, got %d", len(pValues))
	}

	// Test StatusValues
	sValues := StatusValues()
	if len(sValues) != 5 {
		t.Errorf("expected 5 Status values, got %d", len(sValues))
	}

	// Test TriggerValues
	tValues := TriggerValues()
	if len(tValues) != 7 {
		t.Errorf("expected 7 Trigger values, got %d", len(tValues))
	}
}

func TestEnumNames(t *testing.T) {
	t.Parallel()
	// Test ActorKindNames
	akNames := ActorKindNames()
	if len(akNames) != 4 {
		t.Errorf("expected 4 ActorKind names, got %d", len(akNames))
	}
	if !slices.Contains(akNames, "User") {
		t.Error("ActorKindNames should include 'User'")
	}

	// Test PriorityNames
	pNames := PriorityNames()
	if len(pNames) != 4 {
		t.Errorf("expected 4 Priority names, got %d", len(pNames))
	}

	// Test StatusNames
	sNames := StatusNames()
	if len(sNames) != 5 {
		t.Errorf("expected 5 Status names, got %d", len(sNames))
	}

	// Test TriggerNames
	tNames := TriggerNames()
	if len(tNames) != 7 {
		t.Errorf("expected 7 Trigger names, got %d", len(tNames))
	}
}

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

	// Note: The generated Scan method accepts many types but silently ignores unsupported ones
	// This is behavior of go-enum generated code, not something we should test
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
