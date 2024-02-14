package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

type MBR struct {
	MbrTamano        int64
	MbrFechaCreacion [20]byte
	MbrDiskSignature int64
}

var archivoBinarioDisco string = "disk.bin"

func main() {
	limpiarConsola()
	PrintCopyright()
	fmt.Println("Sistema de Archivos - Tarea2")

	var input string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Ingrese el comando:")
	scanner.Scan()
	input = scanner.Text()

	comando, path := parseCommand(input)
	if comando != "EXECUTE" || path == "" {
		fmt.Println("Comando no reconocido o ruta de archivo faltante. Uso: EXECUTE <ruta_al_archivo_de_scripts>")
		return
	}

	path = strings.Trim(path, `"'`)

	fmt.Printf("Leyendo el archivo de scripts de: %s\n", path)

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error leyendo el archivo de scripts: %v\n", err)
		return
	}

	contentStr := string(content)
	contentStr = strings.Replace(contentStr, "\r\n", "\n", -1) // Convertir CRLF a LF
	commands := strings.Split(contentStr, "\n")

	for _, command := range commands {
		cmd := string(command)
		if cmd == "MKDISK" {
			mkdisk(archivoBinarioDisco)
		} else if cmd == "REP" {
			printMBR(archivoBinarioDisco)
		}
	}
}

func parseCommand(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) < 2 {
		return "", ""
	}

	command := parts[0]
	var path string

	for _, part := range parts[1:] {
		if strings.HasPrefix(part, "-path=") {
			path = strings.TrimPrefix(part, "-path=")
			break
		}
	}

	return command, path
}

func mkdisk(filename string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error obteniendo el directorio actual: %v\n", err)
		return
	}
	fmt.Printf("Directorio actual: %s\n", wd)

	var mbr MBR
	mbr.MbrTamano = 5 * 1024 * 1024 // 5 MB
	currentTime := time.Now()
	copy(mbr.MbrFechaCreacion[:], currentTime.Format("2006-01-02T15:04:05"))
	mbr.MbrDiskSignature = 123456789 // Signature de ejemplo

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creando el disco: %v\n", err)
		return
	}
	defer file.Close()

	// Escribir el caracter inicial en el archivo
	_, err = file.Write([]byte{'\x00'})
	if err != nil {
		fmt.Printf("Error escribiendo el caracter inicial: %v\n", err)
		return
	}

	err = binary.Write(file, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Printf("Error escribiendo el MBR: %v\n", err)
		return
	}

	err = file.Truncate(int64(mbr.MbrTamano))
	if err != nil {
		fmt.Printf("Error al asignar espacio en disco: %v\n", err)
		return
	}

	fmt.Println("Disco creado correctamente.")
}

func printMBR(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error abriendo el disco: %v\n", err)
		return
	}
	defer file.Close()

	// Omitir el carácter nulo inicial
	_, err = file.Seek(1, 0)
	if err != nil {
		fmt.Printf("Error al buscar el byte nulo inicial: %v\n", err)
		return
	}

	// Leer y decodificar el MBR
	var mbr MBR
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		fmt.Printf("Error leyendo el MBR del disco: %v\n", err)
		return
	}

	fmt.Printf("Tamaño: %d bytes\nFecha de creación: %s\nDisk Signature: %d\n",
		mbr.MbrTamano,
		string(mbr.MbrFechaCreacion[:]),
		mbr.MbrDiskSignature)
}

func limpiarConsola() {
	fmt.Print("\033[H\033[2J")
}

func lineaDoble(longitud int) {
	fmt.Println(strings.Repeat("=", longitud))
}

func PrintCopyright() {
	lineaDoble(60)
	fmt.Println("Lino Antonio Garcia Vallejo")
	fmt.Println("Carné: 9017323")
	lineaDoble(60)
}
