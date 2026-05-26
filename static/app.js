// Базовый адрес API. Пусто = тот же сервер, что отдал страницу.
const API = "/tasks";

let allTasks = [];      // все задачи с сервера
let currentFilter = "all";

// ---------- Загрузка задач ----------
async function loadTasks() {
    const res = await fetch(API);
    const data = await res.json();
    // Go возвращает null если задач нет — превращаем в пустой массив
    allTasks = data || [];
    render();
}

// ---------- Отрисовка списка ----------
function render() {
    const list = document.getElementById("taskList");
    list.innerHTML = "";

    // Применяем фильтр по приоритету
    const tasks = allTasks.filter(t =>
        currentFilter === "all" ? true : t.priority === Number(currentFilter)
    );

    if (tasks.length === 0) {
        list.innerHTML = `<div class="empty">Задач нет</div>`;
        return;
    }

    const priorityNames = { 1: "Низкий", 2: "Средний", 3: "Высокий" };

    for (const task of tasks) {
        const el = document.createElement("div");
        el.className = `task priority-${task.priority}` + (task.done ? " done" : "");

        el.innerHTML = `
            <input type="checkbox" class="task-checkbox" ${task.done ? "checked" : ""}>
            <div class="task-body">
                <div class="task-title">${escapeHtml(task.title)}</div>
                ${task.description ? `<div class="task-desc">${escapeHtml(task.description)}</div>` : ""}
                <div class="task-meta">
                    <span class="badge priority-${task.priority}">${priorityNames[task.priority]}</span>
                    ${task.deadline ? `<span>📅 ${task.deadline}</span>` : ""}
                </div>
            </div>
            <button class="delete-btn">×</button>
        `;

        // Чекбокс — отметить выполненной
        el.querySelector(".task-checkbox").addEventListener("change", () => {
            toggleDone(task);
        });

        // Кнопка удаления
        el.querySelector(".delete-btn").addEventListener("click", () => {
            deleteTask(task.id);
        });

        list.appendChild(el);
    }
}

// ---------- Создание задачи ----------
async function addTask() {
    const title = document.getElementById("title").value.trim();
    if (!title) {
        alert("Введите название задачи");
        return;
    }

    const newTask = {
        title: title,
        description: document.getElementById("description").value.trim(),
        priority: Number(document.getElementById("priority").value),
        deadline: document.getElementById("deadline").value,
    };

    await fetch(API, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(newTask),
    });

    // Очищаем форму
    document.getElementById("title").value = "";
    document.getElementById("description").value = "";
    document.getElementById("deadline").value = "";

    loadTasks();
}

// ---------- Переключение статуса ----------
async function toggleDone(task) {
    // PUT обновляет ВСЕ поля, поэтому отправляем задачу целиком,
    // меняя только done
    const updated = {
        title: task.title,
        description: task.description,
        priority: task.priority,
        deadline: task.deadline,
        done: !task.done,
    };

    await fetch(`${API}/${task.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(updated),
    });

    loadTasks();
}

// ---------- Удаление ----------
async function deleteTask(id) {
    await fetch(`${API}/${id}`, { method: "DELETE" });
    loadTasks();
}

// ---------- Защита от вставки HTML в текст ----------
function escapeHtml(str) {
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
}

// ---------- Навешиваем обработчики ----------
document.getElementById("addBtn").addEventListener("click", addTask);

document.querySelectorAll(".filter-btn").forEach(btn => {
    btn.addEventListener("click", () => {
        document.querySelector(".filter-btn.active").classList.remove("active");
        btn.classList.add("active");
        currentFilter = btn.dataset.filter;
        render();
    });
});

// ---------- Старт ----------
loadTasks();
