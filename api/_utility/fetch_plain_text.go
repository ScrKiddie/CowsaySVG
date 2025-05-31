package _utility

import (
	"io"
	"log/slog"
	"net/http"
)

func FetchPlainText(apiURL string) (string, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Warn(err.Error())
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}
