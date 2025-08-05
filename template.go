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

// Register{{.ServiceType}}HTTPHandler define http router handle by gin.
func Register{{.ServiceType}}HTTPHandler(g *gin.RouterGroup, srv {{.ServiceType}}HTTPHandler) {
{{- range .Methods}}
    g.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
{{- end}}
}

{{if $validate}}
type Validator interface {
    Validate() error
}
{{end}}

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
        v,ok := interface{}(in).(Validator)
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

type {{$serviceType}}HTTPHandler struct {
    {{$package}}.Unimplemented{{$serviceType}}Server
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
