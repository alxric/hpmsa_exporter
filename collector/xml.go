package collector

// Response is the xml response we get from the HPMSA api
type Response struct {
	Object []object `xml:"OBJECT"`
}

type object struct {
	Property []property `xml:"PROPERTY"`
	Object   []object   `xml:"OBJECT"`
	BaseType string     `xml:"basetype,attr"`
	Name     string     `xml:"name,attr"`
}

type property struct {
	Value string `xml:",chardata"`
	Name  string `xml:"name,attr"`
}
