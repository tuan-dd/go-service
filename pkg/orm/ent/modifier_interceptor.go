package entOrm

import (
	"context"

	"entgo.io/ent"

	"github.com/tuan-dd/go-pkg/common/request"
	"github.com/tuan-dd/go-pkg/common/utils"
)

func (d ModifierMixin) Hook() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
			client := next

			if skip, ok := ctx.Value(ModifierKey{}).(bool); ok && skip {
				return client.Mutate(ctx, mutation)
			}

			userID := request.GetUserCtx[any](ctx).ID
			switch mutation.Op() {
			case ent.OpCreate:
				utils.CallMethod("SetOp", mutation, ent.OpCreate)
				utils.CallMethod("SetCreatedBy", mutation, userID)
				client = utils.CallMethodWithValue[mutateClient]("Client", mutation)
			case ent.OpUpdateOne, ent.OpUpdate:
				utils.CallMethod("SetOp", mutation, ent.OpUpdate)
				utils.CallMethod("SetUpdatedBy", mutation, userID)
				client = utils.CallMethodWithValue[mutateClient]("Client", mutation)
			}

			return client.Mutate(ctx, mutation)
		})
	}
}
