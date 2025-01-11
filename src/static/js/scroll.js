document.addEventListener('DOMContentLoaded', () => {
    const arrowBlock = document.querySelector('.arrow-block');

    arrowBlock.addEventListener('click', () => {
        const currentScrollPosition = window.scrollY;
        const viewportHeight = window.innerHeight;

        // Прокрутити вниз на висоту одного екрану
        window.scrollTo({
            top: currentScrollPosition + viewportHeight,
            behavior: 'smooth',
        });
    });
});