const emailLabel = document.getElementById("email")
const passowrdLabel = document.getElementById("password")
const btn = document.getElementById("registerButton")

btn.addEventListener("click", async () => {
    if (emailLabel.value == "" || passowrdLabel.value == "") {
        alert("Пожалуйста, заполните все поля!")
        return
    }
    if (emailLabel.value.indexOf("@realnoevremya.ru") == -1) {
        alert("Пожалуйста, введите корректный email")
        return
    }
    const response = await fetch("/login", {
        credentials: "same-origin",
        mode: "same-origin",
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            'email': emailLabel.value,
            'password': passowrdLabel.value,
        })
    })

    const data = await response.json()
    if (data.message === "Successfully logged in") {
        window.location.href = "/protected/"
    } else {
        alert(data.error)
    }
})