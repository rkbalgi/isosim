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

            postData += 'msg=' + JSON.stringify(constructReqObj());
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
