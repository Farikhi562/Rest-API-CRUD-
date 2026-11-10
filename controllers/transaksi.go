package controllers
import (
	"errors"
	"strconv"
	"time"

	"keuangan_app/database"
	"keuangan_app/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func validasi(t *models.Transaksi) string {
	if t.Deskripsi == "" {
		return "deskripsi wajib diisi"
	}
	if t.Jumlah <= 0 {
		return "jumlah harus lebih dari 0"
	}
	if t.Tipe != models.TipePemasukan && t.Tipe != models.TipePengeluaran {
		return "tipe harus 'pemasukan' atau 'pengeluaran'"
	}
	return ""
}

func cariByID(c *fiber.Ctx) (t models.Transaksi, ok bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID harus berupa angka"})
		return t, false
	}

	if err := database.DB.First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Transaksi tidak ditemukan"})
		} else {
			_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return t, false
	}
	return t, true
}


// GET /transaksi?tipe=pengeluaran&kategori=Makanan 
func GetTransaksi(c *fiber.Ctx) error {
	var list []models.Transaksi
	query := database.DB.Order("tanggal DESC")

	if tipe := c.Query("tipe"); tipe != "" {
		query = query.Where("tipe = ?", tipe)
	}
	if kategori := c.Query("kategori"); kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}

	if err := query.Find(&list).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// GET /transaksi/:id
func GetTransaksiByID(c *fiber.Ctx) error {
	t, ok := cariByID(c)
	if !ok {
		return nil
	}
	return c.JSON(t)
}

// POST /transaksi
func CreateTransaksi(c *fiber.Ctx) error {
	var t models.Transaksi
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "format JSON tidak valid"})
	}
	if msg := validasi(&t); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}
	if t.Tanggal.IsZero() {
		t.Tanggal = time.Now()
	}

	if err := database.DB.Create(&t).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(t)
}

// PUT /transaksi/:id
func UpdateTransaksi(c *fiber.Ctx) error {
	t, ok := cariByID(c)
	if !ok {
		return nil
	}

	var input models.Transaksi
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "format JSON tidak valid"})
	}
	if msg := validasi(&input); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}


	t.Deskripsi = input.Deskripsi
	t.Jumlah = input.Jumlah
	t.Tipe = input.Tipe
	t.Kategori = input.Kategori
	if !input.Tanggal.IsZero() {
		t.Tanggal = input.Tanggal
	}

	if err := database.DB.Save(&t).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(t)
}

// DELETE /transaksi/:id
func DeleteTransaksi(c *fiber.Ctx) error {
	t, ok := cariByID(c)
	if !ok {
		return nil
	}
	if err := database.DB.Delete(&t).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}


// GET /ringkasan
func GetRingkasan(c *fiber.Ctx) error {
	var pemasukan, pengeluaran int64

	hitung := func(tipe string, hasil *int64) error {
		return database.DB.Model(&models.Transaksi{}).
			Where("tipe = ?", tipe).
			Select("COALESCE(SUM(jumlah), 0)").
			Scan(hasil).Error
	}

	if err := hitung(models.TipePemasukan, &pemasukan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := hitung(models.TipePengeluaran, &pengeluaran); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"total_pemasukan":   pemasukan,
		"total_pengeluaran": pengeluaran,
		"saldo":             pemasukan - pengeluaran,
	})
}