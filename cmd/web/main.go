package main
import (
    "github.com/labstack/echo/v4"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
    addr = "127.0.0.1:4000"
    resumePath = "/sendResume"
    vacancyPath = "/putVacancy"
    dashboardPath = "/dashboard"
)

var (
    db *sql.DB
)

func main() {
    e := echo.New()

    connStr := "host=imperium.org.ru user=constellation password=aeCha2thaM4chaej dbname=cst_backend sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        e.Logger.Fatal("Couldn't to open PostgreSQL database:", err)
    }
    defer db.Close()

    e.POST(resumePath, resumeHandler)
    e.POST(vacancyPath, vacancyHandler)
    e.POST(dashboardPath, dashboardHandler)
    e.Logger.Fatal(e.Start(addr))
}
