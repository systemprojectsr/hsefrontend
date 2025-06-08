const API_BASE = 'http://176.57.215.221:8080/';
const cardsContainer = document.getElementById('cardsContainer');
if (cardsContainer) {
  const token = localStorage.getItem('authToken');
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
            div.innerHTML = `<p>${card.title}</p><p>${card.description}</p>`;
            cardsContainer.appendChild(div);
          });
        }
      })
      .catch(err => console.error('Card list error:', err));
  }
}

