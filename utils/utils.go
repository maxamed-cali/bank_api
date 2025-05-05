package utils
import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("BankSecretKeyzcfsdwgy6euiiw10") // Replace with your own secret key
func GenerateJWT(userID uint, roleName string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
         "role":roleName,
         "exp":     time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}


func ValidateToken(tokenString string) (uint, string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return 0,"", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := uint(claims["user_id"].(float64))
        role := claims["role"].(string)
        return userID, role, nil
    }

    return 0,"", err
}