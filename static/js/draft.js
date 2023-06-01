$(document).ready(function(){
  let colorRange = {
      '91-99' :'overall1-text',
      '86-90' :'overall2-text',
      '81-85' :'overall3-text',
      '1-80'  :'overall4-text',
      '0-0'   :'overall5-text',
  };
  
  function between(value, min, max) {
      return value >= min && value <= max
  }
  
  let color
  let first;
  let second
  let overall
  
  $('.overall-text').each(function(index){
      
      overall = $(this)
      color = parseInt($(this).attr('overall-color'),10)
      
      $.each(colorRange, function(name, value){
          
          first = parseInt(name.split('-')[0],10)
          second = parseInt(name.split('-')[1],10)
          
          if( between(color, first, second) ){
              overall.addClass(value)
          }
      
      });
      
  });
  });
  
  $(document).ready(function() {
  $("#draft-search").on("keyup", function() {
    let value = $(this).val().toLowerCase()
    $("#draft-players-list li:not(.hidden)").filter(function() {
      $(this).toggle($(this).text().toLowerCase().indexOf(value) > -1)
    })
  })
})