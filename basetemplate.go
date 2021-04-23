package main

func GetStaticPath() string {
	staticPath := "/static/" + getVersion()

	if getVersion() == "Unknown" {
		staticPath = "/static"
	}
	return staticPath
}

type BaseTemplateArgs struct {
	StaticPath string
	Version string
}

func GetBaseTemplateArgs() BaseTemplateArgs {


	return BaseTemplateArgs{
		StaticPath: GetStaticPath(),
		Version: getVersion(),
	}
}