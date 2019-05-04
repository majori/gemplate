package parser

type Settings struct {
	Static  Static
	Dynamic Dynamic
}

type Static map[string]interface{}

type Dynamic map[string]interface{}
