package entOrm

import (
	"fmt"
	"slices"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/samber/lo"
	"github.com/tuan-dd/go-pkg/common/pagination"
	"github.com/tuan-dd/go-pkg/common/utils"
)

type Op int32

const (
	OP_LT          Op = 0
	OP_LTE         Op = 1
	OP_GT          Op = 2
	OP_GTE         Op = 3
	OP_EQ          Op = 4
	OP_NEQ         Op = 5
	OP_LIKE        Op = 6
	OP_IS_NULL     Op = 7
	OP_IS_NOT_NULL Op = 8
	OP_IN          Op = 9
	OP_NOT_IN      Op = 10
)

func ParseSort(columns []string, table string, sorts []*pagination.SortMethod) func(s *sql.Selector) {
	var values []string

	for _, sort := range sorts {
		key := utils.ToSnakeCase(sort.Key)
		if !slices.Contains(columns, key) {
			return nil
		}

		col := fmt.Sprintf("`%s`.`%s`", table, key)
		switch sort.Value {
		case pagination.SortType_ASC:
			values = append(values, sql.Asc(col))
		case pagination.SortType_DESC:
			values = append(values, sql.Desc(col))
		default:
			return nil
		}
	}

	return func(s *sql.Selector) {
		if len(values) > 0 {
			s.OrderBy(values...)
		}
	}
}

func ApplyFilter(filters []*pagination.Filter, selector *sql.Selector, op Op, columns ...string) {
	if len(filters) == 0 {
		return
	}
	if len(columns) == 0 {
		columns = selector.Columns()
	}

	for _, f := range filters {
		key := utils.ToSnakeCase(f.Key)
		if !slices.Contains(columns, key) {
			continue
		}
		switch op {
		case OP_LT:
			selector.Where(sql.LT(key, f.Value))
		case OP_LTE:
			selector.Where(sql.LTE(key, f.Value))
		case OP_GT:
			selector.Where(sql.GT(key, f.Value))
		case OP_GTE:
			selector.Where(sql.GTE(key, f.Value))
		case OP_EQ:
			selector.Where(sql.EQ(key, f.Value))
		case OP_NEQ:
			selector.Where(sql.NEQ(key, f.Value))
		case OP_LIKE:
			selector.Where(sql.ContainsFold(key, f.Value))
		case OP_IS_NULL:
			selector.Where(sql.IsNull(key))
		case OP_IS_NOT_NULL:
			selector.Where(sql.NotNull(key))
		case OP_IN:
			parts := ParseInFilterValues(f)
			if len(parts) == 1 {
				selector.Where(sql.EQ(key, parts[0]))
			} else if len(parts) > 1 {
				selector.Where(sql.In(key, lo.ToAnySlice(parts)...))
			}
		case OP_NOT_IN:
			parts := ParseInFilterValues(f)
			if len(parts) == 1 {
				selector.Where(sql.NEQ(f.Key, parts[0]))
			} else if len(parts) > 1 {
				selector.Where(sql.NotIn(f.Key, lo.ToAnySlice(parts)...))
			}
		default:
			selector.Where(sql.ContainsFold(key, f.Value))
		}
	}
}

func ApplyProtoFilter(filters []*pagination.Filter, selector *sql.Selector, columns ...string) {
	if len(filters) == 0 {
		return
	}
	if len(columns) == 0 {
		columns = selector.Columns()
	}

	for _, f := range filters {
		key := utils.ToSnakeCase(f.Key)
		if !slices.Contains(columns, key) {
			continue
		}
		switch f.Op {
		case pagination.OP_OP_LT:
			selector.Where(sql.LT(key, f.Value))
		case pagination.OP_OP_LTE:
			selector.Where(sql.LTE(key, f.Value))
		case pagination.OP_OP_GT:
			selector.Where(sql.GT(key, f.Value))
		case pagination.OP_OP_GTE:
			selector.Where(sql.GTE(key, f.Value))
		case pagination.OP_OP_EQ:
			selector.Where(sql.EQ(key, f.Value))
		case pagination.OP_OP_NEQ:
			selector.Where(sql.NEQ(key, f.Value))
		case pagination.OP_OP_LIKE:
			selector.Where(sql.ContainsFold(key, f.Value))
		case pagination.OP_OP_IS_NULL:
			selector.Where(sql.IsNull(key))
		case pagination.OP_OP_IS_NOT_NULL:
			selector.Where(sql.NotNull(key))
		case pagination.OP_OP_IN:
			parts := ParseInFilterValues(f)
			if len(parts) == 1 {
				selector.Where(sql.EQ(key, parts[0]))
			} else if len(parts) > 1 {
				selector.Where(sql.In(key, lo.ToAnySlice(parts)...))
			}
		case pagination.OP_OP_NOT_IN:
			parts := ParseInFilterValues(f)
			if len(parts) == 1 {
				selector.Where(sql.NEQ(f.Key, parts[0]))
			} else if len(parts) > 1 {
				selector.Where(sql.NotIn(f.Key, lo.ToAnySlice(parts)...))
			}
		default:
			selector.Where(sql.ContainsFold(key, f.Value))
		}
	}
}

func ParseInFilterValues(f *pagination.Filter) []string {
	return lo.FilterMap(strings.Split(f.Value, ","), func(part string, _ int) (string, bool) {
		part = strings.TrimSpace(part)
		return part, part != ""
	})
}

func DefaultPagination(p *pagination.PaginationReq) *pagination.PaginationReq {
	if p == nil {
		return &pagination.PaginationReq{
			PageIndex: 1,
			PageSize:  25,
		}
	}
	return p
}
