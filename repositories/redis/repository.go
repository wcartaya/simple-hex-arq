package redis

import (
	"fmt"
	"strconv"
	"users-ms/domain"
	"github.com/go-redis/redis"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedisRepository(redisURL string) (domain.UserRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, fmt.Errorf("repository.NewRedisRepository %w", err)
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(userId string) (*domain.User, error) {
	user := &domain.User{}
	data, err := r.client.HGetAll(userId).Result()
	if err != nil {
		return nil, fmt.Errorf("repository.Find %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("repository.Find %w", domain.ErrUserNotFound)
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("repository.Find %w", err)
	}
	user.UserId = data["userId"]
	user.Name = data["name"]
	user.CreatedAt = createdAt
	return user, nil
}

func (r *redisRepository) Store(user *domain.User) error {
	data := map[string]interface{}{
		"user_id":       user.UserId,
		"name":        user.Name,
		"created_at": user.CreatedAt,
	}
	_, err := r.client.HMSet(user.UserId, data).Result()
	if err != nil {
		return fmt.Errorf("repository.Store %w", err)
	}
	return nil
}
