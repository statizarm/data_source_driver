package sender

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ProQSoftware/httputil"

	"github.com/ProQSoftware/data_source_driver/request"
)

// ESSender структура источника данных
type ESSender struct {
	url string
}

// SendGetDataRequest функция для GET запроса
func (s *ESSender) SendGetDataRequest(w http.ResponseWriter, dataReq *request.GetRequest) (err error) {
	// Транслируем запрос
	if r, e := request.TranslateToESForm(dataReq); e != nil {
		err = httputil.NewError(http.StatusBadRequest, "text/html", e.Error())
	} else if resp, e := httputil.SendRequest("GET", s.url+dataReq.Source, "application/json", r); e != nil {
		err = e
	} else {
		// Обработка ответа
		w.Header().Add("Content-Type", "application/json")

		var b []byte
		if b, e = getSourceFromESAns(resp.Body); e != nil {
			err = httputil.NewError(http.StatusInternalServerError, "text/html", e.Error())
		} else if _, e = w.Write(b); e != nil {
			err = httputil.NewError(http.StatusInternalServerError, "text/html", e.Error())
		}

		if e = resp.Body.Close(); e != nil {
			panic(e)
		}
	}
	return
}

func getSourceFromESAns(r io.Reader) ([]byte, error) {
	dec := json.NewDecoder(r)

	var data map[string]interface{}
	if err := dec.Decode(&data); err != nil {
		return []byte{}, err
	}

	return json.Marshal(getSource(data))
}

func getSource(jstr interface{}) interface{} {
	switch str := jstr.(type) {
	case map[string]interface{}:
		if v, ok := str["_source"]; ok {
			return v
		}

		res := make([]interface{}, 0)
		for _, v := range str {
			if r := getSource(v); r != nil {
				res = append(res, r)
			}
		}

		if len(res) == 0 {
			return nil
		} else if len(res) == 1 {
			return res[0]
		}

		return res
	case []interface{}:
		var res = make([]interface{}, 0)
		for _, v := range str {
			if r := getSource(v); r != nil {
				res = append(res, r)
			}
		}

		if len(res) == 0 {
			return nil
		} else if len(res) == 1 {
			return res[0]
		}

		return res
	default:
		break
	}
	return nil
}

//NewEsSender некая сущность для интерфейса
func NewEsSender(url string) *ESSender {
	return &ESSender{
		url: url,
	}
}
