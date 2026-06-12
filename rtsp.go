package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"time"
)

/*
	Parametros ffmpeg

-f dshow	Usa DirectShow, la interfaz de Windows para acceder a cámaras y micrófonos.
-rtbufsize 100M	Reserva un buffer de 100 MB para evitar pérdida de frames cuando la captura va más rápido que el procesamiento.
-framerate 30	Captura 30 imágenes por segundo (30 FPS).
-video_size 640x480	Resolución del video. Más resolución = más calidad pero más consumo de CPU y red.
-i video="HD camera "	Dispositivo de entrada. En este caso la webcam llamada "HD camera".
-c:v libx264	Codifica el video usando H.264 mediante la biblioteca x264.
-preset ultrafast	Usa la configuración más rápida posible del codificador para reducir la latencia.
-tune zerolatency	Optimiza el encoder para streaming en tiempo real, reduciendo buffers internos.
-g 6	Genera un frame clave (I-Frame) cada 6 frames. Ayuda a que los clientes se conecten más rápido y reduce la latencia.
-bf 0	Desactiva B-Frames. Mejora la latencia porque evita que FFmpeg tenga que esperar frames futuros.
-pix_fmt yuv420p	Convierte el formato de color a uno compatible con casi todos los reproductores y navegadores.
-rtsp_transport tcp	Usa TCP para transportar RTSP. Más confiable que UDP aunque puede agregar algo de latencia.
-f rtsp	Indica que el formato de salida será RTSP.
rtsp://127.0.0.1:8554/webcam	Dirección donde se publica el stream en MediaMTX.
*/
func main() {

	// FFmpeg (publicador RTSP)
	ffmpeg := exec.Command(
		"ffmpeg",
		"-f", "dshow",
		"-rtbufsize", "10M",
		"-framerate", "30",
		"-video_size", "640x480",
		"-i", `video=HD camera `,
		"-c:v", "libx264", // codificacion H.264
		"-preset", "ultrafast", // es la configuración más rápida posible del codificador
		"-tune", "zerolatency", // Optimiza el encoder para streaming en tiempo real
		"-g", "30", // cada x frames, genera un I-Frame
		"-bf", "0", // B-Frames
		"-pix_fmt", "yuv420p",
		"-rtsp_transport", "udp",
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
