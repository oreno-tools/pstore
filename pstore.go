package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	AppVersion = "0.0.1"
)

var (
	argProfile   = flag.String("profile", "", "Profile 名を指定.")
	argRole      = flag.String("role", "", "Role ARN を指定.")
	argRegion    = flag.String("region", "ap-northeast-1", "Region 名を指定.")
	argEndpoint  = flag.String("endpoint", "", "AWS API のエンドポイントを指定.")
	argVersion   = flag.Bool("version", false, "バージョンを出力.")
	argCsv       = flag.Bool("csv", false, "CSV 形式で出力する")
	argJson      = flag.Bool("json", false, "JSON 形式で出力する")
	argPut       = flag.Bool("put", false, "パラメータを追加する")
	argGet       = flag.Bool("get", false, "パラメータの値を取得する")
	argName      = flag.String("name", "", "パラメータの名前を指定する")
	argValue     = flag.String("value", "", "パラメータ名を値を指定する")
	argOverwrite = flag.Bool("overwrite", false, "パラメータを上書きする")
	argSecure    = flag.Bool("secure", false, "SecureString でパラメータを追加する")
	argInsecure  = flag.Bool("insecure", false, "SecureString を出力する")
	argList      = flag.Bool("list", false, "StringList でパラメータを追加する")
	argDel       = flag.Bool("del", false, "パラメータを削除する")
)

