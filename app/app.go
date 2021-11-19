package app

import "github.com/urfave/cli/v2"

func GenerateCLI() *cli.App {
	usageText := "gobase64 [OPTION]... [COMMAND] [OPTION]... [FILE]\n\n" +
		"   With no FILE, read standard input."
	app := &cli.App{
		Name:      "gobase64",
		Usage:     "Encode/Decode data from FILE or standard input, to standard output",
		UsageText: usageText,
		Version:   "1.0.0",
		Commands:  []*cli.Command{},
	}
	app.Flags = append(app.Flags, &cli.StringFlag{
		Name:    "file",
		Aliases: []string{"f"},
		Usage:   "Read data from a file",
	})

	encodeCMD := &cli.Command{
		Name:   "encode",
		Usage:  "Encode data",
		Action: encode,
	}
	encodeCMD.Flags = append(encodeCMD.Flags, &cli.StringFlag{
		Name:    "format",
		Aliases: []string{"fmt"},
		Usage:   "Encode format, URL|URL_NOPADDING|STD|STD_NOPADDING",
		Value:   "URL_NOPADDING",
	})
	app.Commands = append(app.Commands, encodeCMD)

	decodeCMD := &cli.Command{
		Name:   "decode",
		Usage:  "Decode data",
		Action: decode,
	}
	decodeCMD.Flags = append(decodeCMD.Flags, &cli.StringFlag{
		Name:    "format",
		Aliases: []string{"fmt"},
		Usage:   "Decode format, URL|URL_NOPADDING|STD|STD_NOPADDING",
		Value:   "URL_NOPADDING",
	})
	app.Commands = append(app.Commands, decodeCMD)

	return app
}

func encode(c *cli.Context) error {
	return nil
}

func decode(c *cli.Context) error {
	return nil
}
