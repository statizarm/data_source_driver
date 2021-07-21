package request

import "fmt"

func (f Field) getSQLName() string {
	if f.Child != nil {
		return f.Name + "." + f.Child.getSQLName()
	} else {
		return f.Name
	}
}

// UnifiedRequestToSql method. Transforms a UnifiedRequest structure into a full SQL statement.
func (u GetRequest) translateToSQLSpecificForm() string {
	var req string
	var fields string

	req = "SELECT "

	if len(u.Fields) == 0 {
		fields = "*"
	} else {
		for _, i := range u.Fields {
			fields += i.getSQLName()
			if i != u.Fields[len(u.Fields)-1] {
				fields += ", "
			}

		}
	}

	req += fields + " FROM " + u.Source + " WHERE " + u.Requirements.translateToSQL()

	return req
}

// translateToSQL method. Transforms a RequirementExpression into an SQL condition (after WHERE)
func (expr *RequirementExpression) translateToSQL() string {
	var temp string
	switch expr.Type {
	case "or":
		temp = "(" + expr.OrAnd[0].translateToSQL()
		for _, i := range expr.OrAnd[1:] {
			temp += " OR " + i.translateToSQL()
		}
		return temp + ")"

	case "and":
		temp = "(" + expr.OrAnd[0].translateToSQL()
		for _, i := range expr.OrAnd[1:] {
			temp += " AND " + i.translateToSQL()
		}
		return temp + ")"

	case "requirement":
		return expr.Req.translateToSQL()
	}
	return ""
}

// translateToSQL method. Translates a requirement structure into an SQL requirement
func (r *Requirement) translateToSQL() string {
	var req string
	req = "(" + r.Field.getSQLName()
	switch r.Operator {
	case "in": // Can be done with errors, but I'm not sure!
		values := r.Operand.Spec.([]interface{})
		req += " IN (" + "'" + fmt.Sprint(values[0]) + "'"
		for _, i := range values[1:] {
			req += ", " + "'" + fmt.Sprint(i) + "'"
		}
		return req + "))"

	case "not in":
		values := r.Operand.Spec.([]interface{})
		req += " NOT IN (" + "'" + fmt.Sprint(values[0]) + "'"
		for _, i := range values[1:] {
			req += ", " + "'" + fmt.Sprint(i) + "'"
		}
		return req + "))"

	case "ne":
		return req + " <> '" + fmt.Sprint(r.Operand.Spec) + "')"

	case "eq":
		return req + " = '" + fmt.Sprint(r.Operand.Spec) + "')"

	case "ge":
		return req + " >= " + fmt.Sprint(r.Operand.Spec) + ")"

	case "le":
		return req + " <= " + fmt.Sprint(r.Operand.Spec) + ")"

	case "gt":
		return req + " > " + fmt.Sprint(r.Operand.Spec) + ")"

	case "lt":
		return req + " < " + fmt.Sprint(r.Operand.Spec) + ")"
	}

	return ""
}
