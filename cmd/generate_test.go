package cmd

import (
	"testing"

	"github.com/bitrise-io/go-utils/envutil"
	"github.com/stretchr/testify/require"
)

func Test_generateContent(t *testing.T) {
	t.Log("Simple string template, empty inventory - no substitution")
	{
		genCont, err := generateContent(`Test Content`, nil, "{{", "}}")
		require.NoError(t, err)
		require.Equal(t, `Test Content`, genCont)
	}

	t.Log("Missing inventory key")
	{
		genCont, err := generateContent(`Test {{ .KeyOne }} Content`, nil, "{{", "}}")
		require.EqualError(t, err, `template: :1:8: executing "" at <.KeyOne>: map has no entry for key "KeyOne"`)
		require.Equal(t, ``, genCont)
	}

	t.Log("Simple substitution")
	{
		genCont, err := generateContent(
			`Test {{ .KeyOne }} Content`,
			map[string]interface{}{"KeyOne": "Value 1"},
			"{{", "}}",
		)
		require.NoError(t, err)
		require.Equal(t, `Test Value 1 Content`, genCont)
	}

	t.Log("Template function: var: key found")
	{
		genCont, err := generateContent(
			`Test {{ var "KeyOne" }} Content`,
			map[string]interface{}{"KeyOne": "Value 1"},
			"{{", "}}",
		)
		require.NoError(t, err)
		require.Equal(t, `Test Value 1 Content`, genCont)
	}
	t.Log("Template function: var: key NOT found")
	{
		genCont, err := generateContent(
			`Test {{ var "NonExistingKey" }} Content`,
			map[string]interface{}{"KeyOne": "Value 1"},
			"{{", "}}",
		)
		require.EqualError(t, err, `template: :1:8: executing "" at <var "NonExistingKey">: error calling var: No value found for key: NonExistingKey`)
		require.Equal(t, ``, genCont)
	}

	t.Log("Template function: getenv")
	{
		revokeFn, err := envutil.RevokableSetenv("Test_generateContent_KEY", "Test Env Value")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, revokeFn())
		}()

		genCont, err := generateContent(
			`Test {{ getenv "Test_generateContent_KEY" }} Content`,
			map[string]interface{}{"KeyOne": "Value 1"},
			"{{", "}}",
		)
		require.NoError(t, err)
		require.Equal(t, `Test Test Env Value Content`, genCont)
	}

	t.Log("Template function: getenvRequired: found")
	{
		revokeFn, err := envutil.RevokableSetenv("Test_generateContent_KEY", "Test Env Value")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, revokeFn())
		}()

		genCont, err := generateContent(
			`Test {{ getenvRequired "Test_generateContent_KEY" }} Content`,
			map[string]interface{}{"KeyOne": "Value 1"},
			"{{", "}}",
		)
		require.NoError(t, err)
		require.Equal(t, `Test Test Env Value Content`, genCont)
	}

	t.Log("Template function: getenvRequired: NOT found")
	{
		genCont, err := generateContent(
			`Test {{ getenvRequired "Test_generateContent_KEY_2" }} Content`,
			map[string]interface{}{"KeyOne": "Value 1"},
			"{{", "}}",
		)
		require.EqualError(t, err, "template: :1:8: executing \"\" at <getenvRequired \"Test_generateContent_KEY_2\">: error calling getenvRequired: No environment variable value found for key: Test_generateContent_KEY_2")
		require.Equal(t, ``, genCont)
	}
}

func Test_indentWithSpaces(t *testing.T) {
	require.Equal(t, "", indentWithSpaces(2, ""), "Empty string should result in an empty string")

	require.Equal(t, "  a", indentWithSpaces(2, "a"))

	t.Log("Multiline test")
	{
		orig := `a
b
 c`
		expected := `  a
  b
   c`
		require.Equal(t, expected, indentWithSpaces(2, orig))
	}

	t.Log("Multiline test - ending with newline at the end")
	{
		orig := `a
b
 c
`
		expected := `  a
  b
   c
  `
		require.Equal(t, expected, indentWithSpaces(2, orig))
	}
}

func Test_yaml(t *testing.T) {
	t.Log("Simple no error")
	{
		obj := map[string]string{"key1": "value one"}
		expected := "key1: value one\n"

		s, err := yamlFn(obj)
		require.NoError(t, err)
		require.Equal(t, expected, s)
	}

	t.Log("Simple error")
	{
		// I don't know any way to make `yaml.Marshal` to return an error
	}

	t.Log("More complex - multiline")
	{
		obj := map[string]interface{}{"key1": "value one", "key2": 2}
		expected := "key1: value one\nkey2: 2\n"

		s, err := yamlFn(obj)
		require.NoError(t, err)
		require.Equal(t, expected, s)
	}
}
