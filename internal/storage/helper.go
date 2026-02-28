package storage

import (
	"bytes"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"

	"github.com/adrium/goheif"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
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
		// Проверяем сигнатуру HEIC (ftyp)
		if len(buff) > 12 && string(buff[4:8]) == "ftyp" {
			// --- ОБРАБОТКА HEIC С УЧЕТОМ EXIF ---

			// 1. Сначала пытаемся достать EXIF (данные о повороте)
			// goheif.ExtractExif читает файл, поэтому курсор сдвинется
			exifData, err := goheif.ExtractExif(file)
			
			// Если ошибка при извлечении EXIF, мы её игнорируем, 
			// так как нам все равно нужно декодировать саму картинку.
			// Но сбрасываем курсор обратно в начало файла для декодера.
			if _, err := file.Seek(0, 0); err != nil {
				return nil, fmt.Errorf("failed to seek to start after exif extract: %w", err)
			}

			// 2. Декодируем само изображение
			img, err := goheif.Decode(file)
			if err != nil {
				return nil, fmt.Errorf("failed to decode heic: %w", err)
			}

			// 3. Если EXIF был найден, применяем поворот
			if exifData != nil {
				img = applyOrientation(img, exifData)
			}

			return img, nil
		}

		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}
}

// Вспомогательная функция для применения поворота на основе EXIF данных
func applyOrientation(img image.Image, exifData []byte) image.Image {
	r := bytes.NewReader(exifData)
	x, err := exif.Decode(r)
	if err != nil {
		return img // Если не удалось распарсить EXIF, возвращаем как есть
	}

	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return img
	}

	orient, err := tag.Int(0)
	if err != nil {
		return img
	}

	// Значения ориентации EXIF и необходимые действия:
	// 1: Top-Left (Нормально)
	// 3: Bottom-Right (Перевернуто на 180)
	// 6: Right-Top (Повернуто 90 CW -> нужно повернуть 270 CCW или 90 CW)
	// 8: Left-Bottom (Повернуто 270 CW -> нужно повернуть 90 CCW)
	
	switch orient {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Rotate270(imaging.FlipH(img))
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Rotate90(imaging.FlipH(img))
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
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
