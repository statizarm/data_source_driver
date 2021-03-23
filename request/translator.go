package request

import (
	"bytes"
	"encoding/json"
	"io"
)

// TranslateToESForm function to translate unified get request to es get search request
func TranslateToESForm(req *GetRequest) (io.Reader, error) {
	if simForm, err := translateToSimpleForm(req); err != nil {
		return nil, err
	} else if esForm, err := translateToESSpecificForm(simForm); err != nil {
		return nil, err
	} else if b, err := json.Marshal(esForm); err != nil {
		return nil, err
	} else {
		return bytes.NewReader(b), nil
	}
}
