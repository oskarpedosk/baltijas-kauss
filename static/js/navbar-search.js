let delay
let prevQuery
const handleSearch = () => {
  const query = document.getElementById("navbar-search").value;
  if (query !== prevQuery) {
    clearTimeout(delay);
    delay = setTimeout(() => {
      if (query && query.length > 2) {
        $.ajax({
          type: "GET",
          url: "/search",
          data: { query: query },
          success: function (data) {
            if (data) {
              displaySearchResults(data)
            }
          },
          error: function (jqXHR, textStatus, errorThrown) {
            notie.alert({ type: 3, text: 'Error searching players: ' + errorThrown });
          }
        });
      } else {
        const searchResults = document.getElementById("navbar-search-results");
        searchResults.style.display = "none";
        searchResults.innerHTML = "";
      }
    }, 250);
  }
  prevQuery = query;
}

const displaySearchResults = (players) => {
  console.log(players)
  const searchResults = document.getElementById("navbar-search-results");
  searchResults.style.display = "block";
  searchResults.innerHTML = "";

  if (players.length > 0) {
    let i = 1
    players.forEach((player) => {
      const border = (i !== players.length) ? "border-bottom" : ""
      const borderColor = (player.legend === 1) ? "legend-border" : "border"
      const textColor = (player.legend === 1) ? "legend-text" : "text-black"
      const img = (player.img_url === "") ? "/static/images/players/default_player.png" : player.img_url
      let color
      let colorRange = {
        '91-99': 'overall1',
        '86-90': 'overall2',
        '81-85': 'overall3',
        '1-80': 'overall4',
        '0-0': 'overall5',
      }
      for (const range in colorRange) {
        const [min, max] = range.split('-').map(Number);

        if (player.overall >= min && player.overall <= max) {
          color = colorRange[range]
          break
        }
      }
      const listItem = document.createElement("li");
      listItem.innerHTML = `
        <a href="/players/${player.player_id}" class="d-flex p-2 px-3 justify-content-between ${border}" style="text-decoration: none;">
        <div class="d-flex">
            <img src="${img}" alt="" class="header-image rounded-circle inline-block ${borderColor}">
            <div class="ps-2 d-flex align-items-center">
              <div>
                <p class="m-0 ${textColor}">${player.first_name} ${player.last_name}</p>
                <p class="text-muted my-0" style="font-size: 12px;">
                  ${player.nba_team}
                </p>
              </div>
            </div>
        </div>
        <div class="d-flex align-items-center">
          <span class="overall-badge ${color}" overall-color="${player.overall}">${player.overall}</span>
        </div>
        </a>
        `
      listItem.classList.add("p-0");
      listItem.classList.add("dropdown-item");
      searchResults.appendChild(listItem);
      i++
    });
  } else {
    searchResults.style.display = "none";
    searchResults.innerHTML = "";
  }
};

const searchInput = document.getElementById("navbar-search");
const searchResults = document.getElementById('navbar-search-results');

document.addEventListener('click', function (event) {
  if (!searchResults.contains(event.target) && !searchInput.contains(event.target)) {
    searchResults.style.display = 'none';
  } else if (searchInput.contains(event.target) && searchResults.hasChildNodes()) {
    searchResults.style.display = "block";
  }
});
searchInput.addEventListener("keypress", event => {
  if (event.key === "Enter" && searchResults.hasChildNodes()) {
    const firstChild = searchResults.firstChild;
    if (firstChild.nodeName === "LI") {
      const link = firstChild.querySelector("a");
      if (link && link.href) {
        window.location.href = link.href;
      }
    }
  }
})