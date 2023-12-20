package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"go-web-scrapper/entity"
	"go-web-scrapper/entity/model"
	"go-web-scrapper/framework/logger"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

const URL = "https://www.tokopedia.com/discovery/productlist-clp_handphone-tablet_65-120"
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
	numGoroutines := 1
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

	var productNodes []*cdp.Node
	scrollingScript := `
        //    const scrolls = 10
        //    let scrollCount = 0
           
        //    const scrollInterval = setInterval(() => {
        //      window.scrollTo(0, document.body.scrollHeight)
        //      scrollCount++
           
        //      if (scrollCount === numScrolls) {
        //       clearInterval(scrollInterval)
        //      }
        //    }, 500)

		async function scrollAndClick() {
			const productNodes = document.querySelectorAll('.pcv3__info-content');
			for (const node of productNodes) {
				node.click();
				await new Promise(resolve => setTimeout(resolve, 1000));
			}
		}
	
		scrollAndClick();
        `
	err := chromedp.Run(ctx,
		chromedp.Navigate(URL),
		// chromedp.Navigate("https://scrapingclub.com/exercise/list_infinite_scroll/"),
		chromedp.Evaluate(scrollingScript, nil),
		chromedp.Sleep(2*time.Minute),
		chromedp.Nodes(".css-6bc98m", &productNodes, chromedp.ByQueryAll))
	// chromedp.Nodes(".post", &productNodes, chromedp.ByQueryAll))
	if err != nil {
		u.Logger.Logger.Error(err)
	}

	if err != nil {
		log.Fatal("Error while performing the automation logic:", err)
	}

	fmt.Println("total_data: ", len(productNodes))

	for i := 0; i < len(productNodes); i += numGoroutines {
		for j := 0; j < numGoroutines; j++ {
			wg.Add(1)
			go func(ctx context.Context, a int, b int) {
				defer wg.Done()
				var name, description, imageLink, price, ratings, merchant string
				err = chromedp.Run(ctx,
					// chromedp.Text("h4", &name, chromedp.ByQuery, chromedp.FromNode(productNodes[a+b])),
					chromedp.Text(".prd_link-product-name", &name, chromedp.ByQuery, chromedp.FromNode(productNodes[a+b])),
					// chromedp.Text(".prd_link-product-name", &description, chromedp.ByQuery, chromedp.FromNode(node)),
					// chromedp.Text(".prd_label-product-price", &price, chromedp.ByQuery, chromedp.FromNode(node)),
					// chromedp.Text(".prd_rating-average-text", &ratings, chromedp.ByQuery, chromedp.FromNode(node)),
					// chromedp.Text(".prd_link-product-name", &merchant, chromedp.ByQuery, chromedp.FromNode(node)),
				)

				if err != nil {
					log.Fatal("Error:", err)
				}

				d := &model.Data{
					Id:          uuid.New().String(),
					Name:        name,
					Description: description,
					ImageLink:   imageLink,
					Price:       price,
					Ratings:     ratings,
					MerchatName: merchant,
				}
				data = append(data, d)
			}(ctx, i, j)
		}
	}
	wg.Wait()
	// if err := u.Entity.Insert(ctx, data); err != nil {
	// 	u.Logger.Logger.Error(err)
	// }

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
