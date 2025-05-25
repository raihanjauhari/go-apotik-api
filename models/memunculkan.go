package models

type Memunculkan struct {
	KodeObat    string `json:"kode_obat"`    // KODE_OBAT
	IDEresep    string `json:"id_eresep"`    // ID_ERESEP
	IDDetail    string `json:"id_detail"`    // ID_DETAIL
	Kuantitas   int    `json:"kuantitas"`    // KUANTITAS
	AturanPakai string `json:"aturan_pakai"` // ATURAN_PAKAI
}
