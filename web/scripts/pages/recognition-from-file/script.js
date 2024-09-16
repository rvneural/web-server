const recognizeButton = document.getElementById("recognizeButton") // Кнопка Recognize

const resetButton = document.getElementById("resetButton") // Кнопка Reset
const copyTextButton = document.getElementById("copyText") // Кнопка copyText
const saveFileButton = document.getElementById("saveFile") // Кнопка saveFile

const dialogCheckBox = document.getElementById('dialog') // Чекбокс с указателем того, есть ли в записи диалог
const normalizeCheckBox = document.getElementById('normResult') // Чекбос с указанием на то, какой вид ответа нужен: нормализованный или нет

const outputArea = document.getElementById("outputArea") // Поле вывода результата
const fileInput = document.getElementById('fileInput') // Поле выбора файла
const language = document.getElementById("languageSelect") // Поле выбора языка


var normText = "" // Глобальная переменная, хранящая оригинал нормализованного текста с сервера
var rawText = "" // Глобальная переменная, хранящая оригинал не нормализованного текста с сервера

var currentNormText = "" // Глобальная переменная, хранящая измененный нормализованный текст
var currentRawText = "" // Глобальная переменная, хранящая исходный не нормализованный текст

window.onload = function() {
    document.getElementById('filePage').style.backgroundColor = "#104d2a"

    // Check if normText is in localStorage
    if (localStorage.getItem('normText')!== null) {
        normText = localStorage.getItem('normText')
        currentNormText = normText
        outputArea.value = normText

        //Unlock outputArea
        unlockElements()
    }

    // Check if rawText is in localStorage
    if (localStorage.getItem('rawText')!== null) {
        rawText = localStorage.getItem('rawText')
        currentRawText = rawText
    }
}

// Функция расшифровки
async function recognize(){

    // Проверяем, выбран ли необходимый файл
    if (fileInput.files.length == 0){
        alert('Вы не выбрали файл')
        return
    }

    // Блокируем элементы управления
    lockElements()

    // Включаем анимацию в поле вывода результата
    outputArea.classList.add("loader");
    outputArea.value = "Идет расшифровка"

    // Получаем файл и создаем структуру запроса
    const file = fileInput.files[0];
    const formData = new FormData();

    // Добавляем в структуру запроса необходимые данные
    // TODO: Сделать компрессию файла перед отправкой на сервер
    formData.append('file', file); // Файл
    formData.append('language', language.value) // Язык
    formData.append('dialog', dialogCheckBox.checked) // Указатель того, стоит ли диалог
    console.log(fileInput.value)
    parts = fileInput.files[0].name.split('.')
    console.log("Имя файла:", fileInput.files[0].name)
    console.log(parts)
    formData.append('fileType', parts.at(-1)) // Тип файла

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

    // Сохраняем в глобальных переменных исходные результаты расшифровки
    normText = await data.normText
    rawText = await data.rawText

    // update normText и rawText в localStorage
    localStorage.setItem('normText', normText)
    localStorage.setItem('rawText', rawText)

    console.log("Raw:", rawText)
    console.log("Norm:", normText)

    // Первоначально инициализируем измененные тексты начальными результатами
    currentNormText = await data.normText
    currentRawText = await data.rawText

    console.log("CurRaw:", currentRawText)
    console.log("CurNorm:", currentNormText)
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
    normalizeCheckBox.setAttribute('disabled', '')
    dialogCheckBox.setAttribute('disabled', '')
    resetButton.setAttribute('disabled', '')
    copyTextButton.setAttribute('disabled', '')
    saveFileButton.setAttribute('disabled', '')
}

// Разблокировка элементов
function unlockElements(){
    recognizeButton.removeAttribute('disabled')
    outputArea.removeAttribute('readonly')
    resetButton.removeAttribute('disabled')
    fileInput.removeAttribute('disabled')
    dialogCheckBox.removeAttribute('disabled')
    normalizeCheckBox.removeAttribute('disabled')
    copyTextButton.removeAttribute('disabled')
    saveFileButton.removeAttribute('disabled')
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

    if (normalizeCheckBox.checked){
        // Update normText в localStorage
        localStorage.setItem('normText', text)
    } else{
        // Update rawText в localStorage
        localStorage.setItem('rawText', text)
    }

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