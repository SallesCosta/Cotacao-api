package escrequest

import (
	"context"
	"io"
	"net/http"
)

func Escrequest(method string, url string, ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// json.NewEncoder(w).Encode(body)
	return body
}
