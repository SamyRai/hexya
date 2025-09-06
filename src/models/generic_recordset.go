package models

import "encoding/json"

// GenericRecordSet is a generic type for a collection of records.
type GenericRecordSet[T any] struct {
	*RecordCollection
}

// First returns a pointer to the first record's data.
// It returns nil if the GenericRecordSet is empty.
func (rs GenericRecordSet[T]) First() *T {
	if rs.IsEmpty() {
		return nil
	}
	modelData := rs.RecordCollection.First()
	jsonData, _ := json.Marshal(modelData)
	res := new(T)
	json.Unmarshal(jsonData, res)
	return res
}

// All returns a slice of all records' data.
func (rs GenericRecordSet[T]) All() []T {
	allData := rs.RecordCollection.All()
	res := make([]T, len(allData))
	for i, data := range allData {
		jsonData, _ := json.Marshal(data)
		json.Unmarshal(jsonData, &res[i])
	}
	return res
}

// Records returns a slice of singleton GenericRecordSets.
func (rs GenericRecordSet[T]) Records() []GenericRecordSet[T] {
	recs := rs.RecordCollection.Records()
	res := make([]GenericRecordSet[T], len(recs))
	for i, rec := range recs {
		res[i] = GenericRecordSet[T]{rec.Collection()}
	}
	return res
}

// Sorted returns a new GenericRecordSet sorted by the given less function.
func (rs GenericRecordSet[T]) Sorted(less func(a, b T) bool) GenericRecordSet[T] {
	res := rs.RecordCollection.Sorted(func(rc1, rc2 RecordSet) bool {
		d1 := rc1.Collection().First()
		d2 := rc2.Collection().First()
		jsonData1, _ := json.Marshal(d1)
		jsonData2, _ := json.Marshal(d2)
		t1 := new(T)
		t2 := new(T)
		json.Unmarshal(jsonData1, t1)
		json.Unmarshal(jsonData2, t2)
		return less(*t1, *t2)
	})
	return GenericRecordSet[T]{res.Collection()}
}

// Filtered returns a new GenericRecordSet with only the records for which test is true.
func (rs GenericRecordSet[T]) Filtered(test func(T) bool) GenericRecordSet[T] {
	res := rs.RecordCollection.Filtered(func(rc RecordSet) bool {
		d := rc.Collection().First()
		jsonData, _ := json.Marshal(d)
		t := new(T)
		json.Unmarshal(jsonData, t)
		return test(*t)
	})
	return GenericRecordSet[T]{res.Collection()}
}

// GenericRecordData is a generic type for a single record's data.
type GenericRecordData[T any] struct {
	*ModelData
}
