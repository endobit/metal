package data

import (
	"context"
	"database/sql"
	"errors"

	"endobit.io/metal/internal/data/db"
)

// CreateModel creates a new model in the database. It returns true if
// the model already exists.
func (s Store) CreateModel(ctx context.Context, name string, scope MakeScope) (bool, error) {
	if _, err := s.ReadModels(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if _, err := s.CreateMake(ctx, scope.Make); err != nil { // create make if it doesn't exist
				return false, err
			}

			return false, s.db.CreateModel(ctx, db.CreateModelParams{Make: scope.Make, Name: name})
		}

		return false, err
	}

	return true, nil
}

// CreateModelAttr creates a new model attribute in the database. If
// returns true if the attribute already exists.
func (s Store) CreateModelAttr(ctx context.Context, name string, scope ModelScope) (bool, error) {
	if _, err := s.ReadModelAttrs(ctx, FilterByName(name), scope); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, s.db.CreateModelAttribute(ctx, db.CreateModelAttributeParams{
				Model: scope.Model,
				Make:  scope.Make,
				Name:  name,
			})
		}

		return false, err
	}

	return true, nil
}

func newModel(src db.ReadModelRow) Model {
	return Model{
		MakeScope:    MakeScope{Make: src.Make},
		id:           src.ID,
		Name:         src.Name,
		Architecture: value(src.Architecture),
	}
}

func (s Store) ReadModel(ctx context.Context, model string, scope MakeScope) (*Model, error) {
	models, err := s.ReadModels(ctx, FilterByName(model), scope)
	if err != nil {
		return nil, err
	}

	switch len(models) {
	case 0:
		return nil, nil
	case 1:
		return &models[0], nil
	default:
		return nil, errors.New("multiple models found")
	}
}

// ReadModels reads models from the database.
func (s Store) ReadModels(ctx context.Context, filter Filter, scope MakeScope) ([]Model, error) {
	type from = db.ReadModelRow

	var (
		models []Model
		err    error
	)

	switch {
	case scope.Make != "" && filter.Type == NameFilter:
		var row from

		row, err = s.db.ReadModel(ctx, db.ReadModelParams{
			Make: scope.Make,
			Name: filter.Pattern,
		})
		models = append(models, newModel(row))

	case scope.Make != "" && filter.Type == GlobFilter:
		var rows []db.ReadModelsByGlobRow

		rows, err = s.db.ReadModelsByGlob(ctx, db.ReadModelsByGlobParams{
			Make: scope.Make,
			Glob: filter.Pattern,
		})
		for i := range rows {
			models = append(models, newModel(from(rows[i])))
		}

	case scope.Make != "":
		var rows []db.ReadModelsByMakeRow

		rows, err = s.db.ReadModelsByMake(ctx, db.ReadModelsByMakeParams{
			Make: scope.Make,
		})
		for i := range rows {
			models = append(models, newModel(from(rows[i])))
		}

	default:
		var rows []db.ReadModelsRow

		rows, err = s.db.ReadModels(ctx, db.ReadModelsParams{})
		for i := range rows {
			models = append(models, newModel(from(rows[i])))
		}
	}

	if err != nil {
		return nil, err
	}

	return models, nil
}

func newModelAttr(src db.ReadModelAttributeRow) ModelAttr {
	return ModelAttr{
		ModelScope: ModelScope{
			Make:  src.Make,
			Model: src.Model,
		},
		Attr: Attr{
			id:          src.ID,
			Name:        src.Name,
			Value:       value(src.Value),
			IsProtected: src.IsProtected == 1,
		},
	}
}

