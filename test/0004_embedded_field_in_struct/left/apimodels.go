package apimodels

type Todo struct {
	myType
	ID          int
	Description string
}

type myType struct {
	extraField int
}
