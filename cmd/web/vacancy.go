package main

type VacancyFormated struct {
    VacancyID           int         `json:"vacancy_id"`
    // Опыт работы: [0;+inf]
    Experience          int         `json:"experience"`
    ExperienceWeight    float64     `json:"experience_weight"`
    // Формат работы: "Офис" | "Удалённо" | "Гибрид" | "Разъездной" | "Любой"
    WorkFormat          []string    `json:"work_format"`
    WorkFormatWeight    float64     `json:"work_format_weight"`
    // График работы: "5/2" | "2/2" | "3/3" | "3/2" | "4/3" | "4/2" | "6/1" | "1/3" | "2/1" | "1/2" | "4/4" | "Свободный" | "Другой" | "Любой"
    WorkSchedule        []string    `json:"work_schedule"`
    WorkScheduleWeight  float64     `json:"work_schedule_weight"`
    // Тип занятости: "Полная" | "Частичная" | "Проектная" | "Вахта" | "ГПХ или по совместительству" | "Стажировка" | "Любой"
    Employment          []string    `json:"employment"`
    EmploymentWeight    float64     `json:"employment_weight"`
    // Рабочие часы: [WorkTimeMin;WorkImeMax]
    WorkTimeMin         int         `json:"work_time_min"`
    WorkTimeMax         int         `json:"work_time_max"`
    WorkTimeWeight      float64     `json:"work_time_weight"`
    // Обязанности: map["название_обязанности"] = число от 0 до 2, где 0 - неможет выполнить, 1 - может научиться, 2 - уже выполнял такие задачи
    Responsibilities map[string]int `json:"responsibilities"`
    ResponsibilitiesWeight float64  `json:"responsibilities_weight"`
    // Зарплата: [SalaryMin;SalaryMax]
    SalaryMin           int         `json:"salary_min"`
    SalaryMax           int         `json:"salary_max"`
    SalaryWeight        float64     `json:"salary_weight"`
    // Hard skills: map["название_скилла"] = число от 0 до 2, где 0 - не владеет, 1 - решал похожие задачи, 2 - владеет навыком, имеет опыт
    HardSkills       map[string]int `json:"hard_skills"`
    HardSkillsWeight    float64     `json:"hard_skills_weight"`
    // Soft skills: просто список скиллов
    SoftSkills          []string    `json:"soft_skills"`
    SoftSkillsWeight    float64     `json:"soft_skills_weight"`
    // Командировки: да/нет (готов ехать/не готов)
    BusinessTrips       bool        `json:"business_trips"`
    BusinessTripsWeight float64     `json:"business_trips_weight"`
    // Наличие образования по специальности да/нет
    Education           bool        `json:"education"`
    EducationWeight     float64     `json:"education_weight"`
}

func (*VacancyFormated) addToDB() {

}

func (*VacancyFormated) getFromDB() {

}
