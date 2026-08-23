package ceneocatcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"main/internal/application/email"
	"main/internal/domain/contracts"
	"main/internal/domain/errs/api"
	"main/internal/domain/models"

	"github.com/PuerkitoBio/goquery"
)

type CeneoCatcher struct {
	emailComposer   email.Composer
	mailer          email.IMailer
	tasksRepo       contracts.ITasks
	interval        time.Duration
	maxLimitPrice   float64
	ceneoProductID  int
	ceneoDomain     string
	ceneoProductTag string
}

func New(
	emailComposer email.Composer,
	mailer email.IMailer,
	tasksRepo contracts.ITasks,
	maxLimitPrice float64,
	ceneoProductID int,
	interval time.Duration,
) *CeneoCatcher {
	return &CeneoCatcher{
		emailComposer:   emailComposer,
		mailer:          mailer,
		tasksRepo:       tasksRepo,
		maxLimitPrice:   maxLimitPrice,
		interval:        interval,
		ceneoDomain:     "https://ceneo.pl",
		ceneoProductID:  ceneoProductID,
		ceneoProductTag: ".product-offer__container",
	}
}

func (t *CeneoCatcher) Name() string {
	return "ceneo_catcher"
}

func (t *CeneoCatcher) Run(ctx context.Context) error {
	task, err := t.getOrInsertTask(ctx)
	if err != nil {
		return fmt.Errorf("could get or insert task: %s: %w", t.Name(), err)
	}

	if !task.Active {
		return nil
	}

	url := fmt.Sprintf("%s/%d", t.ceneoDomain, t.ceneoProductID)

	product, err := t.getLowestPriceProduct(url)
	if err != nil {
		return fmt.Errorf("could get suject and content for: %s: %w", t.Name(), err)
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
		return fmt.Errorf("could not request task %v: %w", t.Name(), err)
	}

	err = t.mailer.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send msg %v: %w", t.Name(), err)
	}

	_, err = t.tasksRepo.Update(ctx, task.ID, map[string]any{"last_run_at": time.Now()})
	if err != nil {
		return fmt.Errorf("could not update task %v: %w", t.Name(), err)
	}

	return nil
}

func (t *CeneoCatcher) getOrInsertTask(ctx context.Context) (models.Task, error) {
	task, err := t.tasksRepo.GetByName(ctx, t.Name())
	if err != nil {
		if errors.Is(err, api.ErrTaskNotFound) {
			task, err = t.tasksRepo.Insert(ctx, t.Name())
			if err != nil {
				return models.Task{}, fmt.Errorf("could insert task with name %s: %w", t.Name(), err)
			}
		} else {
			return models.Task{}, fmt.Errorf("could not get task by %s: %w", t.Name(), err)
		}
	}

	return task, nil
}

type Product struct {
	Name    string
	Price   string
	Company string
	URL     string
}

func (t *CeneoCatcher) getLowestPriceProduct(domain string) (*Product, error) {
	resp, err := http.Get(domain)
	if err != nil {
		return nil, fmt.Errorf("could insert task with name %s: %w", t.Name(), err)
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

	if lowestPrice <= t.maxLimitPrice {
		slog.Info("ceneo_catcher - found desired price", "price", lowestPrice)

		return &productWithLowestPrice, nil
	}

	slog.Info("ceneo_catcher - prices are too high", "lowest", lowestPrice)

	return nil, nil
}
