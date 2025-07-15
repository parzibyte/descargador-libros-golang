# descargador-libros-golang
 Descargador de libros de CONALITEG pero ahora con Golang

# Compilar

Para 64 bits:

```bash

go build
```

32 bits suponiendo que ya tienes Go de 32:

```bash
SET PATH=C:\Go32\go\bin;C:\MinGW\bin;%PATH% && SET GOROOT=C:\Go32\go\
go build -o descargador-libros-texto-32.exe
```