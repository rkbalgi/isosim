package iso_http

import "net/http"

type IsoHttpHandler struct {

}

func SetRoutes() {

	//default route
	http.HandleFunc("/iso/", func(rw http.ResponseWriter, req *http.Request) {

		http.ServeFile(rw, req, "C://users//132968//Desktop//iso.html")

		//rw.Write([]byte("<html><body><h2>Welcome to Web ISO Simulator </h2></body></html>"));
	});

	AllSpecsHandler();
	GetSpecMessagesHandler();
	GetMessageTemplateHandler();
}

func sendError(rw http.ResponseWriter, errorMsg string) {
	rw.WriteHeader(http.StatusBadRequest);
	rw.Write([]byte(errorMsg));


}