$(document).ready(function () {
  let colorRange = {
    '91-99': 'overall1',
    '86-90': 'overall2',
    '81-85': 'overall3',
    '1-80': 'overall4',
    '0-0': 'overall5',
  };

  function between(value, min, max) {
    return value >= min && value <= max;
  }

  let color;
  let first;
  let second;
  let overall;

  $('.overall-badge').each(function (index) {

    overall = $(this);
    color = parseInt($(this).attr('overall-color'), 10);

    $.each(colorRange, function (name, value) {

      first = parseInt(name.split('-')[0], 10);
      second = parseInt(name.split('-')[1], 10);

      if (between(color, first, second)) {
        overall.addClass(value);
      }

    });

  });
});

function changeTeam(playerID, firstName, lastName, teamID, abbreviation, color, csrftoken) {
  $.ajax({
    type: "POST",
    url: "/players",
    data: {
      csrf_token: csrftoken,
      action: "change_team",
      player_id: playerID,
      team_id: teamID,
    },
    success: function (response) {
      const team = document.getElementById('player_' + playerID + '_abbreviation');
      team.textContent = abbreviation;
      team.style.color = color;
      notie.alert({ type: 1, text: firstName + ' ' + lastName + ' added to ' + abbreviation });
    },
    error: function (jqXHR, textStatus, errorThrown) {
      notie.alert({ type: 3, text: 'Error adding player: ' + errorThrown });
    }
  });
}

function filterPlayers() {
  let URLWithoutQueries = window.location.origin + window.location.pathname
  let url = new URL(URLWithoutQueries)
  let inputs = document.querySelectorAll('#filterForm input')
  inputs.forEach(input => {
    if ((!input.checked) && input.type === 'checkbox') {
      url.searchParams.set(input.name, 0)
    } else if (input.value && input.type !== 'checkbox') {
      url.searchParams.set(input.name, input.value)
    }
  })
  let teams = document.getElementById('teamSelect')
  if (teams.value !== "0") {
    url.searchParams.set("team", teams.value)
  }
  let limit = document.getElementById('limitSelect')
  if (limit.value !== "20") {
    url.searchParams.set("limit", limit.value)
  }
  let era = document.getElementById('legendSelect')
  if (era.value !== "both") {
    url.searchParams.set("era", era.value)
  }
  window.location.href = url.toString()
}

let sortByCol = document.querySelectorAll(".sortByCol");

let currentURL = window.location.href
let url = new URL(currentURL)

sortByCol.forEach(function (col) {
  col.addEventListener("click", function () {

    sortByCol.forEach(function (otherCol) {
      if (otherCol !== col) {
        otherCol.innerText = otherCol.innerText.replace(" ↓", "").replace(" ↑", "")
        otherCol.removeAttribute("sorted")
      }
    })

    if ((this.getAttribute("sorted") === "desc")) {
      this.innerText = this.innerText.replace(" ↓", "").replace(" ↑", "") + " ↑"
      this.setAttribute("sorted", "asc")
    } else if ((this.getAttribute("sorted") === "asc")) {
      this.innerText = this.innerText.replace(" ↓", "").replace(" ↑", "") + " ↓"
      this.setAttribute("sorted", "desc")
    } else {
      this.innerText += " ↓"
      this.setAttribute("sorted", "desc")
    }
    url.searchParams.set("col", this.getAttribute("data-col"))
    url.searchParams.set("sort", this.getAttribute("sorted"))
    window.location.href = url.toString()
  })
})


currentURL = new URL(window.location.href);
const searchParams = new URLSearchParams(currentURL.search)
const inputs = document.querySelectorAll('#filterForm input')

inputs.forEach(input => {
  if (input.type === 'checkbox') {
    input.checked = !searchParams.has(input.name)
  } else {
    if (input.name === 'search') {
      input.value = searchParams.get(input.name) ? searchParams.get(input.name).replaceAll('+', ' ') : ''
    } else {
      input.value = searchParams.get(input.name)
    }
  }
});

const teams = document.getElementById('teamSelect')
if (searchParams.has('team')) {
  teams.value = searchParams.get('team')
}

const limit = document.getElementById('limitSelect')
if (searchParams.has('limit')) {
  limit.value = searchParams.get('limit')
}

const era = document.getElementById('legendSelect')
if (searchParams.has('era')) {
  era.value = searchParams.get('era')
}

// Handle the table sorting
sortByCol = document.querySelectorAll('.sortByCol')
if (searchParams.has('col')) {
  sortByCol.forEach(col => {
    col.innerText = col.innerText.replace(/[↑↓]/g, '')
    col.removeAttribute('sorted')

    if (col.getAttribute('data-col') === searchParams.get('col')) {
      col.innerText += searchParams.get('sort') === 'asc' ? ' ↑' : ' ↓'
      col.setAttribute('sorted', searchParams.get('sort'))
    }
  })
}

const search = document.getElementById("search")
search.addEventListener("keypress", event => {
  if (event.key === "Enter") {
    filterPlayers()
  }
})
