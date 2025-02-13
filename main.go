package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func fromStdin() []byte {
	var data []byte
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}
	return data
}

func getdata(c *cli.Context) (data []byte, err error) {
	if c.Args().Len() > 1 {
		return nil, fmt.Errorf("file count can not be more than 1")
	}
	name := c.Args().First()
	if len(name) == 0 || name == "-" {
		data = fromStdin()
	} else {
		data, err = os.ReadFile(name)
	}
	return
}

func do(c *cli.Context) error {
	data, err := getdata(c)
	if err != nil {
		return cli.Exit(err, 1)
	}
	format := STD
	if c.IsSet("url") {
		format = URL
	}
	wrap := c.Int("wrap")
	if wrap < 0 {
		return cli.Exit(fmt.Errorf("invalid wrap value %d", wrap), 1)
	}

	if c.IsSet("decode") {
		output, err := Decode(format, data)
		if err != nil {
			return cli.Exit(err, 1)
		}
		if len(output) != 0 {
			fmt.Fprintln(os.Stdout, string(output))
		}
		return nil
	}

	output, err := Encode(format, c.Bool("no-padding"), data)
	if err != nil {
		return cli.Exit(err, 1)
	}
	if len(output) == 0 {
		return nil
	}
	if wrap == 0 {
		fmt.Fprintln(os.Stdout, string(output))
		return nil
	}
	for index := 0; index < len(output); index += wrap {
		fmt.Fprintln(os.Stdout, string(output[index:min(index+wrap, len(output))]))
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
		Version:   "v0.9.0",
		Action:    do,
	}
	app.Flags = append(app.Flags, &cli.BoolFlag{
		Name:    "decode",
		Aliases: []string{"d"},
		Usage:   "decode data",
	})
	app.Flags = append(app.Flags, &cli.BoolFlag{
		Name:    "url",
		Aliases: []string{"u"},
		Usage:   "encode/decode with URL mode, defined in RFC 4648",
	})
	app.Flags = append(app.Flags, &cli.BoolFlag{
		Name:    "no-padding",
		Aliases: []string{"n"},
		Usage:   "when encoding, omit padding characters",
	})
	app.Flags = append(app.Flags, &cli.BoolFlag{
		Name:    "ignore-garbage",
		Aliases: []string{"i"},
		Usage:   "when decoding, ignore new line characters \\r and \\n",
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
