package main

import (
    "math"
    "slices"
)

type ResumeFormated struct {
    ID               int         `json:"id"`
    VacancyID        int         `json:"vacancy_id"`
    // Имя
    Name             string      `json:"name"`
    // Фамилия
    Surname          string      `json:"surname"`
    // Опыт работы: [0;+inf]
    Experience       int         `json:"experience"`
    // Формат работы: "Офис" | "Удалённо" | "Гибрид" | "Разъездной" | "Любой"
    WorkFormat       []string    `json:"work_format"`
    // График работы: "5/2" | "2/2" | "3/3" | "3/2" | "4/3" | "4/2" | "6/1" | "1/3" | "2/1" | "1/2" | "4/4" | "Свободный" | "Другой" | "Любой"
    WorkSchedule     []string    `json:"work_schedule"`
    // Тип занятости: "Полная" | "Частичная" | "Проектная" | "Вахта" | "ГПХ или по совместительству" | "Стажировка" | "Любой"
    Employment       []string    `json:"employment"`
    // Рабочие часы: [WorkTimeMin;WorkImeMax]
    WorkTimeMin      int         `json:"work_time_min"`
    WorkTimeMax      int         `json:"work_time_max"`
    // Обязанности: map["название_обязанности"] = число от 0 до 2, где 0 - неможет выполнить, 1 - может научиться, 2 - уже выполнял такие задачи
    Responsibilities map[string]int `json:"responsibilities"`
    // Зарплата: [SalaryMin;SalaryMax]
    SalaryMin        int         `json:"salary_min"`
    SalaryMax        int         `json:"salary_max"`
    // Hard skills: map["название_скилла"] = число от 0 до 2, где 0 - не владеет, 1 - решал похожие задачи, 2 - владеет навыком, имеет опыт
    HardSkills       map[string]int `json:"hard_skills"`
    // Soft skills: просто список скиллов
    SoftSkills       []string    `json:"soft_skills"`
    // Командировки: да/нет (готов ехать/не готов)
    BusinessTrips    bool        `json:"business_trips"`
    // Наличие образования по специальности да/нет
    Education        bool        `json:"education"`
}

func (*ResumeFormated) addToDB() {
    // Логика работы с БД должна быть здесь, но пока оставлена пустой как в оригинале
}

func (*ResumeFormated) getFromDB() {
    // Логика работы с БД должна быть здесь, но пока оставлена пустой как в оригинале
}

