package techbookfest

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

type Client struct {
	httpClient *http.Client
	graphqlURL *url.URL
	xsrfToken  string
}

func NewTechBookFest(ctx context.Context) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	c := &http.Client{Jar: jar}

	marketURL, err := url.Parse("https://techbookfest.org/market")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, marketURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if _, err := c.Do(req); err != nil {
		return nil, err
	}

	var xsrfToken string
	for _, cookie := range jar.Cookies(marketURL) {
		if cookie.Name == "XSRF-TOKEN" {
			xsrfToken = cookie.Value
		}
	}

	u, err := url.Parse("https://techbookfest.org/api/graphql")
	if err != nil {
		return nil, err
	}

	return &Client{httpClient: c, graphqlURL: u, xsrfToken: xsrfToken}, nil
}

func (c *Client) GetMarketDashboard(ctx context.Context, first int) (*MarketDashboard, error) {
	mdq := newMarketDashboardQuery(first)
	var md MarketDashboard
	if err := c.do(ctx, mdq, &md); err != nil {
		return nil, err
	}

	return &md, nil
}

func (c *Client) GetProductInfo(ctx context.Context, productID string) (*ProductInfo, error) {
	piq := newProductInfoQuery(productID)
	var pi ProductInfo
	if err := c.do(ctx, piq, &pi); err != nil {
		return nil, err
	}

	return &pi, nil
}

