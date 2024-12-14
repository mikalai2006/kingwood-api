package utils

// type LocalesStruct struct {
// 	Ru LocaleItem
// 	En LocaleItem
// 	Fr LocaleItem
// 	Pl LocaleItem
// }

type LocaleItem struct {
	Code        string
	Name        string
	Description string
	Country     string
	Created     string
	Authors     string
	Coordinates string
	Updated     string
}

var (
	Locales = []LocaleItem{
		{
			Code:        "en",
			Name:        "Name",
			Description: "Description",
			Country:     "Country",
			Created:     "Created",
			Updated:     "Updated",
			Authors:     "Authors",
			Coordinates: "Coordinates",
		},
		{
			Code:        "ru",
			Name:        "Название",
			Description: "Описание",
			Country:     "Страна",
			Created:     "Добавлено",
			Updated:     "Обновлено",
			Authors:     "Авторы",
			Coordinates: "Координаты",
		},
		{
			Code:        "pl",
			Name:        "Nazwa",
			Description: "Opis",
			Country:     "Kraj",
			Created:     "Utworzony",
			Updated:     "Zaktualizowano",
			Authors:     "Autorski",
			Coordinates: "Współrzędne",
		},
		{
			Code:        "fr",
			Name:        "Nom",
			Description: "Description",
			Country:     "État",
			Created:     "Créé",
			Updated:     "Mis à jour",
			Authors:     "Auteurs",
			Coordinates: "Coordonnés",
		},
		{
			Code:        "ua",
			Name:        "Ім'я",
			Description: "Опис",
			Country:     "Країна",
			Created:     "Створено",
			Updated:     "Оновлено",
			Authors:     "Автори",
			Coordinates: "Координати",
		},
	}
)
