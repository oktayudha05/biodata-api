package models

type MahasiswaPub struct {
	Nama string `json:"nama" validate:"required"`
	Npm string `json:"npm" validate:"required"`
	Angkatan int16 `json:"angkatan" validate:"required"`
	Alamat string `json:"alamat" validate:"required"`
	NoHp string `json:"nomor_hp" validate:"required"`
	NoHpKeluarga string `json:"nomor_keluarga" validate:"required"`
	AsalSekolah string `json:"asal_sekolah"`
}

type MahasiswaPriv struct {
	MahasiswaPub
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}