package id

import (
	"context"
	"testing"
	"time"
)

func TestSonyflake(t *testing.T) {
	sf := NewSonyflake(WithSonyflakeMachineId(1))
	t.Log(sf.Id(context.Background()))
}

func TestNewSonyflake_DefaultOptions(t *testing.T) {
	sf := NewSonyflake()
	if sf == nil {
		t.Fatal("expected non-nil Sonyflake instance")
	}
	if sf.Error != nil {
		t.Fatalf("unexpected error: %v", sf.Error)
	}
	id, err := sf.NextID()
	if err != nil {
		t.Fatalf("NextID failed: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero id")
	}
}

func TestNewSonyflake_WithCustomMachineID(t *testing.T) {
	sf := NewSonyflake(WithSonyflakeMachineId(2))
	if sf == nil {
		t.Fatal("expected non-nil Sonyflake instance")
	}
	if sf.Error != nil {
		t.Fatalf("unexpected error: %v", sf.Error)
	}
	id, err := sf.NextID()
	if err != nil {
		t.Fatalf("NextID failed: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero id")
	}
}

func TestNewSonyflake_InvalidStartTime(t *testing.T) {
	// Set startTime to future to trigger "invalid start time"
	sf := NewSonyflake(func(opts *SonyflakeOptions) {
		opts.startTime = time.Now().Add(time.Hour)
	})
	if sf == nil {
		t.Fatal("expected non-nil Sonyflake instance")
	}
	if sf.Error == nil {
		t.Error("expected error due to invalid start time, got nil")
	}
}
