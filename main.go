package main

import (
	"extract-cli/config"
	"extract-cli/commands"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "extract-cli"
	app.Usage = "クエリを実行して指定のフォーマットに出力します。"
	app.Version = "0.0.1"
	app.Compiled = time.Now()

	app.Commands = []cli.Command{
		{
			Name:  "config",
			Usage: "config.tomlに設定されているDBの一覧を表示する",
			Action: func(c *cli.Context) error {
				fmt.Println("### config.tomlに設定されているDBの一覧を表示する")
				fmt.Println("/--------------------------------------/")
				for _, db := range config.GetConfig().Databases {
					fmt.Println(db.ToString())
				}

				fmt.Println("/--------------------------------------/")
				return nil
			},
		},
		{
			Name:   "csv",
			Usage:  "csv形式として出力する。csv -dn [config_key] -i [input_filepath] -o [output_path]",
			Flags: []cli.Flag{
				commands.DbNameFlag,
				commands.QueryFlag,
				commands.OutputPathFlag,
				commands.InputPathFlag,
			},
			Action: commands.CsvHandler,
		},
		{
			Name:   "xml",
			Usage:  "xml形式として出力する。TODO:",
			Action: func(c *cli.Context) error {
				fmt.Println("not implemented ....")
				return nil
			},
		},
		{
			Name:   "json",
			Usage:  "JSON形式として出力する。TODO:",
			Action: func(c *cli.Context) error {
				fmt.Println("not implemented ....")
				return nil
			},
		},
		{
			Name:   "excel",
			Usage:  "Excel形式として出力する。TODO:",
			Action: func(c *cli.Context) error {
				fmt.Println("not implemented ....")
				return nil
			},
		},
	}

	app.Run(os.Args)
}
