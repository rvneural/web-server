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


var normText = "" // Глобальная переменная, хранящая оригинал нормализованного текста с сервера
var rawText = "" // Глобальная переменная, хранящая оригинал не нормализованного текста с сервера

var currentNormText = "" // Глобальная переменная, хранящая измененный нормализованный текст
var currentRawText = "" // Глобальная переменная, хранящая исходный не нормализованный текст

var progress = false

fileTypeSelect.addEventListener('change', function() {
    if (fileTypeSelect.value == "file") {
        fileInput.classList.remove('none-display')
        urlInput.classList.add('none-display')
    } else {
        fileInput.classList.add('none-display')
        urlInput.classList.remove('none-display')
    }
})

window.onbeforeunload = function() {
    if (progress) {
        return 'Расшифровка еще не закончена. Вы уверены, что хотите закрыть страницу?'
    }
}

window.onload = function() {
    document.getElementById('filePage').style.backgroundColor = "#494E56"
}

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
    progress = true

    // Включаем анимацию в поле вывода результата
    outputArea.classList.add("loader");
    outputArea.value = "Идет расшифровка"
    const formData = new FormData();

    // Получаем файл и создаем структуру запроса
    if (fileTypeSelect.value == "file"){
        const file = fileInput.files[0];
        formData.append('file', file); // Файл
        parts = file.name.split('.')
        formData.append('fileType', parts.at(-1)) // Тип файла
    } else {
        formData.append('url', urlInput.value) // Ссылка
    }
    
    formData.append('language', language.value) // Язык
    formData.append('auto', autoCheckBox.checked) 

    // Показываем окно о начале расшифровки
    alert("Началась расшифровка файла. В зависимости от его размера, процесс может занять доительное время. В среднем 1 минута расшифровывается 10 секунд")

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
    console.log(data)

    // Разблокирем элементы управления
    unlockElements()
    progress = false

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
})

// Обработка копирования текста
copyTextButton.addEventListener('click', async (event)=>{

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

// Обработка сохранения файла
saveFileButton.addEventListener('click', (event)=>{

    // Создаем BLOB и убираем двойные переносы строк в тексте
    const blob = new Blob([outputArea.value.trim().replaceAll('\n\n', '\n')], {type: 'text/plain'})

    // Создаем документ и ссылку на скачиваение
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);

    // Указываем название файла
    link.download = 'recognition.txt';

    // Отправляем файл клиенту
    link.click();
    URL.revokeObjectURL(link.href);
})

async function sendRequestURL() {
    try {
        // Отправляем GET-запрос на сервер
        const resp = await fetch('get/operation/audio', {
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

document.getElementById('recognizeButton').onclick = async function() {  
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