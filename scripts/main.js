"use strict";

let btn = document.querySelector(".btnLogin");
let modal = document.querySelector(".popUpLogin");
let btnClose = document.querySelector(".close-modal");
let overlay = document.querySelector(".overlay");
let regBtn = document.querySelector(".regBtn");

btn.addEventListener('click', () => {
    overlay.classList.remove("hidden");
    modal.classList.remove("hidden");
});

if (modal.classList.contains("hidden")) {
    btnClose.addEventListener('click', () => {
        modal.classList.add("hidden");
        overlay.classList.add("hidden");
    })
    overlay.addEventListener('click', () => {
        modal.classList.add("hidden");
        overlay.classList.add("hidden");
    })

}

// Загрузка карточек на главной странице
const API_BASE_SEARCH = 'http://176.57.215.221:8070';
const resultsContainer = document.getElementById('results');

async function loadCards() {
  if (!resultsContainer) return;
  try {
    const response = await fetch(`${API_BASE_SEARCH}/search`);
    const data = await response.json();
    resultsContainer.innerHTML = '';
    data.forEach(item => {
      const card = document.createElement('div');
      card.className = 'card';
      card.innerHTML = `
        <img src="img/uborka.png" alt="Услуга">
        <div class="card-content">
          <p>${item.name}</p>
          <p>${item.description}</p>
          <p>${item.price} ₽</p>
        </div>`;
      resultsContainer.appendChild(card);
    });
  } catch (err) {
    console.error('Load cards error:', err);
  }
}

document.addEventListener('DOMContentLoaded', loadCards);
  