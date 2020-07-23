package tlog

import (
	"context"
	"github.com/go-pg/pg/v9"
	"testing"
)

type ShowQuery struct {
	t *testing.T
}

func NewShowQuery(t *testing.T) *ShowQuery {
	return &ShowQuery{t: t}
}

func(s *ShowQuery) BeforeQuery(c context.Context,q *pg.QueryEvent) (context.Context, error){
	s.t.Log(q.FormattedQuery())

	return c, nil
}

func(s *ShowQuery) AfterQuery(c context.Context, q *pg.QueryEvent) error{
	return nil
}
