const input = document.getElementById("search-input")

const button = document.getElementById("search-button")

const from = document.getElementById("search-input-start")
const to = document.getElementById("search-input-end")
const outDiv = document.getElementById("fips-res")

const newUrl = "/fips"

const limit = 15
var offset = 0
var searchLine = ""
var fromLine = ""
var toLine = ""

async function updateOut(){
    let data = new FormData()
    data.append("limit", limit)
    data.append("offset", offset)
    if(searchLine.length > 0) data.append("search", searchLine)
    if(fromLine.length > 0) data.append("from", fromLine)
    if(toLine.length) data.append("to", toLine)

    response = await fetch(newUrl, {
        method: "POST",
        body: data
    })

    if(!response.ok){
        alert("Не получилось загрузить данные")
        return
    }
    const json = await response.json()

    if((json["fips"] === null) || (json["fips"].length === 0)){
        return
    }

    for (var line of json["fips"]){
        outDiv.innerHTML += `<a href="${line["url"]}" target="_blank" class="description">
        <div class="image-element">
            <div class="el-data">
                <img src="${line["image_url"]}" class="image">
            </div>
            <div class="el-desc">
                <p>(${line["registration_number"]}) ${line["author"]}</p>
                <p class="date">${line["mail"]} — ${line["registration_date"]}</p>
            </div>
        </div>
        </a>`
    }
}

window.onload = async function(){
    await updateOut()
}

outDiv.addEventListener("scroll", async function(){
    if(outDiv.scrollTop + outDiv.clientHeight >= outDiv.scrollHeight - 1){
        offset += limit
        await updateOut()
    }
})

button.addEventListener("click", async function(){
    searchLine = ""
    fromLine = ""
    toLine = ""
    outDiv.innerHTML = ""
    offset = 0

    if (input.value.length === 0){
                search = ""
    } else {
        let search_value = input.value
        search_value = search_value.trim()
        search_value = search_value.toLowerCase()
        search_value = search_value.replace(/\s+/g, " ")
        search_value = search_value.replaceAll(", ", ",")
        search_value = search_value.replaceAll(" ", ",")
        searchLine += search_value
    }
        
    if(from.value.length !== 0){
        fromLine = from.value.split("-").reverse().join(".")
    }

    if(to.value.length !== 0){
        toLine = to.value.split("-").reverse().join(".")
    }

    await updateOut()
})