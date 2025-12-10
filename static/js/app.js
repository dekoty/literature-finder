document.addEventListener('DOMContentLoaded', () => {
    // Находим все формы "Избранное"
    const favoriteForms = document.querySelectorAll('.favorite-form');

    favoriteForms.forEach(form => {
        // 2. На каждую форму навешиваем слушатель события 'submit'
        form.addEventListener('submit', async (event) => {
            // Предотвращаем стандартное поведение формы (перезагрузку страницы)
            event.preventDefault();

            const formButton = form.querySelector('.favorite-btn');
            const originalText = formButton.textContent;
            
            // Получаем данные для отправки
            const formData = new FormData(form);
            const bookId = form.dataset.bookId;

            // Блокируем кнопку на время запроса
            formButton.disabled = true;
            formButton.textContent = '...Сохранение';
            const dataToSend = new URLSearchParams(formData);
            try {
                // Отправляем асинхронный POST-запрос на наш Go-хендлер
                const response = await fetch(form.action, {
                    method: 'POST',
                    body: dataToSend, // Отправляем данные формы
                });

                // Проверяем HTTP-статус ответа
                if (response.ok) {
                    // Успех (статус 201 Created)
                    const result = await response.json();
                    
                    console.log('Успех:', result.message);
                    
                    // Обновляем UI
                    formButton.textContent = '✅ В избранном';
                    formButton.classList.add('saved');
                    form.classList.add('is-saved');

                } else {
                    // Ошибка (например, 400 Bad Request или 500 Internal Server Error)
                    const errorText = await response.text();
                    console.error('Ошибка сохранения:', errorText);
                    
                    formButton.textContent = `❌ Ошибка`;
                    alert(`Не удалось добавить в избранное: ${errorText}`);
                }

            } catch (error) {
                // Ошибка сети или другая JS-ошибка
                console.error('Сетевая ошибка:', error);
                formButton.textContent = `❌ Ошибка сети`;
                alert('Произошла сетевая ошибка.');

            } finally {
                // В любом случае разблокируем кнопку через некоторое время
                setTimeout(() => {
                    formButton.disabled = false;
                    // Если не было успеха, возвращаем исходный текст, чтобы можно было попробовать снова
                    if (!form.classList.contains('is-saved')) {
                        formButton.textContent = originalText;
                    }
                }, 1500); 
            }
        });
    });

});
