const querystring = require('querystring');
const fetch = require('node-fetch');

exports.handler = async function (event, context) {
    console.log(`Event body: ${event.body}`)

    let film = "1"
    if (event.queryStringParameters.film) {
        console.log("Received parameters")
        film = event.queryStringParameters.film
        if (film < 1 || film > 8)
        {
            console.log(`film parameter out of range: ${film}`)
            const outOfRange = {
                response_type: "in_channel",
                text: `Film number must be between 1 and 7.`
            }
            const response = {
                statusCode: 200,
                    headers: { "Content-Type": "application/json" },
                body: JSON.stringify(outOfRange)
            }
            return response
        }
    } else {
        console.log("film parameter not received")
    }

    console.log(`Invoking film: ${film}`)

    return fetch(`https://swapi.co/api/films/${film}/`, { headers: { "Accept": "application/json" } })
        .then((resp) => resp.json())
        .then((data) => ({
            statusCode: 200,
            body: `${data.opening_crawl}`
        }))
        .catch(error => ({ statusCode: 400, body: String(error) }));

}
