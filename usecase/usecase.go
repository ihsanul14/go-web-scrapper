package usecase

import (
	"encoding/csv"
	"fmt"
	"go-web-scrapper/entity"
	"go-web-scrapper/entity/model"
	"go-web-scrapper/framework/logger"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const URL = "https://pemilu2024.kpu.go.id/pilpres/hitung-suara/11/1105/110507/1105072002/1105072002001"
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
	// var wg sync.WaitGroup
	var data []*model.Election

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
	defer driver.Close()

	err = driver.Get(URL)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	err = driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		lastProduct, _ := driver.FindElement(selenium.ByTagName, "tbody")
		if lastProduct != nil {
			return lastProduct.IsDisplayed()
		}
		return false, nil
	}, 30*time.Second)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	productElements, err := driver.FindElements(selenium.ByTagName, "tr")
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	productElements = []selenium.WebElement{
		productElements[1],
		productElements[2],
		productElements[3],
		productElements[6],
		productElements[7],
		productElements[8],
		productElements[10],
		productElements[11],
		productElements[12],
	}

	defer service.Stop()
	fmt.Printf("total data %d \n", len(productElements))
	// wg.Add(1)
	// go func(i int) {
	// 	defer wg.Done()

	t01, err := extractElementElection(productElements[3])
	if err != nil {
		return
	}

	t02, err := extractElementElection(productElements[4])
	if err != nil {
		return
	}

	t03, err := extractElementElection(productElements[5])
	if err != nil {
		return
	}

	totalSah, err := extractElementElection(productElements[6])
	if err != nil {
		return
	}

	totalTidakSah, err := extractElementElection(productElements[7])
	if err != nil {
		return
	}

	total, err := extractElementElection(productElements[8])
	if err != nil {
		return
	}

	d := &model.Election{
		Id:            getTpsID(URL),
		T01:           *t01,
		T02:           *t02,
		T03:           *t03,
		TotalSah:      *totalSah,
		TotalTidakSah: *totalTidakSah,
		Total:         *total,
	}
	data = append(data, d)
	// }(i)
	// time.Sleep(2 * time.Minute)
	fmt.Printf("%+v", data[0])
	// wg.Wait()
	// if err := u.Entity.Insert(data); err != nil {
	// 	u.Logger.Logger.Error(err)
	// }
	// u.WriteCsv(data)
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

func extractElementElection(v selenium.WebElement) (*string, error) {
	res, err := v.Text()
	if err != nil {
		return nil, err
	}
	resArr := strings.Split(res, " ")
	return &resArr[len(resArr)-1], nil
}

func getTpsID(url string) string {
	urlArr := strings.Split(url, "/")
	return urlArr[len(urlArr)-1]
}
