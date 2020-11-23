fetch('/now').then(res => res.json())
  .then(data => {
    document.getElementById('current-time').innerText = data.time;
  })
  .catch(() => {
    console.log('Error occurred.')
  });