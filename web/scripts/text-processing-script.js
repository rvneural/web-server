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

var id = ""

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
  await showPopupWithLink()

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
    formData.append('id', id);
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
        id = data.id

        return url; // Возвращаем URL
    } catch (error) {
        console.error('Ошибка при выполнении запроса:', error);
    }
} 

function resetPopup() {
    const popup = document.getElementById('popup');
    const popupMessage = document.getElementById('popupMessage');

    // Скрыть всплывающее окно
    popup.style.display = 'none'; 
    
    // Удалить класс анимации, если он есть
    popup.classList.remove('slide-out'); 
    
    // Очистить сообщение
    popupMessage.innerHTML = ''; 
}

async function showPopupWithLink() {
    const popupContainer = document.getElementById('popupContainer');

    // Показать контейнер, если он неактивен
    if (!popupContainer.classList.contains('popup-active')) {
        popupContainer.classList.add('popup-active');
    }

    // Применяем анимацию к существующим всплывающим окнам
    

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

    const existingPopups = popupContainer.getElementsByClassName('popup');
    for (let popup of existingPopups) {
        popup.classList.add('slide-up'); // Добавляем класс для анимации
    }
    // Ждем завершения анимации (0.5s)
    await new Promise(resolve => setTimeout(resolve, 500));
    for (let popup of existingPopups) {
        popup.classList.remove('slide-up'); // Добавляем класс для анимации
    }
    // Добавляем новое всплывающее окно в контейнер
    popupContainer.appendChild(popup);

    popup.classList.add('slide-in'); // Добавляем класс для анимации появления
    await new Promise(resolve => setTimeout(resolve, 500));
    popup.classList.remove('slide-in');

    // Закрываем всплывающее окно через 20 секунд
    setTimeout(() => closePopup(popup), 20000);
}

// Функция для закрытия всплывающего окна
async function closePopup(popup) {
    popup.classList.add('slide-out'); // Добавляем класс для анимации исчезновения

    // Удаляем элемент после завершения анимации
    popup.addEventListener('animationend', async function() {
       

        // Проверяем, есть ли еще всплывающие окна
        const popupContainer = document.getElementById('popupContainer');
        const popups = Array.from(popupContainer.children);
        
        if (popupContainer.children.length === 0) {
            popupContainer.classList.remove('popup-active'); // Скрываем контейнер, если нет активных окон
        }
            // Спускаем все окна, находящиеся выше закрытого
            let popupIndex = popups.indexOf(popup);
            popup.remove(); // Удаляем всплывающее окно из DOM
            for (let i = popupIndex - 1; i > -1; i--) {
                popups[i].classList.add('slide-down'); // Добавляем анимацию спуска
            }
            await new Promise(resolve => setTimeout(resolve, 500));
            for (let i = popupIndex - 1; i > -1; i--) {
                popups[i].classList.remove('slide-down'); // Добавляем анимацию спуска
            }
    }, { once: true });
}