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

func CacheLink(link Link,NX bool) error {
	//todo
	if NX {
		err := RedisClient.SetNX(Ctx, "fastlink:link:"+link.ShortCode, link, config.Redis().LinkTTL).Err()
		return err
	} else {
		err := RedisClient.Set(Ctx, "fastlink:link:"+link.ShortCode, link, config.Redis().LinkTTL).Err()
		return err
	}
}

func FetchLink(shortCode string) (Link, error) {
	//todo
	return Link{}, nil
}


func UsernameBlomFilterAdd(username string) error {
	//todo
	// Use roaring bitmap

	err := RedisClient.BFAdd(Ctx, "fastlink:username:bloom", username).Err()
	if err != nil {
		return err
	}

	return nil
}

func UsernameBlomFilterCheck(username string) (bool, error) {
	//todo


	
	return false, nil
}


func ShortCodeBlomFilterAdd(shortCode string) error {
	//todo
	return nil
}

func ShortCodeBlomFilterCheck(shortCode string) (bool, error) {
	//todo
	return false, nil
}	


