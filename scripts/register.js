const API_BASE = 'http://176.57.215.221:8080/';
const registerForm = document.querySelector('.registrationForm');
const registerBtn = document.querySelector('.registerSubmit');

// Функция для обработки регистрации
async function handleRegister(event) {
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
        const { token, type } = data.user;
        // Сохраняем токен и тип аккаунта
        localStorage.setItem('authToken', token);
        localStorage.setItem('accountType', type);
        document.cookie = `authToken=${token}; path=/`;

        // Перенаправляем на страницу для авторизованных пользователей
        if (type === 'company') {
          window.location.href = 'companyProfile.html';
        } else {
          window.location.href = 'customerProfile.html';
        }

      } else {
      }

    } catch (error) {
      console.error('Ошибка:', error);
    }
  }

if (registerForm) {
  registerForm.addEventListener('submit', handleRegister);
}
if (registerBtn) {
  registerBtn.addEventListener('click', handleRegister);
}

