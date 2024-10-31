const intupArea = document.getElementById('inputArea')
const outputArea = document.getElementById('outputArea')

const processButton = document.getElementById('rewriteButton')

const resetButton = document.getElementById("resetButton") // Кнопка Reset
const copyTextButton = document.getElementById("copyText") // Кнопка copyText
const saveFileButton = document.getElementById("saveFile") // Кнопка saveFile

var currentText = ""

var id = ""

var progress = false

window.onbeforeunload = function() {
    if (progress) {
        return "Текст все еще обраатывается. Вы уверены, что хотите закрыть страницу?";
    }
};

window.onload = function() {
    document.getElementById('rewritePage').style.backgroundColor = "#494E56"
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

  // Подготовка
  lockElements()
  progress = true
  outputArea.value = 'Переписываем текст...'
  outputArea.classList.add("loader");
  await showPopupWithLink()
  
  // Переписывание текста
  await rewriteText(intupArea.value.trim())

  // Финал
  outputArea.classList.remove("loader");
  outputArea.value = currentText

  unlockElements()
  progress = false
})

// Функция переписывания текста, возвращает строку
async function rewriteText(userText){
    const formData = new FormData();
    formData.append('id', id);
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
        const resp = await fetch('operation/get', {
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
        id = data.id;

        return url; // Возвращаем URL
    } catch (error) {
        console.error('Ошибка при выполнении запроса:', error);
    }
}

async function showPopupWithLink() {
    const popupContainer = document.getElementById('popupContainer'); // Контейнер для всплывающих окон

    // Показать контейнер, если он неактивен
    if (!popupContainer.classList.contains('popup-active')) {
        popupContainer.classList.add('popup-active');
    }

    // Создаем новое всплывающее окно
    const popup = document.createElement('div');
    popup.className = 'popup'; // Добавляем класс для стилей

    const url_page = await sendRequestURL();
    
    // Создаем элемент ссылки
    const link = document.createElement('a');
    link.href = url_page; // Устанавливаем URL
    link.target = '_blank'; // Открываем в новом окне
    link.textContent = 'по этой ссылке';
    
    // Устанавливаем содержимое всплывающего окна
    popup.innerHTML = `Результат операции будет доступен<br>`;
    popup.appendChild(link); // Добавляем ссылку в сообщение

    // Создаем контейнер для кнопок
    const buttonContainer = document.createElement('div');
    buttonContainer.className = 'button-container'; // Класс для стилизации кнопок

    // Создаем кнопку "Закрыть"
    const closeButton = document.createElement('button');
    closeButton.textContent = 'Закрыть';
    closeButton.className = 'closePopup'; // Добавляем класс для стилей
    closeButton.onclick = () => closePopup(popup); // Обработчик для закрытия
    buttonContainer.appendChild(closeButton); // Добавляем кнопку в контейнер

    // Создаем кнопку "Скопировать"
    const copyButton = document.createElement('button');
    copyButton.textContent = 'Скопировать';
    copyButton.className = 'copy-link-button'; // Добавляем класс для стилей
    copyButton.onclick = () => {
        navigator.clipboard.writeText(link.href) // Копируем URL в буфер обмена
            .then(() => {
                copyButton.innerText = 'Скопировано'; // Меняем текст кнопки
                setTimeout(() => {
                    copyButton.innerText = 'Скопировать'; // Возвращаем текст кнопки
                }, 2000);
            })
            .catch(err => {
                console.error('Ошибка при копировании: ', err);
            });
    };
    buttonContainer.appendChild(copyButton); // Добавляем кнопку в контейнер

    // Добавляем контейнер кнопок в всплывающее окно
    popup.appendChild(buttonContainer);

    // Добавляем новое всплывающее окно в контейнер
    popupContainer.appendChild(popup);

    // Закрываем всплывающее окно через 20 секунд
    setTimeout(() => closePopup(popup), 20000);
}

// Функция для закрытия всплывающего окна
function closePopup(popup) {
    popup.classList.add('slide-out'); // Добавляем класс для анимации исчезновения

    // Удаляем элемент после завершения анимации
    popup.addEventListener('animationend', function() {
        popup.remove(); // Удаляем всплывающее окно из DOM

        // Проверяем, есть ли еще всплывающие окна
        const popupContainer = document.getElementById('popupContainer');
        if (popupContainer.children.length === 0) {
            popupContainer.classList.remove('popup-active'); // Скрываем контейнер, если нет активных окон
        }
    }, { once: true });
}