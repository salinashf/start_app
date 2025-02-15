
# start_app

## ¿Qué es start_app?

Esta aplicación es un comando que permite abrir aplicaciones previamente instaladas mediante sus accesos directos. Busca el acceso directo en **C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs** y, una vez encontrado, abre la aplicación a través del archivo ejecutable. La creé con el objetivo de abrir cualquier aplicación desde un comando.

> Ejemplo: **start_app -name brave -exactly true**

Parámetros requeridos
| Parámetro | valor |
| ------ | ------ |
| -name  | Aplicación a abrir  |
| -exactly | Determina si la búsqueda es exacta, **true** o  **false**|


No quería llenar la variable `PATH` del sistema operativo con todas las aplicaciones instaladas. Con esta solución, solo agrego una sola variable dentro del `PATH`.

## Compilación

Compilé el programa con estos parámetros para que no muestre nada por consola:

```sh
go build -ldflags -H=windowsgui -o start_app.exe .
```
## Como instalar

Para instalar el comando se deber colocar `start_app.exe` en una carpeta determinada

> Ejemplo c:/utils
> 
```sh
c:/utils/start_app.exe
```


Esta ruta hay que agregar dentro de las variables de entorno de `PATH`

Terminado esto se puede ejecutar desde la consola el comando  `start_app -name brave -exactly true`

## Como Configurar

El comando, permite configurar la ruta donde se buscará las aplicaciones, 
esta configuración se da través de un archivo llamado `.env`  , existe dos rutas básicas donde se buscara la aplicación

| Clave | Valor  | 
| ------ | ------ |
|FOLDER_IN_PROGRAM_DATA| C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs|
|FOLDER_IN_APP_DATA| C:\\Users\\mi_usuario\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs|

Los dos archivos deben de estar en el mismo directorio  `.env`  y  `start_app.exe` 

>Ejemplo asumiendo que se coloque dentro de la carpeta   **C:\utils** , dentro de esta misma carpeta se creara el archivo de Log **start_app.log**
 
```md
C:\utils
├── start_app.exe
├── .env
└── start_app.log

```
## Control ejecutar  

Para abrir la aplicación de Discord,
>  **start_app  -name 'Discord' -exactly 'true'**
> 

## Control de errores  

Los errores que genera a la aplicación lo registran en un archivo de log, porque no muestra nada por consola  

## Gracias  ChatGPT
No se mucha programación en lenguaje **Goland** he leído lo básico, quien me ayudo fue Chatgpt, ¡¡Gracias ChatGpt!! 

La comencé escribiendo en PowerShell y Python, pero no me terminaron convenciendo, no quería que muestre la consola principal problema, este comando lo necesito para abrir aplicaciones desde otra aplicación




