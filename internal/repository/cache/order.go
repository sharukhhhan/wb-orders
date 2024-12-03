package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/repository/repoerrors"
)

type OrderCache struct {
	*bigcache.BigCache
}

func NewOrderCache(bigCache *bigcache.BigCache) *OrderCache {
	return &OrderCache{BigCache: bigCache}
}

func (c *OrderCache) GetCache(orderUID string) (*entity.Order, error) {
	data, err := c.Get(orderUID)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, repoerrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get order from cache: %w", err)
	}

	var order entity.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order from cache: %w", err)
	}

	return &order, nil
}

func (c *OrderCache) SaveCache(order *entity.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order for cache: %w", err)
	}

	if err := c.Set(order.OrderUID, data); err != nil {
		return fmt.Errorf("failed to save order to cache: %w", err)
	}

	return nil
}
