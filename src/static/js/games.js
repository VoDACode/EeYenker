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

    // Встановлення даних гри
    const setGameDetails = (index) => {
        const card = gameCards[index % gameCards.length]; // Зациклення індекса
        gameNameElement.textContent = card.dataset.name;
        gamePriceElement.textContent = `Ціна: ${card.dataset.price}`;
        changeBackground(index);
    };

    // Функція для зміни фону
    const changeBackground = (index) => {
        document.body.className = "search-result";
        document.body.classList.add(`game-bg-${index % gameCards.length}`);
    };

    initCarousel();
});
