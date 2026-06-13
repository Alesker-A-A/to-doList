// Страж авторизации. Подключается в <head> защищённых страниц ДО отрисовки
// контента. Спрашивает у сервера "кто я" через /api/me:
//   - если сессия валидна — сохраняем пользователя и показываем страницу
//   - если нет (401) — редирект на /login, контент не показываем
//
// Чтобы контент не мелькал перед редиректом, прячем <html> до проверки
// и показываем только после успешного ответа.

(function () {
    // Прячем страницу немедленно (до отрисовки).
    document.documentElement.style.visibility = "hidden";
})();

// Глобально доступный текущий пользователь — другие скрипты (навбар) могут
// прочитать window.currentUser после события "auth:ready".
window.currentUser = null;

async function checkAuth() {
    try {
        const res = await fetch("/api/me");

        if (!res.ok) {
            // Не авторизован — на страницу входа.
            window.location.href = "/login";
            return;
        }

        window.currentUser = await res.json();

        // Показываем страницу и сообщаем остальным скриптам, что юзер готов.
        document.documentElement.style.visibility = "";
        window.dispatchEvent(new CustomEvent("auth:ready", {
            detail: window.currentUser
        }));
    } catch (err) {
        // Сервер недоступен — безопаснее увести на логин.
        window.location.href = "/login";
    }
}

checkAuth();
