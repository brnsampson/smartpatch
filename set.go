package patch

import (
	"fmt"

	"github.com/brnsampson/optional"
)

type PatchSet interface {
	IsNoop(any) (bool, error)
	Apply(any) error
	IsAllNoop() bool
	ApplyAll() error
}

type Patcher interface {
	IsNoop() bool
	Apply() error
}

type ValuePatcher[T comparable] struct {
	applied bool
	target  *T
	patch   Patch[T]
}

func NewValuePatcher[T comparable](target *T, patch Patch[T]) *ValuePatcher[T] {
	return &ValuePatcher[T]{applied: false, target: target, patch: patch}
}

func (p ValuePatcher[T]) IsNoop() bool {
	if p.applied {
		// We are assuming that patches are idempotent and nothing else is modifying p.target
		return true
	} else {
		return p.patch.IsNoop(*p.target)
	}
}

func (p *ValuePatcher[T]) Apply() error {
	if p.applied {
		return nil
	} else {
		return p.patch.Apply(p.target)
	}
}

type OptionalPatcher[T comparable] struct {
	applied bool
	target  optional.Optional[T]
	patch   Patch[T]
}

func NewOptionalPatcher[T comparable](target optional.Optional[T], patch Patch[T]) *OptionalPatcher[T] {
	return &OptionalPatcher[T]{applied: false, target: target, patch: patch}
}

func (p OptionalPatcher[T]) IsNoop() bool {
	if p.applied {
		// We are assuming that patches are idempotent and nothing else is modifying p.target
		return true
	} else {
		return p.patch.IsNoopOption(p.target)
	}
}

func (p *OptionalPatcher[T]) Apply() error {
	if p.applied {
		return nil
	} else {
		return p.patch.ApplyOption(p.target)
	}
}

type PatchMap struct {
	patchMap map[string]Patcher
}

func NewPatchMap() *PatchMap {
	m := make(map[string]Patcher)
	return &PatchMap{m}
}

func (m *PatchMap) Append(label string, patch Patcher) error {
	_, exists := m.patchMap[label]
	if exists {
		return fmt.Errorf("PatchMap Error: label %s was already used for another patch", label)
	}

	m.patchMap[label] = patch
	return nil
}

func (m *PatchMap) getPatchPair(label any) optional.Option[Patcher] {
	tmp, ok := label.(string)
	if ok != true {
		return optional.None[Patcher]()
	}

	p, ok := m.patchMap[tmp]
	if ok != true {
		return optional.None[Patcher]()
	}

	return optional.NewOption(p)
}

func (m *PatchMap) IsNoop(label any) (bool, error) {
	opt := m.getPatchPair(label)
	p, err := opt.Unwrap()
	if err != nil {
		return p.IsNoop(), nil
	}
	return false, fmt.Errorf("Patch not found in PatchMap: %v", label)
}

func (m *PatchMap) Apply(label any) error {
	opt := m.getPatchPair(label)
	p, err := opt.Unwrap()
	if err != nil {
		return p.Apply()
	}
	return fmt.Errorf("Patch not found in PatchMap: %v", label)
}

func (m *PatchMap) IsAllNoop() bool {
	for _, p := range m.patchMap {
		if !p.IsNoop() {
			return false
		}
	}
	return true
}

func (m *PatchMap) ApplyAll() error {
	for lookup, p := range m.patchMap {
		err := p.Apply()
		if err != nil {
			return fmt.Errorf("Patch failed for entry: %v", lookup)
		}
	}
	return nil
}
