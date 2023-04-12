package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c EpicClient) ProfileOperationStr(credentials UserCredentialsResponse, operation string, profileId string, body string) (io.ReadCloser, error) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

	// https://fngw-mcp-gc-livefn.ol.epicgames.com/fortnite/api/game/v2/profile/:accountId/:route/:operation

	resp, err := c.Request("POST", fmt.Sprintf("https://fngw-mcp-gc-livefn.ol.epicgames.com/fortnite/api/game/v2/profile/%s/client/%s?profileId=%s&rvn=-1", credentials.AccountId, operation, profileId), headers, body)
	if err != nil {
		return nil, err
	}

	// defer resp.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", err
	// }

	return resp.Body, nil
}

func (c EpicClient) ProfileOperation(credentials UserCredentialsResponse, operation string, profileId string, body any) (any, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return c.ProfileOperationStr(credentials, operation, profileId, string(bodyBytes))
}

// struct methods cannot use generic types, so this is seperated
// func ProfileOperation[T interface{}](c EpicClient, credentials UserCredentialsResponse, operation string, profileId string, body T) (any, error) {
// 	bodyBytes, err := json.Marshal(body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c.ProfileOperationBasic(credentials, operation, profileId, string(bodyBytes))
// }

// func (c EpicClient) AddFavoriteMnemonic(credentials UserCredentialsResponse, mnemonic string) error {
// 	headers := http.Header{}
// 	headers.Set("Content-Type", "application/json")
// 	headers.Set("Authorization", "Bearer "+credentials.AccessToken)

// 	resp, err := c.Request("POST", fmt.Sprintf("https://fn-service-discovery-live-public.ogs.live.on.epicgames.com/api/v1/links/favorites/%s/%s", credentials.AccountId, mnemonic), headers, "{}")
// 	if err != nil {
// 		return err
// 	}

// 	if resp.StatusCode != http.StatusNoContent {
// 		return fmt.Errorf("invalid mnemonic")
// 	}

// 	return nil
// }

// export async function profileOperationRequest(client: Client, profileOperation: string, profileId: ProfileId, payload?: object) {
//     const { body, statusCode } = await httpRequest(client, `https://fortnite-public-service-prod11.ol.epicgames.com/fortnite/api/game/v2/profile/${client.accountId}/client/${profileOperation}?profileId=${profileId}&rvn=-1`, {
//         method: "POST",
//         headers: {
//             "Content-Type": "application/json"
//         },
//         body: JSON.stringify(payload ?? {})
//     });

//     return { body, statusCode };
// };
