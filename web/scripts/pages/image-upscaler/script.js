const inputFile = document.getElementById("fileInput")
const outputImage = document.getElementById("image")
const upscaleButton = document.getElementById("upscaleButton")

window.onload = function () {
    document.getElementById('upscalePage').style.backgroundColor = "#104d2a"
    upscaleButton.setAttribute("disabled", "")
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

    // Put inputFile.files[0 to outputImage
    outputImage.src = URL.createObjectURL(inputFile.files[0])

})

upscaleButton.addEventListener("click", async function () {
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
        return
    }
    outputImage.src = data.url
})