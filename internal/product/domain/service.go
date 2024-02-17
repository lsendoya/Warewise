package domain

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type UseCase interface {
	Add(formData *FormDataProduct) (*Product, error)
	Update(id uuid.UUID, product *Product) (*Product, error)
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*Product, error)
	List() (Products, error)
	UploadImageFile(ImageFiles []*multipart.FileHeader) ([]string, error)
}

type Storage interface {
	Add(product *Product) (*Product, error)
	Update(id uuid.UUID, product *Product) (*Product, error)
	Delete(id uuid.UUID) error
	Get(id uuid.UUID) (*Product, error)
	List() (Products, error)
}
