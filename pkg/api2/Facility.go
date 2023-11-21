package api2

import (
	"io"
	"net/http"
)

func callOtherAPI(url string) ([]byte, error) {
	// Создаем запрос
	var responseBytes []byte
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return responseBytes, err
		//return "", err
	}
	// Отправляем запрос
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return responseBytes, err
		//return "", err
	}
	defer res.Body.Close()

	// Читаем ответ
	responseBytes, err = io.ReadAll(res.Body)
	if err != nil {
		return responseBytes, err
		//return "", err
	}
	return responseBytes, err
	//return string(responseBytes), nil
}
