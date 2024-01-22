/*
Copyright 2024 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v3

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type NodeTemplateHandler func(string, *v3.NodeTemplate) (*v3.NodeTemplate, error)

type NodeTemplateController interface {
	generic.ControllerMeta
	NodeTemplateClient

	OnChange(ctx context.Context, name string, sync NodeTemplateHandler)
	OnRemove(ctx context.Context, name string, sync NodeTemplateHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() NodeTemplateCache
}

type NodeTemplateClient interface {
	Create(*v3.NodeTemplate) (*v3.NodeTemplate, error)
	Update(*v3.NodeTemplate) (*v3.NodeTemplate, error)
	UpdateStatus(*v3.NodeTemplate) (*v3.NodeTemplate, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.NodeTemplate, error)
	List(namespace string, opts metav1.ListOptions) (*v3.NodeTemplateList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.NodeTemplate, err error)
}

type NodeTemplateCache interface {
	Get(namespace, name string) (*v3.NodeTemplate, error)
	List(namespace string, selector labels.Selector) ([]*v3.NodeTemplate, error)

	AddIndexer(indexName string, indexer NodeTemplateIndexer)
	GetByIndex(indexName, key string) ([]*v3.NodeTemplate, error)
}

type NodeTemplateIndexer func(obj *v3.NodeTemplate) ([]string, error)

type nodeTemplateController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewNodeTemplateController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) NodeTemplateController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &nodeTemplateController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromNodeTemplateHandlerToHandler(sync NodeTemplateHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.NodeTemplate
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.NodeTemplate))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *nodeTemplateController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.NodeTemplate))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateNodeTemplateDeepCopyOnChange(client NodeTemplateClient, obj *v3.NodeTemplate, handler func(obj *v3.NodeTemplate) (*v3.NodeTemplate, error)) (*v3.NodeTemplate, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *nodeTemplateController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *nodeTemplateController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *nodeTemplateController) OnChange(ctx context.Context, name string, sync NodeTemplateHandler) {
	c.AddGenericHandler(ctx, name, FromNodeTemplateHandlerToHandler(sync))
}

func (c *nodeTemplateController) OnRemove(ctx context.Context, name string, sync NodeTemplateHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromNodeTemplateHandlerToHandler(sync)))
}

func (c *nodeTemplateController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *nodeTemplateController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *nodeTemplateController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *nodeTemplateController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *nodeTemplateController) Cache() NodeTemplateCache {
	return &nodeTemplateCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *nodeTemplateController) Create(obj *v3.NodeTemplate) (*v3.NodeTemplate, error) {
	result := &v3.NodeTemplate{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *nodeTemplateController) Update(obj *v3.NodeTemplate) (*v3.NodeTemplate, error) {
	result := &v3.NodeTemplate{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *nodeTemplateController) UpdateStatus(obj *v3.NodeTemplate) (*v3.NodeTemplate, error) {
	result := &v3.NodeTemplate{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *nodeTemplateController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *nodeTemplateController) Get(namespace, name string, options metav1.GetOptions) (*v3.NodeTemplate, error) {
	result := &v3.NodeTemplate{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *nodeTemplateController) List(namespace string, opts metav1.ListOptions) (*v3.NodeTemplateList, error) {
	result := &v3.NodeTemplateList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *nodeTemplateController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *nodeTemplateController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.NodeTemplate, error) {
	result := &v3.NodeTemplate{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type nodeTemplateCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *nodeTemplateCache) Get(namespace, name string) (*v3.NodeTemplate, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.NodeTemplate), nil
}

func (c *nodeTemplateCache) List(namespace string, selector labels.Selector) (ret []*v3.NodeTemplate, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.NodeTemplate))
	})

	return ret, err
}

func (c *nodeTemplateCache) AddIndexer(indexName string, indexer NodeTemplateIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.NodeTemplate))
		},
	}))
}

func (c *nodeTemplateCache) GetByIndex(indexName, key string) (result []*v3.NodeTemplate, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.NodeTemplate, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.NodeTemplate))
	}
	return result, nil
}

type NodeTemplateStatusHandler func(obj *v3.NodeTemplate, status v3.NodeTemplateStatus) (v3.NodeTemplateStatus, error)

type NodeTemplateGeneratingHandler func(obj *v3.NodeTemplate, status v3.NodeTemplateStatus) ([]runtime.Object, v3.NodeTemplateStatus, error)

func RegisterNodeTemplateStatusHandler(ctx context.Context, controller NodeTemplateController, condition condition.Cond, name string, handler NodeTemplateStatusHandler) {
	statusHandler := &nodeTemplateStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromNodeTemplateHandlerToHandler(statusHandler.sync))
}

func RegisterNodeTemplateGeneratingHandler(ctx context.Context, controller NodeTemplateController, apply apply.Apply,
	condition condition.Cond, name string, handler NodeTemplateGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &nodeTemplateGeneratingHandler{
		NodeTemplateGeneratingHandler: handler,
		apply:                         apply,
		name:                          name,
		gvk:                           controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterNodeTemplateStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type nodeTemplateStatusHandler struct {
	client    NodeTemplateClient
	condition condition.Cond
	handler   NodeTemplateStatusHandler
}

func (a *nodeTemplateStatusHandler) sync(key string, obj *v3.NodeTemplate) (*v3.NodeTemplate, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type nodeTemplateGeneratingHandler struct {
	NodeTemplateGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *nodeTemplateGeneratingHandler) Remove(key string, obj *v3.NodeTemplate) (*v3.NodeTemplate, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.NodeTemplate{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *nodeTemplateGeneratingHandler) Handle(obj *v3.NodeTemplate, status v3.NodeTemplateStatus) (v3.NodeTemplateStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.NodeTemplateGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
