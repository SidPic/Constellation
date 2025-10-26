package main
import (
    "net/http"
    "github.com/labstack/echo/v4"
    "strconv"
    "slices"
    "math"
    "sort"
    "io"
    "encoding/json"
    "bytes"
    "log"
)

// Приём резюме, запись в БД
func resumeHandler(c echo.Context) (err error) {

    // читаем тело (можно ограничить размер через io.LimitReader при необходимости)
    b, err := io.ReadAll(c.Request().Body)
    if err != nil {
        return c.String(http.StatusBadRequest, "cannot read request body")
    }
    // валидируем JSON
    if !json.Valid(b) {
        return c.String(http.StatusBadRequest, "invalid json")
    }

    // компактная запись (удаляет лишние пробелы)
    var dst bytes.Buffer
    if err := json.Compact(&dst, b); err != nil {
    // на случай редкой ошибки компактизации — вернуть 400
        return c.String(http.StatusBadRequest, "cannot compact json: "+err.Error())
    }

    // Вставка в БД: placeholder должен быть без кавычек ($1), передаём []byte или string
    result, err := db.Exec("INSERT INTO resumes_jsonb (resume) VALUES ($1)", dst.Bytes())
    if err != nil {
        return c.String(http.StatusInternalServerError, "couldn't insert resume into database: "+err.Error())
    }

    rowsAffected, _ := result.RowsAffected()
    return c.String(http.StatusOK, strconv.FormatInt(rowsAffected, 10))
}

// Приём вакансии, запись в БД
func vacancyHandler(c echo.Context) (err error) {
    // читаем тело (можно ограничить размер через io.LimitReader при необходимости)
    b, err := io.ReadAll(c.Request().Body)
    if err != nil {
        return c.String(http.StatusBadRequest, "cannot read request body")
    }
    // валидируем JSON
    if !json.Valid(b) {
        return c.String(http.StatusBadRequest, "invalid json")
    }

    // компактная запись (удаляет лишние пробелы)
    var dst bytes.Buffer
    if err := json.Compact(&dst, b); err != nil {
    // на случай редкой ошибки компактизации — вернуть 400
        return c.String(http.StatusBadRequest, "cannot compact json: "+err.Error())
    }

    // Вставка в БД: placeholder должен быть без кавычек ($1), передаём []byte или string
    result, err := db.Exec("INSERT INTO vacancies_jsonb (vacancy) VALUES ($1)", dst.Bytes())
    if err != nil {
        return c.String(http.StatusInternalServerError, "couldn't insert vacancy into database: "+err.Error())
    }

    rowsAffected, _ := result.RowsAffected()
    return c.String(http.StatusOK, strconv.FormatInt(rowsAffected, 10))
}

