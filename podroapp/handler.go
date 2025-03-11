package podroapp

import (
	"interview/pkg"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) GetProivedersWeaklyReport(ctx *fiber.Ctx) error {
	resp, err := h.svc.GetProvidersWeaklyReport(ctx.Context(), GetProvidersWeaklyReportReqeust{})
	if err != nil {
		pkg.Logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "somthing went wrotng",
		})
	}
	return ctx.JSON(resp)
}
