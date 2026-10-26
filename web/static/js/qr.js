document.querySelector('[data-print-stands]')?.addEventListener('click', async event => {
  event.currentTarget.disabled = true;
  try {
    await document.fonts.ready;
    await Promise.all([...document.querySelectorAll('.qr-sheets img')].map(image => image.decode()));
    window.print();
  } finally { event.currentTarget.disabled = false; }
});
