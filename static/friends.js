const API = "/api/friends";

function escapeHtml(str) {
    if (!str) return "";
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
}

function initials(name) {
    return name ? name.charAt(0) : "?";
}

function showMessage(el, text, kind) {
    el.textContent = text;
    el.className = `form-message ${kind}`;
}

// --- Отправка заявки ---
async function sendRequest() {
    const input = document.getElementById("addUsername");
    const username = input.value.trim();
    if (!username) return;

    const res = await fetch(`${API}/requests`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username }),
    });

    const msg = document.getElementById("addMessage");
    if (res.ok) {
        showMessage(msg, `Заявка отправлена пользователю «${username}»`, "success");
        input.value = "";
        loadRequests(); // на случай обратной заявки заявка из входящих исчезнет
        loadFriends();  // ...а человек сразу появится в друзьях
    } else {
        const text = await res.text();
        showMessage(msg, text.trim() || "Не удалось отправить заявку", "error");
    }
}

// --- Принять / отклонить ---
async function acceptRequest(requesterId) {
    await fetch(`${API}/accept`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ requester_id: requesterId }),
    });
    loadRequests();
    loadFriends();
}

async function declineRequest(requesterId) {
    await fetch(`${API}/decline`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ requester_id: requesterId }),
    });
    loadRequests();
}

// --- Входящие заявки ---
async function loadRequests() {
    const res = await fetch(`${API}/requests`);
    const requests = await res.json() || [];

    const listEl  = document.getElementById("requestsList");
    const emptyEl = document.getElementById("requestsEmpty");
    const countEl = document.getElementById("requestsCount");

    listEl.innerHTML = "";
    countEl.textContent = requests.length ? `(${requests.length})` : "";

    if (requests.length === 0) {
        emptyEl.classList.remove("hidden");
        return;
    }
    emptyEl.classList.add("hidden");

    for (const user of requests) {
        const row = document.createElement("div");
        row.className = "friend-row";
        row.innerHTML = `
            <div class="friend-avatar">${escapeHtml(initials(user.username))}</div>
            <div class="friend-name">${escapeHtml(user.username)}</div>
            <div class="friend-actions">
                <button class="btn-sm btn-accept">Принять</button>
                <button class="btn-sm btn-decline">Отклонить</button>
            </div>
        `;
        row.querySelector(".btn-accept").addEventListener("click", () => acceptRequest(user.id));
        row.querySelector(".btn-decline").addEventListener("click", () => declineRequest(user.id));
        listEl.appendChild(row);
    }
}

// --- Доступ к календарю ---
// Карта grantee_id -> уровень ("view" | "free_busy"). Кого нет в карте — доступа нет.
let accessMap = {};

async function loadAccessMap() {
    const res = await fetch("/api/access/my-grants");
    const grants = await res.json() || [];
    accessMap = {};
    for (const g of grants) {
        accessMap[g.user.id] = g.level;
    }
}

// level: "none" | "free_busy" | "view"
async function setAccess(granteeId, level) {
    if (level === "none") {
        await fetch("/api/access/revoke", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ grantee_id: granteeId }),
        });
    } else {
        await fetch("/api/access/grant", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ grantee_id: granteeId, level }),
        });
    }
    await loadFriends(); // перерисуем, чтобы подсветка обновилась
}

// --- Список друзей ---
async function loadFriends() {
    await loadAccessMap(); // сначала узнаём текущие уровни доступа

    const res = await fetch(API);
    const friends = await res.json() || [];

    const listEl  = document.getElementById("friendsList");
    const emptyEl = document.getElementById("friendsEmpty");
    const countEl = document.getElementById("friendsCount");

    listEl.innerHTML = "";
    countEl.textContent = friends.length ? `(${friends.length})` : "";

    if (friends.length === 0) {
        emptyEl.classList.remove("hidden");
        return;
    }
    emptyEl.classList.add("hidden");

    for (const user of friends) {
        const current = accessMap[user.id] || "none"; // что сейчас выдано

        const row = document.createElement("div");
        row.className = "friend-row";
        row.innerHTML = `
            <div class="friend-avatar">${escapeHtml(initials(user.username))}</div>
            <div class="friend-name">${escapeHtml(user.username)}</div>
            <div class="access-toggle">
                <button class="access-btn ${current === "none" ? "active" : ""}" data-level="none">Нет</button>
                <button class="access-btn ${current === "free_busy" ? "active" : ""}" data-level="free_busy">Занятость</button>
                <button class="access-btn ${current === "view" ? "active" : ""}" data-level="view">Детали</button>
            </div>
        `;

        // Вешаем обработчик на каждую из трёх кнопок.
        for (const btn of row.querySelectorAll(".access-btn")) {
            btn.addEventListener("click", () => setAccess(user.id, btn.dataset.level));
        }

        listEl.appendChild(row);
    }
}

// --- Инициализация ---
document.getElementById("addBtn").addEventListener("click", sendRequest);
document.getElementById("addUsername").addEventListener("keydown", (e) => {
    if (e.key === "Enter") sendRequest();
});

loadRequests();
loadFriends();