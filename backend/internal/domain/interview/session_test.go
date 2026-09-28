package interview

import "testing"

func TestInterviewSessionFiveRoundLifecycle(t *testing.T) {
	session, err := NewInterviewSession("user-1", "analysis-1")
	if err != nil {
		t.Fatalf("NewInterviewSession() error = %v", err)
	}
	if err := session.Start(); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if session.CurrentRound() != 1 || !session.CanAcceptAnswer() {
		t.Fatalf("started session state = round %d, canAccept %v", session.CurrentRound(), session.CanAcceptAnswer())
	}
	for expectedRound := 2; expectedRound <= MaxRounds; expectedRound++ {
		if err := session.AdvanceRound(); err != nil {
			t.Fatalf("AdvanceRound() to %d error = %v", expectedRound, err)
		}
	}
	if session.CanContinue() {
		t.Fatal("round five must not generate a sixth question")
	}
	if err := session.Complete(); err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if session.Status() != InterviewStatusCompleted || session.CanAcceptAnswer() {
		t.Fatalf("completed session state = %q, canAccept %v", session.Status(), session.CanAcceptAnswer())
	}
}

func TestInterviewSessionRejectsEarlyCompletionAndExtraRound(t *testing.T) {
	session, err := NewInterviewSession("user-1", "analysis-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := session.Start(); err != nil {
		t.Fatal(err)
	}
	if err := session.Complete(); err == nil {
		t.Fatal("Complete() before round five should fail")
	}
	for round := 2; round <= MaxRounds; round++ {
		if err := session.AdvanceRound(); err != nil {
			t.Fatal(err)
		}
	}
	if err := session.AdvanceRound(); err == nil {
		t.Fatal("AdvanceRound() after round five should fail")
	}
}
