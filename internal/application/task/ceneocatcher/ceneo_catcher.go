package ceneocatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"main/internal/application/email"
	"main/internal/domain/contracts"
	"main/internal/domain/models"

	"github.com/PuerkitoBio/goquery"
)

type CeneoCatcherTaskRunner struct {
	emailComposer   email.Composer
	mailer          email.IMailer
	tasksRepo       contracts.ITasks
	ceneoDomain     string
	ceneoProductTag string
}

func NewTaskRunner(
	emailComposer email.Composer,
	mailer email.IMailer,
	tasksRepo contracts.ITasks,
) *CeneoCatcherTaskRunner {
	return &CeneoCatcherTaskRunner{
		emailComposer:   emailComposer,
		mailer:          mailer,
		tasksRepo:       tasksRepo,
		ceneoDomain:     "https://ceneo.pl",
		ceneoProductTag: ".product-offer__container",
	}
}

func (t *CeneoCatcherTaskRunner) Run(ctx context.Context, taskID int, config []byte) error {
	var cfg models.CeneoCatcherConfig

	if err := json.Unmarshal(config, &cfg); err != nil {
		return fmt.Errorf("invalid ceneo catcher config: %w", err)
	}

	url := fmt.Sprintf("%s/%d", t.ceneoDomain, cfg.ProductID)

	product, err := t.getLowestPriceProduct(url, cfg.MaxPrice)
	if err != nil {
		return fmt.Errorf("could get lowest price for task with id %d: %w", taskID, err)
	}

	if product == nil {
		return nil
	}

	data := email.CeneoCatcher{
		ProductName:    product.Name,
		ProductPrice:   product.Price,
		ProductCompany: product.Company,
		ProductURL:     product.URL,
	}

	msg, err := t.emailComposer.ComposeForCeneoCatcher(data)
	if err != nil {
		return fmt.Errorf("could not compose msg for ceneo catcher task: %w", err)
	}

	err = t.mailer.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send msg - %v: %w", msg, err)
	}

	_, err = t.tasksRepo.Update(ctx, taskID, map[string]any{"last_run_at": time.Now()})
	if err != nil {
		return fmt.Errorf("could not update task: %w", err)
	}

	return nil
}

type Product struct {
	Name    string
	Price   string
	Company string
	URL     string
}

func (t *CeneoCatcherTaskRunner) getLowestPriceProduct(domain string, maxPrice int) (*Product, error) {
	resp, err := http.Get(domain)
	if err != nil {
		return nil, fmt.Errorf("could not get domain %s: %w", domain, err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could read response body: %w", err)
	}

	baseURL, _ := url.Parse(domain)

	products := make([]Product, 0)

	doc.Find(t.ceneoProductTag).Each(func(i int, s *goquery.Selection) {
		name := s.Find(".short-name__txt").Text()
		price := s.Find(".price").Text()

		product := Product{
			Name:  name,
			Price: price,
		}

		company, ok := s.Find("img").Attr("alt")
		if ok {
			product.Company = company
		}

		URL, ok := s.Find("button.add-to-basket-no-popup").Attr("data-basket-click-url")
		if !ok {
			URL, ok = s.Attr("data-click-url")
			if ok {
				relativeURL, _ := url.Parse(URL)
				fullURL := baseURL.ResolveReference(relativeURL)
				product.URL = fullURL.String()
			}
		} else {
			relativeURL, _ := url.Parse(URL)
			fullURL := baseURL.ResolveReference(relativeURL)
			product.URL = fullURL.String()
		}

		products = append(products, product)
	})

	var lowestPrice float64

	var productWithLowestPrice Product

	for _, product := range products {
		price, err := strconv.ParseFloat(
			strings.ReplaceAll(
				strings.ReplaceAll(product.Price, " ", ""), ",", "."), 64,
		)
		if err != nil {
			return nil, fmt.Errorf("could parse string to float: %w", err)
		}

		if lowestPrice == 0 {
			lowestPrice = price
			productWithLowestPrice = product

			continue
		}

		if price < lowestPrice {
			lowestPrice = price
			productWithLowestPrice = product
		}
	}

	if lowestPrice <= float64(maxPrice) {
		slog.Info("ceneo_catcher - found desired price", "price", lowestPrice)

		return &productWithLowestPrice, nil
	}

	slog.Info("ceneo_catcher - prices are too high", "lowest", lowestPrice)

	return nil, nil
}
