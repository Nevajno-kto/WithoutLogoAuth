package sms

import (
	"net/http"
	"net/url"
)

const (
	SMSSIMPLE_SCHEME  = "https"
	SMSSIMPLE_HOST    = "smsimple.ru"
	SMSSIMPLE_SEND    = "http_send.php"
	SMSSIMPLE_CHECK   = "http_check.php"
	SMSSIMPLE_ORIGIN  = "http_origins.php"
	SMSSIMPLE_BALANCE = "http_balance.php"

	SMSSIMPLE_LOGIN    = "Somebody"
	SMSSIMPLE_PASSWORD = "123"
)

func Send(phone, msg string) (*http.Response, error) {
	v := url.Values{
		"user":    []string{SMSSIMPLE_LOGIN},
		"pass":    []string{SMSSIMPLE_PASSWORD},
		"or_id":   []string{""},
		"phone":   []string{phone},
		"message": []string{msg},
	}

	u := &url.URL{
		Scheme:   SMSSIMPLE_SCHEME,
		Host:     SMSSIMPLE_HOST,
		Path:     SMSSIMPLE_SEND,
		RawQuery: v.Encode(),
	}

	resp, err := http.Get(u.String())

	if err != nil {

	}

	defer resp.Body.Close()

	// body, _ := ioutil.ReadAll(resp.Body)

	return resp, nil
}
