(this["webpackJsonpisosim-reactjs"]=this["webpackJsonpisosim-reactjs"]||[]).push([[0],{116:function(e,t,a){e.exports=a(148)},121:function(e,t,a){},122:function(e,t,a){},148:function(e,t,a){"use strict";a.r(t);var s=a(0),n=a.n(s),l=a(8),i=a.n(l),o=(a(121),a(122),a(11)),r=a(13),c=a(17),d=a(16),h=a(5),u=a(18),p=a(22),g=a.n(p),m=a(206),f=function(e){function t(e){var a;return Object(o.a)(this,t),(a=Object(c.a)(this,Object(d.a)(t).call(this,e))).state={show:a.props.show,value:a.props.value},a.closeThis=a.closeThis.bind(Object(h.a)(a)),a.valueChanged=a.valueChanged.bind(Object(h.a)(a)),a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"valueChanged",value:function(e){this.setState({value:e.target.value})}},{key:"closeThis",value:function(){this.setState({show:!0}),this.props.onClose(this.state.value)}},{key:"componentDidUpdate",value:function(e,t,a){!1===e.show&&!0===this.props.show&&!1===this.state.show?this.setState({show:!0,value:this.props.value}):!1===this.props.show&&!0===this.state.show&&this.setState({show:!1,value:this.props.value})}},{key:"render",value:function(){return!0===this.state.show?n.a.createElement(n.a.Fragment,null,n.a.createElement("div",{style:{borderBottom:"solid",borderColor:"red"}},n.a.createElement("textarea",{style:{fontFamily:"courier new",width:"100%",minHeight:"80px",maxHeight:"200px"},onChange:this.valueChanged,disabled:this.props.readOnly,value:this.state.value}),n.a.createElement("div",{style:{height:"25px"}},n.a.createElement(m.a,{size:"sm",style:{float:"right",fontSize:"10px"},onClick:this.closeThis}," OK ")))):null}}]),t}(n.a.Component),v=a(102),b=a(97),E=a(203),y=a(205),S=a(209),C=function e(){Object(o.a)(this,e),this.baseUrl="",this.sendMsgUrl=this.baseUrl+"/iso/v0/send",this.loadMsgUrl=this.baseUrl+"/iso/v0/loadmsg",this.allSpecsUrl=this.baseUrl+"/iso/v0/specs",this.templateUrl=this.baseUrl+"/iso/v0/template",this.parseTraceUrl=this.baseUrl+"/iso/v0/parse",this.saveMsgUrl=this.baseUrl+"/iso/v0/save"};C.FixedField="Fixed",C.VariableField="Variable",C.BitmappedField="Bitmapped";var w=new C,k=new(function(){function e(t){Object(o.a)(this,e),this.validate=this.validate.bind(this)}return Object(r.a)(e,[{key:"validate",value:function(e,t,a){console.log("validate",e,t,a);var s=!1;e.Type===C.FixedField&&("ASCII"===e.DataEncoding||"EBCDIC"===e.DataEncoding?t.length!==e.FixedSize&&(a.push('\u2b55 "'.concat(e.Name,'" should have a fixed size of ').concat(e.FixedSize," but has ").concat(t.length)),s=!0):t.length!==2*e.FixedSize&&(a.push('\u2b55 "'.concat(e.Name,'" should have a fixed size of ').concat(e.FixedSize," but has ").concat(t.length/2)),s=!0));var n=!1;if("BCD"!==e.DataEncoding&&"BINARY"!==e.DataEncoding||(t.length%2!==0&&(a.push('\u2b55 "'.concat(e.Name,'" should have even number of characters!')),s=!0,n=!0),"BINARY"!==e.DataEncoding||t.match("^[0-9,a-f,A-F]+$")||(a.push('\u2b55 "'.concat(e.Name,'" supports only hex i.e 0-9,a-z,A-Z')),s=!0),"BCD"!==e.DataEncoding||t.match("^[0-9]+$")||(a.push('\u2b55 "'.concat(e.Name,'" supports only bcd i.e 0-9')),s=!0)),!n&&e.Type===C.VariableField){var l=t.length;"BCD"!==e.DataEncoding&&"BINARY"!==e.DataEncoding||(l=t.length/2),e.MinSize>0&&t.length<e.MinSize&&(a.push('\u2b55 "'.concat(e.Name," size of ").concat(l," is less than required min of ").concat(e.MinSize,'" ')),s=!0),e.MaxSize>0&&t.length>e.MaxSize&&(a.push('\u2b55 "'.concat(e.Name," size  of ").concat(l," is greater than required max of ").concat(e.MinSize,'" ')),s=!0)}return s}}]),e}()),D=function(e){function t(e){var a;Object(o.a)(this,t),(a=Object(c.a)(this,Object(d.a)(t).call(this,e))).fieldValueChanged=a.fieldValueChanged.bind(Object(h.a)(a)),a.fieldSelectionChanged=a.fieldSelectionChanged.bind(Object(h.a)(a)),a.onFieldUpdate=a.onFieldUpdate.bind(Object(h.a)(a)),a.appendFieldContent=a.appendFieldContent.bind(Object(h.a)(a)),a.setSelected=a.setSelected.bind(Object(h.a)(a)),a.setNewValue=a.setNewValue.bind(Object(h.a)(a)),a.showExpanded=a.showExpanded.bind(Object(h.a)(a)),a.closeExpanded=a.closeExpanded.bind(Object(h.a)(a)),a.getBgColor=a.getBgColor.bind(Object(h.a)(a)),a.setError=a.setError.bind(Object(h.a)(a)),a.toggleExpanded=a.toggleExpanded.bind(Object(h.a)(a));if(a.selectable=!0,a.props.readOnly){a.selectable=!1;var s=!1,n=a.props.id2Value.get(a.props.field.Id);n&&(s=!0),a.state={fieldEditable:!0,bgColor:"white",hasError:!1,selected:s,id2Value:a.props.id2Value,fieldValue:n,expandBtnLabel:"+",showExpanded:!1}}else{var l="";if(["Message Type","MTI","Bitmap"].includes(a.props.field.Name)){a.selectable=!1;var i=!0;"Bitmap"===a.props.field.Name&&(l=Array(128).fill("0").reduce((function(){var e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:"",t=arguments.length>1?arguments[1]:void 0;return e+t})),i=!1),a.state={fieldEditable:i,bgColor:"white",hasError:!1,selected:!0,fieldValue:l,expandBtnLabel:"+",showExpanded:!1}}else a.state={fieldEditable:!0,bgColor:"white",selected:!1,hasError:!1,fieldValue:l,expandBtnLabel:"+",showExpanded:!1};a.props.isoMsg.set(a.props.field.Id,Object(h.a)(a))}return a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"getBgColor",value:function(){return this.state.hasError?"red":"white"}},{key:"setError",value:function(e){this.setState({hasError:e})}},{key:"showExpanded",value:function(){this.setState({showExpanded:!0})}},{key:"toggleExpanded",value:function(){this.state.showExpanded?this.setState({showExpanded:!1,expandBtnLabel:"+"}):this.setState({showExpanded:!0,expandBtnLabel:"-"})}},{key:"closeExpanded",value:function(){this.setState({showExpanded:!1})}},{key:"setNewValue",value:function(e){this.setState({fieldValue:e,showExpanded:!1})}},{key:"componentDidUpdate",value:function(e,t,a){e.id2Value!==this.props.id2Value&&this.setState({fieldValue:this.props.id2Value.get(this.props.field.Id),id2Value:this.props.id2Value})}},{key:"onFieldUpdate",value:function(e){var t=this;if(this.props.field.Type===C.BitmappedField)this.props.field.Children.forEach((function(a){if(a.Name===e.fieldName){var s=t.state.fieldValue,n=Array.from(s);if("FieldSelected"===e.ChangeType)n[a.Position-1]="1",a.Position>64&&(n[0]="1");else if("FieldDeselected"===e.ChangeType){n[a.Position-1]="0";for(var l=!0,i=65;i<=128;i++)if("1"===n[i-1]){l=!1;break}l&&(n[0]="0")}var o=n.reduce((function(){var e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:"",t=arguments.length>1?arguments[1]:void 0;return e+t}));t.setState({fieldValue:o})}}));else{var a={fieldName:this.props.field.Name};"FieldSelected"===e.ChangeType?(this.setState({selected:!0}),a.ChangeType="FieldSelected"):"FieldDeselected"===e.ChangeType&&(this.setState({selected:!1}),a.ChangeType="FieldDeselected"),this.props.field.Children.forEach((function(a){"FieldSelected"===e.ChangeType?t.props.isoMsg.get(a.Id).setSelected(!0):"FieldDeselected"===e.ChangeType&&t.props.isoMsg.get(a.Id).setSelected(!1)})),this.props.onFieldUpdate(a)}}},{key:"setSelected",value:function(e){var t=this;this.setState({selected:e}),this.props.field.Children.forEach((function(a){t.props.isoMsg.get(a.Id).setSelected(e)}))}},{key:"fieldSelectionChanged",value:function(e){var t={fieldName:this.props.field.Name},a=!1;e.target.checked?(t.ChangeType="FieldSelected",a=!0):t.ChangeType="FieldDeselected",this.props.field.Type!==C.BitmappedField&&this.setSelected(a),this.props.onFieldUpdate(t)}},{key:"fieldValueChanged",value:function(e){var t=[];if(k.validate(this.props.field,e.target.value,t))this.setState({hasError:!0,errMsg:t[0],fieldValue:e.target.value});else{this.setState({hasError:!1,errMsg:null,fieldValue:e.target.value});var a={fieldName:this.props.field.Name,ChangeType:"ValueChanged"};this.props.onFieldUpdate(a)}}},{key:"appendFieldContent",value:function(e,a,s,l,i){var o=a.Id;this.props.readOnly&&(o="response_seg_"+a.Id),e.push(n.a.createElement(t,{key:o,field:a,id2Value:l,readOnly:this.props.readOnly,parentField:s,isoMsg:this.props.isoMsg,level:i,onFieldUpdate:this.onFieldUpdate}))}},{key:"render",value:function(){var e,t=this;e=this.selectable?n.a.createElement("td",{align:"center"},n.a.createElement(y.a,{type:"checkbox",size:"small",color:"primary",checked:this.state.selected,onChange:this.fieldSelectionChanged})):n.a.createElement("td",{align:"center"},n.a.createElement(y.a,{type:"checkbox",size:"small",color:"primary",disabled:!0,checked:this.state.selected,onChange:this.fieldSelectionChanged}));var a="";this.props.field.ParentId>0&&(a="\u2937"+this.props.field.Position+" ");var s=a+" Type: "+this.props.field.Type+" / ";this.props.field.Type===C.FixedField?s+="Length: "+this.props.field.FixedSize+" / Encoding: "+this.props.field.DataEncoding:this.props.field.Type===C.VariableField?s+="Length Indicator: "+this.props.field.LengthIndicatorSize+" / Length Encoding: "+this.props.field.LengthEncoding+" / Data Encoding: "+this.props.field.DataEncoding:this.props.field.Type;var l=[];this.props.field.Children.forEach((function(e){return t.appendFieldContent(l,e,t.props.field,t.state.id2Value,t.props.level+1)}));for(var i="",o=0;o<this.props.level;o++)i+="\u2193";return n.a.createElement(n.a.Fragment,null,n.a.createElement("tr",null,e,n.a.createElement(v.a,{overlay:n.a.createElement(b.a,{id:"hi",style:{fontSize:"10px"}},s),placement:"top"},n.a.createElement("td",{style:{width:"100px",fontFamily:"lato-regular",fontSize:"14px"}},n.a.createElement(S.a,null,i+" "+this.props.field.Name))),n.a.createElement("td",null,n.a.createElement(E.a,{margin:"dense",size:"small",value:this.state.fieldValue,error:this.state.hasError,helperText:this.state.errMsg,onChange:this.fieldValueChanged,style:{width:"70%"},disabled:this.props.readOnly||!this.state.fieldEditable,key:this.props.key,ondblclick:this.showExpanded}),n.a.createElement(m.a,{size:"sm",style:{float:"right",fontSize:"10px",marginRight:"10px"},onClick:this.toggleExpanded}," ",this.state.expandBtnLabel," ")," ")),n.a.createElement("tr",null,n.a.createElement("td",{colSpan:"3"},n.a.createElement(f,{show:this.state.showExpanded,value:this.state.fieldValue,readOnly:this.props.readOnly,onClose:this.setNewValue}))),l)}}]),t}(n.a.Component),M=a(202),O=function(e){function t(e){var a;return Object(o.a)(this,t),(a=Object(c.a)(this,Object(d.a)(t).call(this,e))).state={show:e.show,selectedMsg:""},a.closeDialogSuccess=a.closeDialogSuccess.bind(Object(h.a)(a)),a.closeDialogFail=a.closeDialogFail.bind(Object(h.a)(a)),a.selectedMsgChanged=a.selectedMsgChanged.bind(Object(h.a)(a)),a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"selectedMsgChanged",value:function(e){this.setState({selectedMsg:e.target.value})}},{key:"componentDidUpdate",value:function(e,t,a){var s=this;!0===this.props.show&&!1===t.show&&g.a.get(w.loadMsgUrl,{params:{specId:this.props.specId,msgId:this.props.msgId}}).then((function(e){s.setState({savedMsgs:e.data,selectedMsg:e.data[0],show:!0})})).catch((function(e){console.log(e),s.setState({show:!0,errorMessage:e.response.data})}))}},{key:"closeDialogSuccess",value:function(){this.setState({show:!1}),this.props.closeLoadMsgDialog(this.state.selectedMsg)}},{key:"closeDialogFail",value:function(){this.setState({show:!1}),this.props.closeLoadMsgDialog(null)}},{key:"render",value:function(){var e;return this.state.show&&(e=this.state.errorMessage?n.a.createElement("div",null,this.state.errorMessage):n.a.createElement(n.a.Fragment,null,n.a.createElement("label",{style:{fontFamily:"lato-regular"}}," Saved Message "),"  ",n.a.createElement("select",{style:{fontFamily:"lato-regular",width:"200px"},value:this.state.selectedMsg,onChange:this.selectedMsgChanged},this.state.savedMsgs.map((function(e){return n.a.createElement("option",{key:e,value:e},e)}))))),n.a.createElement(M.a,{show:this.state.show,onHide:this.closeDialogFail},n.a.createElement(M.a.Header,{closeButton:!0},n.a.createElement(M.a.Title,null,"Load Saved Message")),n.a.createElement(M.a.Body,null,e),n.a.createElement(M.a.Footer,null,n.a.createElement(m.a,{variant:"primary",onClick:this.closeDialogSuccess},"OK"),n.a.createElement(m.a,{variant:"secondary",onClick:this.closeDialogFail},"Close")))}}]),t}(n.a.Component),I=a(150),F=a(192),j=a(210),T=a(194),x=a(195),N=a(196),V=a(98),U=a.n(V),R=function(e){function t(e){var a;return Object(o.a)(this,t),(a=Object(c.a)(this,Object(d.a)(t).call(this,e))).state={show:e.show,data:a.props.data,msgTemplate:a.props.msgTemplate},a.hideResponseSegment=a.hideResponseSegment.bind(Object(h.a)(a)),a.copyToClipboard=a.copyToClipboard.bind(Object(h.a)(a)),a.textAreaRef=n.a.createRef(),a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"appendFieldContent",value:function(e,t,a,s){return e.push(n.a.createElement(D,{key:"response_seg_"+t.Id,field:t,id2Value:a,readOnly:!0,level:s,onFieldUpdate:this.onFieldUpdate})),""}},{key:"hideResponseSegment",value:function(){this.setState({show:!1}),this.props.onClose()}},{key:"collectData",value:function(e,t,a){var s=this;if(t.get(e.Id)){var n=t.get(e.Id);a.push("".concat(e.Name,": ").concat(n))}e.Children.length>0&&e.Children.forEach((function(e){s.collectData(e,t,a)}))}},{key:"componentDidUpdate",value:function(e,t,a){!1===e.show&&!0===this.props.show&&this.setState({show:!0,data:this.props.data,msgTemplate:this.props.msgTemplate})}},{key:"copyToClipboard",value:function(){this.textAreaRef.current.select(),document.execCommand("copy")||alert("Failed to copy to clipboard!")}},{key:"render",value:function(){var e=this,t=[],a=[];if(this.state.show){var s=new Map;this.state.data.forEach((function(e){s.set(e.Id,e.Value)})),this.state.msgTemplate.Fields.forEach((function(t){e.collectData(t,s,a)}));var l="ISO Response  \n|---------------|\n"+a.reduce((function(e,t,a){return 1===a?e+"\n"+t+"\n":e+t+"\n"}));return l="ISO Request  \n|---------------|\n"+this.props.reqData+"\n\n"+l+"\n\n",this.state.msgTemplate.Fields.forEach((function(a){e.appendFieldContent(t,a,s,0)})),n.a.createElement(n.a.Fragment,null,this.state.show?n.a.createElement(j.a,{open:this.state.show,onClose:this.hideResponseSegment,scroll:"paper",PaperComponent:B,"aria-labelledby":"draggable-dialog-title",maxWidth:"sm",fullWidth:!0,disableBackdropClick:!0},n.a.createElement(T.a,{style:{cursor:"move"},id:"draggable-dialog-title"},this.props.dialogTitle),n.a.createElement(x.a,{dividers:!0},n.a.createElement(F.a,null,n.a.createElement("textarea",{ref:this.textAreaRef,style:{opacity:"0.01",position:"absolute",zIndex:-9999,height:0}},l),n.a.createElement("table",{border:"0",align:"center"},n.a.createElement("thead",null,n.a.createElement("tr",{style:{fontFamily:"lato-regular",backgroundColor:"#eed143",fontSize:"15px",align:"center",borderBottom:"solid",borderColor:"blue"}},n.a.createElement("td",{colSpan:"3",align:"center"},"Response Segment")),n.a.createElement("tr",{style:{fontFamily:"lato-regular",backgroundColor:"#3effba",fontSize:"14px"}},n.a.createElement("td",{align:"center"},"Selection"),n.a.createElement("td",{align:"center",style:{width:"35%"}},"Field"),n.a.createElement("td",{align:"center",style:{width:"50%"}},"Field Data"))),n.a.createElement("tbody",null,t)))),n.a.createElement(N.a,null,n.a.createElement(I.a,{onClick:this.copyToClipboard,size:"small",color:"primary",variant:"contained"},"Copy To Clipboard"),n.a.createElement(I.a,{onClick:this.hideResponseSegment,size:"small",color:"primary",variant:"contained"},"Close"))):null)}return null}}]),t}(n.a.Component);function B(e){return n.a.createElement(U.a,{handle:"#draggable-dialog-title",cancel:'[class*="MuiDialogContent-root"]'},n.a.createElement(F.a,e))}var z=function(e){function t(e){var a;return Object(o.a)(this,t),(a=Object(c.a)(this,Object(d.a)(t).call(this,e))).state={show:e.show,traceMsg:""},a.closeDialogSuccess=a.closeDialogSuccess.bind(Object(h.a)(a)),a.closeDialogFail=a.closeDialogFail.bind(Object(h.a)(a)),a.traceChanged=a.traceChanged.bind(Object(h.a)(a)),a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"traceChanged",value:function(e){this.setState({traceMsg:e.target.value})}},{key:"componentDidUpdate",value:function(e,t,a){!0===this.props.show&&!1===t.show&&this.setState({show:!0,traceMsg:""})}},{key:"closeDialogSuccess",value:function(){this.setState({show:!1}),this.props.setTrace(this.state.traceMsg)}},{key:"closeDialogFail",value:function(){this.setState({show:!1}),this.props.setTrace(null)}},{key:"render",value:function(){var e;return this.state.show&&(e=this.state.errorMessage?n.a.createElement("div",null,this.state.errorMessage):n.a.createElement(n.a.Fragment,null,n.a.createElement("label",{style:{fontFamily:"lato-regular"}}," Trace "),"  ",n.a.createElement("textarea",{key:"trace_input",value:this.state.traceMsg,onChange:this.traceChanged,style:{fontFamily:"courier new",width:"100%"}}))),n.a.createElement(M.a,{show:this.state.show,onHide:this.closeDialogFail},n.a.createElement(M.a.Header,{closeButton:!0},n.a.createElement(M.a.Title,null,"Parse Trace")),n.a.createElement(M.a.Body,null,e),n.a.createElement(M.a.Footer,null,n.a.createElement(m.a,{variant:"primary",onClick:this.closeDialogSuccess},"OK"),n.a.createElement(m.a,{variant:"secondary",onClick:this.closeDialogFail},"Close")))}}]),t}(n.a.Component),L=function(e){function t(e){var a;return Object(o.a)(this,t),(a=Object(c.a)(this,Object(d.a)(t).call(this,e))).state={show:e.show,msgName:""},a.closeDialogSuccess=a.closeDialogSuccess.bind(Object(h.a)(a)),a.closeDialogFail=a.closeDialogFail.bind(Object(h.a)(a)),a.msgNameChanged=a.msgNameChanged.bind(Object(h.a)(a)),a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"msgNameChanged",value:function(e){this.setState({errorMessage:"",msgName:e.target.value})}},{key:"componentDidUpdate",value:function(e,t,a){!0===this.props.show&&!1===t.show&&this.setState({show:!0,msgName:this.props.msgName})}},{key:"closeDialogSuccess",value:function(){var e=this;if(this.state.msgName&&""!==this.state.msgName&&this.props.data){var t="specId="+this.props.specId+"&msgId="+this.props.msgId+"&dataSetName="+this.state.msgName+"&msg="+JSON.stringify(this.props.data);g.a.post(w.saveMsgUrl,t).then((function(t){console.log(t),e.props.msgSaveSuccess(e.state.msgName),e.setState({show:!1})})).catch((function(t){e.props.msgSaveFailed(t),e.setState({show:!1})}))}else this.setState({errorMessage:"Please specify a message!"})}},{key:"closeDialogFail",value:function(){this.props.msgSaveCancelled(),this.setState({show:!1})}},{key:"render",value:function(){var e,t;return this.state.show&&(console.log("before sending",this.props),this.state.errorMessage&&(t=n.a.createElement("div",{style:{color:"red"}},this.state.errorMessage)),e=this.props.msgId&&this.props.specId?n.a.createElement(n.a.Fragment,null,n.a.createElement("label",{style:{fontFamily:"lato-regular"}}," Message Name "),"  ",n.a.createElement("input",{type:"text",key:"msg_name_save",value:this.state.msgName,onChange:this.msgNameChanged}),t):n.a.createElement("div",null,"Error: Please load a spec/msg, set data before attempting to save")),n.a.createElement(M.a,{show:this.state.show,onHide:this.closeDialogFail},n.a.createElement(M.a.Header,{closeButton:!0},n.a.createElement(M.a.Title,null,"Save Message")),n.a.createElement(M.a.Body,null,e),n.a.createElement(M.a.Footer,null,n.a.createElement(m.a,{variant:"primary",onClick:this.closeDialogSuccess},"OK"),n.a.createElement(m.a,{variant:"secondary",onClick:this.closeDialogFail},"Close")))}}]),t}(n.a.Component),P=a(99),q=a.n(P),A=(a(145),a(103)),H=a(193),_=a(207),W=a(191),$=function(e){function t(e){var a;return Object(o.a)(this,t),(a=Object(c.a)(this,Object(d.a)(t).call(this,e))).state={targetServerIp:"127.0.0.1",targetServerPort:"6666",mliType:"2i"},a.serverIpChanged=a.serverIpChanged.bind(Object(h.a)(a)),a.serverPortChanged=a.serverPortChanged.bind(Object(h.a)(a)),a.mliTypeChanged=a.mliTypeChanged.bind(Object(h.a)(a)),a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"mliTypeChanged",value:function(e){this.setState({mliType:e.target.value}),this.props.onChange(this.state.targetServerIp,this.state.targetServerPort,e.target.value)}},{key:"serverIpChanged",value:function(e){this.setState({targetServerIp:e.target.value}),this.props.onChange(e.target.value,this.state.targetServerPort,this.state.mliType)}},{key:"serverPortChanged",value:function(e){this.setState({targetServerPort:e.target.value}),this.props.onChange(this.state.targetServerIp,e.target.value,this.state.mliType)}},{key:"render",value:function(){return n.a.createElement("div",{align:"left",style:{align:"left",height:"80px",verticalAlign:"baseline",margin:"10px"}},n.a.createElement("table",{style:{fontFamily:"lato-regular",fontSize:"14px"}},n.a.createElement("tr",null,n.a.createElement("td",null,n.a.createElement(E.a,{id:"outlined-basic",label:"IP",size:"small",variant:"outlined",defaultValue:"127.0.0.1",onChange:this.serverIpChanged})),n.a.createElement("td",null,n.a.createElement(E.a,{id:"outlined-basic",label:"Port",size:"small",variant:"outlined",defaultValue:"6666",onChange:this.serverPortChanged})),n.a.createElement("td",null,n.a.createElement(E.a,{select:!0,size:"small",value:this.state.mliType,variant:"outlined",label:"MLI",onChange:this.mliTypeChanged},n.a.createElement(_.a,{value:"2i"},"2I"),n.a.createElement(_.a,{value:"2e"},"2E"))))))}}]),t}(n.a.Component),J=a(197),K=function(e){function t(e){var a;return Object(o.a)(this,t),a=Object(c.a)(this,Object(d.a)(t).call(this,e)),console.log(a.props),console.log("$msg_structure$",a.props.specs,a.props.spec,a.props.msg),a.state={msgTemplate:null,loaded:!1,spec:e.spec,msg:e.msg,shouldShow:e.showMsgTemplate,targetServerIp:"127.0.0.1",targetServerPort:"6666",mliType:"2I",currentDataSet:"",errDialogVisible:!1,errorMessage:"",showLoadMessagesDialog:!1,showTraceInputDialog:!1,showSaveMsgDialog:!1,showResponse:!1,responseData:null,reqMenuVisible:!1,selectedReqMenuItem:null,reqClipboardData:null},a.onFieldUpdate=a.onFieldUpdate.bind(Object(h.a)(a)),a.appendFieldContent=a.appendFieldContent.bind(Object(h.a)(a)),a.sendToHost=a.sendToHost.bind(Object(h.a)(a)),a.addFieldContent=a.addFieldContent.bind(Object(h.a)(a)),a.showErrorDialog=a.showErrorDialog.bind(Object(h.a)(a)),a.closeErrorDialog=a.closeErrorDialog.bind(Object(h.a)(a)),a.processError=a.processError.bind(Object(h.a)(a)),a.showLoadMessagesDialog=a.showLoadMessagesDialog.bind(Object(h.a)(a)),a.closeLoadMsgDialog=a.closeLoadMsgDialog.bind(Object(h.a)(a)),a.showUnImplementedError=a.showUnImplementedError.bind(Object(h.a)(a)),a.setTrace=a.setTrace.bind(Object(h.a)(a)),a.showTraceInputsDialog=a.showTraceInputsDialog.bind(Object(h.a)(a)),a.showSaveMsgDialog=a.showSaveMsgDialog.bind(Object(h.a)(a)),a.msgSaveSuccess=a.msgSaveSuccess.bind(Object(h.a)(a)),a.msgSaveFailed=a.msgSaveFailed.bind(Object(h.a)(a)),a.msgSaveCancelled=a.msgSaveCancelled.bind(Object(h.a)(a)),a.showInfoDialog=a.showInfoDialog.bind(Object(h.a)(a)),a.showMenu=a.showMenu.bind(Object(h.a)(a)),a.hideMenu=a.hideMenu.bind(Object(h.a)(a)),a.handleMenuClick=a.handleMenuClick.bind(Object(h.a)(a)),a.showResponseDialog=a.showResponseDialog.bind(Object(h.a)(a)),a.getTemplateLabel=a.getTemplateLabel.bind(Object(h.a)(a)),a.networkSettingsChanged=a.networkSettingsChanged.bind(Object(h.a)(a)),a.hideResponse=a.hideResponse.bind(Object(h.a)(a)),a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"networkSettingsChanged",value:function(e,t,a){this.setState({targetServerIp:e,targetServerPort:t,mliType:a})}},{key:"showMenu",value:function(e){this.setState({selectedReqMenuItem:e.currentTarget,reqMenuVisible:!0})}},{key:"hideMenu",value:function(){this.setState({reqMenuVisible:!1}),this.setState({selectedReqMenuItem:null})}},{key:"showResponseDialog",value:function(){this.hideMenu(),this.setState({showResponse:!0})}},{key:"hideResponse",value:function(){this.setState({showResponse:!1})}},{key:"handleMenuClick",value:function(e){alert(e.currentTarget),this.setState({selectedReqMenuItem:e.currentTarget}),this.hideMenu()}},{key:"setTrace",value:function(e){var t=this;null!=e&&g.a.post(w.parseTraceUrl+"/"+this.state.spec.Id+"/"+this.state.msg.Id,e).then((function(e){console.log("parsed msg data",e.data),e.data.forEach((function(e){t.state.isoMsg.get(e.Id).setState({selected:!0,fieldValue:e.Value})}))})).catch((function(e){console.log(e),t.processError(e)})),this.setState({showTraceInputDialog:!1})}},{key:"showUnImplementedError",value:function(){this.setState({errorMessage:"This functionality has not been implemented. Please try the old version of application.",errDialogVisible:!0})}},{key:"closeLoadMsgDialog",value:function(e){var t=this;this.setState({showLoadMessagesDialog:!1,currentDataSet:e}),null!=e&&g.a.get(w.loadMsgUrl,{params:{specId:this.state.spec.Id,msgId:this.state.msg.Id,dsName:e}}).then((function(e){e.data.forEach((function(e){t.state.isoMsg.get(e.Id).setState({selected:!0,fieldValue:e.Value})}))})).catch((function(e){console.log(e),t.processError(e)}))}},{key:"showInfoDialog",value:function(e){this.setState({errDialogVisible:!0,errorMessage:e})}},{key:"msgSaveSuccess",value:function(e){this.showInfoDialog("Message ".concat(e," saved successfully.")),this.setState({showSaveMsgDialog:!1})}},{key:"msgSaveFailed",value:function(e){this.processError(e),this.setState({showSaveMsgDialog:!1})}},{key:"msgSaveCancelled",value:function(){this.setState({showSaveMsgDialog:!1})}},{key:"showSaveMsgDialog",value:function(){var e=this,t=[];this.state.msgTemplate.Fields.forEach((function(a){e.addFieldContent(a,t)})),this.setState({saveData:t,showSaveMsgDialog:!0})}},{key:"showTraceInputsDialog",value:function(){this.hideMenu(),this.setState({showTraceInputDialog:!0})}},{key:"showLoadMessagesDialog",value:function(){this.hideMenu(),this.setState({showLoadMessagesDialog:!0})}},{key:"closeErrorDialog",value:function(){this.setState({errDialogVisible:!1})}},{key:"showErrorDialog",value:function(){this.setState({errDialogVisible:!0})}},{key:"addFieldContent",value:function(e,t,a){var s=this,n=this.state.isoMsg.get(e.Id);n.state.selected&&(k.validate(e,n.state.fieldValue,a)?n.setError(!0):n.setError(!1),t.push({Id:e.Id,Name:e.Name,Value:n.state.fieldValue})),e.Children.forEach((function(e){n.state.selected&&s.addFieldContent(e,t,a)}))}},{key:"sendToHost",value:function(){var e=this;this.hideMenu();var t=[],a=[];if(this.state.msgTemplate.Fields.forEach((function(s){e.addFieldContent(s,t,a)})),a.length>0){var s="";return a.forEach((function(e){return s+=e+"\n"})),this.setState({errorMessage:s}),void this.showErrorDialog()}console.log(t);var n=t.reduce((function(e,t,a){return 1===a?e.Name+":"+e.Value+"\n"+t.Name+":"+t.Value+"\n":e+t.Name+":"+t.Value+"\n"}));this.setState({showResponse:!1,responseData:null,reqClipboardData:n});var l="host="+this.state.targetServerIp+"&port="+this.state.targetServerPort+"&mli="+this.state.mliType+"&specId="+this.state.spec.Id+"&msgId="+this.state.msg.Id+"&msg="+JSON.stringify(t);console.log(l),g.a.post(w.sendMsgUrl,l).then((function(t){console.log("Response from server",t),e.setState({showResponse:!0,responseData:t.data})})).catch((function(t){console.log("error = ",t),e.processError(t)}))}},{key:"processError",value:function(e){if(!e.response)return console.log(e),void this.setState({errorMessage:"Error: Unable to reach API server",errDialogVisible:!0});400===e.response.status?this.setState({errorMessage:e.response.data,errDialogVisible:!0}):this.setState({errorMessage:"Unexpected error from server - "+e.response.status,errDialogVisible:!0})}},{key:"getTemplateLabel",value:function(){return this.state.spec.Name+" // "+this.state.msg.Name}},{key:"onFieldUpdate",value:function(e){}},{key:"componentDidMount",value:function(){this.getMessageTemplate(this.props.spec,this.props.msg)}},{key:"getMessageTemplate",value:function(e,t){var a=this,s=this.props.specs.find((function(t){return t.Name===e?t:null})),n=s.Messages.find((function(e){return e.Name===t?e:null})),l=w.templateUrl+"/"+s.Id+"/"+n.Id;console.log(l),g.a.get(l).then((function(e){console.log(e.data);var t=new Map;t.set("msg_template",e.data),a.setState({spec:s,msg:n,msgTemplate:e.data,loaded:!0,isoMsg:t}),console.log("MsgTemplate = ",a.state.msgTemplate)})).catch((function(e){return a.setState({errorMessage:e,errDialogVisible:!0})}))}},{key:"appendFieldContent",value:function(e,t,a,s){e.push(n.a.createElement(D,{key:t.Id,field:t,isoMsg:a,level:s,onFieldUpdate:this.onFieldUpdate}))}},{key:"render",value:function(){var e=this,t=[];return!0===this.state.loaded&&this.state.msgTemplate.Fields.map((function(a){e.appendFieldContent(t,a,e.state.isoMsg,0)})),n.a.createElement("div",{style:{fontFamily:"lato-regular",fontSize:"12px",fill:"aqua"}},n.a.createElement(M.a,{show:this.state.errDialogVisible,onHide:this.closeErrorDialog},n.a.createElement(M.a.Header,{closeButton:!0},n.a.createElement(M.a.Title,null,"Error")),n.a.createElement(M.a.Body,null,n.a.createElement("pre",{style:{font:"Lato",fontSize:"14px"}},this.state.errorMessage)),n.a.createElement(M.a.Footer,null,n.a.createElement(I.a,{variant:"secondary",onClick:this.closeErrorDialog},"Close"))),n.a.createElement(O,{show:this.state.showLoadMessagesDialog,specId:this.state.spec.Id,msgId:this.state.msg.Id,closeLoadMsgDialog:this.closeLoadMsgDialog}),n.a.createElement(z,{show:this.state.showTraceInputDialog,setTrace:this.setTrace}),n.a.createElement(L,{show:this.state.showSaveMsgDialog,msgId:this.state.msg.Id,specId:this.state.spec.Id,data:this.state.saveData,msgName:this.state.currentDataSet,msgSaveSuccess:this.msgSaveSuccess,msgSaveFailed:this.msgSaveFailed,msgSaveCancelled:this.msgSaveCancelled}),n.a.createElement($,{onChange:this.networkSettingsChanged}),n.a.createElement("div",{align:"left",style:{align:"left",display:"inline-block",width:"40%",float:"left",fill:"aqua"}},n.a.createElement("div",null,n.a.createElement(J.a,{size:"small",color:"primary",fullWidth:!0,variant:"contained"},n.a.createElement(I.a,{onClick:this.showTraceInputsDialog},"Parse"),n.a.createElement(I.a,{onClick:this.showLoadMessagesDialog},"Load"),n.a.createElement(I.a,{onClick:this.showSaveMsgDialog},"Save"),n.a.createElement(I.a,{onClick:this.sendToHost},"Send"),n.a.createElement(I.a,{onClick:this.showResponseDialog,disabled:null==this.state.responseData},"Show Response"))),n.a.createElement(F.a,{variation:"outlined",style:{verticalAlign:"middle"}},n.a.createElement("table",{border:"0",align:"center",style:{align:"center",marginTop:"10px",width:"70%"}},n.a.createElement("thead",null,n.a.createElement("tr",{style:{fontFamily:"lato-regular",backgroundColor:"#ff8f5b",fontSize:"15px",borderBottom:"solid",borderColor:"blue"}},n.a.createElement("td",{colSpan:"3",align:"center"},n.a.createElement("div",{style:{display:"inline-block",float:"left"}},n.a.createElement(W.a,{"aria-label":"more","aria-controls":"long-menu","aria-haspopup":"true",onClick:this.showMenu},n.a.createElement(q.a,null)),n.a.createElement(A.a,{id:"fade-menu",anchorEl:this.state.selectedReqMenuItem,getContentAnchorEl:null,keepMounted:!0,open:this.state.reqMenuVisible,onClose:this.hideMenu,TransitionComponent:H.a},n.a.createElement(_.a,{dense:!0,onClick:this.showTraceInputsDialog},"Parse"),n.a.createElement(_.a,{dense:!0,onClick:this.showLoadMessagesDialog},"Load Message"),n.a.createElement(_.a,{dense:!0,onClick:this.showSaveMsgDialog},"Save Message"),n.a.createElement(_.a,{dense:!0,onClick:this.sendToHost},"Send Message"),n.a.createElement(_.a,{dense:!0,onClick:this.showResponseDialog},"Show Response"))),n.a.createElement("div",{style:{display:"inline-block"}},this.getTemplateLabel()))),n.a.createElement("tr",{style:{fontFamily:"lato-regular",backgroundColor:"#ff8f5b",fontSize:"14px"}},n.a.createElement("td",{align:"center"},"Selection"),n.a.createElement("td",{align:"center",style:{width:"35%"}}," Field"),n.a.createElement("td",{align:"center",style:{width:"70%"}},"Field Data"))),n.a.createElement("tbody",null,t))),n.a.createElement(R,{show:this.state.showResponse,reqData:this.state.reqClipboardData,onClose:this.hideResponse,data:this.state.responseData,dialogTitle:"Response - ["+this.getTemplateLabel()+"]",msgTemplate:this.state.msgTemplate})),n.a.createElement("div",{style:{height:"10px"}}," "))}}]),t}(n.a.Component),Y=a(100),Z=a.n(Y),G=a(101),Q=a.n(G),X=a(208),ee=a(200),te=a(198),ae=a(199),se=a(201),ne=function(e){function t(e){var a;return Object(o.a)(this,t),(a=Object(c.a)(this,Object(d.a)(t).call(this,e))).state={specs:[],loaded:!1,errDialogVisible:!1,errorMessage:""},a.messageClicked=a.messageClicked.bind(Object(h.a)(a)),a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"messageClicked",value:function(e){var t=e.target.parentElement.parentElement.getAttribute("sid"),a=e.target.parentElement.parentElement.getAttribute("mid");this.props.msgSelected(t,a)}},{key:"componentDidMount",value:function(){var e=this;g.a.get(w.allSpecsUrl).then((function(t){console.log(t.data),e.setState({specs:t.data,loaded:!0})})).catch((function(e){return console.log(e)}))}},{key:"buildMessages",value:function(e){var t=this,a=[];return e.Messages.forEach((function(s){a.push(n.a.createElement(X.a,{nodeId:"nodeId_"+e.Id+"_"+s.Id,sid:e.Id,mid:s.Id,label:s.Name,onClick:t.messageClicked}))})),a}},{key:"render",value:function(){var e=this;if(!0===this.state.loaded){var t=[];this.state.specs.forEach((function(a){t.push(n.a.createElement(X.a,{align:"left",nodeId:"nodeId_"+a.Id,icon:n.a.createElement(te.a,{color:"primary"}),label:a.Name},e.buildMessages(a)))}));var a=n.a.createElement(X.a,{nodeId:"nodeId_0",icon:n.a.createElement(ae.a,{color:"primary"}),label:"ISO8583 Specifications"},t);return n.a.createElement(n.a.Fragment,null,n.a.createElement(ee.a,{defaultExpanded:["nodeId_0"],defaultCollapseIcon:n.a.createElement(Z.a,null),defaultExpandIcon:n.a.createElement(Q.a,null),defaultParentIcon:n.a.createElement(te.a,{color:"primary"}),defaultEndIcon:n.a.createElement(se.a,{color:"primary"})},a))}return null}}]),t}(n.a.Component),le=function(e){function t(e){var a;return Object(o.a)(this,t),(a=Object(c.a)(this,Object(d.a)(t).call(this,e))).state={specs:[],currentSpec:"Select",currentSpecMsg:"",showMsgTemplate:!1,loaded:!1,errDialogVisible:!1,errorMessage:""},a.specChanged=a.specChanged.bind(Object(h.a)(a)),a.messageChanged=a.messageChanged.bind(Object(h.a)(a)),a.msgSelected=a.msgSelected.bind(Object(h.a)(a)),a.getSpecByID=a.getSpecByID.bind(Object(h.a)(a)),a.msgTemplateRef=n.a.createRef(),a}return Object(u.a)(t,e),Object(r.a)(t,[{key:"msgSelected",value:function(e,t){console.log(e,t),console.log(this.state.specs);var a=this.getSpecByID(parseInt(e));console.log("spec = ",a);var s=null;a.Messages.forEach((function(e){e.Id===parseInt(t)&&(s=e)})),this.setState({loaded:!0,currentSpec:a.Name,currentSpecMsg:s.Name})}},{key:"closeErrorDialog",value:function(){this.setState({errDialogVisible:!1})}},{key:"showErrorDialog",value:function(){this.setState({errDialogVisible:!0})}},{key:"componentDidMount",value:function(){var e=this;g.a.get(w.allSpecsUrl).then((function(t){console.log(t.data),e.setState({specs:t.data,loaded:!0})})).catch((function(e){return console.log(e)}))}},{key:"render",value:function(){var e,t;return!0===this.state.loaded&&(null==(t=this.getCurrentSpec())&&(t=this.state.specs[0]),e=this.state.currentSpecMsg?this.state.currentSpecMsg:t.Messages[0].Name),n.a.createElement(n.a.Fragment,null,n.a.createElement("div",null,n.a.createElement(M.a,{show:this.state.errDialogVisible,onHide:this.closeErrorDialog},n.a.createElement(M.a.Header,{closeButton:!0},n.a.createElement(M.a.Title,null,"Error")),n.a.createElement(M.a.Body,null,this.state.errorMessage),n.a.createElement(M.a.Footer,null,n.a.createElement(m.a,{variant:"secondary",onClick:this.closeErrorDialog},"Close"))),n.a.createElement("div",{style:{float:"left",display:"inline-block",marginRight:"20px",marginLeft:"20px",backgroundColor:"#fbfff0"}},n.a.createElement(ne,{msgSelected:this.msgSelected})),n.a.createElement("div",{align:"center",style:{backgroundColor:"#fbfff0"}},this.state.loaded&&"Select"!==this.state.currentSpec?n.a.createElement(K,{key:this.state.currentSpec+"_"+e,ref:this.msgTemplateRef,specs:this.state.specs,spec:this.state.currentSpec,msg:this.state.currentSpecMsg}):null)))}},{key:"specChanged",value:function(e){if(this.setState({currentSpec:e.target.value,currentSpecMsg:""}),console.log(e.target.value),this.state.loaded&&"Select"!==e.target.value){console.log("calling update - specChanged");this.getSpecByName(e.target.value)}}},{key:"messageChanged",value:function(e){this.setState({currentSpecMsg:e.target.value}),this.state.loaded&&"Select"!==this.state.currentSpec&&console.log("calling update - msgChanged")}},{key:"specsDropDown",value:function(){return n.a.createElement("select",{style:{fontFamily:"lato-regular",width:"200px"},onChange:this.specChanged},n.a.createElement("option",{key:"Select",value:"Select"},"Select"),this.state.specs.map((function(e){return n.a.createElement("option",{key:e.Name,value:e.Name},e.Name)})))}},{key:"messagesDropDown",value:function(){var e;return this.state.loaded&&(e=this.getCurrentSpec()),"Select"===this.state.currentSpec?n.a.createElement("select",null):n.a.createElement("select",{value:this.state.currentSpecMsg,style:{fontFamily:"lato-regular",width:"150px"},onChange:this.messageChanged},e.Messages.map((function(e){return n.a.createElement("option",{key:e.Id,value:e.Name},e.Name)})))}},{key:"getCurrentSpec",value:function(){var e=this;return this.state.specs.find((function(t,a){return t.Name===e.state.currentSpec?t:null}))}},{key:"getSpecByName",value:function(e){return this.state.specs.find((function(t,a){return t.Name===e?t:null}))}},{key:"getSpecByID",value:function(e){return this.state.specs.find((function(t,a){return t.Id===e?t:null}))}}]),t}(n.a.Component);var ie=function(){return n.a.createElement("div",{style:{backgroundColor:"#fbfff0"}},n.a.createElement("h1",{style:{fontFamily:"shadows-into-light"}},"ISO WebSim - ISO8583 Web Simulator"),n.a.createElement("a",{style:{fontFamily:"lato-regular",fontSize:"12px"},href:"/iso/home",target:"_blank"},"[Non React App]"),n.a.createElement("a",{style:{fontFamily:"lato-regular",fontSize:"12px"},href:"/iso/v0/server",target:"_blank"},"[Manage Servers]"),n.a.createElement("div",{className:"App"},n.a.createElement(le,null)))};Boolean("localhost"===window.location.hostname||"[::1]"===window.location.hostname||window.location.hostname.match(/^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/));a(147);i.a.render(n.a.createElement(ie,{style:{backgroundColor:"#fbfff0"}}),document.getElementById("root")),"serviceWorker"in navigator&&navigator.serviceWorker.ready.then((function(e){e.unregister()})).catch((function(e){console.error(e.message)}))}},[[116,1,2]]]);
//# sourceMappingURL=main.b853c985.chunk.js.map