package isoserv_handlers

import (
	"net/http"
	"github.com/rkbalgi/isosim/data"
	"encoding/json"
	"github.com/rkbalgi/isosim/web/spec"
	"log"
)

//This function will register a handler that will save incoming server definitions into a file

func fetchDefHandler(){

	http.HandleFunc("/iso/v0/server/defs/fetch",func (rw http.ResponseWriter,req *http.Request){

		req.ParseForm();
		strSpecId:=req.Form.Get("specId");
		if len(strSpecId)==0{
			sendError(rw,"Invalid or missing parameter 'specId'")
			return;
		}

		serverDefs,err:=data.DataSetManager().GetServerDefs(strSpecId);
		if spec.DebugEnabled{
			log.Print("Server Defs = ",len(serverDefs),serverDefs);
		}
		if err!=nil{
			sendError(rw,err.Error());
			return;
		}
		json.NewEncoder(rw).Encode(serverDefs);



	});

	http.HandleFunc("/iso/v0/server/defs/get",func (rw http.ResponseWriter,req *http.Request){

		req.ParseForm();
		strSpecId:=req.Form.Get("specId");
		fileName:=req.Form.Get("name");
		if len(strSpecId)==0 || len(fileName)==0{
			sendError(rw,"Invalid or missing parameter 'specId' or 'name'")
			return;
		}

		data,err:=data.DataSetManager().GetServerDef(strSpecId,fileName);
		log.Print("Def = "+string(data))

		if err!=nil{
			sendError(rw,err.Error());
			return;
		}
		rw.Write(data);



	});
}
