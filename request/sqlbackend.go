package request

import (
	"errors"
	"fmt"
)

func (f Field) getSQLName() (name string, err error) {
	if f.Name == "" {
		err = errors.New("in request/sqlbackend.go: getSQLName - Field name is empty or in an incorrect format")
	} else if f.Child != nil {
		subname, _ := f.Child.getSQLName()
		name = f.Name + "." + subname
	} else {
		name = f.Name
	}
	return
}

// UnifiedRequestToSql method. Transforms a UnifiedRequest structure into a full SQL statement.
func translateToSQLSpecificForm (s *simplifiedRequest) (req string, err error) {
	var fields string

	req = "SELECT "

	if len(s.Fields) == 0 {
		fields = "*"
	} else {
		for _, i := range s.Fields {
			if fld, e := i.getSQLName(); e != nil {
				err = e
				return "", err
			} else {
				fields += fld

				if i != s.Fields[len(s.Fields)-1] {
					fields += ", "
				}
			}
		}
	}
	req += fields + " FROM "
	if src := s.Source; src == "" {
		err = errors.New("in request/sqlbackend.go: translateToSQLSpecificForm - Empty Source field")
		return "", err
	} else {
		req += src + " WHERE "
	}

	if rexp, e := s.ReqExpr.translateToSQL(); e != nil{
		err = e
		return "", err
	} else {
		req += rexp
	}
	return
}

// translateToSQL method. Transforms a RequirementExpression into an SQL condition (after WHERE)
func (expr *RequirementExpression) translateToSQL() (rexp string, err error) {
	var temp string
	switch expr.Type {
	case "or":
		if len(expr.OrAnd) < 2{
			err = errors.New("in request/sqlbackend.go: ReqExpr.translateToSQL - Invalid number of expressions in OrAnd field")
		} else {
			if r, e := expr.OrAnd[0].translateToSQL(); e != nil {
				err = e
				return "", err
			} else {
				temp = "(" + r
			}

			for _, i := range expr.OrAnd[1:] {
				if r, e := i.translateToSQL(); e != nil {
					err = e
					return "", err
				} else {
					temp += " OR " + r
				}
			}
			return temp + ")", nil
		}

	case "and":
		if len(expr.OrAnd) < 2{
			err = errors.New("in request/sqlbackend.go: ReqExpr.translateToSQL - Invalid number of expressions in OrAnd field")
		} else {
			if r, e := expr.OrAnd[0].translateToSQL(); e != nil {
				err = e
				return "", err
			} else {
				temp = "(" + r
			}

			for _, i := range expr.OrAnd[1:] {
				if r, e := i.translateToSQL(); e != nil {
					err = e
					return "", err
				} else {
					temp += " AND " + r
				}
			}
			return temp + ")", nil
		}

	case "requirement":
		if r, e := expr.Req.translateToSQL(); e != nil {
			err = e
			return "", err
		} else{
			rexp = r
		}

	default:
		err = errors.New("in request/sqlbackend.go: ReqExpr.translateToSQL - Wrong type of expression")
	}
	return
}

// translateToSQL method. Translates a requirement structure into an SQL requirement
func (r *Requirement) translateToSQL() (req string, err error) {
	if r, e := r.Field.getSQLName(); e != nil {
		err = e
		return "", err
	} else {
		req = "(" + r
	}

	switch r.Operator {
	case "in":
		values := r.Operand.Spec.([]interface{})
		if len(values) < 1 {
			err = errors.New("in request/sqlbackend.go: Requirement.translateToSQL - Empty array of values")
			return "", err
		} else {
			req += " IN (" + "'" + fmt.Sprint(values[0]) + "'"
			for _, i := range values[1:] {
				req += ", " + "'" + fmt.Sprint(i) + "'"
			}
			return req + "))", nil
		}

	case "not in":
		values := r.Operand.Spec.([]interface{})
		if len(values) < 1 {
			err = errors.New("in request/sqlbackend.go: Requirement.translateToSQL - Empty array of values")
			return "", err
		} else {
			req += " NOT IN (" + "'" + fmt.Sprint(values[0]) + "'"
			for _, i := range values[1:] {
				req += ", " + "'" + fmt.Sprint(i) + "'"
			}
			return req + "))", nil
		}

	case "ne":
		if spc := fmt.Sprint(r.Operand.Spec); spc == "" {
			err = errors.New("in request/sqlbackend.go: Requirement.translateToSQL - Empty value")
			return "", err
		} else {
			return req + " <> '" + spc + "')", nil
		}

	case "eq":
		if spc := fmt.Sprint(r.Operand.Spec); spc == "" {
			err = errors.New("in request/sqlbackend.go: Requirement.translateToSQL - Empty value")
			return "", err
		} else {
			return req + " = '" + spc + "')", nil
		}

	case "ge":
		if spc := fmt.Sprint(r.Operand.Spec); spc == "" {
			err = errors.New("in request/sqlbackend.go: Requirement.translateToSQL - Empty value")
			return "", err
		} else {
			return req + " >= '" + spc + "')", nil
		}

	case "le":
		if spc := fmt.Sprint(r.Operand.Spec); spc == "" {
			err = errors.New("in request/sqlbackend.go: Requirement.translateToSQL - Empty value")
			return "", err
		} else {
			return req + " <= '" + spc + "')", nil
		}

	case "gt":
		if spc := fmt.Sprint(r.Operand.Spec); spc == "" {
			err = errors.New("in request/sqlbackend.go: Requirement.translateToSQL - Empty value")
			return "", err
		} else {
			return req + " > '" + spc + "')", nil
		}

	case "lt":
		if spc := fmt.Sprint(r.Operand.Spec); spc == "" {
			err = errors.New("in request/sqlbackend.go: Requirement.translateToSQL - Empty value")
			return "", err
		} else {
			return req + " < '" + spc + "')", nil
		}

	default:
		err = errors.New("in request/sqlbackend.go: Requirement.translateToSQL - Wrong type of requirement operator")
		return "", err
	}
}
