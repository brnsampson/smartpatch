package patch_test

import (
	"testing"

	"github.com/brnsampson/optional"
	"github.com/brnsampson/smartpatch"
	"gotest.tools/v3/assert"
)

func TestPatchNoop(t *testing.T) {
	original := "original"
	changed := "changed"
	originalOption := optional.NewOption(original)
	changedOption := optional.NewOption(changed)
	none := optional.None[string]()

	p1 := smartpatch.NewNoop[string]()
	p2 := smartpatch.NewReplace(original)
	p3 := smartpatch.NewRemoval[string]()
	p4 := smartpatch.NewReplace("changed")

	assert.Assert(t, p1.IsNoop(original), "Noop patch did not evaluate itself as a noop!")
	assert.Assert(t, p2.IsNoop(original), "Replace patch with same value did not evaluate itself as a noop!")
	assert.Assert(t, !p2.IsNoop(changed), "Replace patch with different value incorrectly evaluates itself as a noop!")
	assert.Assert(t, p3.IsNoop(original), "Removal patch on a non-optional value did not evaluate itself as a noop!")
	assert.Assert(t, !p4.IsNoop(original), "Replace patch with different value incorrectly evaluated itself as a noop!")

	assert.Assert(t, p1.IsNoopOption(&originalOption), "Noop patch did not evaluate itself as a noop on an optional with Some value!")
	assert.Assert(t, p2.IsNoopOption(&originalOption), "Replace patch with same value did not evaluate itself as a noop on an optional value!")
	assert.Assert(t, !p2.IsNoopOption(&changedOption), "Replace patch with different value incorrectly evaluates itself as a noop on an optional value!")
	assert.Assert(t, !p3.IsNoopOption(&originalOption), "Removal patch on option with Some value incorrectly evaluated itself as a noop!")
	assert.Assert(t, !p4.IsNoopOption(&originalOption), "Replace patch with different value incorrectly evaluated itself as a noop!")

	assert.Assert(t, p1.IsNoopOption(&none), "Noop patch did not evaluate itself as a noop on an optional with None value!")
	assert.Assert(t, !p2.IsNoopOption(&none), "Replace patch with value incorrectly evaluated itself as a noop on an optional with None value!")
	assert.Assert(t, p3.IsNoopOption(&none), "Removal patch on option with None value did not evaluate itself as a noop!")
}

func TestPatchApply(t *testing.T) {
	p1 := smartpatch.NewNoop[string]()
	p2 := smartpatch.NewRemoval[string]()
	p3 := smartpatch.NewReplace("changed")

	original := "original"
	tester := original

	assert.NilError(t, p1.Apply(&tester))
	assert.Equal(t, tester, original)
	assert.ErrorContains(t, p2.Apply(&tester), "cannot remove value of non-option type")
	assert.Equal(t, tester, original)
	assert.NilError(t, p3.Apply(&tester))
	assert.Equal(t, tester, "changed")
}
