package models

import "time"

type DetailEresep struct {
    IDDetail      string    `json:"id_detail"`
    IDEresep      string    `json:"id_eresep"`
    TanggalEresep time.Time `json:"tanggal_eresep"`
    Catatan       string    `json:"catatan"`
}
