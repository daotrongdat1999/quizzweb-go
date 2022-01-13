function validateForm() {
    var x = document.forms["loginform"]["username"].value;
    var y = document.forms["loginform"]["password"].value;
    if (x == "" || x == null || y == "" || y == null) {
      alert("All box must be filled out");
      return false;
    }
  }

  function validateForm2() {
    var x = document.forms["registform"]["username"].value;
    var y = document.forms["registform"]["password"].value;
    if (x == "" || x == null || y == "" || y == null) {
      alert("All box must be filled out");
      return false;
    }
  }