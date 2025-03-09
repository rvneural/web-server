Date.prototype.addHours= function(h){
    this.setHours(this.getHours()+h);
    return this;
}

Date.prototype.subHours= function(h){
    this.setHours(this.getHours()-h);
    return this;
}

const limit = 50
var lastDate = new Date(Date.now()).addHours(3).toISOString()
const URL = "/agregator/get"
const IMG_URL = "/web/static/img/news.jpg"
var stop = false

const newsList = document.getElementById("news-list-content")
const content = document.getElementById("content-content")
const sources = document.getElementById("sources-content")
const refresh = document.getElementById("refresh-icon")
const search = document.getElementById("search-input")
const searchBtn = document.getElementById("search-btn")

var ref = true

refresh.addEventListener("click", () => {
    if (!refresh.classList.contains("blocked")){
        refresh.classList.add("blocked")
        ref = false
    } else {
        refresh.classList.remove("blocked")
        ref = true
    }
})

async function loadNews() {
    if(stop){
        return
    }
    let url = URL + `?limit=${limit}&date=${lastDate}`
    let q = getQ()
    if (q !== '' && q !== null && q !== undefined){
        url += `&q=${q}`
    }
    const response = await fetch(url)
    const data = await response.json()

    if(data.items.length === 0){
        stop = true
        return
    }
    for(let i = 0; i < data.items.length; i++){
        lastDate = data.items[i].date
        let date = new Date(data.items[i].date)
        var options = { hour: 'numeric', minute: 'numeric', month: '2-digit', day: '2-digit'};
        let rt = ''
        if(data.items[i].isRT){
            rt = 'tatarstan'
        }
        let str = `<div class="news-item ${rt}" onclick="setContent(${data.items[i].id})"><p class="news-date">${date.subHours(3).toLocaleDateString('ru-RU', options).split(',')[1]}</p> <p class="news-list-title">${data.items[i].title}</p></div>`
        newsList.innerHTML += str
    }
}

async function updateNews(){
    if(!ref){
        return
    }
    lastDate = new Date(Date.now()).addHours(3).toISOString()
    newsList.innerHTML = ""
    await loadNews()
}

async function updateForSearch(q){
    if (q === '' || q === null || q === undefined){
        q = getQ()
    }
    lastDate = new Date(Date.now()).addHours(3).toISOString()
    newsList.innerHTML = ""
    await loadNews(q)
}

async function loadSourcesTexts(data){
    if(data.sources.length == 1){
        return
    }
    content.innerHTML += `<hr class="content-hr">`
    for(let i = 1; i < data.sources.length; i++){
        let subcontent = document.createElement('div')
        subcontent.classList.add('subcontent')

        if(data.sources[i].description){
            let desc = document.createElement('div')
            desc.innerHTML += `<h3>${data.sources[i].description}</h3>`
            subcontent.appendChild(desc)
        }
        if (data.sources[i].enclosure ){
            let figure = document.createElement('figure')
            figure.classList.add('content-figure')
            if (data.sources[i].enclosure.includes('.mp4') || data.sources[i].enclosure.includes('.wav') || data.sources[i].enclosure.includes('=flv')){
                let video = document.createElement('video')
                video.src = data.sources[i].enclosure
                video.controls = true
                figure.appendChild(video)
            } else{
                let img = document.createElement('img')
                img.src = data.sources[i].enclosure
                figure.appendChild(img)
            }
            subcontent.append(figure)
        }

        let fullText = document.createElement('div')
        fullText.classList.add('subcontent-full-text')
        fullText.innerHTML = data.sources[i].fullText

        subcontent.appendChild(fullText)

        let summary = document.createElement('summary')
        summary.innerHTML = '<b>' + data.sources[i].name + '</b>: ' + data.sources[i].title

        details = document.createElement('details')
        details.appendChild(summary)
        details.appendChild(subcontent)

        content.appendChild(details)

    }
    
}

async function loadSources(data) {
    sources.innerHTML = ""
    let isRV = false
    let url = ''
    for(let i = 0; i < data.sources.length; i++){
        let date = new Date(data.sources[i].pubDate)
        var options = { year: 'numeric', month: 'long', day: 'numeric', hour: 'numeric', minute: 'numeric' };
        let rv = ''
        if (data.sources[i].name === 'Реальное Время'){
            rv = 'real-time'
            isRV = true
            url = data.sources[i].link
        }
        sources.innerHTML += `<a href="${data.sources[i].link}" target="_blank" class="source-link ${rv}"><div class="source-item"><b class="source-name">${data.sources[i].name}</b>: ${data.sources[i].title}<p class="news-date">${date.subHours(3).toLocaleDateString('ru-RU', options)}</p></div></a>`
    }
    let map = new Map();
    map.set('isRV', isRV)
    map.set('url', url)

    return map
}

