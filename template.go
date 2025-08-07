package main

var httpCodeTmpl = `
{{/*gotype: github.com/shaun-plumb/protoc-gen-go-gin.serviceDesc*/}}
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
{{$validate := .GenValidate}}

// {{.ServiceType}}HTTPHandler defines {{.ServiceType}}Server http handler
type {{.ServiceType}}HTTPHandler interface {
{{- range .Methods}}
    {{.Name}}(*gin.Context, *{{.Request}}) (*{{.Reply}}, error)
    Validate{{.Request}}(*gin.Context, *{{.Request}}) error
{{- end}}
}


type Unimplemented{{$svrType}}HTTPServer struct{}
{{range .Methods}}
func (Unimplemented{{$svrType}}HTTPServer) {{.Name}}(*gin.Context, *{{.Request}}) (*{{.Reply}}, error) {
	return nil, status.Errorf(codes.Unimplemented, "method {{.Name}} not implemented")
}

func (Unimplemented{{$svrType}}HTTPServer) Validate{{.Request}}(*gin.Context, *{{.Request}}) error { return nil }
{{end}}

// Register{{.ServiceType}}HTTPHandlers associates http router handlers in gin.
func Register{{.ServiceType}}HTTPHandler(g *gin.RouterGroup, srv {{.ServiceType}}HTTPHandler) {
{{- range .Methods}}
    g.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
{{- end}}
}

{{range .Methods}}
// _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler is gin http handler to handle
// http request [{{.Method}}] {{.Path}}.
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPHandler) func(ctx *gin.Context) {
    return func(ctx *gin.Context) {
        var (
            err error
            in  = new({{.Request}})
            out = new({{.Reply}})
        )

        {{ if .HasVars }}
        common.ExtractPathParameters(ctx, &in)
        {{ end }}

        if err = ctx.ShouldBind(in{{.Body}}); err != nil {
            ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

    {{if $validate}}
        // call protovalidate to apply any validation rules in the .proto file
        if err = protovalidate.Validate(in);err != nil {
            ctx.AbortWithStatusJSON(http.StatusBadRequest, common.GenerateErrorsFromProtoViolation(err.(*protovalidate.ValidationError)))
            return
        }
    {{end}}

        // Call any supplied validation routines
        if err = srv.Validate{{.Request}}(ctx, in); err != nil {
            ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // execute
        out, err = srv.{{.Name}}(ctx, in)
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        ctx.JSON(http.StatusOK, out)
    }
}
{{end}}
`

var serviceCodeTmpl = `

{{$package := .PackageName}}
{{$serviceType := .ServiceType}}
{{$sampleMethod := index .Methods 0 }}
{{$validate := .GenValidate}}

/*
{{$serviceType}}HTTPHandler is the service handler where the individual method handlers are implemented for {{$serviceType}}
*/
/* === IMPLEMENTATION INSTRUCTIONS ===

Initially, the service is implemented by {{$package}}.Unimplemented{{$serviceType}}HTTPServer, which means that all unimplemented 
methods will respond with an HTTP 500 status and a JSON formatted error message.

{{if $validate}}
This service has been configured automatically to use the protovalidate library (see: https://buf.build/docs/protovalidate/) to create validation code
based on annotations in the .proto file. This means that the incoming request will be validated against any validation annotations present in the .proto file 
and throw a HTTP 400 (Bad Request) error on failure before the request reaches the handler methods in this file. 

There is no need to validate the request inputs in these methods except when not covered by any validation annotations. 
{{end}}

The following tasks remain to implement this service.

* Firstly - implement each of the individual method handlers like this:

func (s *{{$serviceType}}HTTPHandler) {{$sampleMethod.Name}}(ctx *gin.Context, req *{{$package}}.{{$sampleMethod.Request}}) (*{{$package}}.{{$sampleMethod.Reply}}, error) {
 	
    // do some logic here

    return &{{$package}}.{{$sampleMethod.Reply}}{
		Id:       req.Id,
        ... other data
	}, nil   

}

* Secondly, to register this service handler with Go-Gin, do something like:

import (
	"mymodule/generated/{{$package}}"
	"mymodule/service"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
    {{$package}}Handler := service.New{{$serviceType}}HTTPHandler() 
    {{$package}}.Register{{$serviceType}}HTTPHandler(router.Group("/"), {{$package}}Handler)
}

Once implemented - this message can be deleted.
=== */


type {{$serviceType}}HTTPHandler struct {
    {{$package}}.Unimplemented{{$serviceType}}HTTPServer
}

func New{{$serviceType}}HTTPHandler() *{{$serviceType}}HTTPHandler {
	return &{{$serviceType}}HTTPHandler{}
}

{{range .Methods}}
    {{- $reply := printf "%v.%v" $package .Reply }}
    {{- if contains .Reply "." }}
        {{- $reply = .Reply}}
    {{- end }}
func (s *{{$serviceType}}HTTPHandler) {{.Name}}(ctx *gin.Context, req *{{$package}}.{{.Request}}) (*{{ $reply }}, error) {
    return nil, errors.New("method {{.Name}} not implemented")
}

{{end}}
 

`
