package request

import (
	"errors"
)

func (expr *RequirementExpression) translateToES() (res interface{}, err error) {
	if expr != nil {
		switch expr.Type {
		case "requirement":
			res, err = expr.Req.translateToES()
		case "and":
			must := make([]interface{}, len(expr.OrAnd))

			for i, ex := range expr.OrAnd {
				if must[i], err = ex.translateToES(); err != nil {
					return
				}
			}

			res = map[string]interface{}{"bool": map[string]interface{}{"must": must}}
		case "or":
			should := make([]interface{}, len(expr.OrAnd))

			for i, ex := range expr.OrAnd {
				if should[i], err = ex.translateToES(); err != nil {
					return
				}
			}

			res = map[string]interface{}{"bool": map[string]interface{}{"should": should}}
		default:
			err = errors.New("RequirementExpression.translateToES: unknown expression type: " + expr.Type)
		}
	} else {
		res = map[string]interface{}{"match_all": map[string]interface{}{}}
	}

	return
}

func (r *Requirement) translateToES() (res interface{}, err error) {
	fieldName := r.Field.getESName()
	switch r.Operator {
	case "in":
		res = map[string]interface{}{"terms": map[string]interface{}{fieldName: r.Operand.Spec}}
	case "not in":
		res = map[string]interface{}{"bool": map[string]interface{}{"must_not": map[string]interface{}{"terms": map[string]interface{}{fieldName: r.Operand.Spec.([]interface{})}}}}
	case "ne":
		res = map[string]interface{}{"bool": map[string]interface{}{"must_not": map[string]interface{}{"term": map[string]interface{}{fieldName: r.Operand.Spec}}}}
	case "eq":
		res = map[string]interface{}{"term": map[string]interface{}{fieldName: r.Operand.Spec}}
	case "ge":
		res = map[string]interface{}{"range": map[string]interface{}{r.Field.getESName(): map[string]interface{}{"gte": r.Operand.Spec}}}
	case "le":
		res = map[string]interface{}{"range": map[string]interface{}{r.Field.getESName(): map[string]interface{}{"lte": r.Operand.Spec}}}
	case "gt":
		res = map[string]interface{}{"range": map[string]interface{}{r.Field.getESName(): map[string]interface{}{"gt": r.Operand.Spec}}}
	case "lt":
		res = map[string]interface{}{"range": map[string]interface{}{r.Field.getESName(): map[string]interface{}{"lt": r.Operand.Spec}}}
	default:
		err = errors.New("requirement.translateToES: unknown operator: " + r.Operator)
	}
	return
}

func (f Field) getESName() string {
	name := f.Name
	for it := &f; it.Child != nil; it = it.Child {
		name += "." + it.Name
	}
	return name
}

func translateToESSpecificForm(f *simplifiedRequest) (interface{}, error) {
	var res = make(map[string]interface{})
	var err error

	if res["query"], err = f.ReqExpr.translateToES(); err != nil {
		return nil, err
	}

	if len(f.Fields) > 0 {
		fieldNames := make([]string, len(f.Fields))
		for i, field := range f.Fields {
			fieldNames[i] = field.getESName()
		}
		res["_source"] = fieldNames
	}
	for _, m := range f.Modifiers {
		for k, v := range m {
			res[k] = v
		}
	}

	return res, nil
}
