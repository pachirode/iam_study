package scheme

type GroupResource struct {
	Group    string
	Resource string
}

type GroupVersionResource struct {
	Group    string
	Version  string
	Resource string
}

type GroupKind struct {
	Group string
	Kind  string
}

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

type GroupVersion struct {
	Group   string
	Version string
}

type GroupVersions []GroupVersion

type emptyObjectKind struct{}
