function editResult(index) {
  const str = 'edit_result_' + (index);
  $('#' + str).toggle();
}

function updateResult(index, timestamp) {
  const home_team_id = document.getElementById('home_team' + index).value;
  const home_score = document.getElementById('home_score' + index).value;
  const away_score = document.getElementById('away_score' + index).value;
  const away_team_id = document.getElementById('away_team' + index).value;

  const error = document.getElementById('error')
  const errorToast = new bootstrap.Toast(error)

  const success = document.getElementById('success')
  const successToast = new bootstrap.Toast(success)

  if (away_team_id === home_team_id) {
    msg = "Home and away team have to be different";
    document.getElementById("errorid").innerHTML = msg;
    errorToast.show()
  } else if ((home_score < 0) || (away_score < 0)) {
    msg = "Score can only be positive";
    document.getElementById("errorid").innerHTML = msg;
    errorToast.show()
  } else if (home_score === away_score) {
    msg = "Score can't be a draw";
    document.getElementById("errorid").innerHTML = msg;
    errorToast.show()
  } else {
    $.ajax({
      type: 'GET',
      dataType: 'html',
      data: {
        action: "update",
        home_team_id: home_team_id,
        home_score: home_score,
        away_score: away_score,
        away_team_id: away_team_id,
        timestamp: timestamp,
      },
      success: function (data) {
        msg = "Result successfully updated!";
        document.getElementById("successid").innerHTML = msg;
        successToast.show()
        console.log("success");
        $('#results').load(window.location.href + ' #resultsreload');
        $('#standings').load(window.location.href + ' #standingsreload');
      }
    });
  }
}

function deleteResult(timestamp) {
  const success = document.getElementById('success')
  const successToast = new bootstrap.Toast(success)

  const confirm = document.getElementById('confirm');
  const confirmToast = new bootstrap.Toast(confirm);
  confirmToast.show();

  const confirmButton = document.getElementById('confirmbutton');

  setTimeout(() => {
    document.addEventListener("click", function (event) {
      if (event.target === confirmButton) {
        $.ajax({
          type: 'GET',
          dataType: 'html',
          data: {
            action: "delete",
            timestamp: timestamp,
          },
          success: function (data) {
            $('#results').load(window.location.href + ' #resultsreload');
            $('#standings').load(window.location.href + ' #standingsreload');
            msg = "Result deleted!";
            document.getElementById("successid").innerHTML = msg;
            successToast.show()
          }
        });
      }
    }, { once: true })
  }, 100)
}


const dates = document.querySelectorAll('.date')
dates.forEach(element => {
  let date = new Date(element.innerHTML)
  const options = {
    hour: '2-digit',
    minute: '2-digit',
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  };
  const formattedDate = date.toLocaleDateString('est', options).replace(',', '').replaceAll('.', '/').split(' ', -1)
  element.innerHTML = formattedDate[1] + ' - ' + formattedDate[0]
})