// Отправка данных на дашборд
func dashboardHandler(c echo.Context) (err error) {
    vacancy_id, err := strconv.Atoi(c.QueryParam("vacancy_id"))
    if err != nil {
        return c.String(http.StatusBadRequest, "bad query parameter")
    }

    topN, err := strconv.Atoi(c.QueryParam("top_n"))
    if err != nil {
        return c.String(http.StatusBadRequest, "bad number of rows")
    }

    vacancy := new(VacancyFormated)
    vacancy.VacancyID = vacancy_id

    // получение вакансии

    row := db.QueryRow("SELECT vacancy FROM vacancies_jsonb WHERE vacancy @> jsonb_build_object('vacancy_id', $1::int)", vacancy_id)
    var v_json_raw []byte
    err = row.Scan(&v_json_raw)
    if err != nil{
        log.Println(err.Error())
    }
    err = json.Unmarshal(v_json_raw, &vacancy)
    if err != nil{
        log.Println(err.Error())
    }

    // получение резюме

    var resumes []ResumeFormated

    rows, err := db.Query("SELECT resume FROM resumes_jsonb WHERE resume @> jsonb_build_object('vacancy_id', $1::int)", vacancy_id)
    if err != nil {
        return c.String(http.StatusInternalServerError, "couldn't get resumes from database: "+err.Error())
    }
    defer rows.Close()

    rows_number := 0
    for rows.Next(){
        var json_raw []byte
        err := rows.Scan(&json_raw)
        if err != nil{
            log.Println(err.Error())
            continue
        }
        resume := ResumeFormated{}
        err = json.Unmarshal(json_raw, &resume)
        if err != nil{
            log.Println(err.Error())
            continue
        }
        resumes = append(resumes, resume)
        rows_number++
    }

    type ResumeDashboard struct {
        Name                string      `json:"Name"`
        Surname             string      `json:"Surname"`
        Rating              int         `json:"Rating"`
        MatchWithVacancy    int         `json:"MatchWithVacancy"`
        Experience          int         `json:"Experience"`
        WorkFormat          []string    `json:"WorkFormat"`
        WorkSchedule        []string    `json:"WorkSchedule"`
        Employment          []string    `json:"Employment"`
        WorkTime            string      `json:"WorkTime"`
        Responsibilities    map[string]int `json"Responsibilities"`
        Salary              string      `json:"Salary"`
        HardSkills          map[string]int `json:"HardSkills"`
        SoftSkills          []string    `json:"SoftSkills"`
        BusinessTrips       bool        `json:"BusinessTrips"`
        Education           bool        `json:"Education"`
    }

    resumesForDashboard := make([]ResumeDashboard, 0, len(resumes))

    for r, resume := range resumes {
        resumesForDashboard = append(resumesForDashboard, ResumeDashboard{
            Name:           resume.Name,
            Surname:        resume.Surname,
            Experience:     resume.Experience,
            WorkFormat:     resume.WorkFormat,
            WorkSchedule:   resume.WorkSchedule,
            Employment:     resume.Employment,
            WorkTime:       strconv.Itoa(resume.WorkTimeMin) + " - " + strconv.Itoa(resume.WorkTimeMax),
            Salary:         strconv.Itoa(resume.SalaryMin) + " - " + strconv.Itoa(resume.SalaryMax),
            Responsibilities: resume.Responsibilities,
            HardSkills:     resume.HardSkills,
            SoftSkills:     resume.SoftSkills,
            BusinessTrips:  resume.BusinessTrips,
            Education:      resume.Education,
        })

        // считаем доли соответствия
        var experienceMatch float64
        if resume.Experience >= vacancy.Experience {
            experienceMatch = 1.0
        } else {
            experienceMatch = 1.0 / (float64(vacancy.Experience) - float64(resume.Experience))
        }

        matched := 0
        for _, format := range vacancy.WorkFormat {
            if slices.Contains(resume.WorkFormat, format) {
                matched++
            }
        }
        workFormatMatch := float64(matched) / float64(len(vacancy.WorkFormat))

        matched = 0
        for _, schedule := range vacancy.WorkSchedule {
            if slices.Contains(resume.WorkSchedule, schedule) {
                matched++
            }
        }
        workScheduleMatch := float64(matched) / float64(len(vacancy.WorkSchedule))

        matched = 0
        for _, employment := range vacancy.Employment {
            if slices.Contains(resume.Employment, employment) {
                matched++
            }
        }
        employmentMatch := float64(matched) / float64(len(vacancy.Employment))

        workTimeMatch := 1.0 / math.Abs(float64(resume.WorkTimeMin - vacancy.WorkTimeMin + resume.WorkTimeMax - vacancy.WorkTimeMax)  / 2.0)

        salaryMatch := 1.0 / math.Abs(float64(resume.SalaryMin - vacancy.SalaryMin + resume.SalaryMax - vacancy.SalaryMax)  / 2.0)

        responsibilitiesMatch := 0.0
        for i, resp := range vacancy.Responsibilities {
            respMatch := 0.0
            if (resume.Responsibilities[i] >= resp) {
                respMatch = 1.0
            } else {
                respMatch = float64(resume.Responsibilities[i]) / float64(resp)
            }
            responsibilitiesMatch += float64(respMatch) / float64(len(resume.Responsibilities))
        }

        hardsMatch := 0.0
        for skill, value := range vacancy.HardSkills {
            match := 0.0
            if (resume.HardSkills[skill] >= value) {
                match = 1.0
            } else {
                match = float64(resume.HardSkills[skill]) / float64(value)
            }
            hardsMatch += float64(match) / float64(len(resume.HardSkills))
        }

        matched = 0
        for _, soft := range vacancy.SoftSkills {
            if slices.Contains(resume.SoftSkills, soft) {
                matched++
            }
        }
        softsMatch := float64(matched) / float64(len(vacancy.SoftSkills))

        var tripsMatch float64
        if resume.BusinessTrips || resume.BusinessTrips == vacancy.BusinessTrips {
            tripsMatch = 1.0
        } else {
            tripsMatch = 0.0
        }

        var eduMatch float64
        if resume.Education || resume.Education == vacancy.Education {
            eduMatch = 1.0
        } else {
            eduMatch = 0.0
        }

        // Расчёт доли совпадения резюме с вакансией
        resumesForDashboard[r].MatchWithVacancy = int((experienceMatch + workFormatMatch + workScheduleMatch + employmentMatch + workTimeMatch + salaryMatch + responsibilitiesMatch + hardsMatch + softsMatch + tripsMatch + eduMatch) / 11.0 * 100)

        // считаем рейтинг
        resumesForDashboard[r].Rating = int((experienceMatch*vacancy.ExperienceWeight + workFormatMatch*vacancy.WorkFormatWeight + workScheduleMatch*vacancy.WorkScheduleWeight + employmentMatch*vacancy.EmploymentWeight + workTimeMatch*vacancy.WorkTimeWeight + salaryMatch*vacancy.SalaryWeight + responsibilitiesMatch*vacancy.ResponsibilitiesWeight + hardsMatch*vacancy.HardSkillsWeight + softsMatch*vacancy.SoftSkillsWeight + tripsMatch*vacancy.BusinessTripsWeight + eduMatch*vacancy.EducationWeight) / 11.0 * 100)

        // сортируем
        sort.Slice(resumesForDashboard, func(i, j int) (less bool) {
            return resumesForDashboard[i].Rating > resumesForDashboard[j].Rating
        })
    }

    type DashboardResponse struct {
        Resumes []ResumeDashboard   `json:"resumes"`
    }

    response := DashboardResponse{
        Resumes: resumesForDashboard[:int(math.Min(float64(topN), float64(rows_number)))],
    }

    return c.JSON(http.StatusOK, response)
}
