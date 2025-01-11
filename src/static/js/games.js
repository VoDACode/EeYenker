document.addEventListener("DOMContentLoaded", () => {
    const gameCarousel = document.getElementById("gameCarousel");
    const gameGrid = document.getElementById("gameGrid");
    const gameCards = Array.from(document.querySelectorAll(".game-card"));
    const gameNameElement = document.getElementById("gameName");
    const gamePriceElement = document.getElementById("gamePrice");

    // Ініціалізація
    const initCarousel = () => {
        if (gameCards.length > 0) {
            setGameDetails(0);
        }

        // Додаємо події кліку для кожної картки
        gameCards.forEach((card, index) => {
            card.addEventListener("click", () => {
                selectCard(card, index);
            });
        });

        // Додаємо подію прокрутки колесом миші
        gameCarousel.addEventListener("wheel", (e) => {
            e.preventDefault();
            gameCarousel.scrollLeft += e.deltaY;

            // Зациклення елементів
            const scrollLeft = gameCarousel.scrollLeft;
            const firstCard = gameGrid.firstElementChild;
            const lastCard = gameGrid.lastElementChild;

            if (scrollLeft <= 0) {
                // Переносимо останній елемент на початок
                gameGrid.insertBefore(lastCard, firstCard);
                gameCarousel.scrollLeft += lastCard.offsetWidth;
            } else if (scrollLeft + gameCarousel.clientWidth >= gameCarousel.scrollWidth) {
                // Переносимо перший елемент у кінець
                gameGrid.appendChild(firstCard);
                gameCarousel.scrollLeft -= firstCard.offsetWidth;
            }
        });
    };

    // Вибір картки
    const selectCard = (card, index) => {
        document.querySelector(".game-card.selected")?.classList.remove("selected");
        card.classList.add("selected");
        setGameDetails(index);
    };

    // Функція для зміни фону
    const changeBackground = (name) => {
        const formattedName = name
            .toLowerCase() // Переводимо в нижній регістр
            .replace(/\s+/g, "-") // Замінюємо пробіли на дефіси
            .replace(/[^a-z0-9-]/g, ""); // Видаляємо всі символи, крім букв, цифр і дефісів
    
        document.body.className = "search-result"; // Скидаємо попередній фон
        document.body.classList.add(`game-bg-${formattedName}`); // Додаємо новий фон
    };    

    // Змінюємо виклик у setGameDetails
    const setGameDetails = (index) => {
        const card = gameCards[index % gameCards.length]; // Зациклення індекса
        gameNameElement.textContent = card.dataset.name;
        gamePriceElement.textContent = `Ціна: ${card.dataset.price}`;
        changeBackground(card.dataset.name);
    };

    initCarousel();
});
