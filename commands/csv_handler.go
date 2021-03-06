package commands

import (
	"database/sql"
	"extract-cli/config"
	"extract-cli/data"
	"extract-cli/helpers"

	"github.com/urfave/cli"
	"log"
	"errors"
)

var QueryFlag cli.Flag = cli.StringFlag{
	Name:  "query, q",
	Usage: "queryを直接指定",
}

var InputPathFlag cli.Flag = cli.StringFlag{
	Name: "inputpath, i",
	Usage: "SQLファイルを指定",
}

var OutputPathFlag cli.Flag = cli.StringFlag{
	Name: "outputpath, o",
	Usage: "出力先のファイルパスを指定",
}

var DbNameFlag cli.Flag = cli.StringFlag{
	Name: "dbname, dn",
	Usage: "configに設定されているdatabaseのname",
}

type Args struct {
	DbName     string
	Query      string
	OutputPath string
}

func NewArgs(c *cli.Context) (*Args, error) {

	dbName := c.String("db")
	if dbName == "" {
		return nil, errors.New("not selected db")
	}

	// query を取得
	var query string = ""
	q := c.String("query")
	i := c.String("inputpath")
	var filepath string
	if i != "" {
		filepath = i
	} else {
		filepath = config.GetConfig().Default.Input
	}

	if q != ""{
		query = q
	} else if filepath != "" {
		var err error
		query, err = helpers.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("not selected query")
	}

	// output
	o := c.String("outputpath")
	var output_filepath string
	if o != "" {
		output_filepath = o
	} else {
		output_filepath = config.GetConfig().Default.Output
	}

	if output_filepath == "" {
		output_filepath = "./output.csv"
	}

	return &Args{
		DbName: dbName,
		Query: query,
		OutputPath: output_filepath,
	}, nil
}

func CsvHandler(c *cli.Context) error {
	args, err := NewArgs(c)
	if err != nil {
		return cli.NewExitError(err, 404)
	}

	// config.tomlからDB接続情報を取得
	db, err := config.GetConfig().GetDatabase(args.DbName)
	if err != nil {
		return cli.NewExitError(err, 404)
	}
	log.Println(`
selected database
	` + db.ToString())

	// Connection 生成
	con, err := data.NewConnection(db)
	if err != nil{
		return cli.NewExitError(err, 404)
	}

	// query 実行
	log.Printf("execute query ...\n")
	client := data.NewDbClient(con)
	rows, err := client.Execute(args.Query)
	if err != nil {
		return cli.NewExitError(err, 500)
	}

	columns, err := rows.Columns() // カラム名を取得
	if err != nil {
		return cli.NewExitError(err, 500)
	}
	rawBytes := make([]sql.RawBytes, len(columns))

	//  rows.Scan は引数に `[]interface{}`が必要.
	scanArgs := make([]interface{}, len(rawBytes))
	for i := range rawBytes {
		scanArgs[i] = &rawBytes[i]
	}

	// csv出力するの連想配列
	recordes := [][]string{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return cli.NewExitError(err, 500)
		}

		// 1行分
		r := []string{}
		// カラム数分ループ
		for _, col := range rawBytes {
			var value string
			if col != nil {
				value = string(col)
			}
			r = append(r, value)
		}
		// csv１行追加
		recordes = append(recordes, r)
	}

	// csv 出力
	err = helpers.ToCsv(args.OutputPath, recordes)
	if err != nil {
		return cli.NewExitError(err, 500)
	}

	log.Printf("end\n")

	return nil
}
