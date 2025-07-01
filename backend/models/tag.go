package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID   string `gorm:"primaryKey;type:uuid" json:"id"`
	Name string `gorm: "type:varchar(100);unique;not null" json:"name"`
}

type TagCount struct {
	Name  string `json:"tag"`
	Count int    `json:"count"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}

const (
	AssedioMoral          = "Assédio moral"
	AssedioSexual         = "Assédio sexual"
	AssedioOrganizacional = "Assédio organizacional"
	ConflitoInteresses    = "Conflito de interesses"
	Irregularidades       = "Irregularidades e delitos"
	DesrespeitoNormas     = "Desrespeito às normas de segurança"
	CondutasAntieticas    = "Condutas antiéticas"
	AtosDiscriminatorios  = "Atos discriminatórios"
	Outros                = "Outros"
	OcorrenciaInveridica  = "Ocorrência inverídica"
)
