function closePopup() {
  const popups = document.getElementsByClassName("error-popup");
  if (!!popups?.length) popups[popups.length - 1].remove();
}

document.addEventListener("DOMContentLoaded", (event) => {
  console.log("DOMContentLoaded");
  console.log(event);
  document.body.addEventListener('htmx:beforeSwap', function (evt) {
    console.log("beforeSwap");
    console.log(evt);
    console.log(evt.detail);
    if (evt.detail.xhr.status === 404) {
      evt.detail.isError = false;
      const popup = document.createElement("div");
      popup.classList.add("error-popup");
      popup.innerHTML = evt.detail.serverResponse;
      document.body.append(popup);
    }
    if (evt.detail.xhr.status === 422) {
      evt.detail.shouldSwap = true;
      evt.detail.isError = false;
    }
    if (evt.detail.xhr.status === 500) {
      evt.detail.isError = false;

      let popup = document.createElement("div");
      popup.classList.add("error-popup");
      popup.innerHTML = `
        <dialog open>
          <h3>Errore</h3>
          <p><span style="color:darkred;">[500]: </span>Errore generico del server.</p>
          <div class="popup-buttons">
              <button onclick="closePopup()">ok</button>
          </div>
        </dialog>
      `.trim();
      document.body.append(popup);
    }
  })
})