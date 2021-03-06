// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/miniflux/miniflux/http/handler"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
)

// ShowUnreadPage render the page with all unread entries.
func (c *Controller) ShowUnreadPage(ctx *handler.Context, request *handler.Request, response *handler.Response) {
	user := ctx.LoggedUser()
	offset := request.QueryIntegerParam("offset", 0)

	builder := c.store.NewEntryQueryBuilder(user.ID)
	builder.WithStatus(model.EntryStatusUnread)
	countUnread, err := builder.CountEntries()
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	if offset >= countUnread {
		offset = 0
	}

	builder = c.store.NewEntryQueryBuilder(user.ID)
	builder.WithStatus(model.EntryStatusUnread)
	builder.WithOrder(model.DefaultSortingOrder)
	builder.WithDirection(user.EntryDirection)
	builder.WithOffset(offset)
	builder.WithLimit(nbItemsPerPage)
	entries, err := builder.GetEntries()
	if err != nil {
		response.HTML().ServerError(err)
		return
	}

	response.HTML().Render("unread", tplParams{
		"user":        user,
		"countUnread": countUnread,
		"entries":     entries,
		"pagination":  c.getPagination(ctx.Route("unread"), countUnread, offset),
		"menu":        "unread",
		"csrf":        ctx.CSRF(),
	})
}

// MarkAllAsRead marks all unread entries as read.
func (c *Controller) MarkAllAsRead(ctx *handler.Context, request *handler.Request, response *handler.Response) {
	if err := c.store.MarkAllAsRead(ctx.UserID()); err != nil {
		logger.Error("[MarkAllAsRead] %v", err)
	}

	response.Redirect(ctx.Route("unread"))
}
