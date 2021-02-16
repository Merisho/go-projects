package flags

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type Flags struct {
	StrField  string `flag:"str"`
	TrueField bool   `flag:"true"`
	FalseField bool `flag:"false"`
	IntField  int    `flag:"int"`
}

func TestNonPointerError(t *testing.T) {
	err := Parse("", Flags{})
	assert.EqualError(t, err, "flags argument must be a pointer")

	err = Parse("", nil)
	assert.EqualError(t, err, "flags argument must be a pointer")

	err = Parse("", &Flags{})
	assert.NoError(t, err)
}

func TestInvalidFlagType(t *testing.T) {
	err := Parse("-str test -int asdf", &Flags{})
	assert.EqualError(t, err, "invalid int flag type")
}

func TestSetFieldValueFromFlag(t *testing.T) {
	flags := &Flags{}
	err := Parse("-int 123456789 -str this is a string flag -true", flags)

	assert.NoError(t, err)
	assert.Equal(t, 123456789, flags.IntField)
	assert.Equal(t, "this is a string flag", flags.StrField)
	assert.True(t, flags.TrueField)
	assert.False(t, flags.FalseField)
}

func TestParseFlagsString(t *testing.T) {
	flags := parseFlagsStr("test -flagA flag_a_value -flagB 1 2 3 -flagC -flagD flag d value")

	_, ok := flags["test"]
	assert.False(t, ok)
	assert.Equal(t, "flag_a_value", flags["flagA"])
	assert.Equal(t, "1 2 3", flags["flagB"])
	assert.Equal(t, "", flags["flagC"])
	assert.Equal(t, "flag d value", flags["flagD"])
}

func TestMapFlagFieldTagsToFieldNames(t *testing.T) {
	fieldTags := mapFieldTagsToNames(reflect.TypeOf(Flags{}))

	assert.Empty(t, fieldTags[""])
	assert.Equal(t, "StrField", fieldTags["str"])
	assert.Equal(t, "TrueField", fieldTags["true"])
	assert.Equal(t, "FalseField", fieldTags["false"])
	assert.Equal(t, "IntField", fieldTags["int"])
}
