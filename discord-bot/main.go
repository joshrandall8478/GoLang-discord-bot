package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token      string
	AvatarFile string
	AvatarURL  string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&AvatarFile, "f", "", "Avatar File Name")
	flag.StringVar(&AvatarURL, "u", "https://images-na.ssl-images-amazon.com/images/I/61a0edomNQL._AC_SY879_.jpg", "URL to the avatar image")
	flag.Parse()

	//Token = "Njk1MzA1NTczODA1MTk1MzU0.XoYP4w.3AwCp_bpOKs7Te-HYlutF_ou86M"
	if Token == "" || (AvatarFile == "" && AvatarURL == "") {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	var err error
	var dg *discordgo.Session

	// Create a new Discord session using the provided login information.
	dg, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Declare these here so they can be used in the below two if blocks and
	// still carry over to the end of this function.
	var base64img string
	var contentType string

	var img []byte

	// If we're using a URL link for the Avatar
	if AvatarURL != "" {
		var resp *http.Response
		resp, err = http.Get(AvatarURL)
		if err != nil {
			fmt.Println("Error retrieving the file, ", err)
			return
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		img, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading the response, ", err)
			return
		}

	}

	// If we're using a local file for the Avatar
	if AvatarFile != "" {
		img, err = ioutil.ReadFile(AvatarFile)
		if err != nil {
			fmt.Println(err)
		}
	}

	contentType = http.DetectContentType(img)

	pixels := make([]byte, 472*879) // slice of your gray pixels, size of 100x100

	grayImg := image.NewGray(image.Rect(0, 0, 472, 879))
	grayImg.Pix = pixels

	var uploadImage *image.Gray
	uploadImage = rgbaToGray(grayImg)

	base64img = base64.StdEncoding.EncodeToString([]byte(uploadImage.Pix))

	out, err := os.Create("./output.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = jpeg.Encode(out, uploadImage, nil) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Now lets format our base64 image into the proper format Discord wants
	// and then call UserUpdate to set it as our user's Avatar.
	//avatar := fmt.Sprintf("data:%s;base64,%s", contentType, base64img)
	//_, err = dg.UserUpdate("", "", "", avatar, "")
	//if err != nil {
	//	fmt.Println(err)
	//}
}

// Convert to grayscale
func rgbaToGray(img image.Image) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}
	return gray
}
