package models

import "time"

type DetailEresep struct {
    IDDetail      string    `json:"id_detail"`
    IDEresep      string    `json:"id_eresep"`
    TanggalEresep time.Time `json:"tanggal_eresep"`
    Kuantitas    int16      `json:"kuantitas"`
    AturanPakai  string     `json:"aturan_pakai"`
    Catatan      string     `json:"catatan"`
}
