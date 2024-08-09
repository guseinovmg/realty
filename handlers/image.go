package handlers

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

const (
	uploadDir    = "./uploads"
	optimizedDir = "./optimized"
)

/*
func main() {
// Create directories if they don't exist
os.MkdirAll(uploadDir, os.ModePerm)
os.MkdirAll(optimizedDir, os.ModePerm)

http.HandleFunc("/upload", uploadImage)
http.HandleFunc("/image/", serveImage)

fmt.Println("Server started at http://localhost:8080")
http.ListenAndServe(":8080", nil)
}
*/
func uploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := filepath.Base(handler.Filename)
	filePath := filepath.Join(uploadDir, filename)
	optimizedFilePath := filepath.Join(optimizedDir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	if err := optimizeImage(filePath, optimizedFilePath); err != nil {
		http.Error(w, "Error optimizing the image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Image uploaded and optimized successfully")
}

func optimizeImage(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize image to a smaller size (e.g., 800x600)
	resizedImg := resize.Resize(800, 600, img, resize.Lanczos3)

	// Save the optimized image
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Determine the image format and save accordingly
	switch filepath.Ext(inputPath) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 80})
	case ".png":
		err = png.Encode(outFile, resizedImg)
	default:
		err = fmt.Errorf("unsupported file format")
	}

	return err
}

func serveImage(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/image/"):]
	filePath := filepath.Join(optimizedDir, filename)

	http.ServeFile(w, r, filePath)
}
