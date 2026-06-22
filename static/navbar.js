(function () {
  // Определяем активную страницу по pathname
  const path = window.location.pathname;

function isActive(href) {
    if (href === "/app" && path === "/app") return true;
    if (href === "/stats" && path === "/stats") return true;
    if (href === "/archive" && path === "/archive") return true;
    if (href === "/friends" && path === "/friends") return true;
    if (href === "/shared" && path === "/shared") return true;
    return false;
  }

  const nav = document.createElement("header");
  nav.className = "navbar";
  nav.innerHTML = `
    <div class="navbar-left">
      <a href="/app" class="navbar-brand">
        <span class="brand-mark"></span>
        <span class="brand-name">Task Calendar</span>
      </a>
      <button class="btn-today" id="todayBtn">Сегодня</button>
    </div>
    <nav class="navbar-links">
      <a href="/app"     class="nav-link ${isActive("/app")     ? "active" : ""}">Календарь</a>
      <a href="/stats"   class="nav-link ${isActive("/stats")   ? "active" : ""}">Статистика</a>
      <a href="/archive" class="nav-link ${isActive("/archive") ? "active" : ""}">Архив</a>
      <a href="/friends" class="nav-link ${isActive("/friends") ? "active" : ""}">Друзья</a>
      <a href="/shared" class="nav-link ${isActive("/shared") ? "active" : ""}">Общие</a>
      <span class="navbar-user" id="navbarUser"></span>
      <button class="nav-logout" id="logoutBtn">Выход</button>
    </nav>
  `;

  // Вставляем навбар первым элементом в body
  document.body.insertBefore(nav, document.body.firstChild);

  // Кнопка "Сегодня" — глобальное событие, страницы подписываются сами
  document.getElementById("todayBtn").addEventListener("click", function () {
    window.dispatchEvent(new CustomEvent("navbar:today"));
  });

  // --- Имя пользователя ---
  // Данные кладёт auth-guard.js в window.currentUser. Ответ /api/me
  // асинхронный, поэтому: если данные уже есть — показываем сразу,
  // иначе ждём событие auth:ready.
  const userEl = document.getElementById("navbarUser");

  function showUser(user) {
    if (user && user.username) {
      userEl.textContent = user.username;
    }
  }

  if (window.currentUser) {
    showUser(window.currentUser);
  } else {
    window.addEventListener("auth:ready", function (e) {
      showUser(e.detail);
    });
  }

  // --- Выход ---
  document.getElementById("logoutBtn").addEventListener("click", async function () {
    try {
      await fetch("/api/logout", { method: "POST" });
    } catch (err) {
      // Даже если запрос не прошёл, всё равно уводим на логин —
      // кука на клиенте могла протухнуть, серверная сессия не критична.
    }
    window.location.href = "/login";
  });
})();
