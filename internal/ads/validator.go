package ads

import (
	"context"
	"mime/multipart"
	"strconv"
	appErrors "vietio/internal/errors"

	"github.com/google/uuid"
)

type CategoryChecker interface {
	Exists(context.Context, int) (bool, error)
}

type AdChecker interface {
	Exists(context.Context, uuid.UUID) (bool, error)
}

type Validator struct {
	categoryChecker CategoryChecker
	adChecker       AdChecker
}

func NewValidator(categoryChecker CategoryChecker, adChecker AdChecker) *Validator {
	return &Validator{
		categoryChecker: categoryChecker,
		adChecker:       adChecker,
	}
}

func (v *Validator) createAdValidate(
	ctx context.Context,
	payload CreateAdRequestBody,
	images []*multipart.FileHeader,
) *appErrors.ValidationError {
	errors := appErrors.NewValidationError()

    v.validateCommonFields(
        ctx,
        errors,
        payload.Title,
        payload.Description,
        payload.Price,
        payload.CategoryId,
    )

	if len(images) == 0 || len(images) > 3 {
		errors.Add("images", "images должен быть > 0 и меньше 3")
	}

	return errors
}

func (v *Validator) updateAdValidate(
	ctx context.Context,
	payload UpdateAdRequestBody,
	images []*multipart.FileHeader,
) *appErrors.ValidationError {
	errors := appErrors.NewValidationError()

    v.validateCommonFields(
        ctx,
        errors,
        payload.Title,
        payload.Description,
        payload.Price,
        payload.CategoryId,
    )

	adExists, err := v.adChecker.Exists(ctx, payload.Uuid)
	if err != nil {
		errors.Add("uuid", "ошибка БД при проверки существования объявления")
	}
	if !adExists {
		errors.Add("uuid", "такого объявления не существует не существует")
	}

	// общее количество картинок
	countImages := len(images) + len(payload.OldImages)
	if countImages == 0 || countImages > 3 {
		errors.Add("images", "общее количество изображений должно быть > 0 и <= 3")
	}

	return errors
}

func (v *Validator) validateCommonFields(
    ctx context.Context, 
    errors *appErrors.ValidationError,
    title string, 
    description string,
    price int,
    categoryId int,
) {
    if title == "" {
		errors.Add("title", "title не может быть пустым")
	}
	if description == "" {
		errors.Add("description", "description не может быть пустым")
	}
	if price < 0 {
		errors.Add("price", "price не может быть отрицательным")
	}
	if categoryId < 1 {
		errors.Add("category_id", "category_id должен быть >= 1")
	}

	if categoryId >= 1 {
		categoryExists, err := v.categoryChecker.Exists(ctx, categoryId)
		if err != nil {
			errors.Add("category_id", "ошибка БД при проверки существования категории")
		}
		if !categoryExists {
			errors.Add("category_id", "category_id такой категории не существует")
		}
	}
}

func validateIntField(
	fieldName string,
	fieldValue string, 
	isRequired bool, 
	defaultValue int, 
	errors *appErrors.ValidationError,
) int {
	if isRequired && fieldValue == "" {
		errors.Add(fieldName, "поле обязательно для заполнения")
		return 0
	}

	if !isRequired && fieldValue == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(fieldValue);
	if err != nil {
        errors.Add(fieldName, "в поле должны быть число")
		return 0
    }

	return result
}
