const API_BASE = 'http://176.57.215.221:8080/';
const registerForm = document.querySelector('.registartionForm');

// Функция для обработки регистрации
registerForm.addEventListener('submit', async (event) => {
    event.preventDefault(); // Предотвращаем отправку формы по умолчанию

    const fullName = document.querySelector('.registerName').value;
    const email = document.querySelector('.registerEmail').value;
    const phone = document.querySelector('.registerPhone').value;
    const password = document.querySelector('.registerPassword').value;

    // Формируем тело запроса согласно описанию API для регистрации
    const requestBody = {
      user: {
        register: {
          full_name: fullName,
          email: email,
          phone: phone,
          password: password,
          type: "client"
        }
      }
    };

    try {
      const response = await fetch(`${API_BASE}v1/register/client`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
      });

      const data = await response.json(); // Разбираем JSON-ответ

      if (response.ok) {
        // Сохраняем токен в localStorage (можно сразу перенаправить на вход)
        localStorage.setItem('authToken', data.user.token);

        // Перенаправляем на страницу для авторизованных пользователей (замените 'dashboard.html' на ваш URL)
        window.location.href = 'companyProfile.html';

      } else {
      }

    } catch (error) {
      console.error('Ошибка:', error);
    }
  });

