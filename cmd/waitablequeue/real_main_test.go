package main

import (
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fgm/container/queue"
)

const (
	WQCap    = 60
	WQLow    = 20
	WQHigh   = 30
	MainWait = 4 * time.Second
)

// threadSafeWriter is a thread-safe wrapper around a [strings.Builder].
type threadSafeWriter struct {
	mu sync.Mutex
	w  *strings.Builder
}

func (w *threadSafeWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.w.String()
}

// Write implements io.Writer.
func (w *threadSafeWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.w.Write(p)
}

// NewThreadSafeWriter creates a new threadSafeWriter
func NewThreadSafeWriter() io.Writer {
	return &threadSafeWriter{w: &strings.Builder{}}
}

// TestRealMain tests the main function as implemented
func TestRealMain(t *testing.T) {
	t.Parallel()
	var w = NewThreadSafeWriter()

	// Run realMain with a timeout to prevent the test from running too long
	doneCh := make(chan int)
	go func() {
		exitCode := realMain(w, WQCap, WQLow, WQHigh)
		doneCh <- exitCode
	}()

	// Wait at most MainWait seconds (the original function uses a 3-seconds delay)
	var exitCode int
	select {
	case exitCode = <-doneCh:
		// realMain completed normally
	case <-time.After(MainWait):
		t.Fatalf("Test exceeded %v timeout", MainWait)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Verify output
	output := w.(*threadSafeWriter).String()

	// Check that producer() sent messages
	if !strings.Contains(output, "Sent:") {
		t.Error("No producer output detected")
	}

	// Check that consumer received messages
	if !strings.Contains(output, "Received:") {
		t.Error("No consumer output detected")
	}

	// Check that consumer exited properly
	if !strings.Contains(output, "Consumer exiting") {
		t.Error("Consumer exit message not detected")
	}
}

// Create a custom writer that captures the error
type errorCapturingWriter struct {
	io.Writer
	capturedErr error
	mx          sync.Mutex
}

// Override the Write method to extract the error
func (w *errorCapturingWriter) Write(p []byte) (n int, err error) {
	w.mx.Lock()
	defer w.mx.Unlock()
	// Check if this is an error message
	output := string(p)
	if strings.Contains(output, "Failed to create WaitableQueue:") {
		// Extract the error part
		errMsg := strings.TrimPrefix(output, "Failed to create WaitableQueue: ")
		errMsg = strings.TrimSpace(errMsg)

		// Try to match with known errors
		errs := []error{
			queue.ErrCapacityIsNegative,
			queue.ErrHighWatermarkIsLessThanLowWatermark,
			queue.ErrHighWatermarkIsNegative,
			queue.ErrLowWatermarkIsNegative,
		}

		// Try to match with known errors
		for _, err := range errs {
			if strings.Contains(errMsg, err.Error()) {
				w.capturedErr = err
				break
			}
		}
	}
	return w.Writer.Write(p)
}

// TestRealMainErrorHandling tests error handling in realMain by providing invalid parameters
// that trigger the predefined error values from the queue package
func TestRealMainErrorHandling(t *testing.T) {
	t.Parallel()

	// Test cases with parameters that should trigger specific errors
	tests := []struct {
		name          string
		sizeHint      int
		low           int
		high          int
		expectedError error
	}{
		{
			name:          "NegativeCapacity",
			sizeHint:      -1,
			low:           WQLow,
			high:          WQHigh,
			expectedError: queue.ErrCapacityIsNegative,
		},
		{
			name:          "NegativeLowWatermark",
			sizeHint:      WQCap,
			low:           -10,
			high:          WQHigh,
			expectedError: queue.ErrLowWatermarkIsNegative,
		},
		{
			name:          "NegativeHighWatermark",
			sizeHint:      WQCap,
			low:           WQLow,
			high:          -5,
			expectedError: queue.ErrHighWatermarkIsNegative,
		},
		{
			name:          "HighLessThanLow",
			sizeHint:      WQCap,
			low:           WQHigh + 1,
			high:          WQHigh,
			expectedError: queue.ErrHighWatermarkIsLessThanLowWatermark,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var w strings.Builder
			writer := &errorCapturingWriter{Writer: &w}

			// Call realMain with test case parameters
			exitCode := realMain(writer, test.sizeHint, test.low, test.high)

			// Check that we got an error exit code
			if exitCode != 1 {
				t.Errorf("Expected exit code 1 for error, got %d", exitCode)
			}

			// Check for an error message
			output := w.String()
			if !strings.Contains(output, "Failed to create WaitableQueue") {
				t.Error("Missing error message header")
			}

			// Actual errors were extracted by the errorCapturingWriter.
			if writer.capturedErr == nil {
				t.Errorf("No error was captured")
			} else if !errors.Is(writer.capturedErr, test.expectedError) {
				t.Errorf("Expected error %v, got %v", test.expectedError, writer.capturedErr)
			}
		})
	}
}

