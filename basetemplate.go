package main

type BaseTemplateArgs struct {
	Version string
}

func GetBaseTemplateArgs() BaseTemplateArgs {

	return BaseTemplateArgs{
		Version: getVersion(),
	}
}
