const API_DOMAIN = `localhost:8080`;

const expressionInput = document.querySelector(`.form__input`);
const resultInfoTextarea = document.querySelector(`.result`);


const addInfo = (info) => {
    resultInfoTextarea.style.opacity = `1`;
    resultInfoTextarea.textContent += info + `\r\n`;
    resultInfoTextarea.style.height = resultInfoTextarea.scrollHeight + `px`;
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

    const interval = setInterval(async () => {
        await checkResults(data.id, interval);
    }, 1000);
};

const checkResults = async (id, interval) => {
    const response = await fetch(`http://${API_DOMAIN}/api/v1/expressions/${id}`);

    const data = await response.json();

    if (data.status === `completed`) {
        addInfo(`Result: ${data.result}\n`);
    } else if (data.status === `error`) {
        addInfo(`Error: ${data.reason}\n`);
    }

    if (data.status !== `processing`) {
        clearInterval(interval);
        expressionInput.disabled = false;
    }
}

document.querySelector(`.form`).addEventListener(`submit`, handleFormSubmit);
