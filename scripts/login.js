
const API_BASE = 'http://176.57.215.221:8080/';
const loginForm = document.querySelector('.loginFform');
const loginBtn = document.querySelector('.loginButton');

async function handleLogin(event) {
    event.preventDefault(); // Предотвращаем отправку формы по умолчанию

    const email = document.querySelector('.loginEmail').value;
    const password = document.querySelector('.loginPassword').value;

    // Формируем тело запроса согласно описанию API для входа
    const requestBody = {
      user: {
        login: {
          email: email,
          password: password
        }
      }
    };

    try {
      const response = await fetch(`${API_BASE}v1/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
      });

      const data = await response.json(); // Разбираем JSON-ответ

      if (response.ok) { // Вход прошел успешно
        // Сохраняем токен в localStorage
        localStorage.setItem('authToken', data.user.token);

        window.location.href = 'companyProfile.html';

      } else {
        alert('Ошибка!!!!')
      }

    } catch (error) {
      console.error('Ошибка:', error);
    }
  }

if (loginForm) {
  loginForm.addEventListener('submit', handleLogin);
}
if (loginBtn) {
  loginBtn.addEventListener('click', handleLogin);
}