type Parameters struct {
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Name             string `json:"name"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	Version          string `json:"version"`
	LastModifiedDate string `json:"last_modified_date"`
}

func outputTbl(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Value", "Type", "Version", "LastModifiedDate"})
	for _, value := range data {
		table.Append(value)
	}
	table.Render()
}

func outputCsv(data [][]string) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	for _, record := range data {
		if err := w.Write(record); err != nil {
			fmt.Println("Write error: ", err)
			return
		}
		w.Flush()
	}
	fmt.Println(buf.String())
}

func outputJson(data [][]string) {
	var rs []Parameter
	for _, record := range data {
		r := Parameter{Name: record[0], Value: record[1], Type: record[2],
			Version: record[3], LastModifiedDate: record[4]}
		rs = append(rs, r)
	}
	rj := Parameters{
		Parameters: rs,
	}
	b, err := json.Marshal(rj)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}
	os.Stdout.Write(b)
}

func awsSsmClient(profile string, region string, role string) *ssm.SSM {
	var config aws.Config
	if profile != "" && role == "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region),
			Credentials: creds,
			Endpoint:    aws.String(*argEndpoint)}
	} else if profile == "" && role != "" {
		sess := session.Must(session.NewSession())
		creds := stscreds.NewCredentials(sess, role)
		config = aws.Config{Region: aws.String(region),
			Credentials: creds,
			Endpoint:    aws.String(*argEndpoint)}
	} else if profile != "" && role != "" {
		sess := session.Must(session.NewSessionWithOptions(session.Options{Profile: profile}))
		assumeRoler := sts.New(sess)
		creds := stscreds.NewCredentialsWithClient(assumeRoler, role)
		config = aws.Config{Region: aws.String(region),
			Credentials: creds,
			Endpoint:    aws.String(*argEndpoint)}
	} else {
		config = aws.Config{Region: aws.String(region),
			Endpoint: aws.String(*argEndpoint)}
	}

	sess := session.New(&config)
	ssmClient := ssm.New(sess)
	return ssmClient
}

func putParameter(ssmClient *ssm.SSM, pName string, pType string, pValue string) {
	params := &ssm.PutParameterInput{
		Name:        aws.String(pName),
		Value:       aws.String(pValue),
		Description: aws.String(pName),
		Type:        aws.String(pType),
	}
	if *argOverwrite {
		params.SetOverwrite(true)
	}

	_, err := ssmClient.PutParameter(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func delParameter(ssmClient *ssm.SSM, pName string) {
	params := &ssm.DeleteParameterInput{
		Name: aws.String(pName),
	}
	_, err := ssmClient.DeleteParameter(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func convertDate(d time.Time) (convertedDate string) {
	const layout = "2006-01-02 15:04:05"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	convertedDate = d.In(jst).Format(layout)

	return convertedDate
}

func getParameter(ssmClient *ssm.SSM, pName string) (pValue string, pVersion string) {
	params := &ssm.GetParameterInput{
		Name:           aws.String(pName),
		WithDecryption: aws.Bool(true),
	}
	v, err := ssmClient.GetParameter(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *v.Parameter.Type == "SecureString" && !*argInsecure {
		pValue = "******************"
	} else {
		pValue = *v.Parameter.Value
	}

	pVersion = strconv.FormatInt(*v.Parameter.Version, 10)

	return pValue, pVersion
}

func getParameterValue(ssmClient *ssm.SSM, pName string) (pValue string) {
	pValue, _ = getParameter(ssmClient, pName)
	return pValue
}

func listParameters(ssmClient *ssm.SSM) {
	params := &ssm.DescribeParametersInput{}

	allParameters := [][]string{}
	for {
		res, err := ssmClient.DescribeParameters(params)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for _, r := range res.Parameters {
			convertedDate := convertDate(*r.LastModifiedDate)
			pValue, pVersion := getParameter(ssmClient, *r.Name)
			Parameter := []string{
				*r.Name,
				pValue,
				*r.Type,
				pVersion,
				convertedDate,
			}
			allParameters = append(allParameters, Parameter)
		}
		if res.NextToken == nil {
			break
		}
		params.SetNextToken(*res.NextToken)
		continue
	}

	if *argCsv == true {
		outputCsv(allParameters)
	} else if *argJson == true {
		outputJson(allParameters)
	} else {
		outputTbl(allParameters)
	}
}

func main() {
	flag.Parse()

	if *argVersion {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	ssmClient := awsSsmClient(*argProfile, *argRegion, *argRole)

	var yorN string
	if *argPut {
		if *argName == "" {
			fmt.Println("パラメータの名前を指定して下さい.")
			os.Exit(1)
		}

		// Type を選択する (デフォルトは String とする)
		var pType string
		if *argSecure {
			pType = "SecureString"
		} else if *argList || strings.Contains(*argName, "/") {
			pType = "StringList"
		} else {
			pType = "String"
		}

		// パラメータの値を入力する
		var pValue string
		// 引数から入力
		if *argValue != "" {
			pValue = *argValue
		}
		// 標準入力から入力
		if pValue == "" {
			stat, _ := os.Stdin.Stat()
			// 標準入力があれば処理を継続
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				stdin := bufio.NewScanner(os.Stdin)
				stdin.Scan()
				pValue = stdin.Text()
				if err := stdin.Err(); err != nil {
					fmt.Fprintln(os.Stderr, "標準入力から値を読み込めません: ", err)
				}
			}
		}
		// インタラクティブに入力
		if pValue == "" {
			fmt.Print("パラメータの値を入力しますか?(y/n): ")
			fmt.Scan(&yorN)
			switch yorN {
			case "y", "Y":
				fmt.Println("パラメータの値を入力して下さい: ")
				if pType == "SecureString" {
					maskedValue1, err := terminal.ReadPassword(0)
					if err != nil {
						fmt.Println("入力した値が不正です.")
						os.Exit(1)
					}
					fmt.Println("パラメータの値をもう一度入力して下さい: ")
					maskedValue2, err := terminal.ReadPassword(0)
					if err == nil && string(maskedValue1) == string(maskedValue2) {
						pValue = string(maskedValue2)
					} else {
						fmt.Println("入力した値が不正です.")
						os.Exit(1)
					}
				} else {
					fmt.Scan(&pValue)
				}
			case "n", "N":
				fmt.Println("処理を中止します. パラメータの値を指定して下さい.")
				os.Exit(1)
			default:
				fmt.Println("処理を中止します. パラメータの値を指定して下さい.")
				os.Exit(1)
			}
		}
		putParameter(ssmClient, *argName, pType, pValue)
	} else if *argGet {
		if *argName == "" {
			fmt.Println("パラメータの名前を指定して下さい.")
			os.Exit(1)
		}
		fmt.Println(getParameterValue(ssmClient, *argName))
		os.Exit(0)
	} else if *argDel {
		if *argName == "" {
			fmt.Println("パラメータの名前を指定して下さい.")
			os.Exit(1)
		}
		fmt.Print("パラメータを削除しますか?(y/n): ")
		fmt.Scan(&yorN)
		switch yorN {
		case "y", "Y":
			fmt.Println("パラメータを削除します.")
			delParameter(ssmClient, *argName)
		case "n", "N":
			fmt.Println("処理を中止します.")
			os.Exit(0)
		default:
			fmt.Println("処理を中止します.")
			os.Exit(0)
		}
	} else {
		listParameters(ssmClient)
	}
}
