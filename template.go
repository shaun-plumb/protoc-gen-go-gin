package main

// TODO: support validate
var httpCodeTmpl = `
{{/*gotype: github.com/shaun-plumb/protoc-gen-go-gin.serviceDesc*/}}
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
{{$validate := .GenValidate}}

// {{.ServiceType}}HTTPHandler defines {{.ServiceType}}Server http handler
type {{.ServiceType}}HTTPHandler interface {
{{- range .Methods}}
    {{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}


type Unimplemented{{$svrType}}HTTPServer struct{}
{{range .Methods}}
func (Unimplemented{{$svrType}}HTTPServer) {{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error) {
	return nil, status.Errorf(codes.Unimplemented, "method {{.Name}} not implemented")
}
{{end}}



// Register{{.ServiceType}}HTTPHandler define http router handle by gin.
func Register{{.ServiceType}}HTTPHandler(g *gin.RouterGroup, srv {{.ServiceType}}HTTPHandler) {
{{- range .Methods}}
    g.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
{{- end}}
}

{{range .Methods}}
// _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler is gin http handler to handle
// http request [{{.Method}}] {{.Path}}.
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPHandler) func(c *gin.Context) {
    return func(c *gin.Context) {
        var (
            err error
            in  = new({{.Request}})
            out = new({{.Reply}})
        )

        {{ if .HasVars }}
        common.ExtractPathParameters(c, &in)
        {{ end }}

        if err = c.ShouldBind(in{{.Body}}); err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
            return
        }
        {{if $validate}}
        v,ok := interface{}(in).(common.Validator)
        if ok {
            if err = v.Validate();err != nil {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
                return
            }
        }
        {{end}}
        // execute
        out, err = srv.{{.Name}}(c, in)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
            return
        }

        c.JSON(http.StatusOK, out)
    }
}
{{end}}
`

var serviceCodeTmpl = `

{{$package := .PackageName}}
{{$serviceType := .ServiceType}}
{{$sampleMethod := index .Methods 0 }}

/*
{{$serviceType}}HTTPHandler is the service handler where the individual method handlers are implemented for {{$serviceType}}
*/
/* === IMPLEMENTATION INSTRUCTIONS ===
Initially, the service is implemented by {{$package}}.Unimplemented{{$serviceType}}HTTPServer, which means that all unimplemented 
methods will respond with an HTTP 500 status and a JSON formatted error message.

The following tasks remain to implement this service.

* Firstly - implement each of the individual method handlers like this:

func (s *{{$serviceType}}HTTPHandler) {{$sampleMethod.Name}}(ctx context.Context, req *{{$package}}.{{$sampleMethod.Request}}) (*{{$package}}.{{$sampleMethod.Reply}}, error) {
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
/*
func (s *{{$serviceType}}HTTPHandler) {{.Name}}(ctx context.Context, req *{{$package}}.{{.Request}}) (*{{$package}}.{{.Reply}}, error) {
    panic("not implemented")
}
*/
{{end}}
 

`
