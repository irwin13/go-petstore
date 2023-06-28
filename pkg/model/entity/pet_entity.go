package entity

type GetPet struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type InsertPet struct {
	ID          *string `json:"id"`
	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type UpdatePet struct {
	ID          *string `json:"id" binding:"required"`
	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type DeletePet struct {
	ID *string `form:"id" binding:"required"`
}
