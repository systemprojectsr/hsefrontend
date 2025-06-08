const API_BASE = 'http://176.57.215.221:8080/';
const cardsContainer = document.getElementById('cardsContainer');

function getToken() {
  const match = document.cookie.match(/(?:^|; )authToken=([^;]+)/);
  return match ? match[1] : localStorage.getItem('authToken');
}

if (cardsContainer) {
  const token = getToken();
  if (token) {
    fetch(`${API_BASE}v1/account/card/list`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user: { login: { token: token } }
      })
    })
      .then(res => res.json())
      .then(data => {
        cardsContainer.innerHTML = '';
        if (data.cards) {
          data.cards.forEach(card => {
            const div = document.createElement('div');
            div.className = 'card';
            div.innerHTML = `
              <img src="img/moika.png" alt="Услуга">
              <div class="card-content">
                <p>${card.title}</p>
                <p>${card.description}</p>
              </div>`;
            cardsContainer.appendChild(div);
          });
        }
      })
      .catch(err => console.error('Card list error:', err));
  }
}

