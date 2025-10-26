package main
import (
    "net/http"
    "github.com/labstack/echo/v4"
    "strconv"
    "slices"
    "math"
    "sort"
)

// Приём резюме, запись в БД
func resumeHandler(c echo.Context) (err error) {
    resume := new(ResumeFormated)
    if err := c.Bind(resume); err != nil {
        return c.String(http.StatusBadRequest, "bad request")
    }

    resume.addToDB()

    return c.JSON(http.StatusOK, resume)
}

// Приём вакансии, запись в БД
func vacancyHandler(c echo.Context) (err error) {
    vacancy := new(VacancyFormated)
    if err := c.Bind(vacancy); err != nil {
        return c.String(http.StatusBadRequest, "bad request")
    }

    vacancy.addToDB()

    return c.JSON(http.StatusOK, vacancy)
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
    vacancy.getFromDB()

    var resumes []ResumeFormated

    // ТУТА МЫ ЗАПРАШИВАЕМ ДАННЫЕ ИЗ БД, а именно:
    //  нужно получить все резюме, у которых ID = vacancy_id
    //  и записать это всё безобразие в массив resumes

    type ResumeDashboard struct {
        Name                string
        Surname             string
        Rating              int     // от 0 до 100
        MatchWithVacancy    int     // от 0 до 100
        Experience          int
        WorkFormat          []string
        WorkSchedule        []string
        Employment          []string
        WorkTime            string
        Responsibilities    map[string]int
        Salary              string
        HardSkills          map[string]int
        SoftSkills          []string
        BusinessTrips       bool
        Education           bool
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
        resumes []ResumeDashboard   `json:"resumes"`
    }

    response := DashboardResponse{
        resumes: resumesForDashboard[:topN],
    }

    return c.JSON(http.StatusOK, response)
}
