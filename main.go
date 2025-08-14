package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// Version protoc-gen-go-gin 工具版本
const Version = "v0.0.3"

func main() {
	var flags flag.FlagSet
	omitempty := flags.Bool("omitempty", true, "omit if google.api is empty")

	// Flag to include validations from protoc-gen-validate which uses validations embedded in the .proto file
	genValidateCode := flags.Bool("validate", false, "add validate request params in handler")

	genService := flags.Bool("service", false, "generate service code")
	genPath := flags.String("genpath", "", "directory of generated files")

	gp := &GenParam{
		Omitempty:       omitempty,
		GenValidateCode: genValidateCode,
		GenService:      genService,
		GenPath:         genPath,
	}

	// 这里就是入口，指定 option 后执行 Run 方法 ，我们的主逻辑就是在 Run 方法
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			// 这里是我们的生成代码方法
			generateHTTPFile(gen, f, gp)
			if *gp.GenService {
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
	GenService      *bool
	GenPath         *string
	GenRegister     *bool
}
