const recognizeButton = document.getElementById("recognizeButton") // Кнопка Recognize

const resetButton = document.getElementById("resetButton") // Кнопка Reset
const copyTextButton = document.getElementById("copyText") // Кнопка copyText
const saveFileButton = document.getElementById("saveFile") // Кнопка saveFile

const dialogCheckBox = document.getElementById('dialog') // Чекбокс с указателем того, есть ли в записи диалог
const normalizeCheckBox = document.getElementById('normResult') // Чекбос с указанием на то, какой вид ответа нужен: нормализованный или нет

const outputArea = document.getElementById("outputArea") // Поле вывода результата

const fileInput = document.getElementById('fileInput') // Поле выбора файла
const urlInput = document.getElementById('urlInput') // Поле ввода ссылки на файл

const language = document.getElementById("languageSelect") // Поле выбора языка
const fileTypeSelect = document.getElementById("fileTypeSelect") // Поле выбора типа файла

const whisper = document.getElementById("whisper") // Поле ввода whisper

var id = ""
var file_name = ""

var normText = "" // Глобальная переменная, хранящая оригинал нормализованного текста с сервера
var rawText = "" // Глобальная переменная, хранящая оригинал не нормализованного текста с сервера

var currentNormText = "" // Глобальная переменная, хранящая измененный нормализованный текст
var currentRawText = "" // Глобальная переменная, хранящая исходный не нормализованный текст

var version = 1


fileTypeSelect.addEventListener('change', function() {
    if (fileTypeSelect.value == "file") {
        fileInput.classList.remove('none-display')
        urlInput.classList.add('none-display')
    } else {
        fileInput.classList.add('none-display')
        urlInput.classList.remove('none-display')
    }
})


// Функция расшифровки
async function recognize(){


    if (fileTypeSelect.value == "file" && fileInput.files.length == 0) {
        alert('Вы не выбрали файл')
        return
    }

    if (fileTypeSelect.value == "url" && urlInput.value.length == 0) {
        alert('Вы не ввели ссылку')
        return
    }

    // Блокируем элементы управления
    lockElements()
    await showPopupWithLink()
    
    // Включаем анимацию в поле вывода результата
    outputArea.classList.add("loader");
    outputArea.value = "Идет расшифровка"
    const formData = new FormData();

    // Получаем файл и создаем структуру запроса
    if (fileTypeSelect.value == "file"){
        const file = fileInput.files[0];
        formData.append('file', file); // Файл
        formData.append('filename', fileInput.files[0].name)
        file_name = fileInput.files[0].name
        parts = file.name.split('.')
        formData.append('fileType', parts.at(-1)) // Тип файла
    } else {
        formData.append('url', urlInput.value) // Ссылка
        formData.append('filename', urlInput.value)
        file_name = urlInput.value
    }
    
    formData.append('whisper', whisper.checked);
    formData.append('id', id);
    formData.append('language', language.value) // Язык
    formData.append('dialog', dialogCheckBox.checked) // Диалог
    console.log(formData)

    // Показываем окно о начале расшифровки
    //alert("Началась расшифровка файла. В зависимости от его размера, процесс может занять доительное время. В среднем 1 минута расшифровывается 10 секунд")

    // Отправляет запрос на Web Server с данными из веб-формы
    const resp = await fetch('/recognize', {
        method: 'POST',
        body: formData,
    })

    // Расшифровываем результат в JSON
    const data = await resp.json();

    // Убираем анимацию расшифровки и устанавливаем в поле вывода необхожимый текст
    outputArea.classList.remove("loader");

    if ((data.error) && (data.error !== "")){
        data.normText = data.error + "\n\n" + data.details
        data.rawText = data.error + "\n\n" + data.details
    }

    if (normalizeCheckBox.checked){
        outputArea.value = data.normText
    } else {
        outputArea.value = data.rawText
    }

    // Разблокирем элементы управления
    unlockElements()

    // Сохраняем в глобальных переменных исходные результаты расшифровки
    normText = await data.normText
    rawText = await data.rawText

    // Первоначально инициализируем измененные тексты начальными результатами
    currentNormText = await data.normText
    currentRawText = await data.rawText
}

// Обработка кнопки распознавания текста
recognizeButton.addEventListener('click', async (event)=>{
    await recognize()
})

