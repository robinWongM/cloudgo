fetch('/now').then(res => res.json())
  .then(data => {
    document.getElementById('current-time').innerText = `当前时间：${data.time}`;
  })
  .catch(() => {
    console.log('Error occurred.')
  });