package services

import (
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"image/png"
	"os"

	//"image/png"
	//"os"
	"strings"
	"time"
)

// getOTP is a function that generates TOTP using github.com/pquerna/otp
//
// Parameters:
//   - seed: OTP seed is base32
//   - digits: the number of digits for the OTP
//   - algorithm: "sha1", "sha256", "sha512"
//   - skew: the number of minutes to search around the OTP
//   - period: the number of seconds for the OTP to change
func GetOTP(seed string, digits int, algorithm string, skew uint, period uint) string {
	algo := getOTPAlgorithm(strings.ToLower(algorithm))
	opts := totp.ValidateOpts{
		Algorithm: algo,
		Digits:    otp.Digits(digits),
		Skew:      skew,
		Period:    period,
	}

	pass, err := totp.GenerateCodeCustom(seed, time.Now().UTC(), opts)
	if err != nil {
		core.Trail(core.ERROR, "Unable to generate OTP. %s", err)
		return ""
	}
	return pass
}

func VerifyOTP(pass, seed string, digits int, algorithm string, skew uint, period uint) bool {
	algo := getOTPAlgorithm(strings.ToLower(algorithm))
	opts := totp.ValidateOpts{
		Algorithm: algo,
		Digits:    otp.Digits(digits),
		Skew:      skew,
		Period:    period,
	}

	pass = fmt.Sprintf("%0"+fmt.Sprintf("%d.%d", digits, digits)+"s", pass)

	valid, err := totp.ValidateCustom(pass, seed, time.Now().UTC(), opts)
	if err != nil {
		core.Trail(core.ERROR, "Unable to verify OTP. %s", err)
		return false
	}
	return valid
}

func GenerateOTPSeed(digits int, algorithm string, skew uint, period uint, user core.IUser) (secret string, imagePath string) {
	algo := getOTPAlgorithm(strings.ToLower(algorithm))

	opts := totp.GenerateOpts{
		AccountName: user.GetAccountIdentifier(),
		Issuer:      core.CurrentConfig.D.GoMonolith.SiteName,
		Algorithm:   algo,
		Digits:      otp.Digits(digits),
		Period:      period,
		SecretSize:  64,
	}

	key, _ := totp.Generate(opts)
	img, _ := key.Image(250, 250)

	os.MkdirAll("./media/otp/", 0744)

	fName := "./media/otp/" + key.Secret() + ".png"
	for _, err := os.Stat(fName); os.IsExist(err); {
		key, _ = totp.Generate(opts)
		img, _ = key.Image(450, 450)
	}
	qrImg, _ := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE, 0644)
	defer qrImg.Close()

	png.Encode(qrImg, img)

	return key.Secret(), fName
}

func getOTPAlgorithm(algorithm string) otp.Algorithm {
	var algo otp.Algorithm
	switch algorithm {
	case "sha1":
		algo = otp.AlgorithmSHA1
	case "sha256":
		algo = otp.AlgorithmSHA256
	case "sha512":
		algo = otp.AlgorithmSHA512
	default:
		core.Trail(core.WARNING, "getOTPAlgorithm: Unknown hash algorithm (%s). Defaulting to sha1", algorithm)
		return otp.AlgorithmSHA1
	}
	return algo
}
