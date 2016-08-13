package main

import (
	"flag"
	"fmt"
	"github.com/alyu/configparser"
	"github.com/wsxiaoys/terminal/color"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Byfly struct {
	balance  float64
	password string
	client   string
	login    string
	tariff   string
	status   string
}

type ByflyArgs struct {
	login       string
	password    string
	config      string
	onlyBalance bool
}

type Config struct {
	login    string
	password string
}

func prepareArgs() ByflyArgs {
	loginPtr := flag.String("l", "", "byfly login")
	passwordPtr := flag.String("p", "", "byfly password")
	configPtr := flag.String("f", "", `config file
        Example:
        login = 1111
        password = 1111`)
	onlyBalancPtr := flag.Bool("b", false, "show only balance")
	flag.Parse()
	return ByflyArgs{*loginPtr, *passwordPtr, *configPtr, *onlyBalancPtr}
}

func getUserHome() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func readConfig(args ByflyArgs, filename string) ByflyArgs {
	config, err := configparser.Read(filename)
	if err != nil {
		return args
	}
	sections, err := config.AllSections()
	if err != nil {
		return args
	}
	options := sections[0].Options()
	login, exists := options["login"]
	if !exists {
		login = ""
	}
	if len(login) > 0 {
		args.login = login
	}
	password, exists := options["password"]
	if !exists {
		password = ""
	}
	if len(password) > 0 {
		args.password = password
	}
	args.config = filename
	return args
}

func readConfigs(filename string) ByflyArgs {
	home, err := getUserHome()
	if err != nil {
		log.Fatal(err)
	}
	homeConfig, _ := filepath.Abs(filepath.Join(home, defaultConfigFile))
	currentConfig, _ := filepath.Abs(defaultConfigFile)
	configs := []string{homeConfig, currentConfig}
	if len(filename) > 0 {
		argsConfig, _ := filepath.Abs(filename)
		configs = append(configs, argsConfig)
	}
	args := ByflyArgs{"", "", "", false}
	for _, name := range configs {
		args = readConfig(args, name)
	}
	return args
}

func checkConfig(args *ByflyArgs) error {
	if len(args.config) > 0 {
		filename, _ := filepath.Abs(args.config)
		if stat, err := os.Stat(filename); os.IsNotExist(err) || stat.IsDir() {
			return fmt.Errorf("no such file or file is directory: %s", filename)
		}
	}
	return nil
}

func checkArgs(args *ByflyArgs) error {
	if len(args.login) == 0 {
		return fmt.Errorf("no login argument supplied")
	}
	if len(args.password) == 0 {
		return fmt.Errorf("no password argument supplied")
	}
	return nil
}

func mergeConfigs(args ByflyArgs, configArgs ByflyArgs) ByflyArgs {
	if len(args.login) == 0 && len(configArgs.config) > 0 {
		args.login = configArgs.login
	}
	if len(args.password) == 0 && len(configArgs.password) > 0 {
		args.password = configArgs.password
	}
	if len(args.config) == 0 && len(configArgs.config) > 0 {
		args.config = configArgs.config
	}
	return args
}

func (byfly *Byfly) getPage() (string, error) {
	formValues := []formValue{
		formValue{loginFieldName, byfly.login},
		formValue{passwordFieldName, byfly.password},
		formValue{redirectFieldName, redirectFieldValue},
	}
	resp, err := postFormRequest(urlLogin, formValues, agentHeaders)
	if err != nil {
		return "", err
	}
	return readBody(resp)
}

func (byfly *Byfly) printResult(onlyBalance bool) {
	if onlyBalance {
		fmt.Println(byfly.balance)
		return
	}
	data := [][]string{
		{"Абонент", byfly.client},
		{"Логин", byfly.login},
		{"Баланс", fmt.Sprintf("%.2f", byfly.balance)},
		{"Тариф", byfly.tariff},
		{"Статус", byfly.status},
	}
	maxLength := (func(data *[][]string) int {
		length := utf8.RuneCountInString((*data)[0][0])
		for _, value := range (*data)[:1] {
			if utf8.RuneCountInString(value[0]) > length {
				length = utf8.RuneCountInString(value[0])
			}
		}
		return length
	})(&data) + 1
	for _, value := range data {
		color.Printf(
			"@{g}%s%s=> ",
			value[0],
			strings.Repeat(" ", maxLength-utf8.RuneCountInString(value[0])))
		fmt.Println(value[1])
	}
}

func main() {
	args := prepareArgs()
	err := checkConfig(&args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	configArgs := readConfigs(args.config)
	args = mergeConfigs(args, configArgs)
	err = checkArgs(&args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	byfly := Byfly{password: args.password, login: args.login}
	body, err := byfly.getPage()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = byfly.parsePage(body)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	byfly.printResult(args.onlyBalance)
}
