//globals

var reqFieldsPrefixId0 = 'idJsReqFieldSelectedId_';
var reqFieldsPrefixId1 = 'idJsReqFieldId_';
var reqFieldsPrefixId2 = 'idJsReqFieldValueId_';
var reqFieldsValErrorIdPrefix = 'idJsFieldValErr_';

var respFieldsPrefixId0 = 'idJsRespFieldSelectedId_';
var respFieldsPrefixId1 = 'idJsRespFieldId_';
var respFieldsPrefixId2 = 'idJsRespFieldValueId_';

var alphaRE = /^[a-zA-Z]+$/;
var alphaNumericRE = /^[0-9a-zA-Z]+$/;
var numericRE = /^[0-9]+$/;

var gPageState = {
  layoutLoaded: false
};

//functions

function enableOverlay() {
  document.getElementById('overlay').style['z-index'] = 99;
}

function disableOverlay() {
  document.getElementById('overlay').style['z-index'] = 80;
}

function getRequestFieldValElem(fieldId) {
  //alert(document.getElementById(reqFieldsPrefixId2 + fieldId));
  return document.getElementById(reqFieldsPrefixId2 + fieldId);
}

function getResponseFieldValElem(fieldId) {
  return document.getElementById(respFieldsPrefixId2 + fieldId);
}

//functions related to dealing with user input - 

function jsSaveReqData() {

  enableOverlay();
  document.getElementById('idJsUserInput').value = '<Data Set Name>';

  var elem = document.getElementById('idJsUserInputDiv');
  elem.style.display = "block";
  elem.style.position = "fixed";
  elem.style.left = "35%";
  elem.style.top = "50%";
}

function jsProcessUserInput() {

  var dataSetName = document.getElementById('idJsUserInput').value;
  if (dataSetName) {

    if (dataSetName === "<Data Set Name>") {
      showErrorMessage('Please provide a valid data set name.');
      return;
    } else {

      var postData = '';
      postData += 'specId=' + getCurrentSpecId() + '&';
      postData += 'msgId=' + getCurrentSpecMsgId() + '&';
      postData += 'dataSetName=' + dataSetName + '&';

      var reqObj = constructReqObj();
      if (reqObj.length === 0) {
        //document.getElementById('idJsUserInputDiv').style.display = "none";
        showErrorMessage(
            'There is no data to save, please construct a message.');

        return;
      }

      postData += 'msg=' + JSON.stringify(reqObj);
      //alert(postData);

      doAjaxCall('/iso/v0/save', function (response) {
        jsShowUserMsg('Message [' + dataSetName + '] has been saved.');
        updateCurrentDs(dataSetName)

      }, postData, 'POST', 'application/x-www-form-urlencoded');

    }

  } else {
    showErrorMessage('Please provide a valid data set name.');
    return;

  }

  jsCancelUserInputDialog();

}

function jsCancelUserInputDialog() {
  document.getElementById('idJsUserInputDiv').style.display = "none";
  disableOverlay();

}

function getCurrentSpecId() {

  var specsElem = document.getElementById('idJsSpec');
  var selectedOption = specsElem.options[specsElem.selectedIndex];
  return selectedOption.id;

}

function getCurrentSpecMsgId() {


  //var specsElem = document.getElementById('idJsSpec');
  //var selectedOption = specsElem.options[specsElem.selectedIndex];
  //return selectedOption.id;
  var elem = document.getElementById('idJsSpecMsgs');
  var msgId = elem.options[elem.selectedIndex].id;
  return msgId;

}

function jsLoadReqData() {

  doAjaxCall('/iso/v0/loadmsg?specId=' + getCurrentSpecId() + '&msgId='
      + getCurrentSpecMsgId(), function (response) {

    //alert(response);

    var dataSets = JSON.parse(response);
    var htmlContent = '';
    for (var i = 0; i < dataSets.length; i++) {
      htmlContent += '<option value="' + dataSets[i] + "\">" + dataSets[i]
          + "</option>";

    }

    var selectElem = document.getElementById('idJsUserSelection');
    selectElem.innerHTML = htmlContent;
    selectElem.selectedIndex = 0;

    enableOverlay();
    var elem = document.getElementById('idJsUserSelectionDiv');
    elem.style.display = "block";
    elem.style.position = "fixed";
    elem.style.left = "35%";
    elem.style.top = "50%";

  }, 'GET');

}

function getCurrentDs() {
  return document.getElementById('idJsCurrentDs').value;
}

function updateCurrentDs(dsName) {

  if (dsName) {
    gPageState.currentDataSet = dsName;
    document.getElementById('idJsCurrentDs').value = dsName;
    document.getElementById('idJsUpdateMsgBtn').disabled = false;
    document.getElementById('idJsUpdateMsgBtn').className = 'cl_small_btn';

  } else {
    document.getElementById('idJsCurrentDs').value = "";
    document.getElementById('idJsUpdateMsgBtn').disabled = true;
    document.getElementById(
        'idJsUpdateMsgBtn').className = 'cl_small_btn_disabled';
  }

}

function jsProcessUserSelection() {

  var selectElem = document.getElementById('idJsUserSelection');
  var dsName = selectElem.options[selectElem.selectedIndex].value;

  //alert('Fetching ds - ' + dsName);

  doAjaxCall('/iso/v0/loadmsg?specId=' + getCurrentSpecId() + '&msgId='
      + getCurrentSpecMsgId() + '&dsName=' + dsName, function (response) {

    //alert(response);
    var dataSet = JSON.parse(response);
    updateCurrentDs(dsName);

    jsLoadTemplate(dataSet, function () {
      jsCancelUserSelectionDialog();
    });

  }, 'GET');

}

function jsCancelUserSelectionDialog() {
  document.getElementById('idJsUserSelectionDiv').style.display = "none";
  disableOverlay();
}

function jsShowUserMsg(msg) {

  enableOverlay();
  document.getElementById('idJsUserMsg').innerHTML = msg;

  var elem = document.getElementById('idJsUserMsgDiv');
  elem.style.display = "block";
  elem.style.position = "fixed";
  elem.style.left = "35%";
  elem.style.top = "50%";
}

function jsCloseUserMsgDialog() {
  document.getElementById('idJsUserMsgDiv').style.display = "none";
  disableOverlay();
}

function jsUpdateMessage() {

  var postData = '';
  postData += 'specId=' + getCurrentSpecId() + '&';
  postData += 'msgId=' + getCurrentSpecMsgId() + '&';
  postData += 'updateMsg=' + 'true' + '&';
  var ds = getCurrentDs();
  postData += 'dataSetName=' + ds + '&';

  var reqObj = constructReqObj();
  if (reqObj.length === 0) {
    showErrorMessage('There is no data to save, please construct a message.');

    return;
  }

  postData += 'msg=' + JSON.stringify(reqObj);
  //alert(postData);

  doAjaxCall('/iso/v0/save', function (response) {
    jsShowUserMsg('Message [' + ds + '] has been updated.');
  }, postData, 'POST', 'application/x-www-form-urlencoded');

}
