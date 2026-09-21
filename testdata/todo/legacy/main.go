package main

import (
	"errors"
	"fmt"
	"log"
	"mdcoach"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const version = "2.0.0"

func runInferredOutput(cmd *cobra.Command, args []string) error {
	output, _ := cmd.PersistentFlags().GetString(`output`)
	e := mdcoach.NewEnvironment(output)
	defer e.Close()
	if v, _ := cmd.PersistentFlags().GetBool(`force`); v {
		e.Overwrite = true
	}

	ensureExtensionAndCleanPath := func(ext string) string {
		switch old := filepath.Ext(output); strings.ToLower(old) {
		case ext:
			return filepath.Clean(output)
		case "", ".": // only output directory was specified, derrive name from source
			base := filepath.Base(args[0])
			return filepath.Clean(filepath.Join(output, strings.TrimSuffix(base, filepath.Ext(base))+ext))
		default:
			return filepath.Clean(strings.TrimSuffix(output, old) + ext)
		}
	}
	switch f, _ := cmd.PersistentFlags().GetString(`type`); strings.ToLower(f) {
	case `handout`, `syllabus`:
		e.Output = ensureExtensionAndCleanPath(`.pdf`)
		return mdcoach.Paper(e, `handout.css`, args...)
	case `notes`:
		e.Output = ensureExtensionAndCleanPath(`.pdf`)
		return mdcoach.Paper(e, `notes.css`, args...)
	case `book`:
		e.Output = ensureExtensionAndCleanPath(`.pdf`)
		return mdcoach.Paper(e, `book.css`, args...)
	case `split`:
		e.Output = ensureExtensionAndCleanPath(`.pdf`)
		return mdcoach.Paper(e, `split.css`, args...)
	case `assessment`:
		e.Output = ensureExtensionAndCleanPath(`.pdf`)
		return mdcoach.Assessment(e, args...)
	default:
		switch strings.ToLower(filepath.Ext(args[0])) {
		case `.pdf`:
			if len(args) == 1 {
				_, err := readPdf(args[0])
				return err
			}
			if output == `` {
				output = `joined-` + args[0]
			}
			args = append(args, output)
			return mdcoach.Exec(`pdfunite`, args...)
		case `.yaml`:
			// TODO: notice how I don't even use e.Output here
			// TODO: setting up .cache is not required for this to work!
			if len(args) < 2 {
				return errors.New(`Please specify at least one template file.`)
			}
			temp, handle := e.Create(args[0], true)
			err := mdcoach.YamlWithTemplate(handle, args[0], args[1:]...)
			handle.Close()
			// if err != nil {
			// 	return err
			// }
			mdcoach.PDFManual(temp, ensureExtensionAndCleanPath(`.pdf`))
			return err
		}
	}
	e.Output = ensureExtensionAndCleanPath(`.html`)
	return mdcoach.Presentation(e, args...)
}

func main() {
	var CLI = &cobra.Command{
		Use:     `mdcoach`,
		Version: version,
		Short:   `Presentation builder.`,
		Long:    `Presentation builder. Generate lecture presentation slides, notes, and assessments from markdown files.`,
		Run: func(cmd *cobra.Command, args []string) {
			output, _ := cmd.PersistentFlags().GetString(`output`)
			if v, _ := cmd.PersistentFlags().GetBool(`clear`); v {
				target := filepath.Join(filepath.Dir(output), `.cache`)
				os.RemoveAll(target)
				defer log.Printf(`Removed cache at %s.`, target)
			}
			if v, _ := cmd.PersistentFlags().GetBool(`demo`); v {
				e := mdcoach.NewEnvironment(`.`)
				e.Output = `demo.html`
				defer e.Close()
				mdcoach.Presentation(e, `.cache/test-demo.md`)
				fmt.Println(`Demo compiled to demo.html`)
				// TODO: allow specifying output
			} else if len(args) == 0 {
				cmd.Help()
			} else {
				if err := runInferredOutput(cmd, args); err == nil {
					// TODO: output is probably modified!
					fmt.Printf("%s\n", output) // if -o not set, prints nothing
				} else {
					fmt.Printf("Error: %s: %s.\n", output, err.Error())
				}
			}
		},
	}
	CLI.PersistentFlags().StringP(`output`, `o`, ``, `Output to this file.`)
	CLI.PersistentFlags().BoolP(`verbose`, `v`, false, `Print out additional information.`)
	CLI.PersistentFlags().StringP(`type`, `t`, `presentation`, `Set output type: presentation, digest, notes, handout, assessment, book, index.`)
	CLI.PersistentFlags().BoolP(`force`, `f`, false, `Overwrite all files without a warning.`)
	CLI.PersistentFlags().BoolP(`clear`, `c`, false, `First, empty the cache at the output directory.`)
	// CLI.PersistentFlags().Bool(`popup`, false, `Open generated file with default browser.`)
	CLI.PersistentFlags().Bool(`demo`, false, `Generate a sample presentation with usage examples.`)
	// CLI.PersistentFlags().StringArrayP(`format`)
	// CLI.PersistentFlags().StringSliceP("include", "i", []string{}, "Add provided paths to SASS includes.")
	CLI.Execute()
}
