package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {

	// FFmpeg (publicador RTSP)
	ffmpeg := exec.Command(
		"ffmpeg",
		"-f", "dshow",
		"-rtbufsize", "100M",
		"-framerate", "15",
		"-video_size", "640x480",
		"-i", `video=HD camera `,
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-tune", "zerolatency",
		"-g", "15",
		"-bf", "0",
		"-pix_fmt", "yuv420p",
		"-rtsp_transport", "tcp",
		"-f", "rtsp",
		"rtsp://127.0.0.1:8554/webcam",
	)

	ffmpeg.Stdout = os.Stdout
	ffmpeg.Stderr = os.Stderr

	log.Println("Iniciando FFmpeg...")

	if err := ffmpeg.Start(); err != nil {
		log.Fatal(err)
	}

	// Esperar a que el stream esté disponible
	time.Sleep(3 * time.Second)

	// FFplay (cliente)
	ffplay := exec.Command(
		"ffplay",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-rtsp_transport", "tcp",
		"rtsp://127.0.0.1:8554/webcam",
	)

	ffplay.Stdout = os.Stdout
	ffplay.Stderr = os.Stderr

	log.Println("Iniciando FFplay...")

	if err := ffplay.Start(); err != nil {
		log.Fatal(err)
	}

	log.Println("Streaming activo.")
	log.Println("Presiona ENTER para detener la transmisión.")

	// Esperar ENTER
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	log.Println("Deteniendo FFmpeg...")

	if ffmpeg.Process != nil {
		ffmpeg.Process.Kill()
	}

	log.Println("Deteniendo FFplay...")

	if ffplay.Process != nil {
		ffplay.Process.Kill()
	}

	log.Println("Transmisión finalizada.")
}
