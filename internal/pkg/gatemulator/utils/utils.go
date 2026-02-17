package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
)

func RandomPhoneNumber() string {
	prefixes := []string{"7708", "7777", "7383"}
	prefix := prefixes[rand.Intn(len(prefixes))]

	remainingDigits := rand.Intn(1000000)
	phone := fmt.Sprintf("%s%06d", prefix, remainingDigits)
	return phone
}

var names = []string{"John", "Agneshka", "Alice", "Emma", "Sophia", "Oliver", "Liam", "James", "Amelia", "Mia", "Ella"}

func RandomName() string {
	return names[rand.Intn(len(names))]
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func RandomWhatsappMessengerId() string {
	return fmt.Sprintf("false_%d@c.us_%s", rand.Intn(9999999999), RandomString(32))
}

func GenerateRandomImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(rand.Intn(256)),
				G: uint8(rand.Intn(256)),
				B: uint8(rand.Intn(256)),
				A: 255,
			})
		}
	}

	red := image.NewUniform(color.RGBA{R: 255, A: 255})
	rect := image.Rect(50, 50, width-50, height-50)
	draw.Draw(img, rect, red, image.Point{}, draw.Src)

	return img
}
