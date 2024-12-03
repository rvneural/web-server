const togglePassword = document.getElementById('togglePassword');
const passwordInput = document.getElementById('password');
const passwordIcon = document.getElementById('passwordIcon');

const toggleRepPassword = document.getElementById('toggleRepPassword');
const repPasswordInput = document.getElementById('rep-password');
const repPasswordIcon = document.getElementById('repPasswordIcon');

// Функция для переключения видимости пароля
function togglePasswordVisibility(inputField, icon) {
    const type = inputField.getAttribute('type') === 'password' ? 'text' : 'password';
    inputField.setAttribute('type', type);
    icon.src = type === 'password' ? 'web/static/img/auth/view.png' : 'web/static/img/auth/hide.png'; // Меняем изображение
}

// События для основного пароля
togglePassword.addEventListener('click', function () {
    togglePasswordVisibility(passwordInput, passwordIcon);
});

if (toggleRepPassword !== null){
    // События для повторного пароля
    toggleRepPassword.addEventListener('click', function () {
        togglePasswordVisibility(repPasswordInput, repPasswordIcon);
    });
}

// Добавление эффекта наведения для кнопок
[togglePassword, toggleRepPassword].forEach(button => {
    button.addEventListener('mouseover', function() {
        this.classList.add('hover'); // Добавляем класс при наведении
    });

    button.addEventListener('mouseout', function() {
        this.classList.remove('hover'); // Убираем класс при уходе курсора
    });
});