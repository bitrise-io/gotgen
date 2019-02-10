package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"encoding/json"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/templateutil"
	"github.com/bitrise-tools/gotgen/configs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	RunE: generate,
}

func init() {
	RootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generate(cmd *cobra.Command, args []string) error {
	// Read Inventory
	log.Println(colorstring.Blue("Reading GotGen config ..."))
	ggConfContent, err := fileutil.ReadBytesFromFile(gotgenConfigFileName)
	if err != nil {
		return errors.Wrapf(err, "Failed to read GotGen config (%s) file", gotgenConfigFileName)
	}
	ggConf := configs.Model{}
	if err := json.Unmarshal(ggConfContent, &ggConf); err != nil {
		return errors.Wrap(err, "Failed to parse GotGen config (JSON)")
	}
	log.Println(colorstring.Green("[DONE] Reading GotGen config"))

	//

	fmt.Println()
	log.Println(colorstring.Blue("Searching for templates and generaring files ..."))
	fmt.Println()
	templateFiles, err := filepath.Glob("*.gg")
	if err != nil {
		return errors.Wrap(err, "Failed to scan .gg template files")
	}
	for _, aTemplatePth := range templateFiles {
		generatedFilePath := strings.TrimSuffix(aTemplatePth, ".gg")
		fmt.Println(" * ", aTemplatePth, " => ", generatedFilePath)

		templateCont, err := fileutil.ReadStringFromFile(aTemplatePth)
		if err != nil {
			return errors.Wrapf(err, "Failed to read template content (path: %s)", aTemplatePth)
		}

		generatedContent, err := generateContent(templateCont, ggConf.Inventory, ggConf.Delimiter.Left, ggConf.Delimiter.Right)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate file based on content (%s) - invalid content?", aTemplatePth)
		}

		if err := fileutil.WriteStringToFile(generatedFilePath, generatedContent); err != nil {
			return errors.Wrapf(err, "Failed to write generated content into file (to path: %s)", generatedFilePath)
		}
		fmt.Println("   ", colorstring.Green("[OK]"))
	}
	fmt.Println()
	log.Println(colorstring.Green("[DONE] Searching for templates and generaring files"))
	fmt.Println()

	return nil
}

func generateContent(templateCont string, inventory map[string]interface{}, delimiterLeft, delimiterRight string) (string, error) {
	generatedContent, err := templateutil.EvaluateTemplateStringToStringWithDelimiterAndOpts(
		templateCont,
		inventory, createAvailableTemplateFunctions(inventory),
		delimiterLeft, delimiterRight,
		[]string{"missingkey=error"})
	if err != nil {
		return "", errors.WithStack(err)
	}
	return generatedContent, nil
}

func createAvailableTemplateFunctions(inventory map[string]interface{}) template.FuncMap {
	return template.FuncMap{
		"var": func(key string) (interface{}, error) {
			val, isFound := inventory[key]
			if !isFound {
				return "", errors.Errorf("No value found for key: %s", key)
			}
			return val, nil
		},
		"getenv": func(key string) string {
			return os.Getenv(key)
		},
		"getenvRequired": func(key string) (string, error) {
			if val := os.Getenv(key); len(val) > 0 {
				return val, nil
			}
			return "", errors.Errorf("No environment variable value found for key: %s", key)
		},
	}
}
