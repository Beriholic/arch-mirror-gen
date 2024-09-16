package mirror

import (
	"io"
	"net/http"
)

func GetChinaMirrorList() ([]string, error) {
	url := "https://archlinux.org/mirrorlist/?country=CN&protocol=http&protocol=https&ip_version=4"
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	mirrorList, err := ParseMirrorList(string(body))

	if err != nil {
		return nil, err
	}

	return mirrorList, nil
}
