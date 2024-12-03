const emailLabel = document.getElementById("email")
const passowrdLabel = document.getElementById("password")
const firts_nameLabel = document.getElementById("first_name")
const last_nameLabel = document.getElementById("last_name")
const btn = document.getElementById("registerButton")

btn.addEventListener('click', async () => {
    if(emailLabel.value == "" || passowrdLabel.value == "" || firts_nameLabel.value == "" || last_nameLabel.value == ""){
        alert("Пожалуйста, заполните все поля!")
        return
    }
    if (passowrdLabel.value.length < 8 || passowrdLabel.value.length > 20) {
        alert("Пароль должен содержать не менее 8 символов и не более 20 символов")
        return
    }
    if (emailLabel.value.indexOf("@realnoevremya.ru") == -1) {
        alert("Пожалуйста, введите корректный email")
        return
    }
    const response = await fetch("/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            'email': emailLabel.value,
            'password': passowrdLabel.value,
            'firstName': firts_nameLabel.value,
            'lastName': last_nameLabel.value,
        })
    })

    const data = await response.json()
    if (data.exists == true) {
        alert("Пользователь с таким email уже существует")
    } else if (data.message === "Successfully registered") {
        window.location.href = "/protected/"
    } else {
        alert(data.error)
    }
})