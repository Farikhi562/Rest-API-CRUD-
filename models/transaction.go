package models
import (
	"time"
	"gorm.io/gorm"
)

const (
	TipePemasukan   = "pemasukan"
	TipePengeluaran = "pengeluaran"
)

type Transaksi struct {
	gorm.Model
	Deskripsi string    `json:"deskripsi"`
	Jumlah    int64     `json:"jumlah"` 
	Tipe      string    `json:"tipe"`   
	Kategori  string    `json:"kategori"`
	Tanggal   time.Time `json:"tanggal"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}