// ReadModelAttrs reads model attributes from the database.
func (s Store) ReadModelAttrs(ctx context.Context, filter Filter, scope ModelScope) ([]ModelAttr, error) {
	type from = db.ReadModelAttributeRow

	var (
		attrs []ModelAttr
		err   error
	)

	switch {
	case scope.Make != "" && scope.Model != "" && filter.Type == NameFilter:
		var row db.ReadModelAttributeRow

		row, err = s.db.ReadModelAttribute(ctx, db.ReadModelAttributeParams{
			Make:  scope.Make,
			Model: scope.Model,
			Name:  filter.Pattern,
		})
		attrs = append(attrs, newModelAttr(row))

	case scope.Make != "" && scope.Model != "" && filter.Type == GlobFilter:
		var rows []db.ReadModelAttributesByGlobRow

		rows, err = s.db.ReadModelAttributesByGlob(ctx, db.ReadModelAttributesByGlobParams{
			Make:  scope.Make,
			Model: scope.Model,
			Glob:  filter.Pattern,
		})
		for i := range rows {
			attrs = append(attrs, newModelAttr(from(rows[i])))
		}

	case scope.Make != "" && scope.Model != "":
		var rows []db.ReadModelAttributesByMakeModelRow

		rows, err = s.db.ReadModelAttributesByMakeModel(ctx, db.ReadModelAttributesByMakeModelParams{
			Make:  scope.Make,
			Model: scope.Model,
		})
		for i := range rows {
			attrs = append(attrs, newModelAttr(from(rows[i])))
		}

	case scope.Make != "":
		var rows []db.ReadModelAttributesByMakeRow

		rows, err = s.db.ReadModelAttributesByMake(ctx, db.ReadModelAttributesByMakeParams{
			Make: scope.Make,
		})
		for i := range rows {
			attrs = append(attrs, newModelAttr(from(rows[i])))
		}

	default:
		var rows []db.ReadModelAttributesRow

		rows, err = s.db.ReadModelAttributes(ctx, db.ReadModelAttributesParams{})
		for i := range rows {
			attrs = append(attrs, newModelAttr(from(rows[i])))
		}
	}

	return attrs, err
}

// UpdateModelName updates the name of a model in the database.
func (s Store) UpdateModelName(ctx context.Context, name string, scope ModelScope) error {
	return s.db.UpdateModelName(ctx, db.UpdateModelNameParams{
		Make:  scope.Make,
		Model: scope.Model,
		Name:  name,
	})
}

// UpdateModelArchitecture updates the architecture of a model in the database.
func (s Store) UpdateModelArchitecture(ctx context.Context, arch string, scope ModelScope) error {
	return s.db.UpdateModelArchitecture(ctx, db.UpdateModelArchitectureParams{
		Make:         scope.Make,
		Model:        scope.Model,
		Architecture: Optional(arch),
	})
}

// UpdateModelAttrName updates the name of a model attribute in the database.
func (s Store) UpdateModelAttrName(ctx context.Context, name string, scope ModelAttrScope) error {
	return s.db.UpdateModelAttributeName(ctx, db.UpdateModelAttributeNameParams{
		Make:  scope.Make,
		Model: scope.Model,
		Attr:  scope.Attr,
		Name:  name,
	})
}

// UpdateModelAttrValue updates the value of a model attribute in the database.
func (s Store) UpdateModelAttrValue(ctx context.Context, value string, scope ModelAttrScope) error {
	return s.db.UpdateModelAttributeValue(ctx, db.UpdateModelAttributeValueParams{
		Make:  scope.Make,
		Model: scope.Model,
		Attr:  scope.Attr,
		Value: Optional(value),
	})
}

// UpdateModelAttrProtection updates the protection of a model attribute in the database.
func (s Store) UpdateModelAttrProtection(ctx context.Context, protected bool, scope ModelAttrScope) error {
	return s.db.UpdateModelAttributeProtection(ctx, db.UpdateModelAttributeProtectionParams{
		Make:        scope.Make,
		Model:       scope.Model,
		Attr:        scope.Attr,
		IsProtected: boolean(protected),
	})
}

// DeleteModels deletes models from the database.
func (s Store) DeleteModels(ctx context.Context, filter Filter, scope MakeScope) error {
	rows, err := s.ReadModels(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteModel(ctx, db.DeleteModelParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteModelAttrs deletes models attributes from the database.
func (s Store) DeleteModelAttrs(ctx context.Context, filter Filter, scope ModelScope) error {
	rows, err := s.ReadModelAttrs(ctx, filter, scope)
	if err != nil {
		return err
	}

	for i := range rows {
		err := s.db.DeleteAttribute(ctx, db.DeleteAttributeParams{ID: rows[i].id})
		if err != nil {
			return err
		}
	}

	return nil
}
