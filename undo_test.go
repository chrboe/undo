package undo

import (
	"errors"
	"testing"
)

func assert(t *testing.T, actual, expect int) {
	if actual != expect {
		t.Fatalf("unexpected num value: expected %d, got %d", expect, actual)
	}
}

func TestUndo(t *testing.T) {
	num := 5

	u := New(func() error {
		num++
		return nil
	}, func() {
		num--
	})

	// execute the function
	err := u.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert(t, num, 6)

	// undo it
	u.Undo()
	assert(t, num, 5)

	// it is now committed, can't undo again
	u.Undo()
	assert(t, num, 5)

	// execute again, resets committed flag
	err = u.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert(t, num, 6)

	// this time, commit before undoing
	u.Commit()
	assert(t, num, 6)

	// this should do nothing no matter how often we call it
	u.Undo()
	u.Undo()
	assert(t, num, 6)
}

func TestExecuteChainError(t *testing.T) {
	errText := "oops, an error occurred"

	expectTranscript := `execute action1
execute action2
oops, action3 failed
rollback action2
rollback action1
`

	var transcript string
	log := func(l string) {
		transcript += l + "\n"
	}
	action1 := New(func() error {
		log("execute action1")
		return nil
	}, func() {
		log("rollback action1")
	})

	action2 := New(func() error {
		log("execute action2")
		return nil
	}, func() {
		log("rollback action2")
	})

	action3 := New(func() error {
		log("oops, action3 failed")
		return errors.New(errText)
	}, func() {
		log("rollback action3")
	})

	action4 := New(func() error {
		log("execute action4")
		return nil
	}, func() {
		log("rollback action4")
	})

	err := ExecuteChain(action1, action2, action3, action4)
	if err.Error() != errText {
		t.Fatalf("wrong error, expected '%s' got '%s'", errText, err.Error())
	}

	if transcript != expectTranscript {
		t.Fatalf("unexpected transcript. expected:\n%s\n%s\n%s\nactual:\n%s\n%s",
			"-------------", expectTranscript, "-------------",
			"-------------", transcript)
	}
}

func TestExecuteChainSuccess(t *testing.T) {
	expectTranscript := `execute action1
execute action2
execute action3
execute action4
`

	var transcript string
	log := func(l string) {
		transcript += l + "\n"
	}
	action1 := New(func() error {
		log("execute action1")
		return nil
	}, func() {
		log("rollback action1")
	})

	action2 := New(func() error {
		log("execute action2")
		return nil
	}, func() {
		log("rollback action2")
	})

	action3 := New(func() error {
		log("execute action3")
		return nil
	}, func() {
		log("rollback action3")
	})

	action4 := New(func() error {
		log("execute action4")
		return nil
	}, func() {
		log("rollback action4")
	})

	err := ExecuteChain(action1, action2, action3, action4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transcript != expectTranscript {
		t.Fatalf("unexpected transcript. expected:\n%s\n%s\n%s\nactual:\n%s\n%s",
			"-------------", expectTranscript, "-------------",
			"-------------", transcript)
	}
}
