package handlers

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"vk-bot/models"

	"github.com/SevereCloud/vksdk/5.92/object"
)

type textHandler struct {
	*DefHandler
}

func TextHandle(handler *DefHandler) object.MessageNewFunc {
	return (&textHandler{handler}).Handle
}

func (h *textHandler) Handle(obj object.MessageNewObject, groupID int) {
	userId := obj.FromID

	user, err := models.
		Users(models.UserWhere.VKID.EQ(null.IntFrom(userId))).
		One(context.Background(), h.db)
	if err != nil {
		sentry.CaptureException(errors.WithMessagef(err, "Error when Handling text: %v", err))
		h.logger.Errorf("Error when Handling text: %v", err)
	}

	if user == nil {
		newUser := models.User{
			VKID:       null.IntFrom(userId),
			State:      null.StringFrom(models.UsersStateTEST),
			TestPhase:  null.IntFrom(0),
			TestPoints: null.IntFrom(0),
		}
		if err := newUser.Insert(context.Background(), h.db, boil.Infer()); err != nil {
			h.logger.Errorf("Error when Handling text save user: %v", err)
			sentry.CaptureException(errors.WithMessagef(err, "Error when Handling text save user: %v", err))
		}
		user = &newUser
	}

	if user.State.String == models.UsersStateTEST {
		if user.TestPhase.Int == 0 {
			h.HandleTestStart(obj)
		} else {
			h.HandleTest(obj, groupID)
		}
	}
}

const (
	primaryColor   string = "primary"
	secondaryColor string = "secondary"
	negativeColor  string = "negative"
	positiveColor  string = "positive"
)

type question struct {
	Question string
	Answers  []answer
}

type answer struct {
	Text    string
	IsRight bool
}

var testMap = map[int]question{
	1: {"Твое прохождение теста от Bonch.dev начинается прямо сейчас! \n\n JQuery это:", []answer{
		{"Аналог JS", false},
		{"Библиотека JS", true},
		{"Название раскладки клавиатуры", false},
		{"Библиотека Java", false},
	}},
	2: {"Как долго можно разворачивать проект?", []answer{
		{"15 минут", false},
		{"Максимум 3 часа", false},
		{"Вечность", false},
		{"Не надо разворачивать", true},
	}},
	3: {"Что такое Jira?", []answer{
		{"Язык программирования", false},
		{"Фрейморк JS", false},
		{"Система управления проектами", true},
		{"Библиотека Java", false},
	}},
	4: {"Что такое JS?", []answer{
		{"Java, написанная скриптом", false},
		{"Скрипт от Java-разработчиков", false},
		{"Мультипарадигменный язык", true},
		{"Язык разметки Web-страниц", false},
	}},
	5: {"Каким должен быть дизайн", []answer{
		{"Адаптивным", false},
		{"Удобным", false},
		{"Привлекательным", false},
		{"Все вместе", true},
	}},
	6: {"Самый лучший язык", []answer{
		{"HTML", true},
		{"CSS", false},
		{"TXT", false},
		{"PDF", false},
	}},
	7: {"ВК тупит потому что: ", []answer{
		{"В ВК очень нехорошие люди", false},
		{"У них упал сервер", false},
		{"Telegram не тупит", false},
		{"ВК никогда не тупит", true},
	}},
	8: {"Что такое Гугл-таблица?", []answer{
		{"Блокнот", false},
		{"Способ совместного написания ТЗ", false},
		{"Таблица, которую сделал Google", true},
		{"Крутое место для контент-плана", false},
	}},
	9: {"Почему разработчики не вписываются в дэдлайны?", []answer{
		{"Они очень устали", false},
		{"Это риторический вопрос", true},
		{"Не завершили прошлый проект", false},
		{"Нет нормального ТЗ", false},
	}},
	10: {"Как вычислить человека в Интернете?", []answer{
		{"По MAC-адресу", true},
		{"По индексу", false},
		{"По IP", false},
		{"По почте", false},
	}},
	11: {"Почему упал сервер?", []answer{
		{"Его уронил кот", false},
		{"Баланс отрицательный", false},
		{"Из-за перегруза", true},
		{"Он уже лежал", false},
	}},
	12: {"Дедос это?", []answer{
		{"Атака", true},
		{"DEAD OS", false},
		{"Ласковое прозвище для дедушки", false},
		{"Почти лэндос", false},
	}},
	13: {"Что делает дизайнер интерфейсов?", []answer{
		{"Создает логотипы", false},
		{"Верстает сайты", false},
		{"Все, за что заплатят", false},
		{"Создает дизайн сайтов", true},
	}},
	14: {"Что такое Hackathon?", []answer{
		{"IDE", false},
		{"Состязание", true},
		{"Технология организации", false},
		{"Вписка разработчиков", false},
	}},
	15: {"Самая лучшая школа в мире?", []answer{
		{"Школа жизни", false},
		{"Твоя школа", false},
		{"Bonch.dev", true},
		{"Все школы отстой", false},
	}},
}
