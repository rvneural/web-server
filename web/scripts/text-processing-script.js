const intupArea = document.getElementById('inputArea')
const promtArea = document.getElementById('promtArea')
const outputArea = document.getElementById('outputArea')

const promtSelect = document.getElementById('promtSelect')

const processButton = document.getElementById('rewriteButton')

const resetButton = document.getElementById("resetButton") // Кнопка Reset
const copyTextButton = document.getElementById("copyText") // Кнопка copyText
const saveFileButton = document.getElementById("saveFile") // Кнопка saveFile

var currentText = ""

var progress = false

window.onbeforeunload = function() {
    if(progress){
        return "Текст все еще обрабатывается. Вы уверены, что хотите закрыть страницу?"
    }
}

window.onload = function() {
    document.getElementById('textPage').style.backgroundColor = "#494E56"
}

promtSelect.addEventListener('change', () => {
    console.log(promtSelect.value)
    if (promtSelect.value !== '0') {
        promtArea.value = "{{ " + promtSelect.value + " }}"
    } else {
        promtArea.value = ""
    }
})

// Обработчик нажатия кнопки "Переписать"
processButton.addEventListener('click', async () => {
    promtArea.classList.remove('need-enter')
    intupArea.classList.remove('need-enter')
  if(intupArea.value.trim() === ''){
    alert('Введите текст для обработки')
    // add need-enter class to inputArea
    intupArea.classList.add('need-enter')
    setTimeout(() => {
        intupArea.classList.remove('need-enter')
    }, 3000)
    intupArea.value = ''
    return
  }

  if(promtArea.value.trim() === ''){
    alert('Введите запрос для нейросети')
    promtArea.classList.add('need-enter')
    promtArea.value = ''
    setTimeout(() => {
        promtArea.classList.remove('need-enter')
    }, 3000)
    return
  }

  // Подготовка
  lockElements()
  progress = true

  outputArea.value = 'Обрабатываем текст...'
  outputArea.classList.add("loader");

  // Переписывание текста
  await rewriteText(intupArea.value.trim(), promtArea.value.trim())

  // Финал
  outputArea.classList.remove("loader");
  outputArea.value = currentText

  unlockElements()
  progress = false
})

// Функция переписывания текста
async function rewriteText(userText, userPromt){
    const formData = new FormData();
    formData.append('text', userText)
    formData.append('prompt', userPromt)
    console.log(formData)

    // Отправляет запрос на Web Server с данными из веб-формы
    const resp = await fetch('/processTextFromWeb', {
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

// Функция блокировки элементов управления
function lockElements(){
    intupArea.setAttribute('readonly', '')
    outputArea.setAttribute('readonly', '')
    promtArea.setAttribute('readonly', '')
    processButton.setAttribute('disabled', '')
    resetButton.setAttribute('disabled', '')
    copyTextButton.setAttribute('disabled', '')
    saveFileButton.setAttribute('disabled', '')
}

// Функция разблокировки элементов управления
function unlockElements(){
    intupArea.removeAttribute('readonly')
    outputArea.removeAttribute('readonly')
    promtArea.removeAttribute('readonly')
    processButton.removeAttribute('disabled')
    resetButton.removeAttribute('disabled')
    copyTextButton.removeAttribute('disabled')
    saveFileButton.removeAttribute('disabled')
}

// Обработчик нажатия кнопки "Сбросить"
resetButton.addEventListener('click', () => {
  outputArea.value = currentText

})

// Обработчик кнопки "Копировать текст"
copyTextButton.addEventListener('click', async () => {
   // Убираем из текста двойные переносы строк
   var text = outputArea.value.trim().replaceAll('\n\n', '\n')
   console.log("Text for copying: ", text)
   

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

    // Создаем документ и ссылку на скачиваение
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);

    // Указываем название файла
    link.download = 'rewritedText.txt';

    // Отправляем файл клиенту
    link.click();
    URL.revokeObjectURL(link.href);
})

async function sendRequestURL() {
    try {
        // Отправляем GET-запрос на сервер
        const resp = await fetch('get/operation/text', {
            method: 'GET',
        });

        // Проверяем, успешен ли ответ
        if (!resp.ok) {
            throw new Error(`HTTP error! Status: ${resp.status}`);
        }

        // Расшифровываем результат в JSON
        const data = await resp.json();

        // Получаем URL из ответа
        const url = data.url;

        return url; // Возвращаем URL
    } catch (error) {
        console.error('Ошибка при выполнении запроса:', error);
    }
}

document.getElementById('rewriteButton').onclick = async function() {  
    const url_page = await sendRequestURL();
        // Создаем элемент ссылки
        const link = document.createElement('a');
    link.href = url_page; // Устанавливаем URL
    link.target = '_blank'; // Открываем в новом окне
    link.textContent = 'ТУТ'; // Устанавливаем текст ссылки как "тут"

    // Получаем элемент для всплывающего сообщения
    const popupMessage = document.getElementById('popupMessage');
    
    // Удаляем предыдущие содержимое и добавляем новое
    popupMessage.innerHTML = `Ссылка на ваш запрос `; // Устанавливаем текст до ссылки
    popupMessage.appendChild(link); // Добавляем ссылку в сообщение

    // Показываем всплывающее окно
    document.getElementById('popup').style.display = 'block';    
    };

// Обработчик для кнопки закрытия всплывающего окна
document.getElementById('closePopup').onclick = function() {
    document.getElementById('popup').style.display = 'none'; // Скрываем всплывающее окно
};