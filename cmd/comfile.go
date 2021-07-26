package cmd

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/kopy"
	"github.com/itrepablik/sakto"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var comfileCmd = &cobra.Command{
	Use:   "comfile",
	Short: "Compress a single file",
	Long:  `comfile helps to compress a single file`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//source and destination paths
		srcPath := filepath.FromSlash(args[0])
		dstPath := filepath.FromSlash(args[1])

		filenameWithoutExtn := kopy.FileNameWOExt(filepath.Base(args[0]))
		zipFileName := filenameWithoutExtn + ".zip"
		//zip destination location
		zipDst := filepath.FromSlash(path.Join(args[1], zipFileName))

		//File Validation if the same file exits
		if sakto.IsFileExist(zipDst) {
			color.Red("compressed file already existed:" + zipDst)
			return
		}
		msg := `Started Compressing the file:`
		color.Blue(msg + " " + srcPath)
		itrlog.Infow(msg, "src", srcPath, "dist", dstPath, "log_time", time.Now().Format(logTimeFormat))

		files := []string{srcPath}
		os.MkdirAll(dstPath, os.ModePerm) //create a root folder
		if err := kopy.ComFiles(zipDst, files); err != nil {
			color.Red(err.Error())
			itrlog.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			return
		}
		msg = `Done compressing the file:`
		color.Green(msg + " " + srcPath)
		itrlog.Infow(msg, "dst", zipDst, "log_time", time.Now().Format(itrlog.LogTimeFormat))

	},
}

func init() {
	rootCmd.AddCommand(comfileCmd)
}
