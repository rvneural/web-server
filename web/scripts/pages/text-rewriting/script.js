const intupArea = document.getElementById('inputArea')
const outputArea = document.getElementById('outputArea')

const processButton = document.getElementById('rewriteButton')

const resetButton = document.getElementById("resetButton") // Кнопка Reset
const copyTextButton = document.getElementById("copyText") // Кнопка copyText
const saveFileButton = document.getElementById("saveFile") // Кнопка saveFile

var currentText = ""

window.onload = function() {
    document.getElementById('rewritePage').style.backgroundColor = "#0c087466"

    // Check if currentRewriteText is in localStorage
    if (localStorage.getItem('rewriteText')!== null) {

        //Unlock outputArea
        outputArea.removeAttribute('readonly')


        currentText = localStorage.getItem('rewriteText')
        outputArea.value = currentText

        unlockElements()
    }

    // Check if promtRewriteArea is in localStorage
    if (localStorage.getItem('promtRewriteArea')!== null) {
        intupArea.value = localStorage.getItem('promtRewriteArea')
    }
}

// Обработчик нажатия кнопки "Переписать"
processButton.addEventListener('click', async () => {

    intupArea.classList.remove('need-enter')

  if(intupArea.value.trim() === ''){
    alert('Введите текст для переписывания')
    intupArea.classList.add('need-enter')
    intupArea.value = ''
    setTimeout(() => {
        intupArea.classList.remove('need-enter')
    }, 3000)
    return
  }

  // Update promtRewriteArea в localStorage
  localStorage.setItem('promtRewriteArea', intupArea.value)

  // Подготовка
  lockElements()
  outputArea.value = 'Переписываем текст...'
  outputArea.classList.add("loader");

  // Переписывание текста
  await rewriteText(intupArea.value.trim())

  // Финал
  outputArea.classList.remove("loader");
  outputArea.value = currentText

  // Update rewriteText в localStorage
  localStorage.setItem('rewriteText', outputArea.value)

  unlockElements()
})

// Функция переписывания текста, возвращает строку
async function rewriteText(userText){
    const formData = new FormData();
    formData.append('text', userText)
    console.log(formData)

    // Отправляет запрос на Web Server с данными из веб-формы
    const resp = await fetch('/rewriteFromWeb', {
        method: 'POST',
        body: formData,
    })

    // Расшифровываем результат в JSON
    const data = await resp.json();
    if ((data.error) && (data.error !== "")){
        currentText = "Ошибка: " + data.error + "\n\nДетали: " + data.details
    } else {
        currentText = data.newText
    }
}

function lockElements(){
    intupArea.setAttribute('readonly', '')
    outputArea.setAttribute('readonly', '')
    processButton.setAttribute('disabled', '')
    resetButton.setAttribute('disabled', '')
    copyTextButton.setAttribute('disabled', '')
    saveFileButton.setAttribute('disabled', '')
}

function unlockElements(){
    intupArea.removeAttribute('readonly')
    outputArea.removeAttribute('readonly')
    processButton.removeAttribute('disabled')
    resetButton.removeAttribute('disabled')
    copyTextButton.removeAttribute('disabled')
    saveFileButton.removeAttribute('disabled')
}

// Обработчик нажатия кнопки "Сбросить"
resetButton.addEventListener('click', () => {
  outputArea.value = currentText

  // Update currenRewriteText в localStorage
  localStorage.setItem('rewriteText', currentText)
})

// Обработчик кнопки "Копировать текст"
copyTextButton.addEventListener('click', async () => {
   // Убираем из текста двойные переносы строк
   var text = outputArea.value.trim().replaceAll('\n\n', '\n')
   console.log("Text for copying: ", text)

   // Update rewriteText в localStorage
   localStorage.setItem('rewriteText', text)

   // Пытаемся отправить текст в clipboard
   try{
       // Доступно только по протоколу HTTPS
       await navigator.clipboard.writeText(text)
       console.log('Текст скопирован')
       copyTextButton.innerText = "Скопировано"
       setTimeout(() => {
           copyTextButton.innerText = "Скопировать"
       }, 2000)
   } catch (err){
       // Ошибка вызывается в частности в том случае
       // Если сервер читает по HTTP
       try{
           // Пытаемся выделить текст и напрямую вызвать команду ОС 
           // на копирование выделенного фрагмента
           outputArea.focus()
           outputArea.select()
           document.execCommand('copy')
           copyTextButton.innerText = "Скопировано"
           setTimeout(() => {
               copyTextButton.innerText = "Скопировать"
           }, 2000)
       } catch (err) {
           console.log('Ошибка копирования текста:', err)
       }
       console.log('Ошибка копирования текста:', err)
   }
})

// Обработчик кнопки "Сохранить текст в файл"
saveFileButton.addEventListener('click', async () => {
    // Создаем BLOB и убираем двойные переносы строк в тексте
    const blob = new Blob([outputArea.value.trim().replaceAll('\n\n', '\n')], {type: 'text/plain'})

    // Update rewriteText in localStorage
    localStorage.setItem('rewriteText', outputArea.value)

    // Создаем документ и ссылку на скачиваение
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);

    // Указываем название файла
    link.download = 'rewritedText.txt';

    // Отправляем файл клиенту
    link.click();
    URL.revokeObjectURL(link.href);
})