const LEVEL_LABELS = { view: "детали", free_busy: "занятость" };
const PRIORITY_LABELS = { 1: "Низкий", 2: "Средний", 3: "Высокий" };

function escapeHtml(str) {
    if (!str) return "";
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
}

function initials(name) {
    return name ? name.charAt(0) : "?";
}

const MONTH_GENITIVE = [
    "января", "февраля", "марта", "апреля", "мая", "июня",
    "июля", "августа", "сентября", "октября", "ноября", "декабря"
];

function formatDeadline(ymd) {
    if (!ymd) return "";
    const [y, m, d] = ymd.split("-").map(Number);
    return `${d} ${MONTH_GENITIVE[m - 1]} ${y}`;
}

// --- Список тех, кто открыл мне календарь ---
async function loadShared() {
    const res = await fetch("/api/access/shared-with-me");
    const people = await res.json() || [];

    const listEl  = document.getElementById("sharedPeople");
    const emptyEl = document.getElementById("sharedEmpty");
    const countEl = document.getElementById("sharedCount");

    listEl.innerHTML = "";
    countEl.textContent = people.length ? `(${people.length})` : "";

    if (people.length === 0) {
        emptyEl.classList.remove("hidden");
        return;
    }
    emptyEl.classList.add("hidden");

    for (const item of people) {
        const user = item.user;
        const level = item.level;

        const btn = document.createElement("button");
        btn.className = "shared-person";
        btn.innerHTML = `
            <div class="shared-avatar">${escapeHtml(initials(user.username))}</div>
            <div class="shared-person-name">${escapeHtml(user.username)}</div>
            <div class="level-tag">${LEVEL_LABELS[level] || level}</div>
        `;
        btn.addEventListener("click", () => {
            // снимаем подсветку со всех, ставим на текущий
            for (const el of listEl.querySelectorAll(".shared-person")) {
                el.classList.remove("active");
            }
            btn.classList.add("active");
            openCalendar(user, level);
        });
        listEl.appendChild(btn);
    }
}

// --- Загрузка и показ задач выбранного календаря ---
async function openCalendar(user, level) {
    const card    = document.getElementById("viewerCard");
    const titleEl = document.getElementById("viewerTitle");
    const listEl  = document.getElementById("viewerList");
    const emptyEl = document.getElementById("viewerEmpty");

    card.classList.remove("hidden");
    titleEl.textContent = `Календарь: ${user.username}`;
    listEl.innerHTML = "";

    const res = await fetch(`/api/access/calendar/${user.id}`);
    if (!res.ok) {
        emptyEl.classList.remove("hidden");
        emptyEl.textContent = "Не удалось загрузить календарь";
        return;
    }
    const tasks = await res.json() || [];

    if (tasks.length === 0) {
        emptyEl.classList.remove("hidden");
        emptyEl.textContent = "Нет задач для показа";
        return;
    }
    emptyEl.classList.add("hidden");

    const isBusy = level === "free_busy";

    for (const task of tasks) {
        const date = formatDeadline(task.deadline);
        const time = task.start_time && task.end_time
            ? `${task.start_time}–${task.end_time}`
            : "";

        const meta = isBusy
            ? [date, time].filter(Boolean).join(" · ")
            : [date, PRIORITY_LABELS[task.priority], time].filter(Boolean).join(" · ");
        const el = document.createElement("div");
        el.className = `shared-task priority-${task.priority}` +
            (task.done ? " done" : "") +
            (isBusy ? " busy" : "");
        el.innerHTML = `
            <div class="shared-task-title">${escapeHtml(task.title)}</div>
            ${meta ? `<div class="shared-task-meta">${meta}</div>` : ""}
        `;
        listEl.appendChild(el);
    }
}

loadShared();