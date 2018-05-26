package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var imageRoot string = os.Getenv("IMAGE_ROOT_PATH")

func main() {
	http.HandleFunc("/b/", brandImageHandler)
	http.HandleFunc("/p/", productImageHandler)

	//http.Handle("/", http.FileServer(http.Dir(*root))) wow

	log.Println("Listening on 9000")
	//Serving HTTP
	err := http.ListenAndServe(":9000", nil)
	//Serving HTTPS
	//err := http.ListenAndServeTLS(":9000", "/etc/nginx/ssl/certificates.pem", "/etc/nginx/ssl/private.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func brandImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling:" + r.URL.EscapedPath())
	reg, _ := regexp.Compile("/b/(.*)")

	if match := reg.FindStringSubmatch(r.URL.EscapedPath()); match != nil {
		var (
			width  uint
			height uint
		)
		width, height = getImageSize("brand")
		imagePath := fmt.Sprintf("%vbrand/%v", imageRoot, match[1])
		img, err := resizeImage(imagePath, width, height)
		if err != nil {
			returnErrorMessageToBrowser(w, r, http.StatusNotFound)
		} else {
			sendImageToBrowser(w, img)
			return
		}
	} else {
		returnErrorMessageToBrowser(w, r, http.StatusForbidden)
	}
}

func productImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling:" + r.URL.EscapedPath())
	reg, _ := regexp.Compile("/p/.*-([0-9]{2})([0-9]+)-([0-9]?)-?(.+)(\\.[^.]+)")
	if match := reg.FindStringSubmatch(r.URL.EscapedPath()); match != nil {
		var (
			width  uint
			height uint
		)
		width, height = getImageSize(match[4])
		imagePath := fmt.Sprintf("%vproduct/%v/%v/%v.jpg", imageRoot, match[1], match[2], match[3])
		img, err := resizeImage(imagePath, width, height)
		if err != nil {
			returnErrorMessageToBrowser(w, r, http.StatusNotFound)
		} else {
			sendImageToBrowser(w, img)
			return
		}
	} else {
		returnErrorMessageToBrowser(w, r, http.StatusForbidden)
	}
}

func resizeImage(address string, x uint, y uint) (image.Image, error) {
	//read file from hard disk
	if _, err := os.Stat(address); os.IsNotExist(err) {
		log.Println("Error:" + address + " Not Found!")
		return nil, errors.New("File cannot be opened for resize!")
	}

	file, err := os.Open(address)
	if err != nil {
		log.Println("Error:" + address + " Not Accessible!")
		return nil, errors.New("File cannot be opened for resize!")
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Println("Error:" + address + " Conversion Failure!")
		return nil, errors.New("File cannot be decoded into an image!")
	}
	file.Close()

	// resize to new width and height
	return resize.Resize(x, y, img, resize.NearestNeighbor), nil
}

func sendImageToBrowser(w http.ResponseWriter, img image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")

	}
}

func getImageSize(imgType string) (width, height uint) {
	switch imgType {
	case "cart":
		width = 60
		height = 75
	case "catalog_grid_3":
		width = 236
		height = 295
	case "catalog":
		width = 176
		height = 220
	case "gallery":
		width = 35
		height = 44
	case "product":
		width = 292
		height = 365
	case "zoom":
		width = 680
		height = 850
	case "brand":
		width = 100
		height = 80
	default:
		width = 35
		height = 44
	}
	return
}

func returnErrorMessageToBrowser(w http.ResponseWriter, r *http.Request, errorCode int) {
	w.WriteHeader(errorCode)
	w.Write([]byte("☄ " + strconv.Itoa(errorCode)))
}
