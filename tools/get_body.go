package tools

import (
	"context"
	"io"
	"net/http"
	"time"
)

func GetBody(u string) (int, []byte, error) {

	try := 10
	for i := 0; i != try; i++ {
		code, body, err := GetBodyOne(u)
		if code == 200 || code == 404 {
			return code, body, err
		}
		<-time.NewTimer(30 * time.Second).C
	}

	<-time.NewTimer(60 * time.Second).C
	return GetBodyOne(u)

}
func GetBodyOne(u string) (int, []byte, error) {

	ctx, _ := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//panic(err)
		return 0, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	return res.StatusCode, body, err

}
