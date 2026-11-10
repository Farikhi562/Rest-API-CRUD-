package routes
import (
	"keuangan_app/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	t := app.Group("/transaksi")
	t.Get("/", controllers.GetTransaksi)
	t.Get("/:id", controllers.GetTransaksiByID)
	t.Post("/", controllers.CreateTransaksi)
	t.Put("/:id", controllers.UpdateTransaksi)
	t.Delete("/:id", controllers.DeleteTransaksi)

	app.Get("/ringkasan", controllers.GetRingkasan)
}