package requests

import (
	balancer "blcMod/internal/balancer"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ResponseData data of server pespons
type ResponseData struct {
	Name     string
	ServAddr string
	Body     string
}

func SendRequestToBalancer(lb *balancer.Balancer, path string, timeout time.Duration) (ResponseData, error) {
	var respData ResponseData

	serv := lb.ChooseServer()
	if serv == nil {
		return respData, fmt.Errorf("no available servers")
	}

	url := fmt.Sprintf("http://%s%s", serv.ServerConfig.Adress.Host, path) // Полный URL
	respData.Name = serv.Name                                              // ServerName

	client := &http.Client{
		Timeout: timeout,
	}
	startTime := time.Now() // Засекаем время начала запроса
	resp, err := client.Get(url)
	elapsedTime := time.Since(startTime) // Вычисляем прошедшее время

	if err != nil {
		return respData, fmt.Errorf("error in request %s: %w", serv.ServerConfig.Adress.String(), err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return respData, fmt.Errorf("read body error %s: %w", serv.ServerConfig.Adress.String(), err)
	}

	respData.ServAddr = serv.ServerConfig.Adress.String()
	respData.Body = string(bodyBytes)

	fmt.Printf("request end in %s\n", elapsedTime)
	return respData, nil
}
