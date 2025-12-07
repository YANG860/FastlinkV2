package db

import (
	"encoding/json"
	"fastlink/src/config"

	"github.com/redis/go-redis/v9"
)

func CacheRefreshTokenID(UserID string, ID string, NX bool) error {
	var err error
	if NX {
		err = RedisClient.SetNX(Ctx, "fastlink:user:"+UserID+":refresh_token_id", ID, config.Redis().RefreshTokenTTL).Err()
	} else {
		err = RedisClient.Set(Ctx, "fastlink:user:"+UserID+":refresh_token_id", ID, config.Redis().RefreshTokenTTL).Err()
	}
	return err
}

func FetchRefreshTokenID(UserID string) (string, error) {
	ID, err := RedisClient.Get(Ctx, "fastlink:user:"+UserID+":refresh_token_id").Result()
	if err != nil {
		return "", err
	}

	return ID, nil
}

func UpdateRefreshTokenTTL(UserID string) error {
	err := RedisClient.Expire(Ctx, "fastlink:user:"+UserID+":refresh_token_id", config.Redis().RefreshTokenTTL).Err()
	return err
}

func CacheLink(link Link, NX bool) error {
	// 重构
	// 使用hash

	key := "fastlink:link:" + link.ShortCode
	fields := map[string]interface{}{
		"id":         link.ID,
		"created_at": link.CreatedAt,
		"updated_at": link.UpdatedAt,
		"deleted_at": link.DeletedAt,
		"type":       link.Type,
		"creator_id": link.CreatorID,
		"source_url": link.SourceURL,
		"short_code": link.ShortCode,
		"clicks":     link.Clicks,
	}

	err := RedisClient.HSet(Ctx, key, fields).Err()

	if err != nil {
		return err
	}

	return RedisClient.Expire(Ctx, key, config.Redis().LinkTTL).Err()
}

func FetchLink(shortCode string) (Link, error) {

	key := "fastlink:link:" + shortCode
	fields, err := RedisClient.HGetAll(Ctx, key).Result()
	if err != nil {
		return Link{}, err
	}
	if len(fields) == 0 {
		return Link{}, redis.Nil
	}

	var link Link

	data, err := json.Marshal(fields)
	if err != nil {
		return Link{}, err
	}
	if err := json.Unmarshal(data, &link); err != nil {
		return Link{}, err
	}
	return link, nil
}

func DeleteLinkCache(shortCode string) error {
	err := RedisClient.Del(Ctx, "fastlink:link:"+shortCode).Err()
	return err
}

func UpdateLinkTTL(shortCode string) error {
	err := RedisClient.Expire(Ctx, "fastlink:link:"+shortCode, config.Redis().LinkTTL).Err()
	return err
}

func UpdateLinkClicks(shortCode string) error {
	err := RedisClient.HIncrBy(Ctx, "fastlink:link:"+shortCode, "clicks", 1).Err()
	return err
}

func UsernameBloomFilterAdd(username string) error {

	err := RedisClient.BFAdd(Ctx, "fastlink:username:bloom", username).Err()
	if err != nil {
		return err
	}

	return nil
}

func UsernameBloomFilterExists(username string) (bool, error) {

	exists, err := RedisClient.BFExists(Ctx, "fastlink:username:bloom", username).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func ShortCodeBloomFilterAdd(shortCode string) error {

	err := RedisClient.BFAdd(Ctx, "fastlink:shortcode:bloom", shortCode).Err()
	if err != nil {
		return err
	}

	return nil
}

func ShortCodeBloomFilterExists(shortCode string) (bool, error) {

	exists, err := RedisClient.BFExists(Ctx, "fastlink:shortcode:bloom", shortCode).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}
