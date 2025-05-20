package models

type Dokter struct {
	ID         string `json:"id"`
	Nama       string `json:"nama"`
	Poli       string `json:"poli"`
	FotoDokter string `json:"foto_dokter"`
}