// Обработка изменения состояния чекбокса нормализации
// В зависимости от того, стоит флаг нормализации или нет
// В поле вывода результата будет появляться необходимый текст
normalizeCheckBox.onchange = function() {
    if (normalizeCheckBox.checked){
        currentRawText = outputArea.value
        outputArea.value = currentNormText
    } else{
        currentNormText = outputArea.value
        outputArea.value = currentRawText
    }
}

// Блокировка элементов
function lockElements() {
    outputArea.value = 'Здесь появится расшифровка вашего текста'
    recognizeButton.setAttribute('disabled', '')
    outputArea.setAttribute('readonly', '')
    fileInput.setAttribute('disabled', '')
    urlInput.setAttribute('disabled', '')
    normalizeCheckBox.setAttribute('disabled', '')
    dialogCheckBox.setAttribute('disabled', '')
    resetButton.setAttribute('disabled', '')
    copyTextButton.setAttribute('disabled', '')
    saveFileButton.setAttribute('disabled', '')
    fileTypeSelect.setAttribute('disabled', '')
    language.setAttribute('disabled', '')
}

// Разблокировка элементов
function unlockElements(){
    recognizeButton.removeAttribute('disabled')
    outputArea.removeAttribute('readonly')
    resetButton.removeAttribute('disabled')
    fileInput.removeAttribute('disabled')
    urlInput.removeAttribute('disabled')
    dialogCheckBox.removeAttribute('disabled')
    normalizeCheckBox.removeAttribute('disabled')
    copyTextButton.removeAttribute('disabled')
    saveFileButton.removeAttribute('disabled')
    fileTypeSelect.removeAttribute('disabled')
    language.removeAttribute('disabled')
}

// Обработчик кнопки Reset
resetButton.addEventListener('click', (event)=>{
    // Возвращаем измененный текст в нормальное состояние
    // Отдельно доступно восстановление для нормализованного
    // И не нормализованного текста
    if (normalizeCheckBox.checked){
        outputArea.value = normText
        currentNormText = normText
    } else{
        outputArea.value = rawText
        currentRawText = rawText
    }

    saveFileButton.click()
})

// Обработка копирования текста
copyTextButton.addEventListener('click', async (event)=>{

    // Убираем из текста двойные переносы строк
    var text = outputArea.value.trim().replaceAll('\n\n', '\n')

    // Пытаемся отправить текст в clipboard
    try{
        // Доступно только по протоколу HTTPS
        await navigator.clipboard.writeText(text)
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

// Обработка сохранения файла
saveFileButton.addEventListener('click', async (event)=>{

    var dbVersion = await getVersion()
    console.log(dbVersion)
    if (version === dbVersion){
        version += 1
    } else {
        alert('База данных была обновлена. Перейдите на страницу https://neuron-nexus.ru/operarion/'+id + ", чтобы получить обновленный документ.")
        return
    }

    if (normalizeCheckBox.checked){
        currentNormText = outputArea.value
    } else{
        currentRawText = outputArea.value
    }

    normText = currentNormText
    rawText = currentRawText

    const formData = new FormData();
    formData.append("id", id)
    formData.append("type", 'audio')
    formData.append('file_name', file_name)
    formData.append('raw_text', currentRawText)
    formData.append('norm_text', currentNormText)

    try {
    const resp = await fetch('/operation/saveOperation', {
        method: 'POST',
        body: formData,
    })
    console.log(resp)
    } catch (error) {
        console.error('Ошибка при выполнении запроса:', error);
        alert('Ошибка при выполнении запроса:', error);
        return
    }

    saveFileButton.innerText = "Сохранено"
        setTimeout(() => {
            saveFileButton.innerText = "Сохранить"
        }, 1000)
    
})

async function sendRequestURL() {
    try {
        // Отправляем GET-запрос на сервер
        const resp = await fetch('/operation/get', {
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

    try{
        navigator.clipboard.writeText(link.href)
    } catch (err){
        console.log('Ошибка копирования текста:', err)
    }

    // Устанавливаем содержимое всплывающего окна
    popup.innerHTML = `Результат операции будет доступен<br>`;
    popup.appendChild(link); // Добавляем ссылку в сообщение
    popup.innerHTML += `Она скопирована в буфер обмена`;

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

async function getVersion(){
    const formData = new FormData();
    formData.append('id', id);
    try {
        const resp = await fetch('/operation/getVersion', {
            method: 'POST',
            body: formData,
        })
        const data = await resp.json();
        if (data.error) {
            console.log(data.error);
            return 0
        } else {
            return data.version
        }
    } catch (err) {
        console.log(err);
        return 0
    }
}