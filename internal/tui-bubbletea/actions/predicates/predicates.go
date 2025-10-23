package predicates

import "github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"

// Always returns a predicate that always returns true
func Always() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return true
	}
}

// Never returns a predicate that always returns false
func Never() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return false
	}
}

// HasSelection returns true if at least one ticket is selected
func HasSelection() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return len(ctx.SelectedTickets) > 0
	}
}

// HasSingleSelection returns true if exactly one ticket is selected
func HasSingleSelection() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return len(ctx.SelectedTickets) == 1
	}
}

// HasMultipleSelection returns true if more than one ticket is selected
func HasMultipleSelection() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return len(ctx.SelectedTickets) > 1
	}
}

// IsInWorkspace returns true if a workspace is selected
func IsInWorkspace() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return ctx.SelectedWorkspace != nil
	}
}

// HasUnsavedChanges returns true if there are unsaved changes
func HasUnsavedChanges() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return ctx.HasUnsavedChanges
	}
}

// IsOnline returns true if not in offline mode
func IsOnline() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return !ctx.IsOffline
	}
}

// Not inverts a predicate
func Not(pred actions.PredicateFunc) actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return !pred(ctx)
	}
}

// And combines predicates with logical AND
func And(predicates ...actions.PredicateFunc) actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		for _, pred := range predicates {
			if !pred(ctx) {
				return false
			}
		}
		return true
	}
}

// Or combines predicates with logical OR
func Or(predicates ...actions.PredicateFunc) actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		for _, pred := range predicates {
			if pred(ctx) {
				return true
			}
		}
		return false
	}
}
