const inputArea = document.getElementById('inputArea')
const rationSelect = document.getElementById('rationSelect')
const seedArea = document.getElementById('seedArea')
const randomSeed = document.getElementById('randomSeed')

const generateImageButton = document.getElementById('generateImageButton')

const outputImage = document.getElementById('outputImage')
const seedValue = document.getElementById('seedValue')

var progress = false

var images = new Map()

window.onbeforeunload = function() {
    if (progress) {
        return "Выход из страницы приведёт к прерыванию работы генерации изображения. Вы уверены, что хотите выйти?"
    }
}

window.onload = function() {
    document.getElementById('imagePage').style.backgroundColor = "#494E56"
}

generateImageButton.addEventListener('click', async() => {
    if (inputArea.value.trim() === ''){
        alert('Введите текст для генерации изображения')
        inputArea.classList.add('need-enter')
        setTimeout(() => {
            inputArea.classList.remove('need-enter')
        }, 3000)
        return
    }

    if (!(randomSeed.checked) && (seedArea.value.trim() === '')) {
        alert('Введите seed для генерации изображения')
        return
    }

    if (inputArea.value.trim().length > 256) {
        alert('Длина текста не может превышать 256 символов' + '\nТекущая длина текста: '+ inputArea.value.trim().length + ' символов')
        return
    }

    progress = true
    prompt = inputArea.value.trim()
    seed = randomSeed.checked? 'random' : seedArea.value.trim().replaceAll('-', '').replaceAll('.', '').replaceAll(',', '')
    ratio = rationSelect.value.split('-')
    showPopupWithLink()


    lockElements()
    console.log('Запрос на генерацию изображения...')
    console.log('prompt', prompt)
    console.log('seed', seed)
    console.log('ratio', ratio)

    widthRatio = ratio[0]
    heightRatio = ratio[1]

    const formData = new FormData();
    formData.append('prompt', prompt)
    formData.append('seed', seed)
    formData.append('widthRatio', widthRatio)
    formData.append('heightRatio', heightRatio)

    const resp = await fetch('/generateImage', {
        method: 'POST',
        body: formData,
    })

    let data

    try {
        data = await resp.json();
    } catch (err) {
        console.error('Error:', err)
        alert('Ошибка при обращении к серверу')
        unlockElements()
        progress = false
        return
    }

    console.log('data', data)

    if ((data.error) && (data.error!== '')) {
        alert('Ошибка генерации изображения:'+ data.error + "\n\nДетали: " + data.details)
        unlockElements()
        progress = false
        return
    }

    images.set(rationSelect.value.trim() + "-image", data.image.b64String)
    images.set(rationSelect.value.trim() + "-seed", data.image.seed)
    images.set(rationSelect.value.trim() + "-prompt", prompt)

    unlockElements()

    setImage("data:image/png;base64," + data.image.b64String, prompt, data.image.seed)

    // Add file name to outputImage element
    outputImage.alt = data.image.prompt + ".png"

    seedValue.innerText = data.image.seed
    progress = false
})

function setImage(data, alt = '', seed='image') {
    outputImage.src = data
    outputImage.alt = inputArea.value

    // Add tag fileName to outputImage element
    outputImage.setAttribute('filename', seed.trim() + ".png")
}

function lockElements() {
    inputArea.setAttribute('readonly', '')
    rationSelect.setAttribute('disabled', '')
    randomSeed.setAttribute('disabled', '')
    generateImageButton.setAttribute('disabled', '')
}

function unlockElements() {
    inputArea.removeAttribute('readonly')
    rationSelect.removeAttribute('disabled')
    randomSeed.removeAttribute('disabled')
    generateImageButton.removeAttribute('disabled')
}

// On rationSelect change event change outputImage
rationSelect.addEventListener('change', () => {
    const selectedRation = rationSelect.value
    if (images.has(selectedRation + "-image")) {
        setImage("data:image/png;base64," + images.get(selectedRation + "-image"), images.get(selectedRation + "-prompt"), images.get(selectedRation + "-seed"))
        seedValue.innerText = images.get(selectedRation + "-seed")
        outputImage.alt = images.get(selectedRation + "-prompt") + ".png"
        inputArea.value = images.get(selectedRation + "-prompt")
    } else {
        setImage(`/web/static/img/templates/${selectedRation}.png`, `Example template for ${selectedRation} selection ratio`, 'template')
        seedValue.innerText = ""
    }    
})

// Функция обработки изменения checkbox randomSeed
randomSeed.addEventListener('change', (e) => {
    if(e.target.checked){
        //add readonly to seedArea
        seedArea.setAttribute('readonly', '')
        seedArea.value = ""
    } else{
        // remove readonly from seedArea
        seedArea.removeAttribute('readonly')
    }
})

async function sendRequestURL() {
    try {
        // Отправляем GET-запрос на сервер
        const resp = await fetch('get/operation/image', {
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
    const popup = document.getElementById('popup');
    const popupMessage = document.getElementById('popupMessage');

    // Показываем всплывающее окно
    popup.style.display = 'block';
    popup.classList.remove('slide-out'); // Убедитесь, что класс анимации удален перед показом

    const url_page = await sendRequestURL();
    
    // Создаем элемент ссылки
    const link = document.createElement('a');
    link.href = url_page; // Устанавливаем URL
    link.target = '_blank'; // Открываем в новом окне
    link.textContent = 'по этой ссылке';
    
    // Удаляем предыдущее содержимое и добавляем новое
    popupMessage.innerHTML = `Результат операции будет доступен<br>`; // Устанавливаем текст до ссылки
    popupMessage.appendChild(link); // Добавляем ссылку в сообщение
}
    
    // Обработчик для кнопки закрытия всплывающего окна
    document.getElementById('closePopup').onclick = function() {
        const popup = document.getElementById('popup');
        popup.classList.add('slide-out');
    
        // Удаляем элемент после завершения анимации
        popup.addEventListener('animationend', function() {
            resetPopup(); // Сброс состояния после анимации
        }, { once: true }); // Убедитесь, что обработчик вызывается только один раз
    };
    
    // Обработчик для кнопки копирования
    document.getElementById('copyLinkButton').onclick = function() {
        const url_page = document.querySelector('#popupMessage a').href; // Получаем URL из ссылки
        navigator.clipboard.writeText(url_page) // Копируем URL в буфер обмена
            .then(() => {
                const copyButton = document.getElementById('copyLinkButton');
                copyButton.innerText = 'Скопировано'; // Меняем текст кнопки
                
                // Уведомление об успешном копировании
                setTimeout(() => {
                    copyButton.innerText = 'Скопировать'; // Возвращаем текст кнопки
                }, 2000);
            })
            .catch(err => {
                console.error('Ошибка при копировании: ', err);
            });
    };