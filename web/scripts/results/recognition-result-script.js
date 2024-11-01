const normArea = document.getElementById('normArea');
const rawArea = document.getElementById('rawArea');

const saveButton = document.getElementById('saveButton');

async function getVersion(){
    const formData = new FormData();
    formData.append('id', id);
    try {
        const resp = await fetch('/operation/getVersion', {
            method: 'POST',
            body: formData,
        })
        const data = await resp.json();
        if (data.error) {
            console.log(data.error);
            return 0
        } else {
            return data.version
        }
    } catch (err) {
        console.log(err);
        return 0
    }
}

saveButton.addEventListener('click', async (event)=>{

    var dbVersion = await getVersion()
    console.log(dbVersion)
    if (version === dbVersion){
        version += 1
    } else {
        alert('База данных была обновлена. Обновите страницу, чтобы получить обновленный текст')
        return
    }

    const formData = new FormData();
    formData.append("id", id)
    formData.append("type", 'audio')
    formData.append('file_name', file_name)
    formData.append('raw_text', rawArea.value)
    formData.append('norm_text', normArea.value)

    try {
    const resp = await fetch('/operation/saveOperation', {
        method: 'POST',
        body: formData,
    })
    console.log(resp)
    } catch (error) {
        console.error('Ошибка при выполнении запроса:', error);
        alert('Ошибка при выполнении запроса:', error);
        return
    }

    saveButton.innerText = "Сохранено"
        setTimeout(() => {
            saveButton.innerText = "Сохранить"
        }, 1000)
    
})