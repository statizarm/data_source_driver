package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/ProQSoftware/httputil"
)

type simplifiedRequest struct {
	ReqExpr   *RequirementExpression
	Modifiers []map[string]json.RawMessage
	Fields    []Field
}

func translateToSimpleForm(req *GetRequest) (*simplifiedRequest, error) {
	sRequest := &simplifiedRequest{}
	if e := simplifyRequirementExpression(req.HookURL, req.Requirements); e != nil {
		return nil, e
	}

	sRequest.ReqExpr = req.Requirements
	sRequest.Modifiers = req.Modifiers
	sRequest.Fields = req.Fields

	return sRequest, nil
}

func simplifyRequirementExpression(hookURL string, reqExpr *RequirementExpression) error {
	if reqExpr == nil {
		return nil
	}

	switch reqExpr.Type {
	case "or":
		return errors.New(`simplifyRequirement: requirement expression type "or" will be implemented in 0.1v`)
	case "and":
		return errors.New(`simplifyRequirement: requirement expression type "and" will be implemented in 0.1v`)
	case "requirement":
		return simplifyRequirement(hookURL, &reqExpr.Req)
	default:
		return errors.New(`simplifyRequirement: unsupported requirement expression: ` + reqExpr.Type)
	}
}

func simplifyRequirement(hookURL string, req *Requirement) error {
	var e error

	switch req.Operand.Type {
	case "request":
		if req.Operand.Spec, e = getFieldFromRequest(hookURL, req.Operand.Spec, &req.Field); e == nil {
			req.Operand.Type = "raw_value"
		}
	}

	return nil
}

func getFieldFromRequest(hookURL string, req interface{}, f *Field) (interface{}, error) {
	var err error
	var res interface{}

	if str, e := json.Marshal(req); e != nil {
		err = e
	} else if resp, e := httputil.SendRequest("POST", hookURL, "application/json", bytes.NewReader(str)); e != nil {
		err = e
	} else {
		dec := json.NewDecoder(resp.Body)

		if e := dec.Decode(&res); e != nil {
			err = e
		} else if e := resp.Body.Close(); e != nil {
			err = e
		} else {
			res, err = getFieldFromResult(res, f)
		}
	}

	return res, err
}

func getFieldFromResult(res interface{}, f *Field) (interface{}, error) {
	var err error

	switch s := res.(type) {
	case []interface{}:
		for i, v := range s {
			if s[i], err = getFieldFromMap(v, f); err != nil {
				return nil, err
			}
		}

		res = s
	default:
		res, err = getFieldFromMap(res, f)
	}

	return res, err
}

func getFieldFromMap(m interface{}, f *Field) (interface{}, error) {
	var err error
	if f == nil {
		return nil, errors.New("getFieldFromMap: expected field specifier")
	} else {
		switch mp := m.(type) {
		case map[string]interface{}:
			m = mp[f.Name]
			if f.Child != nil {
				m, err = getFieldFromMap(m, f.Child)
			}
		default:
			err = fmt.Errorf("getFieldFromMap: expected map type, received: %v\n", reflect.TypeOf(m))
		}
	}

	return m, err
}