// CalculateMatchAndRating рассчитывает долю совпадения (MatchWithVacancy) и рейтинг (Rating) резюме с вакансией.
// Возвращает (MatchPercent, RatingPercent)
func (r *ResumeFormated) CalculateMatchAndRating(v *VacancyFormated) (int, int) {
    // --- Расчёт долей соответствия (Match Scores) ---

    // 1. Опыт работы
    experienceMatch := 0.0
    if r.Experience >= v.Experience {
        experienceMatch = 1.0
    } else if v.Experience > 0 {
        // Убрал 1.0 / (float64(v.Experience) - float64(r.Experience)) так как может дать > 1.0 или очень большое число.
        // Более логичный подход: (факт / требование), если факт < требования
        experienceMatch = float64(r.Experience) / float64(v.Experience)
    }

    // 2. Формат работы
    workFormatMatch := calculateListMatch(r.WorkFormat, v.WorkFormat)

    // 3. График работы
    workScheduleMatch := calculateListMatch(r.WorkSchedule, v.WorkSchedule)

    // 4. Тип занятости
    employmentMatch := calculateListMatch(r.Employment, v.Employment)

    // 5. Рабочие часы
    // Использован более безопасный подход для избежания деления на ноль, если разница мала.
    workTimeMatch := 0.0
    diffWorkTime := math.Abs(float64(r.WorkTimeMin-v.WorkTimeMin) + float64(r.WorkTimeMax-v.WorkTimeMax))
    if diffWorkTime == 0 {
        workTimeMatch = 1.0
    } else {
        // Если разница в сумме диапазонов не ноль, использовать обратную пропорцию
        workTimeMatch = 1.0 / (1.0 + diffWorkTime/2.0)
    }

    // 6. Зарплата
    salaryMatch := 0.0
    diffSalary := math.Abs(float64(r.SalaryMin-v.SalaryMin) + float64(r.SalaryMax-v.SalaryMax))
    if diffSalary == 0 {
        salaryMatch = 1.0
    } else {
        // Если разница в сумме диапазонов не ноль, использовать обратную пропорцию
        salaryMatch = 1.0 / (1.0 + diffSalary/2.0)
    }

    // 7. Обязанности (Responsibilities)
    responsibilitiesMatch := 0.0
    totalWeight := 0.0
    for key, vacValue := range v.Responsibilities {
        totalWeight += 1.0
        resumeValue := r.Responsibilities[key]
        if vacValue > 0 {
            // Если требование (vacValue) > 0: match = min(факт/требование, 1.0)
            responsibilitiesMatch += math.Min(float64(resumeValue)/float64(vacValue), 1.0)
        } else {
            // Если требование 0, то любое значение > 0, дает 100% совпадение, 0 дает 0%.
            if resumeValue > 0 {
                responsibilitiesMatch += 1.0
            } else {
                responsibilitiesMatch += 0.0
            }
        }
    }
    if totalWeight > 0 {
        responsibilitiesMatch /= totalWeight
    }

    // 8. Hard skills
    hardsMatch := 0.0
    totalWeight = 0.0
    for key, vacValue := range v.HardSkills {
        totalWeight += 1.0
        resumeValue := r.HardSkills[key]
        if vacValue > 0 {
            // Если требование (vacValue) > 0: match = min(факт/требование, 1.0)
            hardsMatch += math.Min(float64(resumeValue)/float64(vacValue), 1.0)
        } else {
            // Если требование 0, то любое значение > 0, дает 100% совпадение, 0 дает 0%.
            if resumeValue > 0 {
                hardsMatch += 1.0
            } else {
                hardsMatch += 0.0
            }
        }
    }
    if totalWeight > 0 {
        hardsMatch /= totalWeight
    }

    // 9. Soft skills
    softsMatch := calculateListMatch(r.SoftSkills, v.SoftSkills)

    // 10. Командировки
    tripsMatch := 0.0
    // Готовность ехать ИЛИ соответствие требованию
    if r.BusinessTrips || r.BusinessTrips == v.BusinessTrips {
        tripsMatch = 1.0
    }

    // 11. Образование
    eduMatch := 0.0
    // Наличие образования ИЛИ соответствие требованию
    if r.Education || r.Education == v.Education {
        eduMatch = 1.0
    }

    // Сумма всех долей
    matchScores := []float64{
        experienceMatch,
        workFormatMatch,
        workScheduleMatch,
        employmentMatch,
        workTimeMatch,
        salaryMatch,
        responsibilitiesMatch,
        hardsMatch,
        softsMatch,
        tripsMatch,
        eduMatch,
    }

    // --- Расчёт MatchWithVacancy (простой процент совпадения) ---
    totalMatchSum := 0.0
    for _, score := range matchScores {
        totalMatchSum += score
    }
    matchWithVacancy := totalMatchSum / float64(len(matchScores)) * 100.0

    // --- Расчёт Rating (взвешенный процент совпадения) ---
    weights := []float64{
        v.ExperienceWeight,
        v.WorkFormatWeight,
        v.WorkScheduleWeight,
        v.EmploymentWeight,
        v.WorkTimeWeight,
        v.SalaryWeight,
        v.ResponsibilitiesWeight,
        v.HardSkillsWeight,
        v.SoftSkillsWeight,
        v.BusinessTripsWeight,
        v.EducationWeight,
    }

    weightedSum := 0.0
    totalWeightSum := 0.0
    for i, score := range matchScores {
        weightedSum += score * weights[i]
        totalWeightSum += weights[i]
    }

    rating := 0.0
    if totalWeightSum > 0 {
        rating = weightedSum / totalWeightSum * 100.0
    }

    return int(math.Round(matchWithVacancy)), int(math.Round(rating))
}

// calculateListMatch - вспомогательная функция для расчета соответствия списков (WorkFormat, SoftSkills и т.д.)
func calculateListMatch(resumeList []string, vacancyList []string) float64 {
    if len(vacancyList) == 0 {
        return 1.0 // Если вакансия не требует, считаем 100%
    }
    matched := 0
    for _, item := range vacancyList {
        if slices.Contains(resumeList, item) {
            matched++
        }
    }
    // Используем len(vacancyList) для нормализации, так как сравниваем с требованиями вакансии
    return float64(matched) / float64(len(vacancyList))
}
