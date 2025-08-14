package main

import (
	"fmt"
	"path"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// generateServiceFiles generates service files to handle requests.
func generateServiceFiles(gen *protogen.Plugin, file *protogen.File, gp *GenParam) *protogen.GeneratedFile {
	if len(file.Services) == 0 || (*gp.Omitempty && !hasHTTPRule(file.Services)) {
		return nil
	}
	filename := fmt.Sprintf("service/%s_service.go", file.GoPackageName)

	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	// g.P("/*")
	// g.P(file.GoPackageName)
	// g.P(file.GoImportPath)
	// g.P(file.GeneratedFilenamePrefix)
	// g.P(path.Dir(file.GeneratedFilenamePrefix))

	// pth := ""
	// bi, ok := debug.ReadBuildInfo()
	// if !ok {
	// 	pth = "Failed to read build info"
	// } else {
	// 	pth = bi.Path
	// }

	// g.P(pth)

	// g.P("*/")

	printHeaders(g, gen, false, "", "Add your service handlers to this file", "")

	g.P("")

	g.P("package ", "service")

	generateServiceFileContent(gen, file, g, gp)
	return g
}

func generateServiceFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, gp *GenParam) {
	if len(file.Services) == 0 {
		return
	}

	p := string(file.GoImportPath)

	i := strings.LastIndex(p, "/")
	if i > -1 {
		pth := *gp.GenPath
		pth = strings.Trim(pth, "./")
		p = p[:i] + "/" + pth + "/" + path.Dir(file.GeneratedFilenamePrefix) // p[i+1:]
	}

	g.P(fmt.Sprintf("import %s %q", file.GoPackageName, p))

	//generatedPackage := protogen.GoImportPath(p)

	// A bunch of variables are defined here to ensure that these packages are correct when the program is compiled.
	// If the package does not exist or the defined package variables do not exist, the compilation will fail.
	g.P("// This is a compile-time assertion to ensure that generated files are safe and compilable.")

	// As long as the Ident method is called, it will be automatically written to the import, so if there is no
	// special requirement for the import package name, just use Ident directly
	g.P("var _ ", contextPackage.Ident("Context"))
	g.P("var _ =", errorsPackage.Ident("New"))
	// g.P("var _ ", generatedPackage.Ident(fmt.Sprintf("%sHTTPHandler", file.Services[0].GoName)))
	g.P("const _ = ", ginPackage.Ident("Version"))

	for _, service := range file.Services {
		genService(gen, file, g, service, gp, true)
	}

}
