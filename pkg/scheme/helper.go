package scheme

import "strings"

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
