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