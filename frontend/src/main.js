// parseHash retrives the access_token from the URL hash
function parseHash() {
  var accessTokenregex = /access_token=([\s\S]{64})&(?:token_type=\S*)&(uid=\d*)/;
  var regexResult;
  if (regexResult = accessTokenregex.exec(window.location.hash), !!regexResult && regexResult.length === 3) {
    return {
      accessToken: regexResult[1],
      uid: regexResult[2]
    };
  }
}

function openDropBoxAuth() {
  var DropboxAuthURL = 'https://www.dropbox.com/1/oauth2/authorize';
  window.location.href = DropboxAuthURL + '?client_id=pe3r68lcz7jt12e' + '&response_type=token' + '&redirect_uri=' + window.location.href;
}

// entry point
document.addEventListener('DOMContentLoaded', function (evt) {
  var h1 = document.getElementById('jsMessage');
  var authBtn = document.getElementById('chiAuth');
  // In case the page has an hash, check if it's valid
  if (window.location.hash !== '') {
    if (AT = parseHash(), !!AT) {
      h1.textContent = 'Access Token received!';
      // TODO: call server!
      reqwest({
        url: '/api/hello',
        data: {
          dropboxUser: AT.accessToken,
          dropboxUid: AT.uid
        }
      }).then(function (resp) {
        h1.textContent = 'Chinchilla approved!';
      })
      .fail(function (err, msg) {
        h1.textContent = 'Chinchilla disapproves: ' + err;
      });
    } else { // bad hash, clear it
      window.location.hash = '';
      h1.textContent = 'Access Token provided was BAD :(';
      setTimeout(function () { h1.textContent = ''; }, 5000);
    }
  }

  var chiUserID;
  if (chiUserID = localStorage.getItem('user'), !!chiUserID) {
    // TODO: AJAX TO VALIDATE chiUserID
    h1.text('Wow you\'re an actual user, that\'s.. just...WOW!');
  }
  authBtn.addEventListener('click', openDropBoxAuth);
});
