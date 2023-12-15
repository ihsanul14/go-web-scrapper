package usecase

import (
	"context"
	"encoding/csv"
	"go-web-scrapper/entity"
	"go-web-scrapper/entity/model"
	"go-web-scrapper/framework/logger"
	"os"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

const URL = "https://www.tokopedia.com/p/handphone-tablet"
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"

type Usecase struct {
	Entity entity.IEntity
	Logger *logger.Log
}
type IUsecase interface {
	Get()
	WriteCsv(data []*model.Data)
}

func NewUsecase(entity entity.IEntity, logger *logger.Log) IUsecase {
	return &Usecase{Entity: entity, Logger: logger}
}

func (u *Usecase) Get() {
	var data []*model.Data
	var wg sync.WaitGroup
	numGoroutines := 100
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(os.Getenv("CHROME_PATH")),
		chromedp.UserAgent(userAgent),
		chromedp.WindowSize(1920, 1080),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU)
	actx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(actx)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.Navigate(URL))
	if err != nil {
		u.Logger.Logger.Error(err)
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(ctx context.Context, id int) {
			defer wg.Done()
			d := &model.Data{
				Id:          uuid.New().String(),
				Name:        "",
				Description: "",
				ImageLink:   "",
				Price:       "",
				Ratings:     "",
				MerchatName: "",
			}

			err = chromedp.Run(ctx, chromedp.WaitVisible(`document.querySelector(".prd_link-product-name")`))
			if err != nil {
				u.Logger.Logger.Error(err)
			}

			err = chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector(".prd_link-product-name").innerHTML`, &d.Name))
			if err != nil {
				u.Logger.Logger.Errorf("name : %v", err)
			}

			err = chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector(".css-1ifrycw").innerHTML`, &d.Description))
			if err != nil {
				u.Logger.Logger.Errorf("description : %v", err)
			}

			err = chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector(".prd_label-product-slash-price").innerHTML`, &d.Price))
			if err != nil {
				u.Logger.Logger.Errorf("price : %v", err)
			}

			err = chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector(".prd_rating-average-text").innerHTML`, &d.Ratings))
			if err != nil {
				u.Logger.Logger.Errorf("ratings : %v", err)
			}

			err = chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector(".css-1wdzqxj-unf-heading.e1qvo2ff2").innerHTML`, &d.MerchatName))
			if err != nil {
				u.Logger.Logger.Errorf("merchant : %v", err)
			}

			data = append(data, d)
		}(ctx, i)
	}
	wg.Wait()
	if err := u.Entity.Insert(ctx, data); err != nil {
		u.Logger.Logger.Error(err)
	}

	u.WriteCsv(data)
}

func (u *Usecase) WriteCsv(data []*model.Data) {
	file, err := os.Create("./framework/csv/output.csv")
	if err != nil {
		u.Logger.Logger.Error(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Id", "Name", "Description", "ImageLink", "Price", "Ratings", "MerchantName"}
	if err := writer.Write(header); err != nil {
		u.Logger.Logger.Error(err)
	}

	for _, v := range data {
		row := []string{v.Id, v.Name, v.Description, v.ImageLink, v.Price, v.Ratings, v.MerchatName}
		if err := writer.Write(row); err != nil {
			u.Logger.Logger.Error(err)
		}
	}
}