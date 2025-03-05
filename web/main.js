const API_DOMAIN = `localhost:8080`;

const expressionInput = document.querySelector(`.form__input`);
const resultInfoTextarea = document.querySelector(`.result`);

let rowsCount = 1;

const addInfo = (info) => {
    resultInfoTextarea.style.opacity = `1`;
    resultInfoTextarea.textContent += info + `\r\n`;
    resultInfoTextarea.style.height = resultInfoTextarea.scrollHeight + `px`;

    resultInfoTextarea.setAttribute(`rows`, rowsCount.toString());
};

const handleFormSubmit = async (evt) => {
    evt.preventDefault();

    const response = await fetch(`http://${API_DOMAIN}/api/v1/calculate`, {
        method: `POST`,
        headers: { "Content-Type": `application/json` },
        body: JSON.stringify({expression: expressionInput.value})
    });

    if (!response.ok) {
        addInfo(`Error: ${response.statusText}`)
    }

    const data = await response.json();

    addInfo(`ID: ${data.id}`)

    expressionInput.disabled = true;

    await checkResults(data.id);

    expressionInput.disabled = false;
};

const checkResults = async (id) => {
    while (true) {
        const response = await fetch(`http://${API_DOMAIN}/api/v1/expressions/${id}`);

        const data = await response.json();

        if (data.status === `completed`) {
            addInfo(`Result: ${data.result}\n`);
            break;
        } else if (data.status === `error`) {
            addInfo(`Error: ${data.reason}\n`);
            break;
        }
    }
}

document.querySelector(`.form`).addEventListener(`submit`, handleFormSubmit);
