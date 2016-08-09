function enableOverlay() {
    document.getElementById('overlay').style['z-index'] = 99;
}

function disableOverlay() {
    document.getElementById('overlay').style['z-index'] = 80;
}


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

        if (dataSetName == '<Data Set Name>') {
            showErrorMessage('Please provide a valid data set name.');
            return;
        } else {

            var postData = '';
            postData += 'specId=' + getCurrentSpecId() + '&';
            postData += 'msgId=' + getCurrentSpecMsgId() + '&';
            postData += 'dataSetName=' + dataSetName + '&';

            var reqObj = constructReqObj();
            if (reqObj.length == 0) {
                //document.getElementById('idJsUserInputDiv').style.display = "none";
                showErrorMessage('There is no data to save, please construct a message.');

                return;
            }

            postData += 'msg=' + JSON.stringify();
            alert(postData);

            doAjaxCall('/iso/v0/save', function (response) {

                alert('data set saved successfully.');

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

    doAjaxCall('/iso/v0/loadmsg?specId=' + getCurrentSpecId() + '&msgId=' + getCurrentSpecMsgId(), function (response) {

        alert(response);

        var dataSets = JSON.parse(response);
        var htmlContent = '';
        for (var i = 0; i < dataSets.length; i++) {
            htmlContent += '<option value="' + dataSets[i] + "\">" + dataSets[i] + "</option>";

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

function jsProcessUserSelection() {

    var selectElem = document.getElementById('idJsUserSelection');
    var dsName = selectElem.options[selectElem.selectedIndex].value;

    alert('Fetching ds - ' + dsName);

    doAjaxCall('/iso/v0/loadmsg?specId=' + getCurrentSpecId() + '&msgId=' + getCurrentSpecMsgId() + '&dsName=' + dsName, function (response) {

        alert(response);

        var dataSet = JSON.parse(response);

        jsCancelUserSelectionDialog();
    }, 'GET');




}

function jsCancelUserSelectionDialog() {
    document.getElementById('idJsUserSelectionDiv').style.display = "none";
    disableOverlay();
}
