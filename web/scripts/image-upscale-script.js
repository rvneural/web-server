const inputFile = document.getElementById("fileInput")
const outputImage = document.getElementById("image")
const upscaleButton = document.getElementById("upscaleButton")
const remBgButton = document.getElementById("removeBgButton")

const downloadLink = document.getElementById('downloadLink')

var progress = false

window.onbeforeunload = function () {
    if (progress) {
        return "Изображение все еще обрабатывается. Вы уверены, что хотите закрыть страницу?"
    }
}

window.onload = function () {
    document.getElementById('upscalePage').style.backgroundColor = "#494E56"
    upscaleButton.setAttribute("disabled", "")
    remBgButton.setAttribute("disabled", "")
}

inputFile.addEventListener("change", function () {
    if (inputFile.files.length == 0) {
       return
    }
    if (inputFile.files[0].type != "image/png" && inputFile.files[0].type != "image/jpeg") {
        alert("Выберите фотографию!")
        inputFile.value = ""
        return
    }
    upscaleButton.removeAttribute("disabled")
    remBgButton.removeAttribute("disabled")

    // Put inputFile.files[0 to outputImage
    outputImage.src = URL.createObjectURL(inputFile.files[0])

})

upscaleButton.addEventListener("click", async function () {
    progress = true
    remBgButton.setAttribute("disabled", "")
    upscaleButton.setAttribute("disabled", "")
    originalImage = inputFile.files[0]
    formData = new FormData()
    formData.append("image", originalImage)
    const resp = await fetch("/upscaleImage", {
        method: "POST",
        body: formData
    })
    const data = await resp.json()
    if(data.error) {
        alert("Произошла следующая ошибка при обработке изображения:\n" + data.error)
        progress = false
        remBgButton.removeAttribute("disabled")
        upscaleButton.removeAttribute("disabled")
        return
    }
    outputImage.src = data.url
    downloadLink.href = data.url
    remBgButton.removeAttribute("disabled")
    upscaleButton.removeAttribute("disabled")
    progress = false
})

remBgButton.addEventListener("click", async function () {
    progress = true
    remBgButton.setAttribute("disabled", "")
    upscaleButton.setAttribute("disabled", "")
    originalImage = inputFile.files[0]
    formData = new FormData()
    formData.append("image", originalImage)
    const resp = await fetch("/removeBackground", {
        method: "POST",
        body: formData
    })
    const data = await resp.json()
    if(data.error) {
        alert("Произошла следующая ошибка при обработке изображения:\n" + data.error)
        progress = false
        remBgButton.removeAttribute("disabled")
        upscaleButton.removeAttribute("disabled")
        return
    }
    outputImage.src = data.url
    downloadLink.href = data.url
    remBgButton.removeAttribute("disabled")
    upscaleButton.removeAttribute("disabled")
    progress = false
})

