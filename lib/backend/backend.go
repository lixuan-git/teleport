/*
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Package backend provides storage backend abstraction layer
package backend

import (
	"context"
	"fmt"
	"iter"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"

	"github.com/gravitational/teleport/api/types"
)

// Forever means that object TTL will not expire unless deleted
const (
	Forever time.Duration = 0
)

// ErrIncorrectRevision is returned from conditional operations when revisions
// do not match the expected value.
var ErrIncorrectRevision = &trace.CompareFailedError{Message: "resource revision does not match, it may have been concurrently created|modified|deleted; please work from the latest state, or use --force to overwrite"}

// Backend implements abstraction over local or remote storage backend.
// Item keys are assumed to be valid UTF8, which may be enforced by the
// various Backend implementations.
type Backend interface {
	// GetName returns the implementation driver name.
	GetName() string

	// Create creates item if it does not exist
	Create(ctx context.Context, i Item) (*Lease, error)

	// Put puts value into backend (creates if it does not
	// exists, updates it otherwise)
	Put(ctx context.Context, i Item) (*Lease, error)

	// CompareAndSwap compares item with existing item
	// and replaces is with replaceWith item
	CompareAndSwap(ctx context.Context, expected Item, replaceWith Item) (*Lease, error)

	// Update updates value in the backend
	Update(ctx context.Context, i Item) (*Lease, error)

	// Get returns a single item or not found error
	Get(ctx context.Context, key Key) (*Item, error)

	// Items produces an iterator of backend items in the range, and order
	// described in the provided [ItemsParams].
	Items(ctx context.Context, params ItemsParams) iter.Seq2[Item, error]

	// GetRange returns the items between the start and end keys, including both
	// (if present).
	GetRange(ctx context.Context, startKey, endKey Key, limit int) (*GetResult, error)

	// Delete deletes item by key, returns NotFound error
	// if item does not exist
	Delete(ctx context.Context, key Key) error

	// DeleteRange deletes range of items with keys between startKey and endKey
	DeleteRange(ctx context.Context, startKey, endKey Key) error

	// KeepAlive keeps object from expiring, updates lease on the existing object,
	// expires contains the new expiry to set on the lease,
	// some backends may ignore expires based on the implementation
	// in case if the lease managed server side
	KeepAlive(ctx context.Context, lease Lease, expires time.Time) error

	// ConditionalUpdate updates the value in the backend if the revision of the [Item] matches
	// the stored revision.
	ConditionalUpdate(ctx context.Context, i Item) (*Lease, error)

	// ConditionalDelete deletes the item by key if the revision matches the stored revision.
	ConditionalDelete(ctx context.Context, key Key, revision string) error

	// AtomicWrite executes a batch of conditional actions atomically s.t. all actions happen if all
	// conditions are met, but no actions happen if any condition fails to hold. If one or more conditions
	// failed to hold, [ErrConditionFailed] is returned. The number of conditional actions must not
	// exceed [MaxAtomicWriteSize] and no two conditional actions may point to the same key. If successful,
	// the returned revision is the new revision associated with all [Put] actions that were part of the
	// operation (the revision value has no meaning outside of the context of puts).
	AtomicWrite(ctx context.Context, condacts []ConditionalAction) (revision string, err error)

	// NewWatcher returns a new event watcher
	NewWatcher(ctx context.Context, watch Watch) (Watcher, error)

	// Close closes backend and all associated resources
	Close() error

	// Clock returns clock used by this backend
	Clock() clockwork.Clock

	// CloseWatchers closes all the watchers
	// without closing the backend
	CloseWatchers()
}

// ItemsParams are parameters that are provided to
// [BackendWithItems.Items] to alter the iteration behavior.
type ItemsParams struct {
	// StartKey is the minimum key in the range yielded by the iteration. This key
	// will be included in the results if it exists.
	StartKey Key
	// EndKey is the maximum key in the range yielded by the iteration. This key
	// will be included in the results if it exists.
	EndKey Key
	// Descending makes the iteration yield items from the biggest to the smallest
	// key (i.e. from EndKey to StartKey). If unset, the iteration will proceed in the
	// usual ascending order (i.e. from StartKey to EndKey).
	Descending bool
	// Limit is an optional maximum number of items to retrieve during iteration.
	Limit int
}

// New initializes a new [Backend] implementation based on the service config.
func New(ctx context.Context, backend string, params Params) (Backend, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	newbk, ok := registry[backend]
	if !ok {
		return nil, trace.BadParameter("unsupported secrets storage type: %q", backend)
	}
	bk, err := newbk(ctx, params)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return bk, nil
}

// Lease represents a lease on the item that can be used
// to extend item's TTL without updating its contents.
//
// Here is an example of renewing object TTL:
//
//	item.Expires = time.Now().Add(10 * time.Second)
//	lease, err := backend.Create(ctx, item)
//	expires := time.Now().Add(20 * time.Second)
//	err = backend.KeepAlive(ctx, lease, expires)
type Lease struct {
	// Key is the resource identifier.
	Key Key
	// Revision is the last known version of the object.
	Revision string
}

// Watch specifies watcher parameters
type Watch struct {
	// Name is a watch name set for debugging
	// purposes
	Name string
	// Prefixes specifies prefixes to watch,
	// passed to the backend implementation
	Prefixes []Key
	// QueueSize is an optional queue size
	QueueSize int
	// MetricComponent if set will start reporting
	// with a given component metric
	MetricComponent string
}

// String returns a user-friendly description
// of the watcher
func (w *Watch) String() string {
	return fmt.Sprintf("Watcher(name=%v, prefixes=%v)", w.Name, w.Prefixes)
}

// Watcher returns watcher
type Watcher interface {
	// Events returns channel with events
	Events() <-chan Event

	// Done returns the channel signaling the closure
	Done() <-chan struct{}

	// Close closes the watcher and releases
	// all associated resources
	Close() error
}

// GetResult provides the result of GetRange request
type GetResult struct {
	// Items returns a list of items
	Items []Item
}

// Event is a event containing operation with item
type Event struct {
	// Type is operation type
	Type types.OpType
	// Item is event Item
	Item Item
}

// Item is a key value item
type Item struct {
	// Key is a key of the key value item
	Key Key
	// Value is a value of the key value item
	Value []byte
	// Expires is an optional record expiry time
	Expires time.Time
	// Revision is the last known version of the object.
	Revision string
}

func (e Event) String() string {
	val := string(e.Item.Value)
	if len(val) > 20 {
		val = val[:20] + "..."
	}
	return fmt.Sprintf("%v %s=%s", e.Type, e.Item.Key, val)
}

// Config is used for 'storage' config section. It's a combination of
// values for various backends: 'etcd', 'filesystem', 'dynamodb', etc.
type Config struct {
	// Type indicates which backend to use (etcd, dynamodb, etc)
	Type string `yaml:"type,omitempty"`

	// Params is a generic key/value property bag which allows arbitrary
	// values to be passed to backend
	Params Params `yaml:",inline"`
}

// Params type defines a flexible unified back-end configuration API.
// It is just a map of key/value pairs which gets populated by `storage` section
// in Teleport YAML config.
type Params map[string]any

// GetString returns a string value stored in Params map, or an empty string
// if nothing is found
func (p Params) GetString(key string) string {
	v, ok := p[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

// NoLimit specifies no limits
const NoLimit = 0

const noEnd = "\x00"

// RangeEnd returns end of the range for given key.
func RangeEnd(key Key) Key {
	end := make([]byte, len(key.s))
	copy(end, key.s)
	for i := len(end) - 1; i >= 0; i-- {
		if end[i] < 0xff {
			end[i] = end[i] + 1
			end = end[:i+1]
			return KeyFromString(string(end))
		}
	}
	// next key does not exist (e.g., 0xffff);
	return Key{noEnd: true}
}

// HostID is a derivation of a KeyedItem that allows the host id
// to be included in the key.
type HostID interface {
	KeyedItem
	GetHostID() string
}

// KeyedItem represents an item from which a pagination key can be derived.
type KeyedItem interface {
	GetName() string
}

// GetPaginationKey returns the pagination key given item.
// For items that implement HostID, the next key will also
// have the HostID part.
func GetPaginationKey(ki KeyedItem) string {
	if h, ok := ki.(HostID); ok {
		return internalKey(h.GetHostID(), h.GetName()).String()
	}

	return ki.GetName()
}

// MaskKeyName masks the given key name.
// e.g "123456789" -> "******789"
func MaskKeyName(keyName string) string {
	maskedBytes := []byte(keyName)
	hiddenBefore := int(0.75 * float64(len(keyName)))
	for i := range hiddenBefore {
		maskedBytes[i] = '*'
	}
	return string(maskedBytes)
}

// Items is a sortable list of backend items
type Items []Item

// Len is part of sort.Interface.
func (it Items) Len() int {
	return len(it)
}

// Swap is part of sort.Interface.
func (it Items) Swap(i, j int) {
	it[i], it[j] = it[j], it[i]
}

// Less is part of sort.Interface.
func (it Items) Less(i, j int) bool {
	return it[i].Key.Compare(it[j].Key) < 0
}

// TTL returns TTL in duration units, rounds up to one second
func TTL(clock clockwork.Clock, expires time.Time) time.Duration {
	ttl := expires.Sub(clock.Now())
	if ttl < time.Second {
		return time.Second
	}
	return ttl
}

// EarliestExpiry returns first of the
// otherwise returns empty
func EarliestExpiry(times ...time.Time) time.Time {
	if len(times) == 0 {
		return time.Time{}
	}
	sort.Sort(earliest(times))
	return times[0]
}

// Expiry converts ttl to expiry time, if ttl is 0
// returns empty time
func Expiry(clock clockwork.Clock, ttl time.Duration) time.Time {
	if ttl == 0 {
		return time.Time{}
	}
	return clock.Now().UTC().Add(ttl)
}

type earliest []time.Time

func (p earliest) Len() int {
	return len(p)
}

func (p earliest) Less(i, j int) bool {
	if p[i].IsZero() {
		return false
	}
	if p[j].IsZero() {
		return true
	}
	return p[i].Before(p[j])
}

func (p earliest) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// CreateRevision generates a new identifier to be used
// as a resource revision. Backend implementations that provide
// their own mechanism for versioning resources should be
// preferred.
func CreateRevision() string {
	return uuid.NewString()
}

// BlankRevision is a placeholder revision to be used by backends when
// the revision of the item in the backend is empty. This can happen
// to any existing resources that were last written before support for
// revisions was added.
var BlankRevision = uuid.Nil.String()

// NewLease creates a lease for the provided [Item].
func NewLease(item Item) *Lease {
	return &Lease{
		Key:      item.Key,
		Revision: item.Revision,
	}
}
