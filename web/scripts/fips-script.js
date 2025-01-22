const input = document.getElementById("search-input")
const button = document.getElementById("search-button")
const from = document.getElementById("search-input-start")
const to = document.getElementById("search-input-end")

window.addEventListener("load", () => {
    var idx = window.location.href.indexOf('?') + 1
    if(idx === 0){
        return
    }

    var hashes = decodeURI(window.location.href).slice(idx).split('&');
    for(var i = 0; i < hashes.length; i++){
        var hash = hashes[i].split('=')
        if(hash[0] === "search"){
            console.log(hash[1])
            input.value = hash[1].split(",").join(", ")
        }
        if(hash[0] === "from"){
            from.value = hash[1].split(".").reverse().join("-")
        }
        if(hash[0] === "to"){
            to.value = hash[1].split(".").reverse().join("-")
        }
    }
})

button.addEventListener("click", () => {
    var search = "?search="
    var search_value = input.value
    if (search_value.length === 0){
        search = ""
    } else {
        search_value = search_value.trim()
        search_value = search_value.toLowerCase()
        search_value = search_value.replace(/\s+/g, " ")
        search_value = search_value.replaceAll(", ", ",")
        search_value = search_value.replaceAll(" ", ",")
        search += search_value
    }

    var from_date = "&from="
    if(search.length === 0){
        from_date = "?from="
    }
    var from_value = from.value

    if(from_value.length === 0){
        from_date = ""
    } else {
        var result_date = from_value.split("-").reverse().join(".")
        from_date += result_date
    }


     var to_date = "&to="
     if((from_date.length === 0) && (search.length === 0)){
        to_date = "?to="
    }
     var to_value = to.value

    if(to_value.length === 0){
        to_date = ""
    } else {
        var result_date = to_value.split("-").reverse().join(".")
        to_date += result_date
    }

    window.location.href = "/protected/fips" + search + from_date + to_date
})