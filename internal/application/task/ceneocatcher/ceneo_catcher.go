package ceneocatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"main/internal/application"
	"main/internal/application/email"
	"main/internal/domain/contracts"
	"main/internal/domain/models"

	"github.com/PuerkitoBio/goquery"
)

type CeneoCatcherTaskRunner struct {
	emailComposer   email.Composer
	mailer          application.IMailer
	tasksRepo       contracts.ITasks
	httpClient      *http.Client
	ceneoDomain     string
	ceneoProductTag string
	senderEmail     string
	senderSignature string
}

func NewTaskRunner(
	emailComposer email.Composer,
	mailer application.IMailer,
	tasksRepo contracts.ITasks,
	senderEmail string,
	senderSignature string,
) *CeneoCatcherTaskRunner {
	return &CeneoCatcherTaskRunner{
		emailComposer:   emailComposer,
		mailer:          mailer,
		tasksRepo:       tasksRepo,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
		ceneoDomain:     "https://ceneo.pl",
		ceneoProductTag: ".product-offer__container",
		senderEmail:     senderEmail,
		senderSignature: senderSignature,
	}
}

func (t *CeneoCatcherTaskRunner) Run(ctx context.Context, taskID int, config []byte) error {
	_, err := t.tasksRepo.Update(ctx, taskID, map[string]any{
		"last_run_at": time.Now(),
	})
	if err != nil {
		return fmt.Errorf("could not update task: %w", err)
	}

	var cfg models.CeneoCatcherConfig

	if err := json.Unmarshal(config, &cfg); err != nil {
		return fmt.Errorf("invalid ceneo catcher config: %w", err)
	}

	product, err := t.getLowestPriceProduct(ctx, t.ceneoDomain, cfg.ProductID, cfg.MaxPrice)
	if err != nil {
		return fmt.Errorf("could get lowest price for task with id %d: %w", taskID, err)
	}

	if product == nil {
		return nil
	}

	data := email.CeneoCatcher{
		ProductName:     product.Name,
		ProductPrice:    product.Price,
		ProductCompany:  product.Company,
		ProductURL:      product.URL,
		RecipientEmail:  cfg.RecipientEmail,
		SenderEmail:     t.senderEmail,
		SenderSignature: t.senderSignature,
	}

	msg, err := t.emailComposer.ComposeForCeneoCatcher(data)
	if err != nil {
		return fmt.Errorf("could not compose msg for ceneo catcher task: %w", err)
	}

	err = t.mailer.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send msg - %v: %w", msg, err)
	}

	newConfig, err := buildConfigWithReducedMaxPrice(cfg)
	if err != nil {
		return fmt.Errorf("could build config with reduced max price: %w", err)
	}

	_, err = t.tasksRepo.Update(ctx, taskID, map[string]any{
		"config": newConfig,
	})
	if err != nil {
		return fmt.Errorf("could not update task: %w", err)
	}

	return nil
}

func buildConfigWithReducedMaxPrice(cfg models.CeneoCatcherConfig) ([]byte, error) {
	cfg.MaxPrice *= 0.95

	result, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not update task: %w", err)
	}

	return result, nil
}

type Product struct {
	Name    string
	Price   string
	Company string
	URL     string
}

func (t *CeneoCatcherTaskRunner) getLowestPriceProduct(
	ctx context.Context, domain string, productID int, maxPrice float32,
) (*Product, error) {
	targetURL := fmt.Sprintf("%s/%d", domain, productID)

	doc, err := t.fetchProductsDoc(ctx, targetURL)
	if err != nil {
		return nil, err
	}

	baseURL, _ := url.Parse(domain)

	products := t.findProducts(doc, baseURL)

	return findCheapestProduct(products, maxPrice)
}

func (t *CeneoCatcherTaskRunner) fetchProductsDoc(
	ctx context.Context, targetURL string,
) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not build request for %s: %w", targetURL, err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request to %s: %w", targetURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, targetURL)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	return doc, nil
}

func findCheapestProduct(products []Product, maxPrice float32) (*Product, error) {
	if len(products) == 0 {
		slog.Warn("ceneo_catcher - no products found, something might have changed on the ceneo.pl")

		return nil, nil
	}

	lowestPrice := math.Inf(1)

	var cheapest Product

	for _, product := range products {
		price, err := parsePrice(product.Price)
		if err != nil {
			return nil, fmt.Errorf("could not parse price %q: %w", product.Price, err)
		}

		if price < lowestPrice {
			lowestPrice = price
			cheapest = product
		}
	}

	if lowestPrice > float64(maxPrice) {
		slog.Info("ceneo_catcher - prices are too high",
			"product", cheapest.Name, "max_price", maxPrice, "lowest", lowestPrice,
		)

		return nil, nil
	}

	slog.Info("ceneo_catcher - found desired price",
		"product", cheapest.Name, "max_price", maxPrice, "price", lowestPrice,
	)

	return &cheapest, nil
}

func parsePrice(raw string) (float64, error) {
	normalized := strings.ReplaceAll(strings.ReplaceAll(raw, " ", ""), ",", ".")

	price, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse float: %w", err)
	}

	return price, nil
}

func (t *CeneoCatcherTaskRunner) findProducts(doc *goquery.Document, baseURL *url.URL) []Product {
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

	return products
}
