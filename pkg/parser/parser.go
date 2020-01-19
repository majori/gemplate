package parser

type Parser = func(*string) (*Settings, *States)

type Settings map[string]interface{}

type States []map[string]interface{}
