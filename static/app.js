// ============== Константы ==============
const API = "/tasks";
const GRID_START = 8;     // сетка начинается с 8:00
const GRID_HOURS = 14;    // показываем 14 часов (8:00 – 22:00)

const WEEKDAY_NAMES = ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"];
const MONTH_NAMES = [
    "январь", "февраль", "март", "апрель", "май", "июнь",
    "июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь"
];
const MONTH_GENITIVE = [
    "января", "февраля", "марта", "апреля", "мая", "июня",
    "июля", "августа", "сентября", "октября", "ноября", "декабря"
];

// ============== Состояние ==============
let allTasks = [];
let weekStart = mondayOf(new Date());     // понедельник текущей видимой недели
let miniMonth = new Date(weekStart);      // месяц в мини-календаре
let currentFilter = "all";

// ============== Хелперы дат ==============
function mondayOf(date) {
    const d = new Date(date);
    d.setHours(0, 0, 0, 0);
    const day = d.getDay(); // 0=вс
    const diff = day === 0 ? -6 : 1 - day;
    d.setDate(d.getDate() + diff);
    return d;
}
function addDays(date, n) {
    const d = new Date(date);
    d.setDate(d.getDate() + n);
    return d;
}
function formatYMD(date) {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const d = String(date.getDate()).padStart(2, "0");
    return `${y}-${m}-${d}`;
}
function sameDay(a, b) {
    return a.getFullYear() === b.getFullYear()
        && a.getMonth() === b.getMonth()
        && a.getDate() === b.getDate();
}
function inCurrentWeek(date) {
    const end = addDays(weekStart, 7);
    return date >= weekStart && date < end;
}

