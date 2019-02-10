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
		require.EqualError(t, err, `template: :1:8: executing "" at <getenvRequired "Test...>: error calling getenvRequired: No environment variable value found for key: Test_generateContent_KEY_2`)
		require.Equal(t, ``, genCont)
	}
}
