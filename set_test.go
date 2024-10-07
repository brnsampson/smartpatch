package patch_test

import (
	"testing"

	"github.com/brnsampson/optional"
	"gotest.tools/v3/assert"
)

func TestPatchSet(t *testing.T) {
	thing1 := "dog"
	thing2 := optional.NewOption(thing1)
	thing3 := optional.None[int]()

	p1 := smartpatch.NewReplace(thing1)
	// var tmp smartpatch.Patch[string] = &p1
	p2 := smartpatch.NewReplace(thing1)
	p3 := smartpatch.NewRemoval[int]()

	pp1 := smartpatch.NewValuePatcher[string](&thing1, &p1)
	pp2 := smartpatch.NewOptionalPatcher[string](&thing2, &p2)
	pp3 := smartpatch.NewOptionalPatcher[int](&thing3, &p3)

	pm := smartpatch.NewPatchMap()

	err := pm.Append(thing1, pp1)
	assert.NilError(t, err)

	err = pm.Append("dogopt", pp2)
	assert.NilError(t, err)

	err = pm.Append("intNone", pp3)
	assert.NilError(t, err)

	assert.Assert(t, pm.IsAllNoop())
}
