package main

import (
	"testing"
	"time"
)

func TestNewReindexScheduler(t *testing.T) {
	schedule, err := newReindexScheduler("12h", func() {})
	if err != nil {
		t.Fatalf("create scheduler: %v", err)
	}
	t.Cleanup(func() {
		if err := schedule.Shutdown(); err != nil {
			t.Errorf("shutdown scheduler: %v", err)
		}
	})

	if jobs := schedule.Jobs(); len(jobs) != 1 {
		t.Fatalf("expected one reindex job, got %d", len(jobs))
	}
}

func TestNewReindexSchedulerRejectsInvalidDuration(t *testing.T) {
	if _, err := newReindexScheduler("not-a-duration", func() {}); err == nil {
		t.Fatal("expected invalid duration to return an error")
	}
}

func TestNewReindexSchedulerStartsImmediately(t *testing.T) {
	run := make(chan struct{}, 1)
	schedule, err := newReindexScheduler("12h", func() {
		run <- struct{}{}
	})
	if err != nil {
		t.Fatalf("create scheduler: %v", err)
	}
	t.Cleanup(func() {
		if err := schedule.Shutdown(); err != nil {
			t.Errorf("shutdown scheduler: %v", err)
		}
	})

	schedule.Start()

	select {
	case <-run:
	case <-time.After(time.Second):
		t.Fatal("expected reindex job to run immediately")
	}
}
