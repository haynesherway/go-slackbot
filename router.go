package slackbot

import (
	"context"
)

type Router interface {
	Match(context.Context, *RouteMatch) (bool, context.Context)
	SetBotID(botID string)
	Hear(regex string) *Route
	Messages(types ...MessageType) *Route
	AddMatcher(m Matcher) *Route
	TalkToSelf() *Route
	NoTalkToSelf() *Route
	AlwaysTalkToSelf() Router
	NeverTalkToSelf() Router
	Handler(handler Handler) error
	MessageHandler(handler MessageHandler) error
	Err() error
}

type SimpleRouter struct {
	// Routes to be matched, in order.
	routes []*Route
	// Slack UserID of the bot UserID
	botUserID string
	// Holds any error that occurred during the function chain that created this router.  Used by subrouters.
	err error
	// if set, all new routes will be set to allow self-talking
	talkToSelf bool
}

// Match matches registered routes against the request.
func (r *SimpleRouter) Match(ctx context.Context, match *RouteMatch) (bool, context.Context) {
	if r.err != nil {
		return false, ctx
	}

	for _, route := range r.routes {
		if matched, ctx := route.Match(ctx, match); matched {
			return true, ctx
		}
	}

	return false, ctx
}

// NewRoute registers an empty route.
func (r *SimpleRouter) newRoute(selftalk bool) *Route {
	route := &Route{
		err:        r.err,
		talkToSelf: selftalk,
	}
	r.routes = append(r.routes, route)
	return route
}

func (r *SimpleRouter) Hear(regex string) *Route {
	return r.newRoute(r.talkToSelf).Hear(regex)
}

func (r *SimpleRouter) Handler(handler Handler) error {
	return r.newRoute(r.talkToSelf).Handler(handler)
}

func (r *SimpleRouter) MessageHandler(handler MessageHandler) error {
	return r.newRoute(r.talkToSelf).MessageHandler(handler)
}

func (r *SimpleRouter) Messages(types ...MessageType) *Route {
	return r.newRoute(r.talkToSelf).Messages(types...)
}

func (r *SimpleRouter) AddMatcher(m Matcher) *Route {
	return r.newRoute(r.talkToSelf).AddMatcher(m)
}

func (r *SimpleRouter) SetBotID(botID string) {
	r.botUserID = botID
	for _, route := range r.routes {
		route.setBotID(botID)
	}
}

func (r *SimpleRouter) AlwaysTalkToSelf() Router {
	r.talkToSelf = true
	return r
}

func (r *SimpleRouter) NeverTalkToSelf() Router {
	r.talkToSelf = false
	return r
}

func (r *SimpleRouter) TalkToSelf() *Route {
	return r.newRoute(true)
}

func (r *SimpleRouter) NoTalkToSelf() *Route {
	return r.newRoute(false)
}

func (r *SimpleRouter) Err() error {
	return r.err
}
