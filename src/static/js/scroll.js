document.addEventListener('DOMContentLoaded', () => {
    const arrowBlock = document.querySelector('.arrow-block');
    const gradientSection = document.getElementById('gradientSection');

    arrowBlock.addEventListener('click', () => {
        if (gradientSection) {
            gradientSection.scrollIntoView({
                behavior: 'smooth',
            });
        }
    });
});