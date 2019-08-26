package handlers

import (
	"context"
	"math/rand"

	"github.com/SevereCloud/vksdk/5.92/object"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"vk-bot/models"
	"vk-bot/vk_objects"
)

func (h *textHandler) HandleTestStart(obj object.MessageNewObject) {
	userId := obj.FromID

	u, err := models.
		Users(models.UserWhere.VKID.EQ(null.IntFrom(userId))).
		One(context.Background(), h.db)
	if err != nil {
		h.logger.Errorf("Error when Handling text: %v", err)
	}
	if u == nil {
		return
	}

	h.testUser(obj, u)
}

func (h *textHandler) HandleTest(obj object.MessageNewObject, groupID int) {
	userId := obj.FromID

	u, err := models.
		Users(models.UserWhere.VKID.EQ(null.IntFrom(userId))).
		One(context.Background(), h.db)
	if err != nil {
		sentry.CaptureException(errors.WithMessagef(err, "Error when Handling text: %v", err))
		h.logger.Errorf("Error when Handling text: %v", err)
	}

	vMap := testMap[u.TestPhase.Int].Answers

	for _, val := range vMap {
		if val.Text == obj.Text {
			h.processTest(obj, u, val.IsRight)
			return
		}
	}

	var text string
	if obj.Text == "Повторить вопрос" {
		text = testMap[u.TestPhase.Int].Question
	} else {
		text = "Используй клавиатуру снизу для ответов"
	}

	message := vk_objects.NewMessage()
	message.UserID = &(obj.FromID)
	message.RandomId = rand.Int31()
	message.Message = text

	for i, val := range testMap[u.TestPhase.Int].Answers {
		if i%2 == 0 {
			message.Keyboard.AddRow()
		}
		message.Keyboard.AddTextButton(val.Text, "", secondaryColor)
	}
	message.Keyboard.AddRow()
	message.Keyboard.AddTextButton("Повторить вопрос", "", primaryColor)

	messageMap := message.ToMap()

	_, vkErr := h.vkapi.MessagesSend(messageMap)
	if vkErr.Code != 0 {
		h.logger.Error(vkErr)
		sentry.CaptureException(errors.Errorf("Error in sending message VK: %v", vkErr))
	}
}

func (h *textHandler) processTest(obj object.MessageNewObject, u *models.User, isRight bool) {
	if isRight {
		u.TestPoints = null.IntFrom(u.TestPoints.Int + 1)
	}

	if u.TestPhase.Int < 15 {
		h.testUser(obj, u)
		return
	}

	h.testResult(obj, u)
}

func (h *textHandler) testUser(obj object.MessageNewObject, u *models.User) {
	u.TestPhase = null.IntFrom(u.TestPhase.Int + 1)

	message := vk_objects.NewMessage()
	message.UserID = &(obj.FromID)
	message.RandomId = rand.Int31()
	message.Message = testMap[u.TestPhase.Int].Question

	for i, val := range testMap[u.TestPhase.Int].Answers {
		if i%2 == 0 {
			message.Keyboard.AddRow()
		}
		message.Keyboard.AddTextButton(val.Text, "", secondaryColor)
	}
	message.Keyboard.AddRow()
	message.Keyboard.AddTextButton("Повторить вопрос", "", primaryColor)

	messageMap := message.ToMap()

	_, err := h.vkapi.MessagesSend(messageMap)
	if err.Code != 0 {
		h.logger.Error(err)
		sentry.CaptureException(errors.Errorf("Error in sending message VK: %v", err))
	}

	if _, err := u.Update(context.Background(), h.db, boil.Infer()); err != nil {
		h.logger.Error("Error when updating user: %v", err)
		sentry.CaptureException(errors.Errorf("Error when updating user: %v", err))
	}
}

func (h *textHandler) testResult(obj object.MessageNewObject, u *models.User) {
	message := vk_objects.NewMessage()
	message.UserID = &(obj.FromID)
	message.RandomId = rand.Int31()
	message.Message = getText(u) + "\n\nПодойди к столу с плюшками, чтобы забрать заслуженный приз!"
	message.Keyboard.AddRow()
	message.Keyboard.AddTextButton("Я иду к организаторам за призом!", "", positiveColor)

	messageMap := message.ToMap()

	u.State = null.StringFrom(models.UsersStateMAIN)

	_, err := h.vkapi.MessagesSend(messageMap)
	if err.Code != 0 {
		h.logger.Error(err)
		sentry.CaptureException(errors.Errorf("Error in sending message VK: %v", err))
	}

	if _, err := u.Update(context.Background(), h.db, boil.Infer()); err != nil {
		h.logger.Error("Error when updating user: %v", err)
		sentry.CaptureException(errors.Errorf("Error when updating user: %v", err))
	}
}

func getText(u *models.User) string {
	points := u.TestPoints.Int
	if points >= 10 && points < 12 {
		return "Ты крут, ты правильно ответил 70% вопросов нашего теста!"
	} else if points >= 12 && points < 14 {
		return "Ты супер крут, ты правильно ответил на 80% вопросов нашего теста!"
	} else if points == 14 {
		return "Ты мега крут, ты правильно ответил на 90% нашего теста!"
	} else if points == 15 {
		return "Ты просто Breathtaking! Ты смог правильно ответить на все наши вопросы!"
	}
	return "Ты неплох, но ты ответил меньше чем на 70% нашего теста"
}
