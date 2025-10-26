package main

type ResumeFormated struct {
    ID                  int         `json:"id"`
    VacancyID           int         `json:"vacancy_id"`
    // Имя
    Name                string      `json:"name"`
    // Фамилия
    Surname             string      `json:"surname"`
    // Опыт работы: [0;+inf]
    Experience          int         `json:"experience"`
    // Формат работы: "Офис" | "Удалённо" | "Гибрид" | "Разъездной" | "Любой"
    WorkFormat          []string    `json:"work_format"`
    // График работы: "5/2" | "2/2" | "3/3" | "3/2" | "4/3" | "4/2" | "6/1" | "1/3" | "2/1" | "1/2" | "4/4" | "Свободный" | "Другой" | "Любой"
    WorkSchedule        []string    `json:"work_schedule"`
    // Тип занятости: "Полная" | "Частичная" | "Проектная" | "Вахта" | "ГПХ или по совместительству" | "Стажировка" | "Любой"
    Employment          []string    `json:"employment"`
    // Рабочие часы: [WorkTimeMin;WorkImeMax]
    WorkTimeMin         int         `json:"work_time_min"`
    WorkTimeMax         int         `json:"work_time_max"`
    // Обязанности: map["название_обязанности"] = число от 0 до 2, где 0 - неможет выполнить, 1 - может научиться, 2 - уже выполнял такие задачи
    Responsibilities map[string]int `json:"responsibilities"`
    // Зарплата: [SalaryMin;SalaryMax]
    SalaryMin           int         `json:"salary_min"`
    SalaryMax           int         `json:"salary_max"`
    // Hard skills: map["название_скилла"] = число от 0 до 2, где 0 - не владеет, 1 - решал похожие задачи, 2 - владеет навыком, имеет опыт
    HardSkills       map[string]int `json:"hard_skills"`
    // Soft skills: просто список скиллов
    SoftSkills          []string    `json:"soft_skills"`
    // Командировки: да/нет (готов ехать/не готов)
    BusinessTrips       bool        `json:"business_trips"`
    // Наличие образования по специальности да/нет
    Education           bool        `json:"education"`
}

func (*ResumeFormated) addToDB() {

}

func (*ResumeFormated) getFromDB() {

}
