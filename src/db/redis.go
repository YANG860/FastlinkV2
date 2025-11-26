package db

import (
	"fastlink/src/config"
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

func CacheLink(link Link) error {
	//todo
	err := RedisClient.Set(Ctx, "fastlink:link:"+link.ShortCode, link, config.Redis().LinkTTL).Err()
	return err

}

func FetchLink(shortCode string) (Link, error) {
	//todo
	return Link{}, nil
}
