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

// Загрузка карточек на главной странице происходит в search.js
  