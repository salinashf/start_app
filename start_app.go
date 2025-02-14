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
	"github.com/joho/godotenv"
)

func buscarArchivos(nombreArchivo string, directorio string, is_exactly bool) ([]string, error) {
	var coincidencias []string

	err := filepath.Walk(directorio, func(ruta string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Comparar los nombres de archivo en minúsculas
		if compararCadenas(nombreArchivo, info.Name(), is_exactly) {
			log.Println("###--Econtrado--###")
			coincidencias = append(coincidencias, ruta)
		}
		if len(coincidencias) >= 1 {
			// Detener la búsqueda después de encontrar las primeras dos coincidencias
			return fmt.Errorf("encontradas dos coincidencias")
		}
		return nil
	})
	log.Println("Coincidencias:")
	log.Println(coincidencias)

	// Ignorar el error de "encontradas dos coincidencias"
	if err != nil && err.Error() != "encontradas dos coincidencias" {
		return nil, err
	}

	return coincidencias, nil
}
func compararCadenas(input string, path_finder string, is_exactly bool) bool {
	log.Println("Comparando:..", strings.ToLower(input), strings.ToLower(path_finder), is_exactly)
	if is_exactly {
		if strings.Compare(strings.ToLower(input)+".lnk", strings.ToLower(path_finder)) == 0 {
			log.Println("   Encontrado_A:..")
			return true
		}
	} else {
		if strings.Contains(strings.ToLower(path_finder), strings.ToLower(input)) {
			log.Println("   Encontrado_B:..")
			return true
		}
	}
	log.Println(" No Encontrado:..")
	return false
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

func abrirRutaEjecutableDeAccesoDirecto(programa string) error {
	// Comando para ejecutar el programa
	cmd := exec.Command(programa)

	// Iniciar el comando sin esperar a que termine
	err := cmd.Start()

	// Manejar errores
	if err != nil {
		log.Println("Error al abrir el programa:", err)
		return fmt.Errorf("error al abrir el programa: %v", err)
	}
	return nil
}

func busqueda(directorio string, nombreArchivo string, exactly bool) error {
	log.Println("###################################")
	log.Println("----Busca dentro de la carpeta----")
	log.Println("----Directorio:", directorio)
	log.Println("###################################")

	// Buscar archivos
	resultados, err := buscarArchivos(nombreArchivo, directorio, exactly)
	if err != nil {
		log.Println("Error al buscar archivos:", err)
		return fmt.Errorf("error al buscar archivos: %v", err)
	}
	if len(resultados) == 0 {
		log.Println("No se encontraron coincidencias")
		return fmt.Errorf("no se encontraron coincidencias: %v", err)
	}

	// Imprimir resultados
	for i, resultado := range resultados {
		log.Printf("Coincidencia %d: %s\n", i+1, resultado)
		rutaEjecutable, err := obtenerRutaEjecutableDeAccesoDirecto(resultado)
		if err != nil {
			log.Println("Error al obtener la ruta del ejecutable:", err)
			return fmt.Errorf("error al obtener la ruta del ejecutable: %v", err)

		}
		errx := abrirRutaEjecutableDeAccesoDirecto(rutaEjecutable)
		if errx != nil {
			log.Println("Error al abrir el archivo:", err)
			return fmt.Errorf("error al abrir el archivo: %v", err)
		}
	}
	return nil
}

func main() {

	// Crear archivo de log
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error al obtener la ruta del ejecutable:", err)
		return
	}
	exeDir := filepath.Dir(exePath)
	err_env := godotenv.Load(filepath.Join(exeDir, ".env"))
	if err_env != nil {
		log.Fatal("Error al cargar  el archivo  .env ")
	}

	log.Println("---------------------------------")
	log.Println("Inicio de la aplicación")
	log.Println("---------------------------------")

	logFilePath := filepath.Join(exeDir, "start_app.log")

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error al crear archivo de log:", err)
		return
	}
	defer logFile.Close()

	// Configurar logger para escribir en el archivo de log
	log.SetOutput(logFile)

	FOLDER_IN_PROGRAM_DATA := os.Getenv("FOLDER_IN_PROGRAM_DATA")
	log.Println(FOLDER_IN_PROGRAM_DATA)

	FOLDER_IN_APP_DATA := os.Getenv("FOLDER_IN_APP_DATA")
	log.Println(FOLDER_IN_APP_DATA)

	// Definir el flag para el nombre del archivo
	nombreArchivo := flag.String("name", "", "Nombre del archivo a buscar")
	exactly := flag.Bool("exactly", true, "Es busqueda exacta")
	flag.Parse()

	if *nombreArchivo == "" {
		log.Println("Por favor, proporciona un nombre de archivo con el parámetro -name")
		return
	}
	log.Println("*****************************************")
	log.Println("Parametros de entrada")
	log.Println("Nombre del archivo:", *nombreArchivo)
	log.Println("Busqueda exacta:", *exactly)
	log.Println("*****************************************")

	err_prg_data := busqueda(FOLDER_IN_PROGRAM_DATA, *nombreArchivo, *exactly)

	if err_prg_data != nil {
		log.Println("Error al buscar archivos:", err_prg_data)
		err_app_data := busqueda(FOLDER_IN_APP_DATA, *nombreArchivo, *exactly)
		if err_app_data != nil {
			log.Println("Error al buscar archivos:", err_app_data)
			return
		}
	}
	log.Println("---------------------------------")
	log.Println("Termina el cierre de la aplicación")
	log.Println("---------------------------------")
}
