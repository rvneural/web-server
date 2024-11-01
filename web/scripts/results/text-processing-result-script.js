const oldTextArea = document.getElementById('oldTextArea');
const newTextArea = document.getElementById('newTextArea');
const promptArea = document.getElementById('promptArea');

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
    formData.append("type", 'text')
    formData.append('old_text', oldTextArea.value)
    formData.append('new_text', newTextArea.value)
    formData.append('prompt', promptArea.value)

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