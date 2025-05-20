package models

type Eresep struct {
	IDEresep      string `json:"id_eresep" gorm:"primaryKey;column:ID_ERESEP"`
	IDPendaftaran string `json:"id_pendaftaran" gorm:"column:ID_PENDAFTARAN"`
	Status        string `json:"status" gorm:"column:STATUS"`
}
