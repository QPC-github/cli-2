/*
Copyright 2019 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by injection-gen. DO NOT EDIT.

package trigger

import (
	context "context"

	apistriggersv1beta1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1beta1"
	versioned "github.com/tektoncd/triggers/pkg/client/clientset/versioned"
	v1beta1 "github.com/tektoncd/triggers/pkg/client/informers/externalversions/triggers/v1beta1"
	client "github.com/tektoncd/triggers/pkg/client/injection/client"
	factory "github.com/tektoncd/triggers/pkg/client/injection/informers/factory"
	triggersv1beta1 "github.com/tektoncd/triggers/pkg/client/listers/triggers/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Triggers().V1beta1().Triggers()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

func withDynamicInformer(ctx context.Context) context.Context {
	inf := &wrapper{client: client.Get(ctx)}
	return context.WithValue(ctx, Key{}, inf)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1beta1.TriggerInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch github.com/tektoncd/triggers/pkg/client/informers/externalversions/triggers/v1beta1.TriggerInformer from context.")
	}
	return untyped.(v1beta1.TriggerInformer)
}

type wrapper struct {
	client versioned.Interface

	namespace string
}

var _ v1beta1.TriggerInformer = (*wrapper)(nil)
var _ triggersv1beta1.TriggerLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apistriggersv1beta1.Trigger{}, 0, nil)
}

func (w *wrapper) Lister() triggersv1beta1.TriggerLister {
	return w
}

func (w *wrapper) Triggers(namespace string) triggersv1beta1.TriggerNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace}
}

func (w *wrapper) List(selector labels.Selector) (ret []*apistriggersv1beta1.Trigger, err error) {
	lo, err := w.client.TriggersV1beta1().Triggers(w.namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector.String(),
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apistriggersv1beta1.Trigger, error) {
	return w.client.TriggersV1beta1().Triggers(w.namespace).Get(context.TODO(), name, v1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}
