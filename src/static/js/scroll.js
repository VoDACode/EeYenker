document.addEventListener('DOMContentLoaded', () => {
    const arrowBlock = document.querySelector('.arrow-block');
    const showMoreSection = document.querySelector('.arrow-block .arrow-label'); // Блок, до якого скролимо

    arrowBlock.addEventListener('click', () => {
        if (showMoreSection) {
            showMoreSection.scrollIntoView({
                behavior: 'smooth', // Плавна анімація
                block: 'start'     // Блок відображається зверху
            });
        }
    });
});