async function setContent(id) {
    let q = getQ()
    window.history.pushState({}, '', '/news/' + id)
    // window.location.href = window.location.href + `?q=${q}`
    window.history.replaceState('', '', updateURLParameter(window.location.href, "q", q));
    const response = await fetch(URL + `/${id}`)
    if(!response.ok){
        let title = 'Новость не обнаружена'
        let description = 'Возможно, вы указали неправильный ID'
        sources.innerHTML = ""
        document.title = title
        content.innerHTML = `<h1>${title}</h1>`
        content.innerHTML += `<h2>${description}</h2>`
        return
    }
    const data = await response.json()
    let isRV = await loadSources(data)
    
    let enclosureURL = ""
    for(let i = 0; i < data.sources.length; i++){
        if(data.sources[i].enclosure && data.sources[i].enclosure != '' && data.sources[i].enclosure != null){
            enclosureURL = data.sources[i].enclosure
            break
        }
        if (!data.description && data.sources[i].description && data.sources[i].description != '') {
            data.description = data.sources[i].description
        }
    }

    if (enclosureURL === ''){
        enclosureURL = IMG_URL
    }

    let date = new Date(data.sources[0].pubDate)
    var options = { year: 'numeric', month: 'long', day: 'numeric', hour: 'numeric', minute: 'numeric' };

    var div = document.createElement("div");
    div.classList.add("news-header-group");
    if(enclosureURL.includes('.mp4') || enclosureURL.includes('.mov') || enclosureURL.includes('.flv')){
        div.innerHTML = `<figure class="content-figure" id="content-figure"><video src='${enclosureURL}' controls alt='${data.title}'></figure>
                        <h1>${data.title}</h1>
                        <p class="news-date">${date.subHours(3).toLocaleDateString('ru-RU', options)} / <a href="${data.sources[0].link}" target="_blank" class="title-source-link">${data.sources[0].name}</a></p>`
    } else {
        div.innerHTML = `<figure class="content-figure" id="content-figure"><img src='${enclosureURL}' alt='${data.title}'></figure>
                    <h1>${data.title}</h1>
                    <p class="news-date">${date.subHours(3).toLocaleDateString('ru-RU', options)} / <a href="${data.sources[0].link}" target="_blank" class="title-source-link">${data.sources[0].name}</a></p>`
    }

    if(data.description){
        div.innerHTML += `<h2>${data.description}</h2>`
    }

    content.innerHTML = ''
    if(isRV.get('isRV')){
        content.innerHTML += `<h4 class="real-time-header"><a class="real-time-link" href="${isRV.get('url')}" target="_blank">Новость опубликована в «Реальном времени»</a></h4>`
    }
    content.appendChild(div);
    content.innerHTML += `<hr class="content-hr">`
    if(!data.rewrite.includes('<p>')) {
        let parts = data.rewrite.split('\n')
        let txt = ''
        for(i = 0; i < parts.length; i++){
            txt += '<p>' + parts[i] + '</p>'
        }
        data.rewrite = txt
    } 
    const url = data.sources[0].link
    let link = url.split('/')[0] + '//' + url.split('/')[2]
    data.rewrite = data.rewrite.replaceAll(`src="/`, `src="${link}/`)
    content.innerHTML += `<div class="content-p">${data.rewrite}</div>`

    document.title = data.title
    await loadSourcesTexts(data)
}

function updateURLParameter(url, param, paramVal){
    var newAdditionalURL = "";
    var tempArray = url.split("?");
    var baseURL = tempArray[0];
    var additionalURL = tempArray[1];
    var temp = "";
    if (additionalURL) {
        tempArray = additionalURL.split("&");
        for (var i=0; i<tempArray.length; i++){
            if(tempArray[i].split('=')[0] != param){
                newAdditionalURL += temp + tempArray[i];
                temp = "&";
            }
        }
    }
    var rows_txt = temp + "" + param + "=" + paramVal;
    return baseURL + "?" + newAdditionalURL + rows_txt;
}


window.onload = async function(){
    let q = getQ()
    q = q.replaceAll(',', ',' )
    search.value = q
    const id = getId()
    if(id !== null){
        await setContent(id)
    }
    updateNews()
    setInterval(updateNews, 60*1000)
}

newsList.addEventListener("scroll", function() {
    if(newsList.scrollTop + newsList.clientHeight >= newsList.scrollHeight-1){
        loadNews(getQ())
    }
})

function searchNews(){
    let q = getSearctValue()
    window.history.replaceState('', '', updateURLParameter(window.location.href, "q", q));
    updateForSearch(q)
}

function getSearctValue(){
    let searchValue = search.value
    searchValue = searchValue.replaceAll(', ', ',')
    searchValue = searchValue.toLowerCase()
    return searchValue
}

search.addEventListener("keyup", function(event) {
    if (event.code === "Enter") {
        stop = false
        searchNews();
    }
})

searchBtn.addEventListener("click", function() {
    stop = false
    searchNews();
})

function getId(){
    const url = window.location.href.split('?')[0]
    const parts = url.split('/')
    const last = parts[parts.length-1]
    if (last==='news'){
        return null
    } else {
        return last
    }
}

function getQ(){
    const urlParams = new URLSearchParams(window.location.search);
    const qF = urlParams.get('q')
    if(qF !== null && qF !== undefined){
        q = '' + qF
    } else {
        q = ''
    }
    return q
}