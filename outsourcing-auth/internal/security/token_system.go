package security

import (
	"core/internal"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

func CreateToken(isCompany bool, accessID uint, lifetimeSec int) string {
	var (
		key []byte
		t   *jwt.Token
		s   string
		err error
	)
	key = []byte(internal.KeyJWT)
	t = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"isCompany": isCompany,
			"accessID":  accessID,
			"lifetime":  lifetimeSec, // in seconds
			"startTime": time.Now().Unix(),
		})
	s, err = t.SignedString(key)
	if err != nil {
		log.Println(err)
		return ""
	}
	return s
}

func CheckToken(tokenS string) (bool, jwt.MapClaims) {
	secretKey := internal.KeyJWT
	parsedToken, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Println(err)
		return false, nil
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		newTime := time.Now()

		startTimeFloat, ok := claims["startTime"].(float64)
		if !ok {
			log.Println("Invalid startTime format")
			return false, nil
		}
		startTime := time.Unix(int64(startTimeFloat), 0)

		lifetimeFloat, ok := claims["lifetime"].(float64)
		if !ok {
			log.Println("Invalid lifetime format")
			return false, nil
		}

		duration := newTime.Sub(startTime)
		minutes := duration.Minutes()
		elapsedSeconds := minutes * 60

		if elapsedSeconds-lifetimeFloat > lifetimeFloat {
			fmt.Println(elapsedSeconds, lifetimeFloat, "not ok")
			return false, claims
		} else {
			fmt.Println(elapsedSeconds, lifetimeFloat, "ok")
			return true, claims
		}
	} else {
		log.Println("Invalid token")
		return false, nil
	}
}
