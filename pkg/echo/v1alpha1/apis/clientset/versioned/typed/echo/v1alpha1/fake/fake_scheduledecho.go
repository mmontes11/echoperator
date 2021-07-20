/*
MIT License

Copyright (c) 2021 Martín Montes

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/mmontes11/echoperator/pkg/echo/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeScheduledEchos implements ScheduledEchoInterface
type FakeScheduledEchos struct {
	Fake *FakeMmontesV1alpha1
	ns   string
}

var scheduledechosResource = schema.GroupVersionResource{Group: "mmontes.io", Version: "v1alpha1", Resource: "scheduledechos"}

var scheduledechosKind = schema.GroupVersionKind{Group: "mmontes.io", Version: "v1alpha1", Kind: "ScheduledEcho"}

// Get takes name of the scheduledEcho, and returns the corresponding scheduledEcho object, and an error if there is any.
func (c *FakeScheduledEchos) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ScheduledEcho, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(scheduledechosResource, c.ns, name), &v1alpha1.ScheduledEcho{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ScheduledEcho), err
}

// List takes label and field selectors, and returns the list of ScheduledEchos that match those selectors.
func (c *FakeScheduledEchos) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ScheduledEchoList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(scheduledechosResource, scheduledechosKind, c.ns, opts), &v1alpha1.ScheduledEchoList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ScheduledEchoList{ListMeta: obj.(*v1alpha1.ScheduledEchoList).ListMeta}
	for _, item := range obj.(*v1alpha1.ScheduledEchoList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested scheduledEchos.
func (c *FakeScheduledEchos) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(scheduledechosResource, c.ns, opts))

}

// Create takes the representation of a scheduledEcho and creates it.  Returns the server's representation of the scheduledEcho, and an error, if there is any.
func (c *FakeScheduledEchos) Create(ctx context.Context, scheduledEcho *v1alpha1.ScheduledEcho, opts v1.CreateOptions) (result *v1alpha1.ScheduledEcho, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(scheduledechosResource, c.ns, scheduledEcho), &v1alpha1.ScheduledEcho{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ScheduledEcho), err
}

// Update takes the representation of a scheduledEcho and updates it. Returns the server's representation of the scheduledEcho, and an error, if there is any.
func (c *FakeScheduledEchos) Update(ctx context.Context, scheduledEcho *v1alpha1.ScheduledEcho, opts v1.UpdateOptions) (result *v1alpha1.ScheduledEcho, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(scheduledechosResource, c.ns, scheduledEcho), &v1alpha1.ScheduledEcho{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ScheduledEcho), err
}

// Delete takes name of the scheduledEcho and deletes it. Returns an error if one occurs.
func (c *FakeScheduledEchos) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(scheduledechosResource, c.ns, name), &v1alpha1.ScheduledEcho{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeScheduledEchos) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(scheduledechosResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ScheduledEchoList{})
	return err
}

// Patch applies the patch and returns the patched scheduledEcho.
func (c *FakeScheduledEchos) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ScheduledEcho, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(scheduledechosResource, c.ns, name, pt, data, subresources...), &v1alpha1.ScheduledEcho{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ScheduledEcho), err
}
