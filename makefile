GO=go
GOFMT=gofmt
OUTPUT_FILE=./descargador-libros-by-parzibyte
.DEFAULT_GOAL=run

run:
	# Instalando dependencias
	${GO} mod tidy
	# Formateando
	${GOFMT} -w .
	# Compilando
	${GO} build -o ${OUTPUT_FILE}.exe -tags desarrollo
	# Ejecutando
	${OUTPUT_FILE}

prod:
	# Si vas a compilar para 32 bits, probablemente quieras:
	# SET PATH=C:\Go32\go\bin;C:\MinGW\bin;%PATH% && SET GOROOT=C:\Go32\go\
	# Instalando dependencias
	${GO} mod tidy
	# Formateando
	${GOFMT} -w .
	# Compilando
	${GO} build -o ${OUTPUT_FILE}_prod_64.exe -tags produccion -ldflags "-H windowsgui"

prod_32:
	# Si vas a compilar para 32 bits, probablemente quieras:
	# SET PATH=C:\Go32\go\bin;C:\MinGW\bin;%PATH% && SET GOROOT=C:\Go32\go\
	# Instalando dependencias
	${GO} mod tidy
	# Formateando
	${GOFMT} -w .
	# Compilando
	${GO} build -o ${OUTPUT_FILE}_prod_32.exe -tags produccion -ldflags "-H windowsgui"