package main

// CLI colors

const (
	clNoColor      = "\033[0m"
	clWhite        = "\033[1;37m"
	clBlack        = "\033[0;30m"
	clBlue         = "\033[0;34m"
	clLightBlue    = "\033[1;34m"
	clGreen        = "\033[0;32m"
	clLightGreen   = "\033[1;32m"
	clCyan         = "\033[0;36m"
	clLightCyan    = "\033[1;36m"
	clRed          = "\033[0;31m"
	clLightRed     = "\033[1;31m"
	cPurple        = "\033[0;35m"
	clLightPurple  = "\033[1;35m"
	clBrown        = "\033[0;33m"
	clYellow       = "\033[1;33m"
	clGray         = "\033[0;30m"
	clLightGray    = "\033[0;37m"
	clDefault      = clGreen
	clAliasDefault = `default`
)

type ColorizerStruct struct {
	FontColorMap map[string]string // key values for colors in format alias => font code
}

// Wrap wraps text with color
func (c *ColorizerStruct) Wrap(text string, clAlias string) string {
	var fontColor string

	if fontColorVal, ok := c.FontColorMap[clAlias]; !ok { // if alias not found
		fontColor = clDefault
	} else {
		fontColor = fontColorVal
	}

	return fontColor + text + clNoColor
}

// NewColorizer creates new instance of colorizer
func NewColorizer(fontColorMap map[string]string) *ColorizerStruct {
	if len(fontColorMap) == 0 {
		fontColorMap[clAliasDefault] = clDefault
	}

	var colorizer = ColorizerStruct{
		FontColorMap: fontColorMap,
	}

	return &colorizer
}
