const API_BASE = 'http://176.57.215.221:8080/';
const addForm = document.querySelector('.main-container');
if (addForm) {
  const publishBtn = document.querySelector('.button-group .button');
  publishBtn.addEventListener('click', async (e) => {
    e.preventDefault();
    const title = document.querySelector('.cardTitle').value;
    const description = document.querySelector('.textarea').value;
    const token = localStorage.getItem('authToken');
    const body = {
      user: {
        login: {
          token: token
        }
      },
      card: {
        title: title,
        description: description
      }
    };
    try {
      const response = await fetch(`${API_BASE}v1/account/card/create`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(body)
      });
      if (response.ok) {
        window.location.href = 'companyProfile.html';
      } else {
        console.error('Failed to create card');
      }
    } catch (err) {
      console.error('Ошибка:', err);
    }
  });
}

