async function fetchData() {
    const response = await fetch('https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1');
    const data = await response.json();
    return data;
}

function displayData(data) {
    let tableHTML = '<table><tr><th>id</th><th>symbol</th><th>name</th></tr>';
    data.forEach((coin, index) => {
        let bgColor = '';
        if (index < 5) bgColor = 'blue';
        if (coin.symbol === 'usdt') bgColor = 'green';
        tableHTML += `<tr style="background-color: ${bgColor}"><td>${coin.id}</td><td>${coin.symbol}</td><td>${coin.name}</td></tr>`;
    });
    tableHTML += '</table>';
    document.getElementById('cryptoTable').innerHTML = tableHTML;
}

fetchData().then(data => displayData(data));
