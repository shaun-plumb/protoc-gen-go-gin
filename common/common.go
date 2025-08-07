// Package common to provide functions used by all generated files
package common

import (
	"encoding/json"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/gin-gonic/gin"
)

type internalError struct {
	Type    string `json:"type"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type GeneralHTTPError struct {
	Errors     []internalError `json:"errors"`
	StatusCode int             `json:"statuscode"`
}

func CreateHTTPError(stat int, typ string) GeneralHTTPError {
	return GeneralHTTPError{StatusCode: stat}
}

func (s *GeneralHTTPError) AddError(typ string, field string, message string) {
	s.Errors = append(s.Errors, internalError{Type: typ, Field: field, Message: message})
}

func (s *GeneralHTTPError) AddProtoViolations(err *protovalidate.ValidationError) {
	for _, vl := range err.Violations {
		s.AddError("validation", *vl.Proto.Field.Elements[0].FieldName, vl.Proto.GetMessage())
	}
}

func GenerateErrorsFromProtoViolation(viols *protovalidate.ValidationError) (GeneralHTTPError) {
	err := GeneralHTTPError{StatusCode: 400}
	err.AddProtoViolations(viols)
	return err
}

func ExtractPathParameters(c *gin.Context, in any) {
	if len(c.Params) > 0 {
		pMap := "{"
		for _, v := range c.Params {
			pMap += "\"" + v.Key + "\":\"" + v.Value + "\","
		}
		pMap = strings.Trim(pMap, ",")
		pMap += "}"
		err := json.Unmarshal([]byte(pMap), &in)
		if err != nil {
			panic(err)
		}
	}
}
