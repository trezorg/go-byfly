package main

import (
	"net/http"
)

var agentHeaders = http.Header{
	"User-Agent": []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36"},
	"Accept":     []string{"*/*"},
}

const (
	urlLogin           = "https://issa.beltelecom.by/main.html"
	passwordFieldName  = "passwd"
	loginFieldName     = "oper_user"
	redirectFieldName  = "redirect"
	redirectFieldValue = "/main.html"
	balanceImageSrc    = "/data/img/liber/coins.png"
	defaultConfigFile  = "byfly.conf"
	tariffText         = "Тарифный план на услуги"
	clientText         = "Абонент"
	statusText         = "Статус блокировки"
)
