package serviceaccount

import (
	crd "cloud.redhat.com/clowder/v2/apis/cloud.redhat.com/v1alpha1"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/errors"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/object"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/providers"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/utils"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createServiceAccountForClowdObj(cache *providers.ObjectCache, ident providers.ResourceIdent, obj object.ClowdObject, pullSecretNames crd.PullSecrets) error {

	if obj.GetClowdNamespace() == "" {
		err := errors.New("targetNamespace not yet populated")
		err.Requeue = true
		return err
	}

	nn := types.NamespacedName{
		Name:      obj.GetClowdSAName(),
		Namespace: obj.GetClowdNamespace(),
	}

	labeler := utils.GetCustomLabeler(nil, nn, obj)

	return CreateServiceAccount(cache, ident, pullSecretNames, nn, labeler)
}

func CreateServiceAccount(cache *providers.ObjectCache, ident providers.ResourceIdent, pullSecretNames crd.PullSecrets, nn types.NamespacedName, labeler func(v1.Object)) error {

	sa := &core.ServiceAccount{}
	if err := cache.Create(ident, nn, sa); err != nil {
		return err
	}

	sa.ImagePullSecrets = []core.LocalObjectReference{}

	for _, pullSecret := range pullSecretNames {
		sa.ImagePullSecrets = append(sa.ImagePullSecrets, core.LocalObjectReference{
			Name: string(pullSecret),
		})
	}

	labeler(sa)

	if err := cache.Update(ident, sa); err != nil {
		return err
	}

	return nil
}