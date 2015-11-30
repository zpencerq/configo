package configo

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testCbDesignDoc = "api_1"
)

var testPassString string = "pass_String"
var testPassPtrString string = "pass_PtrString"

var testPassInt int = 9812
var testPassPtrInt int = 1289

var testEnvString string = "env_String"
var testEnvPtrString string = "env_PtrString"

var testEnvInt int = 6789
var testEnvPtrInt int = 9876

var testErrString string = "fail_if_you_see_this"
var testErrInt int = 0
var testErrBool bool = false

var testBoolTrue bool = true

type Types struct {
	ZeroString    string
	DefaultString string `default:"default_String"`
	PassString    string `default:"fail_if_you_see_this"`
	TomlString    string `default:"fail_if_you_see_this"`
	EnvString     string `env:"CONFIGO_TEST_ENVSTRING" default:"fail_if_you_see_this"`

	ZeroPtrString    *string
	DefaultPtrString *string `default:"default_PtrString"`
	PassPtrString    *string `default:"fail_if_you_see_this"`
	TomlPtrString    *string `default:"fail_if_you_see_this"`
	EnvPtrString     *string `env:"CONFIGO_TEST_ENVPTRSTRING" default:"fail_if_you_see_this"`

	ZeroInt    int
	DefaultInt int `default:"1234"`
	PassInt    int `default:"0"`
	TomlInt    int `default:"0"`
	EnvInt     int `env:"CONFIGO_TEST_ENVINT" default:"0"`

	ZeroPtrInt    *int
	DefaultPtrInt *int `default:"4321"`
	PassPtrInt    *int `default:"0"`
	TomlPtrInt    *int `default:"0"`
	EnvPtrInt     *int `env:"CONFIGO_TEST_ENVPTRINT" default:"0"`

	ZeroBool    bool
	DefaultBool bool `default:"true"`
	PassBool    bool `default:"false"`
	TomlBool    bool `default:"false"`
	EnvBool     bool `env:"CONFIGO_TEST_ENVBOOL" default:"false"`

	ZeroPtrBool    *bool
	DefaultPtrBool *bool `default:"true"`
	PassPtrBool    *bool `default:"false"`
	TomlPtrBool    *bool `default:"false"`
	EnvPtrBool     *bool `env:"CONFIGO_TEST_ENVPTRBOOL" default:"false"`
}

type SubTypes struct {
	Struct    Types
	StructPtr *Types

	ZeroString    string
	DefaultString string `default:"default_String"`
	PassString    string `default:"fail_if_you_see_this"`
	TomlString    string `default:"fail_if_you_see_this"`
	EnvString     string `env:"CONFIGO_TEST_ENVSTRING" default:"fail_if_you_see_this"`

	ZeroPtrString    *string
	DefaultPtrString *string `default:"default_PtrString"`
	PassPtrString    *string `default:"fail_if_you_see_this"`
	TomlPtrString    *string `default:"fail_if_you_see_this"`
	EnvPtrString     *string `env:"CONFIGO_TEST_ENVPTRSTRING" default:"fail_if_you_see_this"`

	ZeroInt    int
	DefaultInt int `default:"1234"`
	PassInt    int `default:"0"`
	TomlInt    int `default:"0"`
	EnvInt     int `env:"CONFIGO_TEST_ENVINT" default:"0"`

	ZeroPtrInt    *int
	DefaultPtrInt *int `default:"4321"`
	PassPtrInt    *int `default:"0"`
	TomlPtrInt    *int `default:"0"`
	EnvPtrInt     *int `env:"CONFIGO_TEST_ENVPTRINT" default:"0"`

	ZeroBool    bool
	DefaultBool bool `default:"true"`
	PassBool    bool `default:"false"`
	TomlBool    bool `default:"false"`
	EnvBool     bool `env:"CONFIGO_TEST_ENVBOOL" default:"false"`

	ZeroPtrBool    *bool
	DefaultPtrBool *bool `default:"true"`
	PassPtrBool    *bool `default:"false"`
	TomlPtrBool    *bool `default:"false"`
	EnvPtrBool     *bool `env:"CONFIGO_TEST_ENVPTRBOOL" default:"false"`
}

type TestArraySub struct {
	Strings []string
	Ints    []int
}

type TestArray struct {
	Strings []string
	Ints    []int
	Sub1    TestArraySub
}

type TestIgnoreFields struct{}

