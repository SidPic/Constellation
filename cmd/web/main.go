package main
import (
    "github.com/labstack/echo/v4"
)

const (
    addr = "127.0.0.1:4000"
    resumePath = "/sendResume"
    vacancyPath = "/putVacancy"
    dashboardPath = "/dashboard"
)

func main() {
    e := echo.New()
    e.POST(resumePath, resumeHandler)
    e.POST(vacancyPath, vacancyHandler)
    e.POST(dashboardPath, dashboardHandler)
    e.Logger.Fatal(e.Start(addr))
}
