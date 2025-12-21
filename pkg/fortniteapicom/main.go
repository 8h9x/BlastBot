package fortniteapicom

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/8h9x/fortgo/request"
)

type Client struct {
    httpClient *http.Client
    cosmetics  map[string]CosmeticItem
    nameMap    map[string]string
}

func New() *Client {
    return &Client{
        httpClient: &http.Client{},
        cosmetics: make(map[string]CosmeticItem),
        nameMap: make(map[string]string),
    }
}

func (c *Client) fetchCosmetics() ([]CosmeticItem, error) {
    req, err := request.MakeRequest(
        http.MethodGet,
        "https://fortnite-api.com",
        "v2/cosmetics/br",
    )
    if err != nil {
        return []CosmeticItem{}, err
    }

    res, err := c.httpClient.Do(req)
    if err != nil {
        return []CosmeticItem{}, err
    }

    resp, err := request.ParseResponse[CosmeticListResponse](res)
    if err != nil {
        return []CosmeticItem{}, err
    }

    return resp.Data.Data, nil
}

func (c *Client) fetchCosmeticByID(id string) (CosmeticItem, error) {
    req, err := request.MakeRequest(
        http.MethodGet,
        "https://fortnite-api.com",
        fmt.Sprintf("v2/cosmetics/br/%s", strings.ToLower(id)),
    )
    if err != nil {
        return CosmeticItem{}, err
    }

    res, err := c.httpClient.Do(req)
    if err != nil {
        return CosmeticItem{}, err
    }

    resp, err := request.ParseResponse[CosmeticSearchResponse](res)
    if err != nil {
        return CosmeticItem{}, err
    }

    return resp.Data.Data, nil
}

func (c *Client) fetchCosmeticByName(name string) (CosmeticItem, error) {
    req, err := request.MakeRequest(
        http.MethodGet,
        "https://fortnite-api.com",
        fmt.Sprintf("v2/cosmetics/br/search?name=%s", name),
    )
    if err != nil {
        return CosmeticItem{}, err
    }

    res, err := c.httpClient.Do(req)
    if err != nil {
        return CosmeticItem{}, err
    }

    resp, err := request.ParseResponse[CosmeticSearchResponse](res)
    if err != nil {
        return CosmeticItem{}, err
    }

    return resp.Data.Data, nil
}

func (c *Client) PreloadCache() error {
    cosmetics, err := c.fetchCosmetics()
    if err != nil {
        return err
    }

    for _, item := range cosmetics {
        c.cosmetics[strings.ToLower(item.ID)] = item
        c.nameMap[strings.ToLower(item.Name)] = strings.ToLower(item.ID)
    }

    return nil
}

func (c *Client) GetCosmeticByID(id string) (CosmeticItem, error) {
    lowerID := strings.ToLower(id)

    if _, ok := c.cosmetics[lowerID]; !ok {
        cosmetic, err := c.fetchCosmeticByID(id)
        if err != nil {
            return CosmeticItem{}, err
        }

        c.cosmetics[lowerID] = cosmetic
    }

    return c.cosmetics[lowerID], nil
}

func (c *Client) GetCosmeticByName(name string) (CosmeticItem, error) {
    if id, ok := c.nameMap[strings.ToLower(name)]; !ok {
        cosmetic, err := c.GetCosmeticByID(id)
        if err != nil {
            return CosmeticItem{}, err
        }

        c.nameMap[strings.ToLower(cosmetic.Name)] = cosmetic.ID

        return cosmetic, nil
    } else {
        cosmetic, err := c.fetchCosmeticByName(name)
        if err != nil {
            return CosmeticItem{}, err
        }

        c.nameMap[strings.ToLower(cosmetic.Name)] = cosmetic.ID

        return cosmetic, nil
    }
}