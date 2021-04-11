// //go:generate dataloaden AccountLoader string *github.com/orsenkucher/cocopuff/graphql.Account
//go:generate go run github.com/vektah/dataloaden AccountLoader string *github.com/orsenkucher/cocopuff/graphql.Account

package dataloader

import _ "github.com/vektah/dataloaden"
