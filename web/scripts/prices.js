const tableSelect = document.getElementById("modelSelect")
const tableBox = document.getElementById("table-box")
const fromInput = document.getElementById("search-input-start")
const toInput = document.getElementById("search-input-end")
const searchButton = document.getElementById("search-button")
const redoButton = document.getElementById("redo-button")

window.onload = function(){
    if (selectValue === ""){
        selectValue = "none"
    }
    tableSelect.value = selectValue

    if (fromBack !== ""){
        fromInput.value = fromBack
    }
    if (toBack !== ""){
        toInput.value = toBack
    }
}

document.addEventListener('DOMContentLoaded', function() {
    tableBox.innerHTML = tableText
});

redoButton.addEventListener('click', function() {
    fromInput.value = ""
    toInput.value = ""
})

searchButton.addEventListener('click', function() {
    if (tableSelect.value === "none"){
        alert("Выберите таблицу")
        return
    }
    let url = "/protected/price?model="+tableSelect.value

    if (fromInput.value !== null && fromInput.value !== ""){
        url += "&from="+fromInput.value.split("-").reverse().join(".")
    }
    if (toInput.value !== null && toInput.value !== ""){
        url += "&to="+toInput.value.split("-").reverse().join(".")
    }

    window.location.href = url
})
