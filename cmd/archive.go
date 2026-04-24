package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envport/internal/archive"
	"github.com/user/envport/internal/profile"
)

func init() {
	packCmd := &cobra.Command{
		Use:   "archive pack <file> <snapshot> [snapshot...]",
		Short: "Pack snapshots into a portable archive file",
	}
	unpackCmd := &cobra.Command{
		Use:   "archive unpack <file>",
		Short: "Restore snapshots from an archive file",
	}

	archiveCmd := &cobra.Command{
		Use:   "archive",
		Short: "Pack and unpack snapshot archives",
	}

	packCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: archive pack <file> <snapshot> [snapshot...]")
		}
		return runArchivePack(args[0], args[1:])
	}

	unpackCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("usage: archive unpack <file>")
		}
		return runArchiveUnpack(args[0])
	}

	archiveCmd.AddCommand(packCmd, unpackCmd)
	rootCmd.AddCommand(archiveCmd)
}

func runArchivePack(dest string, names []string) error {
	m, err := profile.NewManager()
	if err != nil {
		return err
	}
	ar, err := archive.Pack(m, names)
	if err != nil {
		return err
	}
	data, err := archive.Marshal(ar)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("archive: writing file: %w", err)
	}
	fmt.Printf("packed %s → %s\n", strings.Join(names, ", "), dest)
	return nil
}

func runArchiveUnpack(src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("archive: reading file: %w", err)
	}
	ar, err := archive.Unmarshal(data)
	if err != nil {
		return err
	}
	m, err := profile.NewManager()
	if err != nil {
		return err
	}
	restored, err := archive.Unpack(m, ar)
	if err != nil {
		return err
	}
	fmt.Printf("unpacked %d snapshot(s): %s\n", len(restored), strings.Join(restored, ", "))
	return nil
}
