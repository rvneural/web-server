const inputArea = document.getElementById('inputArea')
const rationSelect = document.getElementById('rationSelect')
const seedArea = document.getElementById('seedArea')
const randomSeed = document.getElementById('randomSeed')

const generateImageButton = document.getElementById('generateImageButton')

const downloadLink = document.getElementById('downloadLink')

const outputImage = document.getElementById('outputImage')
const seedValue = document.getElementById('seedValue')

var id = ""

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

    if (inputArea.value.trim().length > 300) {
        alert('Длина текста не может превышать 300 символов' + '\nТекущая длина текста: '+ inputArea.value.trim().length + ' символов')
        return
    }

    progress = true
    prompt = inputArea.value.trim()
    seed = randomSeed.checked? 'random' : seedArea.value.trim().replaceAll('-', '').replaceAll('.', '').replaceAll(',', '')
    ratio = rationSelect.value.split('-')


    lockElements()
    await showPopupWithLink()

    widthRatio = ratio[0]
    heightRatio = ratio[1]

    const formData = new FormData();
    formData.append('id', id)
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

    imageName = prompt + ".jpg"

    // Устанавливаем изображение
    setImage(imagedata="data:image/jpeg;base64," + data.image.b64String, alt=imageName);

    seedValue.innerText = data.image.seed;
    progress = false;
})

function setImage(imagedata, alt = 'neuron-nexus') {
    outputImage.src = imagedata;
    outputImage.alt = alt; // Используем переданный alt для описания изображения

    downloadLink.href = imagedata;
    downloadLink.download = alt;
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
    try{
        navigator.clipboard.writeText(link.href);
    } catch (err) {
        console.error('Ошибка при копировании ссылки:', err);
    }
    link.target = '_blank'; // Открываем в новом окне
    link.textContent = 'по этой ссылке';

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