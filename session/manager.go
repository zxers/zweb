package session

import (
	"zweb"
	"github.com/google/uuid"
)

type Manager struct {
	Store
	Propagator
	SessCtxKey string
}

func (m *Manager) GetSession(ctx *zweb.Context) (Session, error) {
	sessionId, err := m.Propagator.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	sess, err := m.Store.Get(ctx.Req.Context(), sessionId)
	if err != nil {
		return nil, err
	}
	return sess, err
}

func (m *Manager) InitSession(ctx *zweb.Context) (Session, error) {
	id := uuid.New()
	err := m.Propagator.Inject(id.String(), ctx.Resp)
	if err != nil {
		return nil, err
	}
	return m.Store.Generate(ctx.Req.Context(), id.String())
}

func (m *Manager) RefreshSession(ctx *zweb.Context) (error) {
	sessionId, err := m.Propagator.Extract(ctx.Req)
	if err != nil {
		return err
	}
	err = m.Store.Refresh(ctx.Req.Context(), sessionId)
	if err != nil {
		return err
	}
	return m.Inject(sessionId, ctx.Resp)
}

func (m *Manager) RemoveSession(ctx *zweb.Context) (error) {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), session.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}