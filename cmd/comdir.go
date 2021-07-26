package cmd

import (
	"bytes"
	"io"
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
var comdirCmd = &cobra.Command{
	Use:   "comdir",
	Short: "Compress the directory or a folder",
	Long:  `comdir helps to compress a directory or a folder from your terminal`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//source and destination paths
		srcPath := filepath.FromSlash(args[0])
		dstPath := filepath.FromSlash(args[1])

		//compose a zip filename
		fileNameWithoutExtn := kopy.FileNameWOExt(filepath.Base(args[0])) // Returns a filename without an extension.
		//zip externsion <.tar.gz>
		zipDir := fileNameWithoutExtn + kopy.ComFileFormat

		//zip destination dir
		zipDst := filepath.FromSlash(path.Join(args[1], zipDir))

		//File Validation if the same file exits
		if sakto.IsFileExist(zipDst) {
			color.Red("compressed file already existed:" + zipDst)
			return
		}
		msg := `Started Compressing the folder:`
		color.Blue(msg + " " + srcPath)
		itrlog.Infow(msg, "src", srcPath, "dist", dstPath, "log_time", time.Now().Format(logTimeFormat))

		// Start compressing the entire directory or a folder using the tar + gzip
		var buf bytes.Buffer
		if err := kopy.CompressDIR(srcPath, &buf, ingoreDir); err != nil {
			color.Red(err.Error())
			itrlog.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			return
		}

		// write the .tar.gzip
		os.MkdirAll(dstPath, os.ModePerm) // Create the root folder first
		fileToWrite, err := os.OpenFile(zipDst, os.O_CREATE|os.O_RDWR, os.FileMode(0600))
		if err != nil {
			itrlog.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			panic(err)
		}
		if _, err := io.Copy(fileToWrite, &buf); err != nil {
			itrlog.Errorw("error", "err", err, "log_time", time.Now().Format(itrlog.LogTimeFormat))
			panic(err)
		}
		defer fileToWrite.Close()

		msg = `Done compressing the directory or a folder:`
		color.Green(msg + " " + srcPath)
		itrlog.Infow(msg, "dst", zipDst, "log_time", time.Now().Format(itrlog.LogTimeFormat))

	},
}

func init() {
	rootCmd.AddCommand(comdirCmd)
}
