package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	_ "image/jpeg"
	_ "image/png"

	"github.com/signintech/gopdf"
)

func obtenerCodigoFuenteDeVisualizadorDeLibro(urlLibro string) (string, error) {
	respuesta, err := http.Get(urlLibro)
	if err != nil {
		return "", err
	}
	if respuesta.StatusCode != 200 {
		return "", fmt.Errorf("obteniendo información de libro: código de respuesta %d", respuesta.StatusCode)
	}
	defer respuesta.Body.Close()
	cuerpoRespuesta, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return "", err
	}
	respuestaComoCadena := string(cuerpoRespuesta)
	return respuestaComoCadena, nil
}

func extraerAñoDeLibroSegunUrl(urlLibro string) (string, error) {
	expresionRegularPaginas := regexp.MustCompile(`(?m)https://libros\.conaliteg\.gob\.mx/(\d+)/\w+\.html?`)
	coincidenciasPagina := expresionRegularPaginas.FindStringSubmatch(urlLibro)
	if len(coincidenciasPagina) < 2 {
		return "", fmt.Errorf("no se pudo extraer el año del libro. ¿Es una URL correcta? debe tener la forma https://libros.conaliteg.gob.mx/AÑO/CLAVE.htm. Coincidencias: %v", coincidenciasPagina)
	}
	return coincidenciasPagina[1], nil
}

func extraerClaveAñoYPaginas(urlLibro string) (claveLibro string, año string, cantidadDePaginas int, err error) {
	codigoFuente, err := obtenerCodigoFuenteDeVisualizadorDeLibro(urlLibro)
	if err != nil {
		return
	}
	expresionRegularPaginas := regexp.MustCompile(`(?m)ag_pages = (\d*);`)
	coincidenciasPagina := expresionRegularPaginas.FindStringSubmatch(codigoFuente)
	if len(coincidenciasPagina) < 2 {
		err = fmt.Errorf("no se pudo extraer la cantidad de páginas. Coincidencias: %v", coincidenciasPagina)
		return
	}
	cantidadDePaginasComoCadena := coincidenciasPagina[1]
	cantidadDePaginas, err = strconv.Atoi(cantidadDePaginasComoCadena)
	if err != nil {
		return
	}
	expresionRegularClave := regexp.MustCompile(`(?m)ag_clave = "(\w+)";`)
	coincidenciasClave := expresionRegularClave.FindStringSubmatch(codigoFuente)
	if len(coincidenciasClave) < 2 {
		err = fmt.Errorf("no se pudo extraer la clave del libro. Coincidencias: %v", coincidenciasPagina)
		return
	}
	claveLibro = coincidenciasClave[1]
	año, err = extraerAñoDeLibroSegunUrl(urlLibro)
	return
}

func descargarImagenDeInternetYDevolverReader(url string) (io.Reader, error) {
	respuesta, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if respuesta.StatusCode != 200 {
		return nil, fmt.Errorf("error descargando imagen, el código de respuesta fue %d", respuesta.StatusCode)
	}
	return respuesta.Body, nil
}

func descargarLibro(urlLibro string) error {
	claveLibro, año, cantidadPaginas, err := extraerClaveAñoYPaginas(urlLibro)
	if err != nil {
		return err
	}
	fmt.Printf(`Información del libro %s extraída.
Clave: %s
Año: %s
Cantidad de páginas: %d
`,
		urlLibro, claveLibro, año, cantidadPaginas)
	pdf := gopdf.GoPdf{}
	tamañoDePagina := *gopdf.PageSizeA4
	pdf.Start(gopdf.Config{PageSize: tamañoDePagina})
	var url string
	/*
		El servidor asume que los libros siempre tendrán a lo mucho 999 páginas, así
		que siempre pide 3 dígitos. Si no fuera así y los dígitos dependieran del número
		de página, se podría hacer así:
		cantidadDeDigitosQueConformanLaCantidadDePaginas := len(strconv.Itoa(cantidadPaginas))

		Mientras tanto se ha dejado en 3
	*/
	cantidadDeDigitosQueConformanLaCantidadDePaginas := 3
	for numeroDePaginaActual := 0; numeroDePaginaActual < cantidadPaginas; numeroDePaginaActual++ {
		/*
			La URL de un libro se ve así:
			https://libros.conaliteg.gob.mx/2023/P1LPM.htm

			La de una imagen de ese libro, así:
			https://libros.conaliteg.gob.mx/2023/c/P1LPM/000.jpg
		*/
		numeroDePaginaConCeros := fmt.Sprintf("%0*d", cantidadDeDigitosQueConformanLaCantidadDePaginas, numeroDePaginaActual)
		url = fmt.Sprintf(
			"https://libros.conaliteg.gob.mx/%s/c/%s/%s.jpg",
			año, claveLibro, numeroDePaginaConCeros,
		)
		fmt.Printf("Descargando imagen %d/%d que se encuentra en %s...", numeroDePaginaActual+1, cantidadPaginas, url)
		lectorRespuestaHttp, err := descargarImagenDeInternetYDevolverReader(url)
		fmt.Printf("OK\nAgregando imagen descargada a PDF...")
		if err != nil {
			return fmt.Errorf("al descargar imagen: %s", err.Error())
		}
		pdf.AddPage()
		imageHolder, err := gopdf.ImageHolderByReader(lectorRespuestaHttp)
		if err != nil {
			return err
		}
		err = pdf.ImageByHolder(imageHolder, 0, 0, &gopdf.Rect{
			W: tamañoDePagina.W,
			H: tamañoDePagina.H,
		})
		if err != nil {
			return fmt.Errorf("al agregar la imagen al PDF: %s", err.Error())
		}

		fmt.Printf("OK\n")
	}
	nombrePdf := claveLibro + ".pdf"
	fmt.Printf("Guardando PDF con el nombre %s...", nombrePdf)
	return pdf.WritePdf(nombrePdf)
}

func main() {
	fmt.Printf(`Ingresa la URL del libro. Puede tener la forma 1 o 2. Después de haberlo ingresado, presiona Enter.\n
URL del libro: `)
	var urlLibro string
	fmt.Scanln(&urlLibro)
	err := descargarLibro(urlLibro)
	if err != nil {
		fmt.Printf("Error descargando libro. El error es: %v. Recuerde: la URL debe tener el formato https://libros.conaliteg.gob.mx/AÑO/CLAVE.htm por ejemplo https://libros.conaliteg.gob.mx/2023/P1LPM.htm",
			err)
	} else {
		fmt.Printf("\n\nLibro descargado correctamente en el mismo lugar donde se encuentra este programa")
	}
}
