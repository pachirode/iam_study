package scheme

import (
	"fmt"
	"strings"
)

// resource.group.com -> group=com, version=group, resource=resource
func ParseResourceArg(arg string) (*GroupVersionResource, GroupResource) {
	var gvr *GroupVersionResource
	if strings.Count(arg, ".") >= 2 {
		strList := strings.SplitN(arg, ".", 3)
		gvr = &GroupVersionResource{Group: strList[2], Version: strList[1], Resource: strList[0]}
	}

	return gvr, ParseGroupResource(arg)
}

func ParseGroupResource(arg string) GroupResource {
	if idx := strings.Index(arg, "."); idx >= 0 {
		return GroupResource{Group: arg[idx+1:], Resource: arg[:idx]}
	}

	return GroupResource{Resource: arg}
}

func ParseKindArg(arg string) (*GroupVersionKind, GroupKind) {
	var gvk *GroupVersionKind
	if strings.Count(arg, ".") >= 2 {
		strList := strings.SplitN(arg, ".", 3)
		gvk = &GroupVersionKind{Group: strList[2], Version: strList[1], Kind: strList[0]}
	}

	return gvk, ParseGroupKind(arg)
}

func ParseGroupKind(arg string) GroupKind {
	if idx := strings.Index(arg, "."); idx >= 0 {
		return GroupKind{Group: arg[idx+1:], Kind: arg[:idx]}
	}

	return GroupKind{Kind: arg}
}

// group/version
func ParseGroupVersion(arg string) (GroupVersion, error) {
	if len(arg) == 0 || arg == "/" {
		return GroupVersion{}, nil
	}

	switch strings.Count(arg, "/") {
	case 0:
		return GroupVersion{"", arg}, nil
	case 1:
		idx := strings.Index(arg, "/")
		return GroupVersion{arg[:idx], arg[idx+1:]}, nil
	default:
		return GroupVersion{}, fmt.Errorf("Unexpected groupVersion string: %s", arg)
	}
}

func bestMatch(kinds []GroupVersionKind, targets []GroupVersionKind) GroupVersionKind {
	for _, gvk := range targets {
		for _, k := range kinds {
			if k == gvk {
				return k
			}
		}
	}

	return targets[0]
}

func FormAPIVersionAndKind(apiVersion, kind string) GroupVersionKind {
	if groupVersion, err := ParseGroupVersion(apiVersion); err != nil {
		return GroupVersionKind{Group: groupVersion.Group, Version: groupVersion.Group, Kind: kind}
	}

	return GroupVersionKind{Kind: kind}
}
