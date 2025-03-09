const txtarea = document.getElementById('text');
const audios = document.getElementById('audio-list');

const voiceSelect = document.getElementById('voice');
const roleSelect = document.getElementById('role');
const speedSelect = document.getElementById('speed');
const pitchShiftSelect = document.getElementById('pitchShift');

window.onload = function(){
    for(let[key,value] of voices.entries()){
        voiceSelect.innerHTML += `<option value="${key}">${value}</option>`;
    }

    changeRoles()
}

voiceSelect.addEventListener('change', changeRoles);

function changeRoles(){
    roleSelect.innerHTML = '';
    var role = roles.get(voiceSelect.value).split(',');
    for(let i of role){
        roleSelect.innerHTML += `<option value="${i}">${ampluas.get(i)}</option>`;
    }
}

function add(text){
    var ps = txtarea.selectionStart
    if (ps == 0){
        txtarea.value = txtarea.value + text;
        txtarea.selectionStart = txtarea.value.length;
        txtarea.selectionEnd = txtarea.value.length;
        txtarea.focus();
    } else {
        txtarea.value = txtarea.value.substring(0, ps) + text + txtarea.value.substring(ps, txtarea.value.length);
        txtarea.selectionStart = txtarea.selectionEnd = ps+text.length;
        txtarea.focus();
    }
}

function addSelection(st, nd){
    var start = txtarea.selectionStart;
    var end = txtarea.selectionEnd;
    txtarea.value = txtarea.value.substring(0, start) + st + txtarea.value.substring(start, end) + nd + txtarea.value.substring(end, txtarea.value.length)
    txtarea.selectionStart = txtarea.selectionEnd =end + st.length + nd.length
    txtarea.focus();
}

async function generate(){
    if(txtarea.value.length == 0){
        alert('Вы не ввели текст');
        return;
    }
    var data = new FormData();
    data.append('text', txtarea.value);
    data.append('voice', voiceSelect.value);
    if(roleSelect.value !== '-') data.append('role', roleSelect.value);
    data.append('speed', speedSelect.value);
    data.append('pitchShift', pitchShiftSelect.value);
    var res = await fetch('/tts', {
        method: 'POST',
        body: data
    })
    if(!res.ok){
        alert('Ошибка сервера');
        return;
    }
    var audio = await res.json();
    bytes = audio.audio;
    audios.innerHTML = `<audio controls src="data:audio/wav;base64,${bytes}" class="fade-in"></audio>` + audios.innerHTML;
    setTimeout(function() {
        for(let elem of audios.children){
            elem.classList.remove('fade-in');
        }
    }, 1000)
}

function setDefaultValue(id, val){
    document.getElementById(id).value = val;
}

function saveSettings(){
    localStorage.setItem('voice', voiceSelect.value);
    localStorage.setItem('role', roleSelect.value);
    localStorage.setItem('speed', speedSelect.value);
    localStorage.setItem('pitchShift', pitchShiftSelect.value);
}

function loadSettings(){
    // check if localStorage is empty
    if(localStorage.length == 0) return;

    voiceSelect.value = localStorage.getItem('voice');
    changeRoles();
    roleSelect.value = localStorage.getItem('role');
    speedSelect.value = localStorage.getItem('speed');
    pitchShiftSelect.value = localStorage.getItem('pitchShift');

    document.getElementById('pitchOutput').value = localStorage.getItem('pitchShift');
    document.getElementById('speedOutput').value = localStorage.getItem('speed');
}