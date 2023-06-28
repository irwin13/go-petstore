package dao

import (
	"github.com/irwin13/go-petstore/pkg/model/entity"
	uuid "github.com/satori/go.uuid"
)

func GenerateUuidV4() string {
	return uuid.NewV4().String()
}

type PetDao interface {
	Search(filter string) ([]entity.GetPet, error)
	Insert(request entity.InsertPet) (string, error)
	Update(request entity.UpdatePet) (int64, error)
	Delete(request entity.DeletePet) (int64, error)
}
