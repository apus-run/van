package errorsx

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	e := A()

	fmt.Printf("%+v", e)
	fmt.Println()
	fmt.Println(e.Error())
}

func A() error {
	return B()
}

func B() error {
	return C()
}

func C() error {
	return InternalServer("InternalServer").
		WithCause(fmt.Errorf("db connection error")).
		WithStack()
}

func TestError(t *testing.T) {
	e := New(400, "InvalidParams").
		WithMessage("Invalid username format").
		WithMetadata(map[string]string{
			"field":  "username",
			"format": "letters followed by 6 digits",
		})

	if e.Code != 400 || e.Reason != "InvalidParams" || e.Message != "Invalid username format" {
		t.Errorf("unexpected error: %+v", e)
	}

	if e.Metadata["field"] != "username" || e.Metadata["format"] != "letters followed by 6 digits" {
		t.Errorf("unexpected metadata: %+v", e.Metadata)
	}
}

func TestWithMessage(t *testing.T) {
	e := New(400, "InvalidParams").
		WithMessage("Invalid username format").
		WithMetadata(map[string]string{
			"field":  "username",
			"format": "letters followed by 6 digits",
		})

	if e.Code != 400 || e.Reason != "InvalidParams" || e.Message != "Invalid username format" {
		t.Errorf("unexpected error: %+v", e)
	}

	if e.Metadata["field"] != "username" || e.Metadata["format"] != "letters followed by 6 digits" {
		t.Errorf("unexpected metadata: %+v", e.Metadata)
	}
}

func TestWithMetadata(t *testing.T) {
	e := New(400, "InvalidParams").
		WithMessage("Invalid username format").
		WithMetadata(map[string]string{
			"field":  "username",
			"format": "letters followed by 6 digits",
		})

	if e.Code != 400 || e.Reason != "InvalidParams" || e.Message != "Invalid username format" {
		t.Errorf("unexpected error: %+v", e)
	}

	if e.Metadata["field"] != "username" || e.Metadata["format"] != "letters followed by 6 digits" {
		t.Errorf("unexpected metadata: %+v", e.Metadata)
	}
}

func TestWithCause(t *testing.T) {
	originalErr := fmt.Errorf("db connection error")
	e := New(500, "DBError").
		WithCause(originalErr).
		WithStack()

	if e.Cause != originalErr {
		t.Errorf("unexpected cause: %+v", e.Cause)
	}

	if len(e.Stack) == 0 {
		t.Error("stack trace should not be empty")
	}
}

func TestKV(t *testing.T) {
	e := New(400, "ParseError").
		WithMessage("JSON parsing failed").
		KV("input", "{invalid: json}").
		KV("service", "user-api")

	if e.Metadata["input"] != "{invalid: json}" || e.Metadata["service"] != "user-api" {
		t.Errorf("unexpected metadata: %+v", e.Metadata)
	}
}

func TestGRPCStatus(t *testing.T) {
	e := New(500, "InternalError").
		WithMessage("Something went wrong").
		WithMetadata(map[string]string{
			"key": "value",
		})

	st := e.GRPCStatus()
	if st.Message() != "InternalError: Something went wrong" {
		t.Errorf("unexpected gRPC status message: %s", st.Message())
	}

	details := st.Details()
	if len(details) != 1 {
		t.Errorf("unexpected gRPC status details: %+v", details)
	}
}

func TestFromError(t *testing.T) {
	originalErr := fmt.Errorf("standard error")
	e := FromError(originalErr)

	if e.Code != 500 || e.Reason != "InternalError" || e.Message != "standard error" {
		t.Errorf("unexpected error: %+v", e)
	}

	if e.Cause != originalErr {
		t.Errorf("unexpected cause: %+v", e.Cause)
	}
}
func TestEnableStackCapture(t *testing.T) {
	// Save the original value of EnableStackCapture to restore it later
	originalValue := EnableStackCapture
	defer func() { EnableStackCapture = originalValue }()

	// Test when EnableStackCapture is true
	EnableStackCapture = true
	e := New(500, "TestError").WithStack()
	if len(e.Stack) == 0 {
		t.Error("stack trace should not be empty when EnableStackCapture is true")
	}

	// Test when EnableStackCapture is false
	EnableStackCapture = false
	e = New(500, "TestError").WithStack()
	if len(e.Stack) != 0 {
		t.Error("stack trace should be empty when EnableStackCapture is false")
	}
}
