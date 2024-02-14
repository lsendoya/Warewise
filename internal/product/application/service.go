package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lsendoya/Warewise/internal/product/domain"
	"github.com/lsendoya/Warewise/pkg/aws"
		"os"
	"time"
)

type Product struct {
	storage domain.Storage
	awsSvc  aws.Service
}

func New(s domain.Storage, awsSvc aws.Service) Product {
	return Product{
		storage: s,
		awsSvc:  awsSvc,
	}
}

func (p *Product) Add(formData *domain.FormDataProduct) (*domain.Product, error) {
	product := new(domain.Product)
	product.BeforeCreated()

	urls, err := p.UploadImageFile(formData)
	if err != nil {
		return nil, err
	}

	jsonData, errMarshal := json.Marshal(urls)
	if errMarshal != nil {
		fmt.Printf("Error serializing URLs: %v", err)
		return nil, errMarshal
	}

	product.ImageURLS = jsonData
	product.Name = formData.Name
	product.Price = formData.Price
	product.Color = formData.Color
	product.Description = formData.Description
	product.IsAvailable = formData.IsAvailable
	product.Size = formData.Size

	return p.storage.Add(product)
}
func (p *Product) Update(id uuid.UUID, product *domain.Product) (*domain.Product, error) {
	mdl, err := p.storage.Update(id, product)
	if err != nil {
		return nil, err
	}
	mdl.UpdatedAt = time.Now()

	return mdl, nil
}
func (p *Product) Delete(id uuid.UUID) error {
	return p.storage.Delete(id)
}
func (p *Product) Get(id uuid.UUID) (*domain.Product, error) {
	return p.storage.Get(id)
}
func (p *Product) List() (domain.Products, error) {
	return p.storage.List()
}

func (p *Product) UploadImageFile(formData *domain.FormDataProduct) ([]string, error) {
	var urls []string
	for _, fileHeader := range formData.ImageFiles {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, errors.New("error opening the uploaded file")
		}

				fileName := fileHeader.Filename

		url, errUpload := p.awsSvc.UploadFile(os.Getenv("BUCKET_NAME"), fileName, file)
		if errUpload != nil {
			return nil, errUpload
		}

		errFileClose := file.Close()
		if err != nil {
		return nil, errFileClose
		}

		urls = append(urls, url)

	}

	return urls, nil
}