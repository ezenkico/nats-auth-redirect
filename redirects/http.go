package redirects

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"encoding/json"

	natsauth "github.com/ezenkico/nats-auth-redirect/base"
)

func SetupHttpAuth(server string) {
	natsauth.Listen(func(token *string) (*natsauth.ResponseData, error) {
		log.Printf("(main)Token: %s", *token)

		req, err := http.NewRequest(http.MethodGet, server, nil)
		if err != nil {
			return nil, err
		}

		bearer := "Bearer " + *token

		req.Header.Add("Authorization", bearer)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode > 300 {
			return nil, errors.New("invalid token")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		fmt.Printf("Data returned: %q", string(body))

		var res natsauth.ResponseData

		err = json.Unmarshal(body, &res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}, natsauth.GetConnectionDataEnv())
}
