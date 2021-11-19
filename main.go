package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Armourstill/gobase64/app"
	"github.com/Armourstill/gobase64/input"
)

var encoderMap = map[string]*base64.Encoding{
	"URL":           base64.URLEncoding,
	"URL_NOPADDING": base64.RawURLEncoding,
	"STD":           base64.StdEncoding,
	"STD_NOPADDING": base64.RawStdEncoding,
}

func encode(c *cli.Context) error {
	// 通过c.Args().Get(0)获取文件名
	var encoder = encoderMap["URL_NOPADDING"]
	if c.IsSet("format") {
		fmt.Println("Format: ", c.String("format"))
		encoder = encoderMap[c.String("format")]
	}
	data, _ := input.ReadFromStdin()
	encoded := make([]byte, encoder.EncodedLen(len(data)))
	encoder.Encode(encoded, data)
	fmt.Println(string(encoded))
	return nil
}

func decode(c *cli.Context) error {
	var encoder = encoderMap["URL_NOPADDING"]
	if c.IsSet("format") {
		fmt.Println("Format: ", c.String("format"))
		encoder = encoderMap[c.String("format")]
	}
	data, _ := input.ReadFromStdin()
	decoded := make([]byte, encoder.DecodedLen(len(data)))
	_, err := encoder.Decode(decoded, data)
	if err != nil {
		return err
	}
	fmt.Println(string(decoded))
	return nil
}

func main() {
	if err := app.GenerateCLI().Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
