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
		{"1) Аналог JS", false},
		{"2) Библиотека JS", true},
		{"3) Название раскладки клавиатуры", false},
		{"4) Библиотека Java", false},
	}},
	2: {"Как долго можно разворачивать проект?", []answer{
		{"1) 15 минут", false},
		{"2) Максимум 3 часа", false},
		{"3) Вечность", false},
		{"4) Не надо разворачивать", true},
	}},
	3: {"Что такое Jira?", []answer{
		{"1) Язык программирования", false},
		{"2) Фрейморк JS", false},
		{"3) Система управления проектами", true},
		{"4) Библиотека Java", false},
	}},
	4: {"Что такое JS?", []answer{
		{"1) Java, написанная скриптом", false},
		{"2) Скрипт от Java-разработчиков", false},
		{"3) Мультипарадигменный язык", true},
		{"4) Язык разметки Web-страниц", false},
	}},
	5: {"Каким должен быть дизайн", []answer{
		{"1) Адаптивным", false},
		{"2) Удобным", false},
		{"3) Привлекательным", false},
		{"4) Все вместе", true},
	}},
	6: {"Самый лучший язык", []answer{
		{"1) HTML", true},
		{"2) CSS", false},
		{"3) TXT", false},
		{"4) PDF", false},
	}},
	7: {"ВК тупит потому что: ", []answer{
		{"1) В ВК очень нехорошие люди", false},
		{"2) У них упал сервер", false},
		{"3) Telegram не тупит", false},
		{"4) ВК никогда не тупит", true},
	}},
	8: {"Что такое Гугл-таблица?", []answer{
		{"1) Блокнот", false},
		{"2) Способ совместного написания ТЗ", false},
		{"3) Таблица, которую сделал Google", true},
		{"4) Крутое место для контент-плана", false},
	}},
	9: {"Почему разработчики не вписываются в дэдлайны?", []answer{
		{"1) Они очень устали", false},
		{"2) Это риторический вопрос", true},
		{"3) Не завершили прошлый проект", false},
		{"4) Нет нормального ТЗ", false},
	}},
	10: {"Как вычислить человека в Интернете?", []answer{
		{"1) По MAC-адресу", true},
		{"2) По индексу", false},
		{"3) По IP", false},
		{"4) По почте", false},
	}},
	11: {"Почему упал сервер?", []answer{
		{"1) Его уронил кот", false},
		{"2) Баланс отрицательный", false},
		{"3) Из-за перегруза", true},
		{"4) Он уже лежал", false},
	}},
	12: {"Дедос это?", []answer{
		{"1) Атака", true},
		{"2) DEAD OS", false},
		{"3) Ласковое прозвище для дедушки", false},
		{"4) Почти лэндос", false},
	}},
	13: {"Что делает дизайнер интерфейсов?", []answer{
		{"1) Создает логотипы", false},
		{"2) Верстает сайты", false},
		{"3) Все, за что заплатят", false},
		{"4) Создает дизайн сайтов", true},
	}},
	14: {"Что такое Hackathon?", []answer{
		{"1) IDE", false},
		{"2) Состязание", true},
		{"3) Технология организации", false},
		{"4) Вписка разработчиков", false},
	}},
	15: {"Самая лучшая школа в мире?", []answer{
		{"1) Школа жизни", false},
		{"2) Твоя школа", false},
		{"3) Bonch.dev", true},
		{"4) Все школы отстой", false},
	}},
}
