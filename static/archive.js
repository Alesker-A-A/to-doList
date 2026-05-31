const API = "/tasks";

const MONTH_GENITIVE = [
    "января", "февраля", "марта", "апреля", "мая", "июня",
    "июля", "августа", "сентября", "октября", "ноября", "декабря"
];

const PRIORITY_LABELS = { 1: "Низкий", 2: "Средний", 3: "Высокий" };

function formatDeadline(ymd) {
    if (!ymd) return null;
    const [y, m, d] = ymd.split("-").map(Number);
    return `${d} ${MONTH_GENITIVE[m - 1]} ${y}`;
}

function groupLabel(ymd) {
    if (!ymd) return "Без даты";
    const [y, m] = ymd.split("-").map(Number);
    const now = new Date();
    if (y === now.getFullYear() && m === now.getMonth() + 1) return "Этот месяц";
    const months = [
        "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
        "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"
    ];
    return y === now.getFullYear()
        ? months[m - 1]
        : `${months[m - 1]} ${y}`;
}

function escapeHtml(str) {
    if (!str) return "";
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
}

async function deleteTask(id) {
    await fetch(`${API}/${id}`, { method: "DELETE" });
    renderArchive();
}

async function renderArchive() {
    const res = await fetch(API);
    const tasks = await res.json() || [];

    const done = tasks
        .filter(t => t.done)
        .sort((a, b) => {
            // Сортировка: задачи с датой — по убыванию даты, без даты — в конец
            if (!a.deadline && !b.deadline) return 0;
            if (!a.deadline) return 1;
            if (!b.deadline) return -1;
            return b.deadline.localeCompare(a.deadline);
        });

    const countEl = document.getElementById("archiveCount");
    const emptyEl = document.getElementById("archiveEmpty");
    const listEl  = document.getElementById("archiveList");

    listEl.innerHTML = "";

    if (done.length === 0) {
        countEl.textContent = "";
        emptyEl.classList.remove("hidden");
        return;
    }

    emptyEl.classList.add("hidden");
    countEl.textContent = `${done.length} выполнено`;

    // Группируем по месяцу дедлайна
    const groups = new Map();
    for (const task of done) {
        const key = task.deadline
            ? task.deadline.slice(0, 7)  // "YYYY-MM"
            : "no-date";
        if (!groups.has(key)) groups.set(key, []);
        groups.get(key).push(task);
    }

    for (const [key, groupTasks] of groups) {
        const sampleDeadline = key === "no-date" ? null : key + "-01";
        const label = groupLabel(sampleDeadline);

        const group = document.createElement("div");
        group.className = "archive-group";
        group.innerHTML = `<div class="archive-group-title">${label}</div>`;

        for (const task of groupTasks) {
            const el = document.createElement("div");
            el.className = `archive-task priority-${task.priority}`;

            const deadlineStr = formatDeadline(task.deadline);
            const timeStr = task.start_time && task.end_time
                ? ` · ${task.start_time}–${task.end_time}`
                : "";
            const priorityStr = PRIORITY_LABELS[task.priority] || "";
            const meta = [deadlineStr, priorityStr + timeStr]
                .filter(Boolean).join(" · ");

            el.innerHTML = `
                <span class="archive-task-check">✓</span>
                <div class="archive-task-body">
                    <div class="archive-task-title" title="${escapeHtml(task.title)}">${escapeHtml(task.title)}</div>
                    ${meta ? `<div class="archive-task-meta">${meta}</div>` : ""}
                </div>
                <button class="archive-task-delete" title="Удалить">×</button>
            `;

            el.querySelector(".archive-task-delete").addEventListener("click", () => {
                deleteTask(task.id);
            });

            group.appendChild(el);
        }

        listEl.appendChild(group);
    }
}

renderArchive();
