package storage

import (
	"bytes"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"

	"github.com/adrium/goheif"
	"github.com/disintegration/imaging"
	"golang.org/x/image/webp"
)

// определяет формат по содержимому и декодирует картинку
func decodeImage(file multipart.File) (image.Image, error) {
	// Считываем первые 512 байт для определения типа
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return nil, fmt.Errorf("failed to read file header: %w", err)
	}

	// Возвращаем "курсор" файла в самое начало,
	// иначе декодер начнет читать с 513-го байта и сломается.
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek to start: %w", err)
	}

	contentType := http.DetectContentType(buff)

	switch contentType {
	case "image/jpeg", "image/png":
		// imaging.Decode автоматически обрабатывает EXIF-поворот для Jpeg/Png
		img, err := imaging.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode standard image: %w", err)
		}
		return img, nil

	case "image/webp":
		img, err := webp.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode webp: %w", err)
		}
		return img, nil

	default:
		// Go стандартно определяет HEIC как "application/octet-stream",
		// поэтому проверяем сигнатуру "ftyp" вручную.
		// HEIF/HEIC файлы обычно содержат 'ftyp' на позициях 4-8.
		if len(buff) > 12 && string(buff[4:8]) == "ftyp" {
			// Пробуем декодировать как HEIC
			img, err := goheif.Decode(file)
			if err != nil {
				return nil, fmt.Errorf("failed to decode heic: %w", err)
			}
			return img, nil
		}

		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}
}

func encodeToJPG(img image.Image, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)
	// imaging.Encode пишет результат в буфер (в память)
	err := imaging.Encode(buf, img, imaging.JPEG, imaging.JPEGQuality(quality))
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}
