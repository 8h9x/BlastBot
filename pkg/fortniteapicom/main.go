package fortniteapicom

import (
    "github.com/8h9x/vinderman/request"
    "net/http"
    "strings"
)

func FetchCosmetics() (Cosmetics, error) {
    resp, err := http.Get("https://fortnite-api.com/v2/cosmetics/br")
    if err != nil {
        return Cosmetics{}, err
    }

    defer resp.Body.Close()

    res, err := request.ResponseParser[CosmeticListResponse](resp)
    if err != nil {
        return Cosmetics{}, err
    }

    items := make(Cosmetics, len(res.Body.Data))
    for _, item := range res.Body.Data {
        items[strings.ToLower(item.ID)] = item
    }

    return items, nil
}