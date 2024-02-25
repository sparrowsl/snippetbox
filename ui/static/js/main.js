"use strict"

document.getElementById("year").outerHTML = new Date().getFullYear();
const navLinks = document.querySelectorAll("nav a");
const datesToFormat = document.querySelectorAll("[data-custom-date]");

for (const link of navLinks) {
  if (link.getAttribute('href') === window.location.pathname) {
    link.classList.add("live");
    break;
  }
}

for (const element of datesToFormat) {
  const newDate = new Date(element.textContent)
  const formattedDate = new Intl.DateTimeFormat("en-US", { dateStyle: "long", timeStyle: "short" }).format(newDate)

  element.textContent = formattedDate
}
