package patch

import (
	"fmt"

	"github.com/brnsampson/optional"
)

const (
	Noop PatchAction = iota
	Remove
	Replace
)

type PatchAction int

type PatchFunc[T comparable] func(left, right T) error

type Patchable[Self comparable, P Patch[Self]] interface {
	WithPatch(P) (Self, error)
	Patch(P) error
}

type Patch[T comparable] interface {
	IsNoop(T) bool
	IsNoopOption(optional.Option[T]) bool
	Apply(*T) error
	ApplyOption(optional.Option[T]) error
}

type PrimativePatch[T comparable] interface {
	Patch[T]
	Action() PatchAction
	Peek() optional.Option[T]
}

type PatchField[T comparable] struct {
	action PatchAction
	value  optional.Option[T]
}

func GetFieldPatch[T comparable](old, updated T) PatchField[T] {
	if old == updated {
		return NewNoop[T]()
	} else {
		return NewReplace(updated)
	}
}

func GetOptionFieldPatch[T comparable](old, updated optional.Optional[T]) PatchField[T] {
	//oldNone := old.IsNone()
	//updatedNone := updated.IsNone()
	//if oldNone && updatedNone {
	//	return NewNoop[T]()
	//} else if oldNone && !updatedNone {
	//	newValue, err := updated.Get()
	//	if err != nil {
	//    // We should never get here...
	//		return NewNoop[T]()
	//	}
	//	return NewReplace(newValue)
	//} else if !oldNone && updatedNone {
	//	return NewRemoval[T]()
	//} else if !oldNone && !updatedNone {
	//	newValue, err := updated.Get()
	//	if err != nil {
	//		// Again, we should never get here.
	//		return NewNoop[T]()
	//	}
	//	return NewReplace(newValue)
	//}
	// This should be unreachable

	newValue, err := updated.Get()
	if err == nil {
		return NewNoop[T]()
	} else {
		return NewReplace(newValue)
	}
	// panic("Panic: reached unreachable code when calculating patch for optional field value")
}

func NewNoop[T comparable]() PatchField[T] {
	return PatchField[T]{action: Noop, value: optional.None[T]()}
}

func NewRemoval[T comparable]() PatchField[T] {
	return PatchField[T]{action: Remove, value: optional.None[T]()}
}

func NewReplace[T comparable](updated T) PatchField[T] {
	return PatchField[T]{action: Replace, value: optional.NewOption(updated)}
}

func (p PatchField[T]) Peek() optional.Optional[T] {
	// This SHOULD create a new value because of the non-pointer receiver and the pointer returned should point to that,
	// but we should really have a test case to make sure the user of this cannot modify the internal optional.
	return &p.value
}

func (p PatchField[T]) Action() PatchAction {
	return p.action
}

func (p PatchField[T]) IsNoop(current T) bool {
	if p.action == Remove {
		return true
	} else if p.action == Replace {
		return p.value.Matches(current)
	} else if p.action == Noop {
		return true
	} else {
		return false
	}
}

func (p PatchField[T]) IsNoopOption(current optional.Optional[T]) bool {
	if p.action == Remove {
		if current.IsSome() {
			return false
		} else {
			return true
		}
	} else if p.action == Replace {
		return p.value.Eq(current)
	} else if p.action == Noop {
		return true
	} else {
		// Unknown/unsuppoerted operation? We must have introduced new functionality but didn't update this function...
		return false
	}
}

func (p PatchField[T]) Apply(operand *T) error {
	if p.action == Remove {
		// We can't do something like current = nil because this is not an optional. It would only change the pointer
		// passed to this function.
		return fmt.Errorf("cannot remove value of non-option type")
	} else if p.action == Replace {
		if p.value.Match(*operand) {
			return nil
		} else if p.value.IsSome() {
			// Because we have a non-pointer receiver, it should be a copy in memory and this should not modify the orignal
			// patch object. We should really make a test to make sure though.
			*operand, _ = p.value.Get()
			return nil
		} else {
			return fmt.Errorf("cannot replace non-option type with nil")
		}
	}
	// p.action == Noop
	return nil
}

func (p PatchField[T]) ApplyOption(current optional.Optional[T]) error {
	if current.IsNone() {
		if p.action == Remove {
			return nil
		} else if p.action == Replace {
			if p.value.IsNone() {
				// Set to none, but it is already none??
				return nil
			} else {
				// Again, p should be a copy so this shouldn't affect the original, but we should make a test to be sure.
				current.Set(p.value.UnsafeUnwrap())
				return nil
			}
		}
	}

	// current is Some
	if p.action == Remove {
		current.Clear()
		return nil
	} else if p.action == Replace {
		tmp, err := p.value.Unwrap()
		if err != nil {
			// p.value is none, so we are replacing with none? Same as Remove case.
			current.Clear()
			return nil
		} else {
			current.Set(tmp)
			return nil
		}
	}

	// p.action == Noop
	return nil
}
