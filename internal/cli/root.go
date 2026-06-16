package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"qed/internal/codec"
	"qed/internal/version"
)

func NewRootCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	var encrypt bool
	var decrypt bool

	cmd := &cobra.Command{
		Use:           "qed (-e|-d) password",
		Short:         "Encrypt or decrypt stdin using qed format",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if encrypt == decrypt {
				return fmt.Errorf("set exactly one of -e or -d")
			}

			input, err := io.ReadAll(stdin)
			if err != nil {
				return fmt.Errorf("read input: %w", err)
			}

			password := args[0]
			var output []byte
			if encrypt {
				output, err = codec.Encrypt(input, password)
			} else {
				output, err = codec.Decrypt(input, password)
			}
			if err != nil {
				return err
			}

			if _, err := stdout.Write(output); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	}

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.Flags().BoolVarP(&encrypt, "encrypt", "e", false, "encrypt stdin")
	cmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "decrypt stdin")
	cmd.AddCommand(newVersionCommand(stdout))

	return cmd
}

func newVersionCommand(stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:           "version",
		Short:         "Print qed version information",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(
				stdout,
				"qed %s\ncommit: %s\nbuilt: %s\n",
				version.Version,
				version.Commit,
				version.Date,
			)
			return err
		},
	}
}
