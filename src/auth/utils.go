package auth

import (
	"errors"
	"fastlink/src/db"

	"gorm.io/gorm"
)

func RefreshTokenValidInDB(claims *Token) (bool, error) {

	user, err := gorm.G[db.User](db.MySQLClient).Where("id = ?", claims.UserID).First(db.Ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) || user.AccessTokenID != claims.ID {
		return false, nil
	}

	return true, err
}
