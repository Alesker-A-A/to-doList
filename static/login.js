// Текущий режим: "login" или "register"
let mode = "login";

const tabLogin    = document.getElementById("tabLogin");
const tabRegister = document.getElementById("tabRegister");
const form        = document.getElementById("authForm");
const usernameEl  = document.getElementById("username");
const passwordEl  = document.getElementById("password");
const errorEl     = document.getElementById("authError");
const hintEl      = document.getElementById("authHint");
const submitBtn   = document.getElementById("submitBtn");

// Переключение режима (вход / регистрация)
function setMode(newMode) {
    mode = newMode;
    errorEl.textContent = "";

    if (mode === "login") {
        tabLogin.classList.add("active");
        tabRegister.classList.remove("active");
        submitBtn.textContent = "Войти";
        passwordEl.autocomplete = "current-password";
        hintEl.textContent = "";
    } else {
        tabRegister.classList.add("active");
        tabLogin.classList.remove("active");
        submitBtn.textContent = "Создать аккаунт";
        passwordEl.autocomplete = "new-password";
        hintEl.textContent = "Имя — от 3 символов, пароль — от 6.";
    }
}

tabLogin.addEventListener("click", () => setMode("login"));
tabRegister.addEventListener("click", () => setMode("register"));

// Отправка формы
form.addEventListener("submit", async (e) => {
    e.preventDefault();
    errorEl.textContent = "";

    const username = usernameEl.value.trim();
    const password = passwordEl.value;

    if (!username || !password) {
        errorEl.textContent = "Заполните оба поля.";
        return;
    }

    const endpoint = mode === "login" ? "/api/login" : "/api/register";

    submitBtn.disabled = true;
    try {
        const res = await fetch(endpoint, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password })
        });

        if (res.ok) {
            // Успех — сервер уже поставил куку сессии. Переходим в приложение.
            window.location.href = "/app";
            return;
        }

        // Ошибка — сервер прислал текстовое сообщение, покажем его.
        const text = await res.text();
        errorEl.textContent = text.trim() || "Что-то пошло не так. Попробуйте снова.";
    } catch (err) {
        errorEl.textContent = "Не удалось связаться с сервером.";
    } finally {
        submitBtn.disabled = false;
    }
});

// Стартовый режим
setMode("login");