// TestRealMainWithTimeout tests the behavior of realMain with a timeout
func TestRealMainWithTimeout(t *testing.T) {
	t.Parallel()
	var w = NewThreadSafeWriter()

	// Run realMain with a shorter timeout
	doneCh := make(chan int)
	go func() {
		exitCode := realMain(w, WQCap, WQLow, WQHigh)
		doneCh <- exitCode
	}()

	// Wait at most 500 msec before checking partial results
	// This allows testing behavior during execution without waiting for full completion
	select {
	case <-doneCh:
		t.Log("realMain completed earlier than expected")
	case <-time.After(500 * time.Millisecond):
		// This is expected - realMain normally takes 3+ seconds
	}

	// Check partial output
	output := w.(*threadSafeWriter).String()

	// At this point, the producer should have sent at least a few messages
	if !strings.Contains(output, "Sent:") {
		t.Error("No producer output detected after 500ms")
	}

	// The consumer should have received at least a few messages
	if !strings.Contains(output, "Received:") {
		t.Error("No consumer output detected after 500ms")
	}

	// Cleanup - wait for realMain to finish to avoid zombie goroutines
	select {
	case <-doneCh:
		// realMain completed
	case <-time.After(5 * time.Second):
		t.Log("Timeout waiting for realMain to complete")
	}
}

// TestRealMainOutputContents tests the specific content of realMain's output
func TestRealMainOutputContents(t *testing.T) {
	t.Parallel()
	var w = NewThreadSafeWriter()

	// Run realMain with timeout
	doneCh := make(chan int)
	go func() {
		exitCode := realMain(w, WQCap, WQLow, WQHigh)
		doneCh <- exitCode
	}()

	// Wait for realMain to complete or timeout after 5 seconds
	select {
	case <-doneCh:
		// realMain completed normally
	case <-time.After(MainWait):
		t.Fatal("Test exceeded WQWait timeout", MainWait)
	}

	// Detailed output verification
	output := w.(*threadSafeWriter).String()
	lines := strings.Split(output, "\n")

	// Count sent and received messages
	sentCount := 0
	receivedCount := 0

	for _, line := range lines {
		if strings.Contains(line, "Sent:") {
			sentCount++
		}
		if strings.Contains(line, "Received:") {
			receivedCount++
		}
	}

	// Verify multiple messages were sent (at least 10)
	if sentCount < 10 {
		t.Errorf("Insufficient number of sent messages: %d (expected at least 10)", sentCount)
	}

	// Verify multiple messages were received (at least 5)
	if receivedCount < 5 {
		t.Errorf("Insufficient number of received messages: %d (expected at least 5)", receivedCount)
	}

	// Verify presence of messages indicating different queue states
	states := []string{
		"QueueIsBelowLowWatermark",
		"QueueIsNominal",
		"QueueIsAboveHighWatermark",
	}

	for _, state := range states {
		if !strings.Contains(output, state) {
			t.Errorf("Missing queue state in output: %s", state)
		}
	}

	// Verify consumer properly exited
	if !strings.Contains(output, "Consumer exiting on WaitableQueue closure") {
		t.Error("Proper consumer exit message not detected")
	}
}
