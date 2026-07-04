package cobra

// cobra-doc-start
var docsCobraMapping = map[string]string{
	`field1`: `Doc 1`,
	`replicas`: `Doc Replicas`,
}

func cflagVar[T any](p *T, name string, value T) (pOut *T, nameOut string, valueOut T, reason string) {
	return p, name, value, docsCobraMapping[name]
}

func cflag[T any](name string, value T) (nameOut string, valueOut T, reason string) {
	return name, value, docsCobraMapping[name]
}

func cflagP[T any](name, shorthand string, value T) (nameOut, shorthandOut string, valueOut T, reason string) {
	return name, shorthand, value, docsCobraMapping[name]
}

// cobra-doc-end
