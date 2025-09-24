package main

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version         = "Unknown"
	compilerVersion = "Unknown"
)

// openInput returns an io.ReadCloser for the input source.
// If no FILE is provided or FILE is "-", it returns stdin as a ReadCloser.
func openInput(c *cli.Context) (io.ReadCloser, error) {
	if c.Args().Len() > 1 {
		return nil, fmt.Errorf("file count can not be more than 1")
	}
	name := c.Args().First()
	if len(name) == 0 || name == "-" {
		// Wrap stdin to satisfy the ReadCloser interface.
		return io.NopCloser(os.Stdin), nil
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func do(c *cli.Context) error {
	in, err := openInput(c)
	if err != nil {
		return cli.Exit(err, 1)
	}
	defer in.Close()

	wrap := c.Int("wrap")
	if wrap < 0 {
		return cli.Exit(fmt.Errorf("invalid wrap value %d", wrap), 1)
	}
	// Streaming decode
	if c.IsSet("decode") {
		if err := DecodeStream(c.Bool("url"), c.Bool("ignore-garbage"), in, os.Stdout); err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}
	// Streaming encode
	if err := EncodeStream(c.Bool("url"), c.Bool("no-padding"), wrap, in, os.Stdout); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func main() {
	usageText := "gbase64 [OPTION]... [FILE]\n\n" +
		"   With no FILE, or when FILE is -, read standard input."
	app := &cli.App{
		Name:      "gbase64",
		Usage:     "Encode/Decode data from FILE or standard input, to standard output",
		UsageText: usageText,
		Version:   fmt.Sprint(version, " (", compilerVersion, ")"),
		Action:    do,
	}
	app.Flags = append(app.Flags, &cli.BoolFlag{
		Name:    "decode",
		Aliases: []string{"d"},
		Usage:   "decode data",
	})
	app.Flags = append(app.Flags, &cli.BoolFlag{
		Name:    "ignore-garbage",
		Aliases: []string{"i"},
		Usage:   "when decoding, ignore new line characters \\r and \\n",
	})
	app.Flags = append(app.Flags, &cli.BoolFlag{
		Name:    "no-padding",
		Aliases: []string{"n"},
		Usage:   "when encoding, omit padding characters",
	})
	app.Flags = append(app.Flags, &cli.BoolFlag{
		Name:    "url",
		Aliases: []string{"u"},
		Usage:   "encode/decode with URL mode, defined in RFC 4648",
	})
	app.Flags = append(app.Flags, &cli.IntFlag{
		Name:    "wrap",
		Aliases: []string{"w"},
		Usage:   "when encoding, wrap encoded lines after some characters, 0 to disable line wrapping",
		Value:   76,
	})

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
