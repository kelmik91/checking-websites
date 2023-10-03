package serverGo

import (
	"net/http"
)

func RunServer() error {
	//handler := http.HandlerFunc(handleRequest)
	//http.Handle("/", handler)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		return err
	}
	return nil
}
