package sample

import (
	"context"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/account"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/transaction"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/strings"
	"gopkg.in/yaml.v3"
	stdstrings "strings"
)

type accountWithChildren struct {
	ID                   string
	DisplayName          string                `yaml:"displayName"`
	Icon                 string                `yaml:"icon"`
	Children             []accountWithChildren `yaml:"children"`
	OutgoingTransactions []transactionSpec     `yaml:"outgoingTransactions"`
}

func (a *accountWithChildren) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		ID                   string
		DisplayName          string                `yaml:"displayName"`
		Icon                 string                `yaml:"icon"`
		Children             []accountWithChildren `yaml:"children"`
		OutgoingTransactions []transactionSpec     `yaml:"outgoingTransactions"`
	}
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("failed decoding YAML account: %w", err)
	}
	if raw.DisplayName == "" {
		raw.DisplayName = strings.CamelCaseToHumanReadable(raw.ID)
	}
	if raw.Icon == "" {
		raw.Icon = stdstrings.ToUpper(string(raw.ID[0])) + raw.ID[1:]
	}
	a.ID = raw.ID
	a.DisplayName = raw.DisplayName
	a.Icon = raw.Icon
	a.Children = raw.Children
	a.OutgoingTransactions = raw.OutgoingTransactions
	return nil
}

func (a *accountWithChildren) applyAccount(ctx context.Context, ids map[string]string, h account.Handler, tenantID string, parent *accountWithChildren) error {
	var parentID *string
	if parent != nil {
		parentID = lang.PtrOf(ids[parent.ID])
	}

	if acc, err := h.Create(ctx, account.CreateRequest{
		TenantID:    tenantID,
		DisplayName: a.DisplayName,
		Icon:        strings.NilIfEmpty(a.Icon),
		ParentID:    parentID,
	}); err != nil {
		return fmt.Errorf("failed creating account '%s': %w", a.ID, err)
	} else {
		ids[a.ID] = acc.ID
	}

	for _, child := range a.Children {
		if err := child.applyAccount(ctx, ids, h, tenantID, a); err != nil {
			return fmt.Errorf("failed applying child account '%s' in parent '%s': %w", child.ID, a.ID, err)
		}
	}

	return nil
}

func (a *accountWithChildren) applyOutgoingTransactions(ctx context.Context, ids map[string]string, h transaction.Handler, tenantID, defaultCurrency string) error {
	for _, tx := range a.OutgoingTransactions {
		if err := tx.apply(ctx, ids, h, tenantID, defaultCurrency, a.ID); err != nil {
			return fmt.Errorf("failed applying an outgoing transaction of account '%s': %w", a.ID, err)
		}
	}

	for _, child := range a.Children {
		if err := child.applyOutgoingTransactions(ctx, ids, h, tenantID, defaultCurrency); err != nil {
			return fmt.Errorf("failed applying outgoing transactions of child account '%s' in parent '%s': %w", child.ID, a.ID, err)
		}
	}

	return nil
}
