  /**
          The central method that does all of the AJAX handling..
          
          */
  function doAjaxCall(url, onSuccess, data, method = 'GET', contentType = null) {
      xhr = new XMLHttpRequest();
      xhr.onreadystatechange = function () {


              if (xhr.readyState == 4) {
                  if (xhr.status == 200) {
                      onSuccess(xhr.responseText);
                  } else {
                      console.log('failed. status code =' + xhr.status);
                      showErrorMessage('Request Failed. ' + 'Http Status: ' + xhr.status + ' Message = ' + xhr.responseText);
                  }
              }
          } //end onreadystatechange
      xhr.open(method, url);
      if (contentType) {
          xhr.setRequestHeader('Content-Type', contentType);
      }

      console.log('Making ajax request .. ' + url);
      xhr.send(data);
  } //end doAjaxCall()


  function showErrorMessage(errMsg) {

      enableOverlay();
      document.getElementById('idJsErrorMsg').innerHTML = errMsg;

      var elem = document.getElementById('idJsErrorDiv');
      elem.style.display = "block";
      elem.style.position = "fixed";
      elem.style.left = "35%";
      elem.style.top = "50%";
      //elem.style['z-index'] = "2";

  }

  function closeErrorDialog() {
      document.getElementById('idJsErrorDiv').style.display = "none";
      disableOverlay();
  }