const API = "/tasks";

async function loadStats() {
    const res = await fetch(API);
    const tasks = await res.json() || [];

    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const total   = tasks.length;
    const done    = tasks.filter(t => t.done).length;
    const active  = total - done;
    const overdue = tasks.filter(t =>
        !t.done && t.deadline && new Date(t.deadline) < today
    ).length;
    const noDate  = tasks.filter(t => !t.deadline).length;

    const pct = total === 0 ? 0 : Math.round((done / total) * 100);

    // Карточки
    document.getElementById("statTotal").textContent  = total;
    document.getElementById("statDone").textContent   = done;
    document.getElementById("statActive").textContent = active;
    document.getElementById("statOverdue").textContent = overdue;

    // Прогресс-бар
    document.getElementById("progressFill").style.width = pct + "%";
    document.getElementById("progressPct").textContent  = pct + "%";

    // По приоритетам
    const priorities = [
        { key: 1, label: "Низкий",   cls: "low" },
        { key: 2, label: "Средний",  cls: "medium" },
        { key: 3, label: "Высокий",  cls: "high" },
    ];
    const rows = document.getElementById("priorityRows");
    rows.innerHTML = "";

    for (const p of priorities) {
        const count = tasks.filter(t => t.priority === p.key).length;
        const pctP  = total === 0 ? 0 : Math.round((count / total) * 100);

        const row = document.createElement("div");
        row.className = "priority-row";
        row.innerHTML = `
            <span class="priority-row-label">${p.label}</span>
            <div class="priority-track">
                <div class="priority-fill ${p.cls}" style="width: ${pctP}%"></div>
            </div>
            <span class="priority-count">${count}</span>
        `;
        rows.appendChild(row);
    }

    // Без даты
    document.getElementById("statNoDate").textContent =
        noDate === 0 ? "Нет задач без даты" : `${noDate} задач без назначенной даты`;
}

loadStats();
