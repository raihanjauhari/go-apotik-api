package models

type Obat struct {
	KodeObat    string  `json:"kode_obat"`
	IDUser      string  `json:"id_user"`
	NamaObat    string  `json:"nama_obat"`
	HargaSatuan float64 `json:"harga_satuan"`
	Stok        int     `json:"stok"`
	Deskripsi   string  `json:"deskripsi"`
}
