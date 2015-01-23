// parseHash retrives the access_token from the URL hash
function parseHash() {
  var accessTokenregex = /access_token=([\s\S]{64}).*(uid=\d*)/;
  var regexResult;
  if (regexResult = accessTokenregex.exec(window.location.hash), regexResult.length === 3) {
    return {
      accessToken: regexResult[1],
      uid: regexResult[2]
    };
  }
}

// silly function to sugget the user to authorize chinchilla
function nagAboutAccessToken() {
  var chiButton = document.getElementById('chiAuth');
  chiButton.addEventListener('click', function (evt) {
    var DropboxAuthURL = 'https://www.dropbox.com/1/oauth2/authorize';
    window.location.href = DropboxAuthURL + '?client_id=pe3r68lcz7jt12e' + '&response_type=token' + '&redirect_uri=' + window.location.href;
  });
}

document.addEventListener("DOMContentLoaded", function (event) {
  var messageH1 = document.getElementById('jsMessage');
  var AT;
  // dropbox puts the access token in the hash
  if (window.location.hash !== '') {
    if (AT = parseHash(), !!AT) {
      messageH1.textContent = 'Access Token received!';
    }
  }

  var chiUserID;
  if (chiUserID = localStorage.getItem('user'), !!chiUserID) {
  	console.log('user should be logged in already');
  }

  var chiButton = document.getElementById('chiAuth').style.display = 'none';
});
