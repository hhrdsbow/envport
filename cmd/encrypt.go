package cmd

import (
	"fmt"
	"os"

	"envport/internal/encrypt"
	"envport/internal/profile"

	"github.com/spf13/cobra"
)

func init() {
	var passphrase string
	var decrypt bool

	cmd := &cobra.Command{
		Use:   "encrypt <profile>",
		Short: "Encrypt or decrypt all values in a snapshot profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEncrypt(args[0], passphrase, decrypt)
		},
	}
	cmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase for encryption/decryption (required)")
	cmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "Decrypt instead of encrypt")
	_ = cmd.MarkFlagRequired("passphrase")
	rootCmd.AddCommand(cmd)
}

func runEncrypt(name, passphrase string, decrypt bool) error {
	dir, err := defaultProfileDir()
	if err != nil {
		return err
	}
	mgr, err := profile.New(dir)
	if err != nil {
		return err
	}
	snap, err := mgr.Load(name)
	if err != nil {
		return fmt.Errorf("load profile %q: %w", name, err)
	}
	var transformed map[string]string
	if decrypt {
		transformed, err = encrypt.DecryptMap(snap.Vars, passphrase)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	} else {
		transformed, err = encrypt.EncryptMap(snap.Vars, passphrase)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
	}
	snap.Vars = transformed
	if err := mgr.Save(name, snap); err != nil {
		return fmt.Errorf("save profile %q: %w", name, err)
	}
	action := "encrypted"
	if decrypt {
		action = "decrypted"
	}
	fmt.Fprintf(os.Stdout, "Profile %q %s successfully.\n", name, action)
	return nil
}
