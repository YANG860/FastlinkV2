package db

import (
	"encoding/json"
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

func CacheLink(link Link, NX bool) error {
	key := "fastlink:link:" + link.ShortCode
	data, err := json.Marshal(link)
	if err != nil {
		return err
	}

	// 
	if NX {
		err := RedisClient.SetNX(Ctx, key, data, config.Redis().LinkTTL).Err()
		return err
	} else {
		err := RedisClient.Set(Ctx, key, data, config.Redis().LinkTTL).Err()
		return err
	}
}

func FetchLink(shortCode string) (Link, error) {
	var link Link
	s, err := RedisClient.Get(Ctx, "fastlink:link:"+shortCode).Result()
	if err != nil {
		return Link{}, err
	}
	if err := json.Unmarshal([]byte(s), &link); err != nil {
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
	


func AddUsernameBloomFilter(username string) error {
	
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

func AddShortCodeBloomFilter(shortCode string) error {
	
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
