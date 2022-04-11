// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dao

import (
	"context"
	"time"

	"github.com/goharbor/harbor/src/lib/errors"
	"github.com/goharbor/harbor/src/lib/orm"
	"github.com/goharbor/harbor/src/lib/q"
)

// DAO defines the DAO operations of registry
type DAO interface {
	// Create the registry
	Create(ctx context.Context, as *AccelerationService) (id int64, err error)
	// Count returns the count of registries according to the query
	Count(ctx context.Context, query *q.Query) (count int64, err error)
	// List the registries according to the query
	List(ctx context.Context, query *q.Query) (ases []*AccelerationService, err error)
	// Get the registry specified by ID
	Get(ctx context.Context, id int64) (as *AccelerationService, err error)
	// Update the specified registry
	Update(ctx context.Context, as *AccelerationService, props ...string) (err error)
	// Delete the registry specified by ID
	Delete(ctx context.Context, id int64) (err error)
}

// NewDAO creates an instance of DAO
func NewDAO() DAO {
	return &dao{}
}

type dao struct{}

func (d *dao) Create(ctx context.Context, as *AccelerationService) (int64, error) {
	ormer, err := orm.FromContext(ctx)
	if err != nil {
		return 0, err
	}
	id, err := ormer.Insert(as)
	if e := orm.AsConflictError(err, "AccelerationService %s already exists", as.Name); e != nil {
		err = e
	}
	return id, err
}

func (d *dao) Count(ctx context.Context, query *q.Query) (int64, error) {
	qs, err := orm.QuerySetterForCount(ctx, &AccelerationService{}, query)
	if err != nil {
		return 0, err
	}
	return qs.Count()
}

func (d *dao) List(ctx context.Context, query *q.Query) ([]*AccelerationService, error) {
	registries := []*AccelerationService{}
	qs, err := orm.QuerySetter(ctx, &AccelerationService{}, query)
	if err != nil {
		return nil, err
	}
	if _, err = qs.All(&registries); err != nil {
		return nil, err
	}
	return registries, nil
}

func (d *dao) Get(ctx context.Context, id int64) (*AccelerationService, error) {
	registry := &AccelerationService{
		ID: id,
	}
	ormer, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := ormer.Read(registry); err != nil {
		if e := orm.AsNotFoundError(err, "AccelerationService %d not found", id); e != nil {
			err = e
		}
		return nil, err
	}
	return registry, nil
}

func (d *dao) Update(ctx context.Context, registry *AccelerationService, props ...string) error {
	ormer, err := orm.FromContext(ctx)
	if err != nil {
		return err
	}
	registry.UpdateTime = time.Now()
	n, err := ormer.Update(registry, props...)
	if err != nil {
		if e := orm.AsConflictError(err, "AccelerationService %s already exists", registry.Name); e != nil {
			err = e
		}
		return err
	}
	if n == 0 {
		return errors.NotFoundError(nil).WithMessage("registry %d not found", registry.ID)
	}
	return nil
}

func (d *dao) Delete(ctx context.Context, id int64) error {
	ormer, err := orm.FromContext(ctx)
	if err != nil {
		return err
	}
	n, err := ormer.Delete(&AccelerationService{
		ID: id,
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.NotFoundError(nil).WithMessage("AccelerationService %d not found", id)
	}
	return nil
}
