const input = document.getElementById("search-input")
const button = document.getElementById("search-button")

button.addEventListener("click", () => {
    var value = input.value
    if (value.length === 0){
        return
    }
    value = value.trim()
    value = value.toLowerCase()
    value = value.replace(", ", ",")
    value = value.replace(" ", ",")
    //redirect to /fips/search?q=value
    window.location.href = "/protected/fips?search=" + value
})