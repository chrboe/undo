// undo provides functionality to represent an action that can be rolled
// back at a later point in time.
package undo

// Undo represents an undoable action. It can be executed, undone, and committed.
type Undo struct {
	exec      func() error
	undo      func()
	committed bool
}

// New constructs a new Undo
func New(exec func() error, undo func()) *Undo {
	return &Undo{exec, undo, false}
}

// Execute executes the exec function of the action. The committed flag is reset to
// `false`. This effectively means that an action can be reused after it has
// been committed. If the action fails, the error from the exec function is
// returned.
func (u *Undo) Execute() error {
	u.committed = false
	return u.exec()
}

// Undo executes the undo function of the action. It is intended to be used in
// defer statements. If the action has already been committed, this is a no-op.
func (u *Undo) Undo() {
	if !u.committed {
		u.committed = true
		u.undo()
	}
}

// Commit sets the committed flag on the action. This essentially "defuses" the
// mechanism and all subsequent calls to Undo will be no-ops.
func (u *Undo) Commit() {
	u.committed = true
}

// CommitAll is a helper which can be used to commit several actions at once.
func CommitAll(actions ...*Undo) {
	for _, a := range actions {
		a.Commit()
	}
}

// ExecuteChain represents a common pattern where a chain of actions is executed
// one after the other. Every action gets executed, checked for an error (in
// which case the function returns), and their Undo gets deferred. If all actions
// complete without error, they are all committed. If one action fails, all actions
// up to that point will be undone and the error from the Execute function will be
// returned.
func ExecuteChain(actions ...*Undo) error {
	for _, a := range actions {
		err := a.Execute()
		if err != nil {
			return err
		}
		defer a.Undo()
	}

	CommitAll(actions...)
	return nil
}
