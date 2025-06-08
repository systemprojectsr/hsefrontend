const API_BASE = 'http://176.57.215.221:8080/';
const registerForm = document.querySelector('.registrationForm');
const registerBtn = document.querySelector('.registerSubmit');

async function hashPassword(password) {
  const msgUint8 = new TextEncoder().encode(password);
  const hashBuffer = await crypto.subtle.digest('SHA-256', msgUint8);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

// Функция для обработки регистрации
async function handleRegister(event) {
    event.preventDefault(); // Предотвращаем отправку формы по умолчанию

    const fullName = document.querySelector('.registerName').value;
    const email = document.querySelector('.registerEmail').value;
    const phone = document.querySelector('.registerPhone').value;
    const password = document.querySelector('.registerPassword').value;
    const passwordHash = await hashPassword(password);

    // Формируем тело запроса согласно описанию API для регистрации
    const requestBody = {
      user: {
        register: {
          full_name: fullName,
          email: email,
          phone: phone,
          password: passwordHash,
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
  }

if (registerForm) {
  registerForm.addEventListener('submit', handleRegister);
}
if (registerBtn) {
  registerBtn.addEventListener('click', handleRegister);
}

