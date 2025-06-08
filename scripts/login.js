
const loginForm = document.querySelector('.loginFform'); 


// Функция для обработки входа
loginForm.addEventListener('submit', async (event) => {
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
      const response = await fetch('v1/login', { 
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
      });

      const data = await response.json(); // Разбираем JSON-ответ

      if (response.ok) { //Вход прошел успешно
        // Сохраняем токен в localStorage
        localStorage.setItem('authToken', data.user.token);

        window.location.href = 'companyProfile.html';

      } else {
        alert('Ошибка!!!!')
      }

    } catch (error) {
      console.error('Ошибка:', error);
    }
  });

  