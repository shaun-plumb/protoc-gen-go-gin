// Package common to provide functions used by all generated files
package common

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Type    string
	Field   string
	Message string
}

type GeneralHTTPError struct {
	Errors     Error
	StatusCode int
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
