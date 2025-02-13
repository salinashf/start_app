package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func buscarArchivos(nombreArchivo string, directorio string) ([]string, error) {
	var coincidencias []string

	err := filepath.Walk(directorio, func(ruta string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Comparar los nombres de archivo en minúsculas
		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(nombreArchivo)) {
			coincidencias = append(coincidencias, ruta)
		}
		if len(coincidencias) >= 1 {
			// Detener la búsqueda después de encontrar las primeras dos coincidencias
			return fmt.Errorf("encontradas dos coincidencias")
		}
		return nil
	})

	// Ignorar el error de "encontradas dos coincidencias"
	if err != nil && err.Error() != "encontradas dos coincidencias" {
		return nil, err
	}

	return coincidencias, nil
}

func obtenerRutaEjecutableDeAccesoDirecto(rutaAccesoDirecto string) (string, error) {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// Crear objeto COM WScript.Shell
	shell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return "", fmt.Errorf("error creando objeto WScript.Shell: %w", err)
	}
	defer shell.Release()

	// Crear una instancia de IDispatch
	wshell, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return "", fmt.Errorf("error creando IDispatch: %w", err)
	}
	defer wshell.Release()

	// Obtener el acceso directo como objeto COM
	shortcut := oleutil.MustCallMethod(wshell, "CreateShortcut", rutaAccesoDirecto)
	defer shortcut.Clear()

	// Obtener la propiedad TargetPath (ruta del ejecutable)
	targetPath := oleutil.MustGetProperty(shortcut.ToIDispatch(), "TargetPath").ToString()

	return targetPath, nil
}

func abrirRutaEjecutableDeAccesoDirecto(programa string) {
	// Comando para ejecutar el programa
	cmd := exec.Command(programa)

	// Iniciar el comando sin esperar a que termine
	err := cmd.Start()

	// Manejar errores
	if err != nil {
		log.Println("Error al abrir el programa:", err)
		return
	}
}

func main() {
	// Crear archivo de log
	logFile, err := os.OpenFile("start_app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error al crear archivo de log:", err)
		return
	}
	defer logFile.Close()

	// Configurar logger para escribir en el archivo de log
	log.SetOutput(logFile)

	// Definir el flag para el nombre del archivo
	nombreArchivo := flag.String("name", "", "Nombre del archivo a buscar")
	flag.Parse()

	if *nombreArchivo == "" {
		log.Println("Por favor, proporciona un nombre de archivo con el parámetro -name")
		fmt.Println("Por favor, proporciona un nombre de archivo con el parámetro -name")
		return
	}

	// Directorio donde deseas buscar
	directorio := "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs"

	// Buscar archivos
	resultados, err := buscarArchivos(*nombreArchivo, directorio)
	if err != nil {
		log.Println("Error al buscar archivos:", err)
		fmt.Println("Error al buscar archivos:", err)
		return
	}

	// Imprimir resultados
	for i, resultado := range resultados {
		log.Printf("Coincidencia %d: %s\n", i+1, resultado)
		fmt.Printf("Coincidencia %d: %s\n", i+1, resultado)
		rutaEjecutable, err := obtenerRutaEjecutableDeAccesoDirecto(resultado)
		if err != nil {
			log.Println("Error al obtener la ruta del ejecutable:", err)
			fmt.Println("Error al obtener la ruta del ejecutable:", err)
			return
		}
		abrirRutaEjecutableDeAccesoDirecto(rutaEjecutable)
	}
}
