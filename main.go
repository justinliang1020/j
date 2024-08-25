package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var directories map[string]string = map[string]string{
	"justin": "/Users/justin/Documents/code/justin-blog",
	"notes":  "/Users/justin/Library/'Mobile Documents'/iCloud~md~obsidian/Documents/Justin",
}

func runCommand(command string, args ...string) (string, error) {
	// Create the command
	cmd := exec.Command(command, args...)

	// Run the command and get the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %v\nOutput: %s", err, output)
	}

	// Trim any leading/trailing whitespace and return the output
	return strings.TrimSpace(string(output)), nil
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "image",
				Aliases: []string{"i"},
				Usage:   "options for image maniuplation",
				Subcommands: []*cli.Command{
					{
						Name:  "compress",
						Usage: "compress an image using ffmpeg",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "input",
								Aliases:  []string{"i"},
								Usage:    "input image file",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "output",
								Aliases: []string{"o"},
								Usage:   "output image file",
							},
							&cli.IntFlag{
								Name:    "quality",
								Aliases: []string{"q"},
								Value:   80,
								Usage:   "compression quality (0-100)",
							},
						},
						Action: func(cCtx *cli.Context) error {
							input := cCtx.String("input")
							output := cCtx.String("output")
							quality := cCtx.Int("quality")
							if output == "" {
								// Get the base filename without extension
								baseName := strings.TrimSuffix(input, filepath.Ext(input))
								// Append "_compressed" and the original extension
								output = baseName + "_compressed" + filepath.Ext(input)
							}

							args := []string{
								"-i", input,
								"-q:v", fmt.Sprintf("%d", quality),
								"-y",
								output,
							}

							result, err := runCommand("ffmpeg", args...)
							if err != nil {
								return fmt.Errorf("failed to compress image: %v", err)
							}

							fmt.Printf("Image compressed successfully: %s\n", result)
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove an existing template",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("removed task template: ", cCtx.Args().First())
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
