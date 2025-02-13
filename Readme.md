
# start_app

## ¿Qué es start_app?

Esta aplicación es un comando que permite abrir aplicaciones previamente instaladas mediante sus accesos directos. Busca el acceso directo en **C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs** y, una vez encontrado, abre la aplicación a través del archivo ejecutable. La creé con el objetivo de abrir cualquier aplicación desde un comando.

> Ejemplo: **start_app -name brave**

No quería llenar la variable `PATH` del sistema operativo con todas las aplicaciones instaladas. Con esta solución, solo agrego una sola variable dentro del `PATH`.

## Compilación

Compilé el programa con estos parámetros para que no muestre nada por consola:

```sh
go build -ldflags -H=windowsgui -o start_app.exe .
```

## Control de errores  

Los errores que genera a la aplicación lo registran en un archivo de log, porque no muestra nada por consola  

## Gracias  ChatGPT
No se mucha programación en lenguaje **Goland** he leído lo básico, quien me ayudo fue Chatgpt, ¡¡Gracias ChatGpt!! 

La comencé escribiendo en PowerShell y Python, pero no me terminaron convenciendo, no quería que muestre la consola principal problema, este comando lo necesito para abrir aplicaciones desde otra aplicación




