package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const Version = "dev"

func main() {
	var flags flag.FlagSet
	omitempty := flags.Bool("omitempty", true, "omit if google.api is empty")

	// Flag to include validations from protoc-gen-validate which uses validations embedded in the .proto file
	genValidateCode := flags.Bool("validate", false, "add validate request params in handler")

	genServiceFiles := flags.Bool("service", false, "generate service code")
	genPath := flags.String("genpath", "", "directory of generated files")

	gp := &GenParam{
		Omitempty:       omitempty,
		GenValidateCode: genValidateCode,
		GenServiceFiles: genServiceFiles,
		GenPath:         genPath,
	}

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			// generate the HTTP handlers
			generateHTTPHandlerFile(gen, f, gp)
			// if we are generating the service, do that here
			if *gp.GenServiceFiles {
				if *gp.GenPath != "" {
					generateServiceFiles(gen, f, gp)
				} else {
					panic("Cannot specify service generation without genpath - use '--go-gin_opt=paths=source_relative,service=true,genpath=$OUT_PATH'")
				}
			}
		}
		return nil
	})
}

type GenParam struct {
	Omitempty       *bool
	GenValidateCode *bool
	GenServiceFiles *bool
	GenPath         *string
	GenRegister     *bool
}
