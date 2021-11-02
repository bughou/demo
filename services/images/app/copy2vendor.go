package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func copy2vendorCmd() *cobra.Command {
	var all bool
	var excludeTest bool
	cmd := &cobra.Command{
		Use:   `copy2vendor [package-path] ...`,
		Short: `Copy the specified packages to vendor dir.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 && !all {
				return errors.New(`no package path provided.`)
			}
			return copy2Vendor(args, excludeTest)
		},
	}
	cmd.Flags().BoolVarP(&all, `all`, `a`, false, `copy all dependences not in vendor dir.`)
	cmd.Flags().BoolVarP(&excludeTest, `exclude-test`, `e`, false, `exclude test dependences.`)
	return cmd
}

func copy2Vendor(pkgs []string, excludeTest bool) error {
	if len(pkgs) == 0 {
		var err error
		if pkgs, err = getDeps(false, excludeTest); err != nil {
			return err
		}
	}
	goSrcDir := fs.GoSrcPath()
	vendorDir := filepath.Join(release.ProjectRoot(), `vendor`)
	for _, pkg := range pkgs {
		if err := syncGoFiles(filepath.Join(goSrcDir, pkg), filepath.Join(vendorDir, pkg)); err != nil {
			return err
		}
	}
	return nil
}

func syncGoFiles(srcDir, destDir string) error {
	srcFiles, err := filepath.Glob(srcDir + "/*.go")
	if err != nil {
		return err
	}
	if fs.Exist(destDir) {
		if err := removeRedundantDestFiles(srcFiles, destDir); err != nil {
			return err
		}
	} else if err = os.MkdirAll(destDir, 0775); err != nil {
		return err
	}
	for _, src := range srcFiles {
		dst := strings.Replace(src, srcDir, destDir, 1)
		if err := fs.Copy(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func removeRedundantDestFiles(srcFiles []string, destDir string) error {
	goFiles := map[string]bool{}
	for _, srcFile := range srcFiles {
		goFiles[filepath.Base(srcFile)] = true
	}

	destFiles, err := filepath.Glob(destDir + "/*.go")
	if err != nil {
		return err
	}

	for _, destFile := range destFiles {
		if !goFiles[filepath.Base(destFile)] {
			if err := os.Remove(destFile); err != nil {
				return err
			}
		}
	}
	return nil
}
