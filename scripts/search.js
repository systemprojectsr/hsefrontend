const API_BASE_SEARCH = 'http://176.57.215.221:8070';
const form = document.querySelector('.searchForm');
const results = document.getElementById('results');
if (form && results) {
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const query = form.querySelector('input[type="text"]').value;
    const url = `${API_BASE_SEARCH}/search?` + new URLSearchParams({ q: query });
    try {
      const response = await fetch(url);
      const data = await response.json();
      results.innerHTML = '';
      data.forEach(item => {
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `<p>${item.name}</p><p>${item.price} ₽</p><p>${item.location}</p>`;
        results.appendChild(card);
      });
    } catch(err) {
      console.error('Search error:', err);
    }
  });
}

