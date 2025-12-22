package types

import "github.com/mikestefanello/pagoda/pkg/controller"

type (
	LoginForm struct {
		Username   string `form:"username" validate:"required"`
		Password   string `form:"password" validate:"required"`
		Submission controller.FormSubmission
	}
)
