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