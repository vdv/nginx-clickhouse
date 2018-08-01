package nginx

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/satyrius/gonx"
	"github.com/vdv/nginx-clickhouse/config"
)

func GetParser(config *config.Config) (*gonx.Parser, error) {

	// Use nginx config file to extract format by the name
	nginxConfig := strings.NewReader(fmt.Sprintf("%s%s%s", `
		http {
			log_format   main  '`, config.Nginx.LogFormat, `';
		}
	`))

	return gonx.NewNginxParser(nginxConfig, config.Nginx.LogType)
}

func ParseField(column config.Column, value string) interface{} {
	key := column.VarName

	switch column.VarType {
	case "datetime":
		t, err := time.Parse(config.NginxTimeLayout, value)

		if err == nil {
			return t.UTC().Format(config.CHTimeLayout)
		}

		return value

	case "integer":
		if value == "-" {
			value = "-1"
		}

		val, err := strconv.Atoi(value)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"key":   key,
				"value": value,
			}).Error("Error to convert string to int")
		}

		return val
	case "float":
		if value == "-" {
			value = "-1"
		}

		val, err := strconv.ParseFloat(value, 32)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"key":   key,
				"value": value,
			}).Error("Error to convert string to float32")
		}

		return val
	default:
		return value
	}

	return value
}

func ParseLogs(parser *gonx.Parser, logLines []string) []gonx.Entry {

	logReader := strings.NewReader(strings.Join(logLines, "\n"))
	reader := gonx.NewParserReader(logReader, parser)

	var logs []gonx.Entry

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		// Process the record... e.g.
		logs = append(logs, *rec)
	}

	return logs
}
