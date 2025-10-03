package utils

import (
	"confession-wall-backend/config/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 生成token
func GenerateToken(userID uint) (string, error) {
	claims:=jwt.MapClaims{
        "user_id":userID,
        "exp":time.Now().Add(time.Hour*24).Unix(),
        "iat":time.Now().Unix(),
        "iss":"gin_jwt_demo",
        "nbf":time.Now().Unix(),
    }
    token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
    tokenStr,err:=token.SignedString([]byte(config.Config.GetString("jwt.secret")))
    if err!=nil{
        return "",err
    }
    return tokenStr,nil
}

//验证token
func ParseToken(tokenStr string) (*jwt.Token, error) {
    //检查签名方法是否正确
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil,jwt.ErrSignatureInvalid
        }
        return []byte(config.Config.GetString("jwt.secret")), nil
    })
        
    if err != nil {
        return nil,err
    }
	return token,nil
}
func ExtractClaims(token *jwt.Token)(jwt.MapClaims,error){
    if claims,ok:=token.Claims.(jwt.MapClaims);ok&&token.Valid{
        return claims,nil
    }else{
        return nil,jwt.ErrTokenInvalidClaims 
    }
}