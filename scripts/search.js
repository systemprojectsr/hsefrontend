const API_BASE_SEARCH = 'http://176.57.215.221:8070';
const form = document.querySelector('.searchForm');
const results = document.getElementById('results');

async function renderCards(cards) {
  results.innerHTML = '';
  cards.forEach(item => {
    const card = document.createElement('div');
    card.className = 'card';
    card.innerHTML = `<h3>${item.name}</h3><p>${item.description}</p><p>${item.price} ₽</p>`;
    results.appendChild(card);
  });
}

async function fetchAllCards() {
  const resp = await fetch(`${API_BASE_SEARCH}/search`);
  return resp.json();
}

if (form && results) {
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const query = form.querySelector('input[type="text"]').value.toLowerCase();
    try {
      const data = await fetchAllCards();
      const filtered = data.filter(item => item.name.toLowerCase().includes(query));
      renderCards(filtered);
    } catch(err) {
      console.error('Search error:', err);
    }
  });

  // начальная загрузка всех карточек
  fetchAllCards().then(renderCards).catch(err => console.error('Load error:', err));
}

