(function () {
  // Определяем активную страницу по pathname
  const path = window.location.pathname;

  function isActive(href) {
    if (href === "/app" && path === "/app") return true;
    if (href === "/stats" && path === "/stats") return true;
    if (href === "/archive" && path === "/archive") return true;
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
    </nav>
  `;

  // Вставляем навбар первым элементом в body
  document.body.insertBefore(nav, document.body.firstChild);

  // Экспортируем кнопку "Сегодня" как глобальное событие,
  // чтобы каждая страница могла подписаться на неё
  document.getElementById("todayBtn").addEventListener("click", function () {
    window.dispatchEvent(new CustomEvent("navbar:today"));
  });
})();
