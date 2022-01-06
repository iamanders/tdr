// Toggle non gap on day page
let dayNoGapCheckbox = document.querySelector('.day-no-gap-label input[type="checkbox"]');
if (dayNoGapCheckbox) {
    dayNoGapCheckbox.addEventListener('change', function() {
        let start = document.querySelector('.day-toggle-advanced-form div:first-child');
        start.style.display = start.style.display == 'none' ? 'block' : 'none';
    })
}

// Toggle advanced form on day page
let dayToggleAdvancedFormButton = document.querySelector('.day-toggle-advanced-form-button');
if (dayToggleAdvancedFormButton) {
    // Toggle form
    let dayToggleAdvancedForm = document.querySelector('.day-toggle-advanced-form');
    dayToggleAdvancedFormButton.addEventListener('click', (event) => {
        event.preventDefault();
        dayToggleAdvancedForm.classList.toggle('d-none');
        dayToggleAdvancedFormButton.querySelector('svg').classList.toggle('rotate-90');
    });

    // Insert current time
    let dayToggleInsertCurrentTimes = document.querySelectorAll('.day-toggle-insert-current-time');
    dayToggleInsertCurrentTimes.forEach((a) => {
        a.addEventListener('click', (event) => {
            event.preventDefault();
            let date = new Date();
            let dateString = `${date.toLocaleDateString('sv-SE')} ${date.toLocaleTimeString('sv-SE').substr(0, 5)}`;
            event.target.parentNode.querySelector('input').value = dateString;
        });
    });
}

// Auto grow textarea on day page
function textareaAutoHeight(element) {
    element.style.height = "1rem";
    element.style.height = (element.scrollHeight + 16)+"px";
}
function notesTextAreaFocus(element) {
    // Show submit button
    let submit = document.querySelector('.notes-form button[type="submit"]');
    submit.style.display = 'block';
}
document.addEventListener('DOMContentLoaded', function(event) {
    let textarea = document.querySelector('textarea.auto-height');
    if (textarea) {
        textareaAutoHeight(textarea);
    }
});

// Week page charts
let chartProjectSummaryCanvas = document.querySelector('#chart-project-summary');
if (chartProjectSummaryCanvas) {
    new Chart(chartProjectSummaryCanvas, {
        type: 'pie',
        data: chartProjectSummaryData,
        options: {
            plugins: {
                legend: {
                    position: 'bottom',
                    align: 'start',
                },
            },
        }
    });
}
