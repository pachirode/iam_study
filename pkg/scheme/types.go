package scheme

import (
	"fmt"
	"strings"
)

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

func (groupResource GroupResource) WithVersion(version string) GroupVersionResource {
	return GroupVersionResource{Group: groupResource.Group, Version: version, Resource: groupResource.Resource}
}

func (groupResource GroupResource) Empty() bool {
	return len(groupResource.Group) == 0 && len(groupResource.Resource) == 0
}

func (groupResource GroupResource) String() string {
	if len(groupResource.Group) == 0 {
		return groupResource.Resource
	}
	return groupResource.Resource + "." + groupResource.Group
}

func (gvr GroupVersionResource) Empty() bool {
	return len(gvr.Group) == 0 && len(gvr.Version) == 0 && len(gvr.Resource) == 0
}

func (gvr GroupVersionResource) GroupResource() GroupResource {
	return GroupResource{Group: gvr.Group, Resource: gvr.Resource}
}

func (gvr GroupVersionResource) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvr.Group, Version: gvr.Version}
}

func (gvr GroupVersionResource) String() string {
	return strings.Join([]string{gvr.Group, "/", gvr.Version, ", Resource=", gvr.Resource}, "")
}

func (gk GroupKind) Empty() bool {
	return len(gk.Group) == 0 && len(gk.Kind) == 0
}

func (gk GroupKind) WithVersion(version string) GroupVersionKind {
	return GroupVersionKind{Group: gk.Group, Version: version, Kind: gk.Kind}
}

func (gk GroupKind) String() string {
	if len(gk.Group) == 0 {
		return gk.Kind
	}
	return gk.Kind + "." + gk.Group
}

func (gvk GroupVersionKind) Empty() bool {
	return len(gvk.Group) == 0 && len(gvk.Version) == 0 && len(gvk.Kind) == 0
}

func (gvk GroupVersionKind) GroupKind() GroupKind {
	return GroupKind{Group: gvk.Group, Kind: gvk.Kind}
}

func (gvk GroupVersionKind) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvk.Group, Version: gvk.Version}
}

func (gvk GroupVersionKind) String() string {
	return gvk.Group + "/" + gvk.Version + ", Kind=" + gvk.Kind
}

func (gv GroupVersion) Empty() bool {
	return len(gv.Group) == 0 && len(gv.Version) == 0
}

func (gv GroupVersion) String() string {
	if len(gv.Group) > 0 {
		return gv.Group + "/" + gv.Version
	}
	return gv.Version
}

func (gv GroupVersion) Identifier() string {
	return gv.String()
}

func (gv GroupVersion) KindForGroupVersionKinds(kinds []GroupVersionKind) (target GroupVersionKind, ok bool) {
	for _, gvk := range kinds {
		if gvk.Group == gv.Group && gvk.Version == gv.Version {
			return gvk, true
		}
	}
	for _, gvk := range kinds {
		if gvk.Group == gv.Group {
			return gv.WithKind(gvk.Kind), true
		}
	}
	return GroupVersionKind{}, false
}

func (gv GroupVersion) WithKind(kind string) GroupVersionKind {
	return GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}
}

func (gv GroupVersion) WithResource(resource string) GroupVersionResource {
	return GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource}
}

func (gvs GroupVersions) Identifier() string {
	groupVersions := make([]string, 0, len(gvs))
	for i := range gvs {
		groupVersions = append(groupVersions, gvs[i].String())
	}
	return fmt.Sprintf("[%s]", strings.Join(groupVersions, ","))
}

func (gvs GroupVersions) KindForGroupVersionKinds(kinds []GroupVersionKind) (GroupVersionKind, bool) {
	var targets []GroupVersionKind
	for _, gv := range gvs {
		target, ok := gv.KindForGroupVersionKinds(kinds)
		if !ok {
			continue
		}
		targets = append(targets, target)
	}
	if len(targets) == 1 {
		return targets[0], true
	}
	if len(targets) > 1 {
		return bestMatch(kinds, targets), true
	}
	return GroupVersionKind{}, false
}

func (gvk GroupVersionKind) ToAPIVersionAndKind() (string, string) {
	if gvk.Empty() {
		return "", ""
	}
	return gvk.GroupVersion().String(), gvk.Kind
}
