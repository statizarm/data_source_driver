package main
//
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ProQSoftware/httputil"

	"github.com/ProQSoftware/data_source_driver/request"
	"github.com/ProQSoftware/data_source_driver/sender"
)

// requestSender содержит методы обращения к бд
type requestSender interface {
	SendGetDataRequest(w http.ResponseWriter, getReq *request.GetRequest) error
}

type handler func(w http.ResponseWriter, r *http.Request) error
type getDataHandler func(w http.ResponseWriter, getReq *request.GetRequest) error

var senders = map[string]requestSender {
	"es": sender.NewEsSender("http://10.128.190.49:9200/"),
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("too few args: first arg - source type (e.g. es), second arg - addr")
	} else if len(os.Args) > 2 {
		log.Fatal("too much args")
	}

	addr := os.Args[1]

	for k, v := range senders {
		http.Handle("/" + k, getDataHandler(v.SendGetDataRequest))
	}

	//обработчики
	http.Handle("/status", handler(statusHandler))

	log.Fatal(http.ListenAndServe(addr, nil))

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := h(w, r); e != nil {
		httputil.WriteHttpError(w, e)
	}
}

func (h getDataHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	var getReq request.GetRequest
	var err error

	if r.Method == "POST" {
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			err = httputil.NewError(http.StatusUnsupportedMediaType, "text/html", "Expected JSON")
		} else if wholeReqBody, e := ioutil.ReadAll(r.Body); e != nil {
			err = httputil.NewError(http.StatusBadRequest, "text/html", e.Error())
		} else if e = json.Unmarshal(wholeReqBody, &getReq); e != nil {
			errmsg := fmt.Sprintf(`Unknown type of request:\n%senderObj`, wholeReqBody)
			err = httputil.NewError(http.StatusBadRequest, "text/html", errmsg)
		} else {
			err = h(w, &getReq)
		}
	} else {
		err = httputil.NewError(http.StatusMethodNotAllowed, "test/html",
			`Method "`+r.Method+`" not allowed`)
	}

	if err != nil {
		httputil.WriteHttpError(w, err)
	}
}

//статус работоспособности драйвера
func statusHandler(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Method == "POST" {
		w.Header().Add("Content-Type", "application/json")

		_, _ = fmt.Fprint(w, `{"status":"working"}`)
	} else {
		err = httputil.NewError(http.StatusMethodNotAllowed, "text/html", "Wrong method")
	}
	return
}
