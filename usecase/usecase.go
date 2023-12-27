package usecase

import (
	"encoding/csv"
	"fmt"
	"go-web-scrapper/entity"
	"go-web-scrapper/entity/model"
	"go-web-scrapper/framework/logger"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const URL = "https://www.tokopedia.com/discovery/productlist-clp_handphone-tablet_65-120"
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"

type Usecase struct {
	Entity entity.IEntity
	Logger *logger.Log
}

type Product struct {
	link string
}
type IUsecase interface {
	Get()
	WriteCsv(data []*model.Data)
}

func NewUsecase(entity entity.IEntity, logger *logger.Log) IUsecase {
	return &Usecase{Entity: entity, Logger: logger}
}

func (u *Usecase) Get() {
	var wg sync.WaitGroup
	var data []*model.Data

	var products []Product
	service, err := selenium.NewChromeDriverService(os.Getenv("CHROME_PATH"), 4444)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		fmt.Sprintf("--user-agent=%s", userAgent),
	}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	err = driver.Get(URL)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	scrollingScript := `
	// scroll down the page 10 times
	const scrolls = 10
	let scrollCount = 0
	
	// scroll down and then wait for 5s
	const scrollInterval = setInterval(() => {
	window.scrollTo(0, document.body.scrollHeight)
	scrollCount++
	if (scrollCount === scrolls) {
	clearInterval(scrollInterval)
	}
	}, 5000)
	`
	_, err = driver.ExecuteScript(scrollingScript, []interface{}{})
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	err = driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		lastProduct, _ := driver.FindElement(selenium.ByCSSSelector, ".css-10kdh43:nth-child(10)")
		if lastProduct != nil {
			return lastProduct.IsDisplayed()
		}
		return false, nil
	}, 30*time.Second)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	productElements, err := driver.FindElements(selenium.ByCSSSelector, ".pcv3__info-content")
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	for _, productElement := range productElements {
		linkElementDetail, err := productElement.GetAttribute("href")
		if err != nil {
			log.Fatal("Error:", err)
			return
		}

		product := Product{}
		product.link = linkElementDetail
		products = append(products, product)
	}
	driver.Close()

	defer service.Stop()
	fmt.Printf("total data %d \n", len(productElements))
	for i := 0; i < len(products); i += 5 {
		for j := 0; j < 5; j++ {
			wg.Add(1)
			go func(i int, j int) {
				defer wg.Done()
				fmt.Println(products[i+j].link)

				driver, err := selenium.NewRemote(caps, "")
				if err != nil {
					log.Fatal("Error:", err)
					return
				}

				err = driver.Get(products[i+j].link)
				if err != nil {
					log.Println("Error:", err)
					return
				}

				err = driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
					time.Sleep(10 * time.Second)
					lastProduct, _ := driver.FindElement(selenium.ByCSSSelector, ".css-1os9jjn")
					if lastProduct != nil {
						return lastProduct.IsDisplayed()
					}
					return false, nil
				}, 20*time.Second)
				if err != nil {
					log.Fatal("Error:", err)
					return
				}

				name, err := extractElement(driver, ".css-1os9jjn")
				if err != nil {
					log.Println("Error:", err)
					return
				}

				price, err := extractElement(driver, ".price")
				if err != nil {
					log.Println("Error:", err)
					return
				}

				description, err := extractElement(driver, ".css-16inwn4")
				if err != nil {
					log.Println("Error:", err)
					return
				}

				ratings, err := extractElement(driver, "span.main")
				if err != nil {
					log.Println("Error:", err)
					return
				}

				merchant, err := extractElement(driver, ".e1qvo2ff2")
				if err != nil {
					log.Println("Error:", err)
					return
				}

				imageLinkElement, err := driver.FindElement(selenium.ByCSSSelector, ".css-1c345mg")
				if err != nil {
					log.Println("Error:", err)
					return
				}

				imageLink, err := imageLinkElement.GetAttribute("src")
				if err != nil {
					log.Fatal("Error:", err)
					return
				}

				d := &model.Data{
					Id:          uuid.New().String(),
					Name:        *name,
					Price:       *price,
					Description: *description,
					ImageLink:   imageLink,
					Ratings:     *ratings,
					MerchatName: *merchant,
				}
				data = append(data, d)
				defer driver.Close()
			}(i, j)
		}
		time.Sleep(2 * time.Minute)
	}
	wg.Wait()
	if err := u.Entity.Insert(data); err != nil {
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

func extractElement(v selenium.WebDriver, id string) (*string, error) {
	resElement, err := v.FindElement(selenium.ByCSSSelector, id)
	if err != nil {
		return nil, err
	}
	res, err := resElement.Text()
	if err != nil {
		return nil, err
	}

	return &res, nil
}
