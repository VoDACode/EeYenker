document.addEventListener("DOMContentLoaded", () => {
    const gameCarousel = document.getElementById("gameCarousel");
    const gameGrid = document.getElementById("gameGrid");
    const gameCards = Array.from(document.querySelectorAll(".game-card"));
    const gameNameElements = document.querySelectorAll(".gameName");
    const gamePriceElement = document.getElementById("gamePrice");

    const onlineChartElement = document.getElementById("online-chart");
    const onlineChartFromElement = document.getElementById("online-chart-from");
    const onlineChartToElement = document.getElementById("online-chart-to");
    const onlineChartPeriodElement = document.getElementById("online-chart-period");

    let chart = null;

    let selectedGameId = 730;

    // Ініціалізація
    const initCarousel = () => {
        onlineChartFromElement.addEventListener('change', updateChart);
        onlineChartToElement.addEventListener('change', updateChart);
        onlineChartPeriodElement.addEventListener('change', updateChart);
        selectedGameId = games[0].data.steam_appid;

        console.log('Game by default:', selectedGameId);
        onlineChartFromElement.value = onlineChartFromElement.value || new Date().toISOString().split('T')[0];
        onlineChartToElement.value = onlineChartToElement.value || new Date().toISOString().split('T')[0];
        onlineChartPeriodElement.value = onlineChartPeriodElement.value || 'h';
        updateChart();

        if (gameCards.length > 0) {
            setGameDetails(0);
        }

        // Додаємо події кліку для кожної картки
        gameCards.forEach((card, index) => {
            card.addEventListener("click", () => {
                selectCard(card, index);
                console.log(card.dataset.id);
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
        selectedGameId = games[index % games.length].data.steam_appid;
        setGameDetails(index);
        updateChart();
    };

    // Функція для зміни фону
    const changeBackground = (name) => {
        const formattedName = name
            .toLowerCase() // Переводимо в нижній регістр
            .replace(/\s+/g, "-") // Замінюємо пробіли на дефіси
            .replace(/[^a-z0-9-]/g, ""); // Видаляємо всі символи, крім букв, цифр і дефісів

        const searchResultElement = document.querySelector('.search-result'); // Знаходимо контейнер
        searchResultElement.classList.remove(...searchResultElement.classList); // Скидаємо всі попередні класи
        searchResultElement.classList.add('search-result', `game-bg-${formattedName}`); // Додаємо новий фон
    };


    const setGameDetails = (index) => {
        const card = gameCards[index % gameCards.length];
        const gameInfo = games[index % games.length];
    
        gameNameElements.forEach((element) => {
            element.textContent = card.dataset.name;
        });
        gamePriceElement.textContent = `Ціна: ${card.dataset.price}`;
    
        // Оновлення опису
        const detailedDescriptionElement = document.getElementById("gameDetailedDescription");
        detailedDescriptionElement.innerHTML = gameInfo.data.detailed_description || "Опис гри недоступний.";
    
        changeBackground(card.dataset.name);
    
        // Оновлення зображень
        featured.style.backgroundImage = `url(${gameInfo.data.screenshots[0].path_full})`;
        const screenshotGallery = document.getElementById("screenshot-gallery");
        screenshotGallery.innerHTML = "";
    
        gameInfo.data.screenshots.forEach((screenshot, i) => {
            const itemWrapper = document.createElement("div");
            itemWrapper.className = "item-wrapper";
    
            const galleryItem = document.createElement("figure");
            galleryItem.className = `gallery-item image-holder r-3-2 ${i === 0 ? "active" : ""} transition`;
            galleryItem.style.backgroundImage = `url(${screenshot.path_full})`;
            galleryItem.addEventListener("click", (e) => selectItem(e));
    
            itemWrapper.appendChild(galleryItem);
            screenshotGallery.appendChild(itemWrapper);
        });
    };    

    function updateChart() {
        const from = new Date(onlineChartFromElement.value).toISOString().split('.')[0];
        const to = new Date(onlineChartToElement.value).toISOString().split('.')[0];
        const detail = onlineChartPeriodElement.value;

        fetchGameStats(selectedGameId, from, to, detail)
            .then(data => initChart(data))
            .catch(error => console.error(error));
    }

    async function fetchGameStats(appId, from, to, detail) {
        const baseUrl = '/api/stats';
        const params = new URLSearchParams({
            appid: appId,
            from: from,
            to: to,
            detail: detail,
        });

        try {
            const response = await fetch(`${baseUrl}?${params.toString()}`);
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('Error fetching aggregated data:', error);
            throw error;
        }
    }

    function initChart(data) {
        if (!data || data.length === 0) {
            console.error('No data to display');
            return;
        }

        if (chart) {
            chart.destroy();
        }

        const ctx = onlineChartElement.getContext('2d');
        console.log(data);
        chart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: data.map((item) => item.period),
                datasets: [{
                    label: 'Гравців онлайн',
                    data: data.map((item) => item.avg_count),
                    borderColor: '#00bdd6',
                    backgroundColor: '#00bdd6',
                    tension: 0.2,
                }],
            },
        });
    }

    initCarousel();
});