func testSetEnv(env map[string]string) error {
	for k, v := range env {
		err := os.Setenv("CONFIGO_TEST_"+k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func testUnsetEnv() error {
	vars := []string{
		"ENVSTRING",
		"ENVPTRSTRING",

		"ENVINT",
		"ENVPTRINT",

		"ENVBOOL",
		"ENVPTRBOOL",
	}

	for _, v := range vars {
		err := os.Unsetenv("CONFIGO_TEST_" + v)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestUnmarshalFile1(t *testing.T) {
	err := testUnsetEnv()
	if err != nil {
		t.Fatal(err)
	}

	env := map[string]string{
		"ENVSTRING":    testEnvString,
		"ENVPTRSTRING": testEnvPtrString,
		"ENVINT":       fmt.Sprintf("%d", testEnvInt),
		"ENVPTRINT":    fmt.Sprintf("%d", testEnvPtrInt),
		"ENVBOOL":      "true",
		"ENVPTRBOOL":   "true",
	}

	err = testSetEnv(env)
	if err != nil {
		t.Fatal(err)
	}

	sub1 := Types{
		PassString:    testPassString,
		PassPtrString: &testPassPtrString,
		TomlString:    testErrString,
		TomlPtrString: &testErrString,
		EnvString:     testErrString,
		EnvPtrString:  &testErrString,

		PassInt:    testPassInt,
		PassPtrInt: &testPassPtrInt,
		TomlInt:    testErrInt,
		TomlPtrInt: &testErrInt,
		EnvInt:     testErrInt,
		EnvPtrInt:  &testErrInt,

		PassBool:    testBoolTrue,
		PassPtrBool: &testBoolTrue,
		TomlBool:    testErrBool,
		TomlPtrBool: &testErrBool,
		EnvBool:     testErrBool,
		EnvPtrBool:  &testErrBool,
	}

	got := SubTypes{
		Struct: sub1,
		StructPtr: &Types{
			PassString:    testPassString,
			PassPtrString: &testPassPtrString,
			TomlString:    testErrString,
			TomlPtrString: &testErrString,
			EnvString:     testErrString,
			EnvPtrString:  &testErrString,

			PassInt:    testPassInt,
			PassPtrInt: &testPassPtrInt,
			TomlInt:    testErrInt,
			TomlPtrInt: &testErrInt,
			EnvInt:     testErrInt,
			EnvPtrInt:  &testErrInt,

			PassBool:    testBoolTrue,
			PassPtrBool: &testBoolTrue,
			TomlBool:    testErrBool,
			TomlPtrBool: &testErrBool,
			EnvBool:     testErrBool,
			EnvPtrBool:  &testErrBool,
		},

		PassString:    testPassString,
		PassPtrString: &testPassPtrString,
		TomlString:    testErrString,
		TomlPtrString: &testErrString,
		EnvString:     testErrString,
		EnvPtrString:  &testErrString,

		PassInt:    testPassInt,
		PassPtrInt: &testPassPtrInt,
		TomlInt:    testErrInt,
		TomlPtrInt: &testErrInt,
		EnvInt:     testErrInt,
		EnvPtrInt:  &testErrInt,

		PassBool:    testBoolTrue,
		PassPtrBool: &testBoolTrue,
		TomlBool:    testErrBool,
		TomlPtrBool: &testErrBool,
		EnvBool:     testErrBool,
		EnvPtrBool:  &testErrBool,
	}

	err = UnmarshalFile("testdata/types.toml", &got)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, got.ZeroString)
	assert.Equal(t, "default_String", got.DefaultString)
	assert.Equal(t, "pass_String", got.PassString)
	assert.Equal(t, "toml_String", got.TomlString)
	assert.Equal(t, "env_String", got.EnvString)
	assert.Empty(t, *got.ZeroPtrString)
	assert.Equal(t, "default_PtrString", *got.DefaultPtrString)
	assert.Equal(t, "pass_PtrString", *got.PassPtrString)
	assert.Equal(t, "toml_PtrString", *got.TomlPtrString)
	assert.Equal(t, "env_PtrString", *got.EnvPtrString)

	assert.Empty(t, got.ZeroInt)
	assert.Equal(t, 1234, got.DefaultInt)
	assert.Equal(t, testPassInt, got.PassInt)
	assert.Equal(t, 7878, got.TomlInt)
	assert.Equal(t, testEnvInt, got.EnvInt)
	assert.Empty(t, *got.ZeroPtrInt)
	assert.Equal(t, 4321, *got.DefaultPtrInt)
	assert.Equal(t, testPassPtrInt, *got.PassPtrInt)
	assert.Equal(t, 8787, *got.TomlPtrInt)
	assert.Equal(t, testEnvPtrInt, *got.EnvPtrInt)

	assert.Empty(t, got.ZeroBool)
	assert.Equal(t, true, got.DefaultBool)
	assert.Equal(t, true, got.PassBool)
	assert.Equal(t, true, got.TomlBool)
	assert.Equal(t, true, got.EnvBool)
	assert.Empty(t, *got.ZeroPtrBool)
	assert.Equal(t, true, *got.DefaultPtrBool)
	assert.Equal(t, true, *got.PassPtrBool)
	assert.Equal(t, true, *got.TomlPtrBool)
	assert.Equal(t, true, *got.EnvPtrBool)

	assert.Empty(t, got.Struct.ZeroString)
	assert.Equal(t, "default_String", got.Struct.DefaultString)
	assert.Equal(t, "pass_String", got.Struct.PassString)
	assert.Equal(t, "toml_String", got.Struct.TomlString)
	assert.Equal(t, "env_String", got.Struct.EnvString)
	assert.Empty(t, *got.Struct.ZeroPtrString)
	assert.Equal(t, "default_PtrString", *got.Struct.DefaultPtrString)
	assert.Equal(t, "pass_PtrString", *got.Struct.PassPtrString)
	assert.Equal(t, "toml_PtrString", *got.Struct.TomlPtrString)
	assert.Equal(t, "env_PtrString", *got.Struct.EnvPtrString)

	assert.Empty(t, got.Struct.ZeroInt)
	assert.Equal(t, 1234, got.Struct.DefaultInt)
	assert.Equal(t, testPassInt, got.Struct.PassInt)
	assert.Equal(t, 7878, got.Struct.TomlInt)
	assert.Equal(t, testEnvInt, got.Struct.EnvInt)
	assert.Empty(t, *got.Struct.ZeroPtrInt)
	assert.Equal(t, 4321, *got.Struct.DefaultPtrInt)
	assert.Equal(t, testPassPtrInt, *got.Struct.PassPtrInt)
	assert.Equal(t, 8787, *got.Struct.TomlPtrInt)
	assert.Equal(t, testEnvPtrInt, *got.Struct.EnvPtrInt)

	assert.Empty(t, got.Struct.ZeroBool)
	assert.Equal(t, true, got.Struct.DefaultBool)
	assert.Equal(t, true, got.Struct.PassBool)
	assert.Equal(t, true, got.Struct.TomlBool)
	assert.Equal(t, true, got.Struct.EnvBool)
	assert.Empty(t, *got.Struct.ZeroPtrBool)
	assert.Equal(t, true, *got.Struct.DefaultPtrBool)
	assert.Equal(t, true, *got.Struct.PassPtrBool)
	assert.Equal(t, true, *got.Struct.TomlPtrBool)
	assert.Equal(t, true, *got.Struct.EnvPtrBool)

	assert.Empty(t, got.StructPtr.ZeroString)
	assert.Equal(t, "default_String", got.StructPtr.DefaultString)
	assert.Equal(t, "pass_String", got.StructPtr.PassString)
	assert.Equal(t, "toml_String", got.StructPtr.TomlString)
	assert.Equal(t, "env_String", got.StructPtr.EnvString)
	assert.Empty(t, *got.StructPtr.ZeroPtrString)
	assert.Equal(t, "default_PtrString", *got.StructPtr.DefaultPtrString)
	assert.Equal(t, "pass_PtrString", *got.StructPtr.PassPtrString)
	assert.Equal(t, "toml_PtrString", *got.StructPtr.TomlPtrString)
	assert.Equal(t, "env_PtrString", *got.StructPtr.EnvPtrString)

	assert.Empty(t, got.StructPtr.ZeroInt)
	assert.Equal(t, 1234, got.StructPtr.DefaultInt)
	assert.Equal(t, testPassInt, got.StructPtr.PassInt)
	assert.Equal(t, 7878, got.StructPtr.TomlInt)
	assert.Equal(t, testEnvInt, got.StructPtr.EnvInt)
	assert.Empty(t, *got.StructPtr.ZeroPtrInt)
	assert.Equal(t, 4321, *got.StructPtr.DefaultPtrInt)
	assert.Equal(t, testPassPtrInt, *got.StructPtr.PassPtrInt)
	assert.Equal(t, 8787, *got.StructPtr.TomlPtrInt)
	assert.Equal(t, testEnvPtrInt, *got.StructPtr.EnvPtrInt)

	assert.Empty(t, got.StructPtr.ZeroBool)
	assert.Equal(t, true, got.StructPtr.DefaultBool)
	assert.Equal(t, true, got.StructPtr.PassBool)
	assert.Equal(t, true, got.StructPtr.TomlBool)
	assert.Equal(t, true, got.StructPtr.EnvBool)
	assert.Empty(t, *got.StructPtr.ZeroPtrBool)
	assert.Equal(t, true, *got.StructPtr.DefaultPtrBool)
	assert.Equal(t, true, *got.StructPtr.PassPtrBool)
	assert.Equal(t, true, *got.StructPtr.TomlPtrBool)
	assert.Equal(t, true, *got.StructPtr.EnvPtrBool)
}

func TestUnmarshalFileArrays(t *testing.T) {
	var got TestArray

	err := UnmarshalFile("testdata/arrays.toml", &got)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []string{"a", "b", "c"}, got.Strings)
	assert.Equal(t, []int{1, 2, 3}, got.Ints)
	assert.Equal(t, []string{"a", "b", "c"}, got.Sub1.Strings)
	assert.Equal(t, []int{1, 2, 3}, got.Sub1.Ints)
}

func TestUnmarshalFileIgnoreFields(t *testing.T) {
	var got TestIgnoreFields

	err := UnmarshalFile("testdata/ignore_fields.toml", &got)
	if err != nil {
		t.Fatal(err)
	}
}
