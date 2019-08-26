package handlers

import (
	"database/sql"

	"github.com/SevereCloud/vksdk/5.92/api"
	"github.com/sirupsen/logrus"
)

type DefHandler struct {
	vkapi  api.VK
	db     *sql.DB
	logger *logrus.Logger
}

func NewDefHandler(vk api.VK, db *sql.DB, logger *logrus.Logger) *DefHandler {
	return &DefHandler{vkapi: vk, db: db, logger: logger}
}
