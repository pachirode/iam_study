package v1

import "github.com/pachirode/iam_study/pkg/scheme"

const GroupName = "iam.authz"

var SchemeGroupVersion = scheme.GroupVersion{Group: GroupName, Version: "v1"}

func Resource(resource string) scheme.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
