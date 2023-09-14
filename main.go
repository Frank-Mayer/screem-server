package main

import (
	"fmt"
	"net"
    "bytes"
	"os"
	"encoding/binary"
	"image"
	"image/png"
)

const (
	serverAddr    = "0.0.0.0:12345"
	maxPacketSize = 1400 // Recommended size for IPv4
	imageFilePath = "received_image.png"
)

func main() {
	// Create a connection
	conn, err := net.ListenPacket("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error creating connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("server is listening on", serverAddr)

	// Create a buffer to store received packet data
	buffer := make([]byte, maxPacketSize)

	// Create a buffer to accumulate the image data
	var imageBuffer []byte
	var imageSize int

	for {
		// Read a packet from the connection
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			continue
		}

		// If imageBuffer is empty, the first packet should contain the image size
		if len(imageBuffer) == 0 {
			if n != 4 {
				fmt.Println("Invalid image size packet")
				continue
			}
			imageSize = int(binary.BigEndian.Uint32(buffer[:4]))
			imageBuffer = make([]byte, imageSize)
			fmt.Printf("Receiving image with size: %d bytes\n", imageSize)
			continue
		}

		// Append the packet data to the image buffer
		imageBuffer = append(imageBuffer, buffer[:n]...)

		// If the image buffer is now the same size as the expected image, decode it
		if len(imageBuffer) == imageSize {
			decodedImage, _, err := image.Decode(bytes.NewReader(imageBuffer))
			if err != nil {
				fmt.Println("Error decoding image:", err)
				continue
			}

			// Save the received image to a file
			saveImageToFile(decodedImage.(*image.RGBA), imageFilePath)
			fmt.Println("Received and saved image to", imageFilePath)

			// Reset the image buffer and size
			imageBuffer = nil
			imageSize = 0
		} else {
            fmt.Printf("Received %d bytes of image data\n", len(buffer))
            os.Exit(1)
        }
	}
}

func saveImageToFile(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating image file:", err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Println("Error encoding image to file:", err)
	}
}