// ============== API ==============
async function loadTasks() {
    const res = await fetch(API);
    const data = await res.json();
    allTasks = data || [];
    renderAll();
}
async function createTask(payload) {
    await fetch(API, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
    });
    loadTasks();
}
async function updateTask(id, payload) {
    await fetch(`${API}/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
    });
    loadTasks();
}
async function deleteTask(id) {
    await fetch(`${API}/${id}`, { method: "DELETE" });
    loadTasks();
}

// ============== Фильтр ==============
function passesFilter(t) {
    return currentFilter === "all" || t.priority === Number(currentFilter);
}

// ============== Главная отрисовка ==============
function renderAll() {
    renderWeekHeader();
    renderWeekGrid();
    renderMiniCal();
    renderGlobalList();
}

function renderWeekHeader() {
    const end = addDays(weekStart, 6);
    let label;
    if (weekStart.getMonth() === end.getMonth()) {
        label = `${weekStart.getDate()}–${end.getDate()} ${MONTH_GENITIVE[weekStart.getMonth()]} ${weekStart.getFullYear()}`;
    } else {
        label = `${weekStart.getDate()} ${MONTH_GENITIVE[weekStart.getMonth()]} – ${end.getDate()} ${MONTH_GENITIVE[end.getMonth()]} ${end.getFullYear()}`;
    }
    document.getElementById("weekLabel").textContent = label;
}

function renderWeekGrid() {
    const grid = document.getElementById("weekGrid");
    grid.innerHTML = "";

    const today = new Date();
    today.setHours(0, 0, 0, 0);

    // --- Шапка с днями ---
    const headerRow = document.createElement("div");
    headerRow.className = "day-header-row";
    headerRow.appendChild(document.createElement("div")); // пустой угол

    for (let i = 0; i < 7; i++) {
        const date = addDays(weekStart, i);
        const cell = document.createElement("div");
        cell.className = "day-header";
        if (sameDay(date, today)) cell.classList.add("today");
        cell.innerHTML = `
            <div class="day-header-name">${WEEKDAY_NAMES[i]}</div>
            <div class="day-header-num">${date.getDate()}</div>
        `;
        headerRow.appendChild(cell);
    }
    grid.appendChild(headerRow);

    // --- Строка "Весь день" (задачи с датой, но без времени) ---
    const alldayRow = document.createElement("div");
    alldayRow.className = "allday-row";

    const alldayLabel = document.createElement("div");
    alldayLabel.className = "allday-label";
    alldayLabel.textContent = "Весь день";
    alldayRow.appendChild(alldayLabel);

    for (let i = 0; i < 7; i++) {
        const date = addDays(weekStart, i);
        const ymd = formatYMD(date);
        const cell = document.createElement("div");
        cell.className = "allday-cell";

        const alldayTasks = allTasks.filter(t =>
            passesFilter(t)
            && t.deadline === ymd
            && (!t.start_time || !t.end_time)
        );
        for (const task of alldayTasks) {
            const el = document.createElement("div");
            el.className = `allday-task priority-${task.priority}${task.done ? " done" : ""}`;
            if (task.color) el.style.background = task.color;
            el.textContent = task.title;
            el.title = task.title;
            el.addEventListener("click", () => toggleDone(task));
            cell.appendChild(el);
        }
        alldayRow.appendChild(cell);
    }
    grid.appendChild(alldayRow);

    // --- Основная сетка по часам ---
    const timeGrid = document.createElement("div");
    timeGrid.className = "time-grid";

    // Колонка с подписями часов
    const timeCol = document.createElement("div");
    timeCol.className = "time-col";
    for (let h = GRID_START; h < GRID_START + GRID_HOURS; h++) {
        const lbl = document.createElement("div");
        lbl.className = "time-label";
        lbl.textContent = `${String(h).padStart(2, "0")}:00`;
        timeCol.appendChild(lbl);
    }
    timeGrid.appendChild(timeCol);

    // Колонки дней
    for (let i = 0; i < 7; i++) {
        const date = addDays(weekStart, i);
        const ymd = formatYMD(date);
        const col = document.createElement("div");
        col.className = "day-col";
        if (sameDay(date, today)) col.classList.add("today");

        // Линии часов
        for (let h = 1; h < GRID_HOURS; h++) {
            const line = document.createElement("div");
            line.className = "hour-line";
            line.style.top = (h * 60) + "px";
            col.appendChild(line);
        }

        // Задачи с временем на этот день
        const dayTasks = allTasks.filter(t =>
            passesFilter(t)
            && t.deadline === ymd
            && t.start_time && t.end_time
        );

        for (const task of dayTasks) {
            const [sh, sm] = task.start_time.split(":").map(Number);
            const [eh, em] = task.end_time.split(":").map(Number);
            const startMin = (sh - GRID_START) * 60 + sm;
            const endMin = (eh - GRID_START) * 60 + em;
            // Обрезаем по границам сетки
            const maxMin = GRID_HOURS * 60;
            if (endMin <= 0 || startMin >= maxMin) continue;
            const top = Math.max(0, startMin);
            const height = Math.max(22, Math.min(maxMin, endMin) - top);

            const block = document.createElement("div");
            block.className = `task-block priority-${task.priority}${task.done ? " done" : ""}`;
            block.style.top = top + "px";
            block.style.height = height + "px";
            if (task.color) block.style.background = task.color;

            block.innerHTML = `
                <button class="task-block-delete" type="button" title="Удалить">×</button>
                <div class="task-block-title">${escapeHtml(task.title)}</div>
                <div class="task-block-time">${task.start_time} – ${task.end_time}</div>
            `;

            block.querySelector(".task-block-delete").addEventListener("click", (e) => {
                e.stopPropagation();
                deleteTask(task.id);
            });
            block.addEventListener("click", () => toggleDone(task));

            col.appendChild(block);
        }

        timeGrid.appendChild(col);
    }
    grid.appendChild(timeGrid);
}

function renderMiniCal() {
    document.getElementById("miniMonth").textContent =
        `${MONTH_NAMES[miniMonth.getMonth()]} ${miniMonth.getFullYear()}`;

    const grid = document.getElementById("miniGrid");
    grid.innerHTML = "";

    const firstOfMonth = new Date(miniMonth.getFullYear(), miniMonth.getMonth(), 1);
    const start = mondayOf(firstOfMonth);

    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const datesWithTasks = new Set(
        allTasks.filter(t => t.deadline).map(t => t.deadline)
    );

    // Всегда 6 недель = 42 ячейки для постоянной высоты
    for (let i = 0; i < 42; i++) {
        const date = addDays(start, i);
        const btn = document.createElement("button");
        btn.className = "mini-day";
        if (date.getMonth() !== miniMonth.getMonth()) btn.classList.add("other-month");
        if (inCurrentWeek(date)) btn.classList.add("in-week");
        if (sameDay(date, today)) btn.classList.add("today");
        if (datesWithTasks.has(formatYMD(date))) btn.classList.add("has-tasks");
        btn.textContent = date.getDate();
        btn.addEventListener("click", () => {
            weekStart = mondayOf(date);
            renderAll();
        });
        grid.appendChild(btn);
    }
}

function renderGlobalList() {
    const list = document.getElementById("globalList");
    list.innerHTML = "";

    const globals = allTasks.filter(t => passesFilter(t) && !t.deadline);
    if (globals.length === 0) {
        list.innerHTML = `<div class="empty">Нет задач</div>`;
        return;
    }

    for (const task of globals) {
        const el = document.createElement("div");
        el.className = `global-task priority-${task.priority}${task.done ? " done" : ""}`;
        el.innerHTML = `
            <input type="checkbox" class="global-task-checkbox" ${task.done ? "checked" : ""}>
            <span class="global-task-title" title="${escapeHtml(task.title)}">${escapeHtml(task.title)}</span>
            <button class="global-task-delete" title="Удалить">×</button>
        `;
        el.querySelector(".global-task-checkbox").addEventListener("change", () => toggleDone(task));
        el.querySelector(".global-task-delete").addEventListener("click", () => deleteTask(task.id));
        list.appendChild(el);
    }
}

// ============== Действия ==============
function toggleDone(task) {
    updateTask(task.id, {
        title: task.title,
        description: task.description,
        priority: task.priority,
        start_time: task.start_time,
        end_time: task.end_time,
        color: task.color,
        deadline: task.deadline,
        done: !task.done
    });
}

function escapeHtml(str) {
    if (!str) return "";
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
}

// ============== Обработчики ==============
document.getElementById("addForm").addEventListener("submit", (e) => {
    e.preventDefault();
    const title = document.getElementById("title").value.trim();
    if (!title) return;
    const payload = {
        title,
        description: document.getElementById("description").value.trim(),
        priority: Number(document.getElementById("priority").value),
        deadline: document.getElementById("deadline").value,
        start_time: document.getElementById("startTime").value,
        end_time: document.getElementById("endTime").value,
        color: ""
    };
    createTask(payload);
    e.target.reset();
    document.getElementById("priority").value = "2";
});

document.getElementById("prevWeek").addEventListener("click", () => {
    weekStart = addDays(weekStart, -7);
    renderAll();
});
document.getElementById("nextWeek").addEventListener("click", () => {
    weekStart = addDays(weekStart, 7);
    renderAll();
});
document.getElementById("todayBtn").addEventListener("click", () => {
    weekStart = mondayOf(new Date());
    miniMonth = new Date(weekStart);
    renderAll();
});

document.getElementById("miniPrev").addEventListener("click", () => {
    miniMonth = new Date(miniMonth.getFullYear(), miniMonth.getMonth() - 1, 1);
    renderMiniCal();
});
document.getElementById("miniNext").addEventListener("click", () => {
    miniMonth = new Date(miniMonth.getFullYear(), miniMonth.getMonth() + 1, 1);
    renderMiniCal();
});

document.querySelectorAll(".filter-btn").forEach(btn => {
    btn.addEventListener("click", () => {
        document.querySelector(".filter-btn.active").classList.remove("active");
        btn.classList.add("active");
        currentFilter = btn.dataset.filter;
        renderAll();
    });
});

// ============== Старт ==============
loadTasks();