func (c *Client) do(ctx context.Context, reqQuery, v interface{}) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(reqQuery); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.graphqlURL.String(), &b)
	if err != nil {
		return err
	}
	req.Header.Set("x-xsrf-token", c.xsrfToken)
	req.Header.Set("content-type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP Status: %d.", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

type marketDashboardQuery struct {
	OperationName string                   `json:"operationName"`
	Variables     marketDashboardVariables `json:"variables"`
	Query         string                   `json:"query"`
}

type marketDashboardVariables struct {
	First                       int    `json:"first"`
	SuggestionProductInfo1ID    string `json:"suggestionProductInfo1ID"`
	SuggestionProductInfo2ID    string `json:"suggestionProductInfo2ID"`
	SuggestionProductInfo3ID    string `json:"suggestionProductInfo3ID"`
	SuggestionProductInfo4ID    string `json:"suggestionProductInfo4ID"`
	SuggestionProductInfo5ID    string `json:"suggestionProductInfo5ID"`
	SuggestionProductInfo6ID    string `json:"suggestionProductInfo6ID"`
	SuggestionProductInfo7ID    string `json:"suggestionProductInfo7ID"`
	SuggestionProductInfo8ID    string `json:"suggestionProductInfo8ID"`
	ShowProductInfoThumbnailDLC bool   `json:"showProductInfoThumbnailDLC"`
}

type MarketDashboard struct {
	Data struct {
		AllProductVariants struct {
			PageInfo struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
			Nodes []struct {
				Products struct {
					Nodes []struct {
						ID string `json:"id"`
					} `json:"nodes"`
				} `json:"products"`
			} `json:"nodes"`
		} `json:"allProductVariants"`
	} `json:"data"`
}

func newMarketDashboardQuery(first int) *marketDashboardQuery {
	return &marketDashboardQuery{
		OperationName: "MarketDashboardQuery",
		Variables: marketDashboardVariables{
			First:                       first,
			SuggestionProductInfo1ID:    "ProductInfo:6264122186399744",
			SuggestionProductInfo2ID:    "ProductInfo:6264122186399744",
			SuggestionProductInfo3ID:    "ProductInfo:6264122186399744",
			SuggestionProductInfo4ID:    "ProductInfo:6264122186399744",
			SuggestionProductInfo5ID:    "ProductInfo:6264122186399744",
			SuggestionProductInfo6ID:    "ProductInfo:6264122186399744",
			SuggestionProductInfo7ID:    "ProductInfo:6264122186399744",
			SuggestionProductInfo8ID:    "ProductInfo:6264122186399744",
			ShowProductInfoThumbnailDLC: false,
		},
		Query: "query MarketDashboardQuery($first: Int!, $after: String, $showProductInfoThumbnailDLC: Boolean!, $suggestionProductInfo1ID: ID!, $suggestionProductInfo2ID: ID!, $suggestionProductInfo3ID: ID!, $suggestionProductInfo4ID: ID!, $suggestionProductInfo5ID: ID!, $suggestionProductInfo6ID: ID!, $suggestionProductInfo7ID: ID!, $suggestionProductInfo8ID: ID!) {\n  suggestionProductInfo1: product(id: $suggestionProductInfo1ID) {\n    ...MarketTopSuggestionProductInfoFragment\n    __typename\n  }\n  suggestionProductInfo2: product(id: $suggestionProductInfo2ID) {\n    ...MarketTopSuggestionProductInfoFragment\n    __typename\n  }\n  suggestionProductInfo3: product(id: $suggestionProductInfo3ID) {\n    ...MarketTopSuggestionProductInfoFragment\n    __typename\n  }\n  suggestionProductInfo4: product(id: $suggestionProductInfo4ID) {\n    ...MarketTopSuggestionProductInfoFragment\n    __typename\n  }\n  suggestionProductInfo5: product(id: $suggestionProductInfo5ID) {\n    ...MarketTopSuggestionProductInfoFragment\n    __typename\n  }\n  suggestionProductInfo6: product(id: $suggestionProductInfo6ID) {\n    ...MarketTopSuggestionProductInfoFragment\n    __typename\n  }\n  suggestionProductInfo7: product(id: $suggestionProductInfo7ID) {\n    ...MarketTopSuggestionProductInfoFragment\n    __typename\n  }\n  suggestionProductInfo8: product(id: $suggestionProductInfo8ID) {\n    ...MarketTopSuggestionProductInfoFragment\n    __typename\n  }\n  allProductVariants: productVariants(first: $first, after: $after, input: {route: \"market\"}) {\n    pageInfo {\n      hasNextPage\n      endCursor\n      __typename\n    }\n    nodes {\n      ...MarketTopAllProductVariantFragment\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment MarketTopSuggestionProductInfoFragment on ProductInfo {\n  id\n  ...ProductInfoThumbnailAtomFragment\n  __typename\n}\n\nfragment ProductInfoThumbnailAtomFragment on ProductInfo {\n  id\n  databaseID\n  name\n  organization {\n    id\n    name\n    __typename\n  }\n  images(first: 1) {\n    nodes {\n      id\n      url\n      height\n      width\n      __typename\n    }\n    __typename\n  }\n  downloadContent @include(if: $showProductInfoThumbnailDLC) {\n    id\n    fileName\n    downloadURL\n    __typename\n  }\n  __typename\n}\n\nfragment MarketTopAllProductVariantFragment on ProductVariant {\n  id\n  products(first: 1) {\n    nodes {\n      ...ProductInfoThumbnailAtomFragment\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n",
	}
}

type productInfoQuery struct {
	OperationName string               `json:"operationName"`
	Variables     productInfoVariables `json:"variables"`
	Query         string               `json:"query"`
}

type productInfoVariables struct {
	ProductInfoID               string `json:"productInfoID"`
	ShowProductInfoThumbnailDLC bool   `json:"showProductInfoThumbnailDLC"`
}

type ProductInfo struct {
	Data struct {
		Product struct {
			DatabaseID               string `json:"databaseID"`
			Name                     string `json:"name"`
			Description              string `json:"description"`
			Page                     int    `json:"page"`
			FirstAppearanceEventName string `json:"firstAppearanceEventName"`
			Organization             struct {
				Name string `json:"name"`
			} `json:"organization"`
			ProductVariants struct {
				Nodes []struct {
					Name                   string `json:"name"`
					Price                  int    `json:"price"`
					MarketShippingRequired bool   `json:"marketShippingRequired"`
				} `json:"nodes"`
			} `json:"productVariants"`
		} `json:"product"`
	} `json:"data"`
}

func newProductInfoQuery(productInfoID string) *productInfoQuery {
	return &productInfoQuery{
		OperationName: "ProductInfoQuery",
		Variables: productInfoVariables{
			ProductInfoID:               productInfoID,
			ShowProductInfoThumbnailDLC: false,
		},
		Query: "query ProductInfoQuery($productInfoID: ID!, $showProductInfoThumbnailDLC: Boolean!) {\n  viewer {\n    id\n    __typename\n  }\n  product(id: $productInfoID) {\n    ...ProductInfoFragment\n    ...ProductPurchaseCompleteFragment\n    __typename\n  }\n}\n\nfragment ProductInfoFragment on ProductInfo {\n  id\n  databaseID\n  name\n  description\n  page\n  firstAppearanceEventName\n  loginUserBookShelfItem {\n    id\n    causedAt\n    __typename\n  }\n  images(first: 4) {\n    nodes {\n      ...ProductThumbImageFragment\n      __typename\n    }\n    __typename\n  }\n  organization {\n    ...ProductInfoOrganization\n    __typename\n  }\n  productVariants(first: 20, input: {route: \"market\"}) {\n    nodes {\n      ...ProductInfoProductVariant\n      ...ProductVariantButtonFragment\n      __typename\n    }\n    __typename\n  }\n  recommendedProducts(first: 7, input: {fillInWithRecentlyUpdated: true}) {\n    nodes {\n      ...ProductInfoThumbnailAtomFragment\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment ProductInfoProductVariant on ProductVariant {\n  id\n  name\n  price\n  marketShippingRequired\n  __typename\n}\n\nfragment ProductInfoOrganization on Organization {\n  id\n  name\n  __typename\n}\n\nfragment ProductThumbImageFragment on Image {\n  id\n  databaseID\n  url\n  height\n  width\n  __typename\n}\n\nfragment ProductInfoThumbnailAtomFragment on ProductInfo {\n  id\n  databaseID\n  name\n  organization {\n    id\n    name\n    __typename\n  }\n  images(first: 1) {\n    nodes {\n      id\n      url\n      height\n      width\n      __typename\n    }\n    __typename\n  }\n  downloadContent @include(if: $showProductInfoThumbnailDLC) {\n    id\n    fileName\n    downloadURL\n    __typename\n  }\n  __typename\n}\n\nfragment ProductVariantButtonFragment on ProductVariant {\n  id\n  name\n  price\n  __typename\n}\n\nfragment ProductPurchaseCompleteFragment on ProductInfo {\n  id\n  databaseID\n  name\n  organization {\n    id\n    name\n    __typename\n  }\n  __typename\n}\n",
	}
}
