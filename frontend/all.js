--

## 📂 index.html

```html
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8" />
  <title>HR JSON Editor</title>
  <style>
    body {
      font-family: sans-serif;
      margin: 20px;
      background: #fafafa;
    }
    h2 {
      margin-top: 30px;
    }
    .tree {
      border-left: 2px solid #ccc;
      margin-left: 10px;
      padding-left: 10px;
    }
    .tree label {
      display: block;
      margin: 4px 0;
    }
    input[type="text"], input[type="number"] {
      width: 200px;
    }
    button {
      margin-top: 10px;
      padding: 5px 10px;
      cursor: pointer;
    }
    textarea {
      width: 100%;
      height: 200px;
      margin-top: 20px;
    }
  </style>
</head>
<body>
  <h1>HR Dashboard JSON Editor</h1>

  <div>
    <button onclick="showSection('vacancy')">Вакансия</button>
    <button onclick="showSection('resume')">Резюме</button>
  </div>

  <div id="vacancy" class="section">
    <h2>Создать вакансию</h2>

    <label>Опыт работы: <input id="exp" type="number" value="1"></label>

    <div class="tree">
      <strong>Work Format</strong>
      <label><input type="checkbox" name="wf" value="Офис"> Офис</label>
      <label><input type="checkbox" name="wf" value="Удалённо" checked> Удалённо</label>
      <label><input type="checkbox" name="wf" value="Гибрид"> Гибрид</label>
    </div>

    <div class="tree">
      <strong>График</strong>
      <label><input type="checkbox" name="ws" value="5/2" checked> 5/2</label>
      <label><input type="checkbox" name="ws" value="Свободный"> Свободный</label>
    </div>

    <div>
      <strong>Обязанности</strong>
      <div id="resp-tree" class="tree"></div>
      <button onclick="addNode('resp')">➕ Добавить обязанность</button>
    </div>

    <div>
      <strong>Hard Skills</strong>
      <div id="hard-tree" class="tree"></div>
      <button onclick="addNode('hard')">➕ Добавить скилл</button>
    </div>

    <div>
      <strong>Soft Skills (через запятую)</strong><br>
      <input id="soft" type="text" value="Коммуникация, Работа в команде">
    </div>

    <label>Командировки: <input id="bt" type="checkbox" checked></label><br>
    <label>Образование: <input id="edu" type="checkbox" checked></label><br>

    <label>Зарплата мин <input id="sal_min" type="number" value="50000"> макс <input id="sal_max" type="number" value="120000"></label><br>
    <label>Время мин <input id="wt_min" type="number" value="30"> макс <input id="wt_max" type="number" value="40"></label><br>

    <button onclick="sendVacancy()">📤 Отправить вакансию</button>

    <h3>Сгенерированный JSON:</h3>
    <textarea id="outVacancy" readonly></textarea>
  </div>

  <div id="resume" class="section" style="display:none">
    <h2>Добавить резюме</h2>
    <label>Имя: <input id="r_name" type="text"></label>
    <label>Фамилия: <input id="r_surname" type="text"></label>
    <label>Опыт: <input id="r_exp" type="number" value="1"></label>

    <div>
      <strong>Hard Skills</strong>
      <div id="r-hard-tree" class="tree"></div>
      <button onclick="addNode('rhard')">➕ Добавить скилл</button>
    </div>

    <div>
      <strong>Soft Skills (через запятую)</strong><br>
      <input id="r-soft" type="text" value="Коммуникация">
    </div>

    <label>Готов к командировкам <input id="r_bt" type="checkbox"></label><br>
    <label>Образование по специальности <input id="r_edu" type="checkbox"></label><br>

    <label>Зарплата мин <input id="r_sal_min" type="number" value="60000"> макс <input id="r_sal_max" type="number" value="90000"></label><br>

    <button onclick="sendResume()">📩 Отправить резюме</button>

    <h3>Сгенерированный JSON:</h3>
    <textarea id="outResume" readonly></textarea>
  </div>

<script>
function showSection(id) {
  document.querySelectorAll(".section").forEach(s => s.style.display = "none");
  document.getElementById(id).style.display = "block";
}

function addNode(prefix) {
  const treeId = prefix === "resp" ? "resp-tree" :
                 prefix === "hard" ? "hard-tree" : "r-hard-tree";
  const container = document.getElementById(treeId);
  const div = document.createElement("div");
  div.innerHTML = `
    <label>${prompt("Название")}:
      <input type="number" min="0" max="2" value="1">
    </label>`;
  container.appendChild(div);
}

// --- Vacancy ---
function sendVacancy() {
  const wf = Array.from(document.querySelectorAll('[name=wf]:checked')).map(i=>i.value);
  const ws = Array.from(document.querySelectorAll('[name=ws]:checked')).map(i=>i.value);
  const respEls = document.querySelectorAll('#resp-tree input');
  const hardEls = document.querySelectorAll('#hard-tree input');
  const resp = {}, hard = {};
  respEls.forEach(el => {
    const name = el.parentElement.textContent.replace(':','').trim();
    resp[name] = Number(el.value);
  });
  hardEls.forEach(el => {
    const name = el.parentElement.textContent.replace(':','').trim();
    hard[name] = Number(el.value);
  });
  const obj = {
    vacancy_id: 0,
    experience: Number(document.getElementById('exp').value),
    experience_weight: 0.5,
    work_format: wf,
    work_format_weight: 0.8,
    work_schedule: ws,
    work_schedule_weight: 0.6,
    employment: ["Полная"],
    employment_weight: 0.9,
    work_time_min: Number(document.getElementById('wt_min').value),
    work_time_max: Number(document.getElementById('wt_max').value),
    work_time_weight: 0.7,
    responsibilities: resp,
    responsibilities_weight: 0.85,
    salary_min: Number(document.getElementById('sal_min').value),
    salary_max: Number(document.getElementById('sal_max').value),
    salary_weight: 0.75,
    hard_skills: hard,
    hard_skills_weight: 0.9,
    soft_skills: document.getElementById('soft').value.split(',').map(s=>s.trim()),
    soft_skills_weight: 0.6,
    business_trips: document.getElementById('bt').checked,
    business_trips_weight: 0.4,
    education: document.getElementById('edu').checked,
    education_weight: 0.3
  };

  document.getElementById("outVacancy").value = JSON.stringify(obj, null, 2);
  fetch("http://localhost:4000/putVacancy", {
    method: 'POST',
    headers: { 'Content-Type':'application/json' },
    body: JSON.stringify(obj)
  }).then(r => r.json().catch(()=>({}))).then(d => {
    alert("✅ Вакансия отправлена");
    console.log(d);
  });
}

// --- Resume ---
function sendResume() {
  const hardEls = document.querySelectorAll('#r-hard-tree input');
  const hard = {};
  hardEls.forEach(el => {
    const name = el.parentElement.textContent.replace(':','').trim();
    hard[name] = Number(el.value);
  });

  const obj = {
    name: document.getElementById('r_name').value,
    surname: document.getElementById('r_surname').value,
    experience: Number(document.getElementById('r_exp').value),
    work_format: ["Удалённо"],
    work_schedule: ["5/2"],
    employment: ["Полная"],
    work_time_min: 30,
    work_time_max: 40,
    responsibilities: { "Разработка": 2 },
    salary_min: Number(document.getElementById('r_sal_min').value),
    salary_max: Number(document.getElementById('r_sal_max').value),
    hard_skills: hard,
    soft_skills: document.getElementById('r-soft').value.split(',').map(s=>s.trim()),
    business_trips: document.getElementById('r_bt').checked,
    education: document.getElementById('r_edu').checked
  };

  document.getElementById("outResume").value = JSON.stringify(obj, null, 2);
  fetch("http://localhost:4000/sendResume", {
    method: 'POST',
    headers: { 'Content-Type':'application/json' },
    body: JSON.stringify(obj)
  }).then(r => r.json().catch(()=>({}))).then(d => {
    alert("📩 Резюме отправлено");
    console.log(d);
  });
}
</script>
</body>
</html>
```
