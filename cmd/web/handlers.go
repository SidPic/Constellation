package main

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "math"
    "net/http"
    "sort"
    "strconv"

    "github.com/labstack/echo/v4"
)

// Структура для стандартизированного ответа при ошибке
type ErrorResponse struct {
    Message string `json:"message"`
}

// Универсальный обработчик для приёма JSON и записи в БД
func insertJSONBHandler(tableName string) echo.HandlerFunc {
    return func(c echo.Context) (err error) {
        // Читаем тело запроса
        b, err := io.ReadAll(c.Request().Body)
        if err != nil {
            log.Printf("Error reading request body for table %s: %v", tableName, err)
            return c.JSON(http.StatusBadRequest, ErrorResponse{"cannot read request body"})
        }

        // Валидируем JSON
        if !json.Valid(b) {
            return c.JSON(http.StatusBadRequest, ErrorResponse{"invalid json format"})
        }

        // Компактизация (удаление лишних пробелов)
        var dst bytes.Buffer
        if err := json.Compact(&dst, b); err != nil {
            log.Printf("Error compacting json for table %s: %v", tableName, err)
            return c.JSON(http.StatusBadRequest, ErrorResponse{"cannot compact json"})
        }

        // Вставка в БД
        query := "INSERT INTO " + tableName + " (resume) VALUES ($1)" // Вставка в таблицу резюме
        if tableName == "vacancies_jsonb" {
            query = "INSERT INTO " + tableName + " (vacancy) VALUES ($1)" // Вставка в таблицу вакансий
        }

        result, err := db.Exec(query, dst.Bytes())
        if err != nil {
            log.Printf("Couldn't insert into %s: %v", tableName, err)
            return c.JSON(http.StatusInternalServerError, ErrorResponse{"couldn't insert data into database"})
        }

        rowsAffected, _ := result.RowsAffected()
        return c.String(http.StatusOK, strconv.FormatInt(rowsAffected, 10))
    }
}

// Приём резюме, запись в БД
func resumeHandler(c echo.Context) (err error) {
    return insertJSONBHandler("resumes_jsonb")(c)
}

// Приём вакансии, запись в БД
func vacancyHandler(c echo.Context) (err error) {
    return insertJSONBHandler("vacancies_jsonb")(c)
}

// Структура для данных резюме на дашборде
type ResumeDashboard struct {
    Name             string `json:"Name"`
    Surname          string `json:"Surname"`
    Rating           int    `json:"Rating"`
    MatchWithVacancy int    `json:"MatchWithVacancy"`
    Experience       int    `json:"Experience"`
    WorkFormat       []string `json:"WorkFormat"`
    WorkSchedule     []string `json:"WorkSchedule"`
    Employment       []string `json:"Employment"`
    WorkTime         string `json:"WorkTime"`
    Responsibilities map[string]int `json:"Responsibilities"`
    Salary           string `json:"Salary"`
    HardSkills       map[string]int `json:"HardSkills"`
    SoftSkills       []string `json:"SoftSkills"`
    BusinessTrips    bool `json:"BusinessTrips"`
    Education        bool `json:"Education"`
}

// Структура для ответа дашборда
type DashboardResponse struct {
    Resumes []ResumeDashboard `json:"resumes"`
}

// Отправка данных на дашборд
func dashboardHandler(c echo.Context) (err error) {
    vacancyIDStr := c.QueryParam("vacancy_id")
    topNStr := c.QueryParam("top_n")

    // Парсинг vacancy_id
    vacancyID, err := strconv.Atoi(vacancyIDStr)
    if err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{"invalid vacancy_id: must be an integer"})
    }

    // Парсинг top_n
    topN, err := strconv.Atoi(topNStr)
    if err != nil || topN <= 0 {
        return c.JSON(http.StatusBadRequest, ErrorResponse{"invalid top_n: must be a positive integer"})
    }

    // Получение вакансии
    vacancy, err := getVacancyFromDB(vacancyID)
    if err != nil {
        log.Println("Error getting vacancy:", err)
        return c.JSON(http.StatusInternalServerError, ErrorResponse{"couldn't retrieve vacancy"})
    }

    // Получение резюме
    resumes, err := getResumesForVacancy(vacancyID)
    if err != nil {
        log.Println("Error getting resumes:", err)
        return c.JSON(http.StatusInternalServerError, ErrorResponse{"couldn't retrieve resumes"})
    }

    resumesForDashboard := make([]ResumeDashboard, 0, len(resumes))

    for _, resume := range resumes {
        matchPercent, ratingPercent := resume.CalculateMatchAndRating(vacancy)

        resumesForDashboard = append(resumesForDashboard, ResumeDashboard{
            Name:             resume.Name,
            Surname:          resume.Surname,
            Experience:       resume.Experience,
            WorkFormat:       resume.WorkFormat,
            WorkSchedule:     resume.WorkSchedule,
            Employment:       resume.Employment,
            WorkTime:         strconv.Itoa(resume.WorkTimeMin) + " - " + strconv.Itoa(resume.WorkTimeMax),
            Salary:           strconv.Itoa(resume.SalaryMin) + " - " + strconv.Itoa(resume.SalaryMax),
            Responsibilities: resume.Responsibilities,
            HardSkills:       resume.HardSkills,
            SoftSkills:       resume.SoftSkills,
            BusinessTrips:    resume.BusinessTrips,
            Education:        resume.Education,
            Rating:           ratingPercent,
            MatchWithVacancy: matchPercent,
        })
    }

    // Сортируем по рейтингу в порядке убывания
    sort.Slice(resumesForDashboard, func(i, j int) bool {
        return resumesForDashboard[i].Rating > resumesForDashboard[j].Rating
    })

    // Ограничение по topN
    limit := int(math.Min(float64(topN), float64(len(resumesForDashboard))))
    response := DashboardResponse{
        Resumes: resumesForDashboard[:limit],
    }

    return c.JSON(http.StatusOK, response)
}

// getVacancyFromDB извлекает вакансию по ID
func getVacancyFromDB(vacancyID int) (*VacancyFormated, error) {
    row := db.QueryRow("SELECT vacancy FROM vacancies_jsonb WHERE vacancy @> jsonb_build_object('vacancy_id', $1::int)", vacancyID)
    var vJSONRaw []byte
    err := row.Scan(&vJSONRaw)
    if err != nil {
        return nil, err
    }

    vacancy := new(VacancyFormated)
    err = json.Unmarshal(vJSONRaw, &vacancy)
    if err != nil {
        return nil, err
    }
    return vacancy, nil
}

// getResumesForVacancy извлекает все резюме, привязанные к VacancyID
func getResumesForVacancy(vacancyID int) ([]*ResumeFormated, error) {
    var resumes []*ResumeFormated
    rows, err := db.Query("SELECT resume FROM resumes_jsonb WHERE resume @> jsonb_build_object('vacancy_id', $1::int)", vacancyID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var jsonRaw []byte
        err := rows.Scan(&jsonRaw)
        if err != nil {
            log.Println("Error scanning resume row:", err)
            continue
        }
        resume := &ResumeFormated{}
        err = json.Unmarshal(jsonRaw, &resume)
        if err != nil {
            log.Println("Error unmarshalling resume JSON:", err)
            continue
        }
        resumes = append(resumes, resume)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return resumes, nil
}
