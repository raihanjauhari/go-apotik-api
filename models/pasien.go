package models

type Pasien struct {
	IDPendaftaran string  `json:"id_pendaftaran" db:"ID_PENDAFTARAN"`
	NamaPasien    string  `json:"nama_pasien" db:"NAMA_PASIEN"`
	Umur          int     `json:"umur" db:"UMUR"`
	Diagnosa      string  `json:"diagnosa" db:"DIAGNOSA"`
	BeratBadan    float64 `json:"berat_badan" db:"BERAT_BADAN"`
	FotoPasien    string  `json:"foto_pasien" db:"FOTO_PASIEN"`
}
