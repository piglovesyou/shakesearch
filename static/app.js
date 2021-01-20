(() => {
  const DOMUtils = {
    removeChildren(n) {
      while (n.firstChild) n.removeChild(n.firstChild);
    },
  };

  const Controller = {
    search: (ev) => {
      ev.preventDefault();
      const formData = new FormData(document.forms.form);
      const data = Object.fromEntries(formData);
      fetch(`/search?q=${data.query}`).then((response) => {
        const { status } = response;
        if (200 <= status && status <= 300) {
          response.json().then(Controller.updateTable);
        } else {
          response.text().then(Controller.updateMessage);
        }
      });
    },

    clearContent() {
      const messageEl = document.getElementById('message');
      DOMUtils.removeChildren(messageEl);
      const tBodyEl = document.getElementById('table-body');
      DOMUtils.removeChildren(tBodyEl);
      return [messageEl, tBodyEl];
    },

    updateTable: (results) => {
      const rows = document.createDocumentFragment();
      for (let result of results) {
        const trEl = document.createElement('tr');
        const tdEl = document.createElement('td');

        tdEl.innerHTML = `<p>${result.Highlights.map((h) => {
          return `<span>${h.SubContent.slice(0, h.Start)}<b>${h.SubContent.slice(
            h.Start,
            h.End,
          )}</b>${h.SubContent.slice(h.End)}</span>`;
        })}</p><h6>${result.Name}</h6>`;
        trEl.appendChild(tdEl);
        rows.appendChild(trEl);
      }
      const [messageEl, tBodyEl] = Controller.clearContent();

      if (results.length) tBodyEl.appendChild(rows);
      else messageEl.textContent = 'No result.';
    },

    updateMessage(message) {
      const [messageEl, tBodyEl] = Controller.clearContent();
      messageEl.textContent = message;
    },
  };

  document.addEventListener('DOMContentLoaded', () => {
    const form = document.forms.form;
    form.addEventListener('submit', Controller.search);
    // XXX: The only purpose users visit the site is
    // to start searching. Help focus and start typing.
    form.elements.query.focus();
  });
})();
