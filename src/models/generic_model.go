package models

// GenericModel is a generic type for a model.
type GenericModel[T any] struct {
	*Model
}

// NewSet returns a new GenericRecordSet[T] instance in the given Environment
func (m GenericModel[T]) NewSet(env Environment) GenericRecordSet[T] {
	return GenericRecordSet[T]{
		RecordCollection: env.Pool(m.name),
	}
}

// Create creates a new record and returns the newly created GenericRecordSet[T] instance.
func (m GenericModel[T]) Create(env Environment, data GenericRecordData[T]) GenericRecordSet[T] {
	return GenericRecordSet[T]{
		RecordCollection: m.Model.Create(env, data.Underlying()),
	}
}

// Search searches the database and returns a new GenericRecordSet[T] instance
// with the records found.
func (m GenericModel[T]) Search(env Environment, cond Conditioner) GenericRecordSet[T] {
	return GenericRecordSet[T]{
		RecordCollection: m.Model.Search(env, cond),
	}
}

// Browse returns a new GenericRecordSet with the records with the given ids.
// Note that this function is just a shorcut for Search on a list of ids.
func (m GenericModel[T]) Browse(env Environment, ids []int64) GenericRecordSet[T] {
	return GenericRecordSet[T]{
		RecordCollection: m.Model.Browse(env, ids),
	}
}

// BrowseOne returns a new GenericRecordSet with the record with the given id.
// Note that this function is just a shorcut for Search on the given id.
func (m GenericModel[T]) BrowseOne(env Environment, id int64) GenericRecordSet[T] {
	return GenericRecordSet[T]{
		RecordCollection: m.Model.BrowseOne(env, id),
	}
}

// NewData returns a pointer to a new empty GenericRecordData[T] instance.
//
// Optional field maps if given will be used to populate the data.
func (m GenericModel[T]) NewData(fm ...FieldMap) GenericRecordData[T] {
	return GenericRecordData[T]{
		ModelData: NewModelData(m.Model, fm...),
	}
}
