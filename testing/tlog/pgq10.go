package tlog

import (
	"context"
	"github.com/go-pg/pg/v10"
	"testing"
)

type ShowQuery10 struct {
	t *testing.T
}

func NewShowQuery10(t *testing.T) *ShowQuery10 {
	return &ShowQuery10{t: t}
}

func (s *ShowQuery10) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	b, err := q.FormattedQuery()
	s.t.Log(string(b), err)

	return c, nil
}

func (s *ShowQuery10) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	return nil
}
