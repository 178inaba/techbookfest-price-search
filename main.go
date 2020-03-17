package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"
	"sync"

	"github.com/178inaba/techbookfest-price-search/techbookfest"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	c, err := techbookfest.NewTechBookFest(ctx)
	if err != nil {
		return err
	}

	md, err := c.GetMarketDashboard(ctx, 2000)
	if err != nil {
		return err
	}

	fmt.Printf("All Books: %d\n", len(md.Data.AllProductVariants.Nodes))

	g, ctx := errgroup.WithContext(ctx)
	limit := make(chan struct{}, 50)
	var m sync.Mutex
	ddMap := map[string]displayDetail{}
	for _, node := range md.Data.AllProductVariants.Nodes {
		node := node
		g.Go(func() error {
			limit <- struct{}{}
			defer func() { <-limit }()

			pi, err := c.GetProductInfo(ctx, node.Products.Nodes[0].ID)
			if err != nil {
				return err
			}

			for _, node := range pi.Data.Product.ProductVariants.Nodes {
				if node.Price == 0 {
					u, err := url.Parse("https://techbookfest.org/product/" + pi.Data.Product.DatabaseID)
					if err != nil {
						return err
					}

					m.Lock()
					defer m.Unlock()
					ddMap[pi.Data.Product.DatabaseID] = displayDetail{
						name:                     pi.Data.Product.Name,
						description:              pi.Data.Product.Description,
						organization:             pi.Data.Product.Organization.Name,
						page:                     pi.Data.Product.Page,
						firstAppearanceEventName: pi.Data.Product.FirstAppearanceEventName,
						url:                      u,
						price:                    node.Price,
						physical:                 node.MarketShippingRequired,
					}

					break
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	dds := make([]displayDetail, 0, len(ddMap))
	for _, dd := range ddMap {
		dds = append(dds, dd)
	}

	sort.Slice(dds, func(i, j int) bool {
		if dds[i].organization != dds[j].organization {
			return dds[i].organization < dds[j].organization
		}

		if dds[i].firstAppearanceEventName != dds[j].firstAppearanceEventName {
			return dds[i].firstAppearanceEventName < dds[j].firstAppearanceEventName
		}

		return dds[i].name < dds[j].name
	})

	fmt.Println("| 書名 | サークル名 | 初出イベント | ページ数 |")
	fmt.Println("| -- | -- | -- | -- |")
	for _, dd := range dds {
		fmt.Printf("| [%s](%s) | %s | %s | %d |\n", dd.name, dd.url.String(), dd.organization, dd.firstAppearanceEventName, dd.page)
	}

	return nil
}

type displayDetail struct {
	name                     string
	description              string
	organization             string
	page                     int
	firstAppearanceEventName string
	url                      *url.URL
	price                    int
	physical                 bool
}
