package models

// GenericCondition is a generic type for a query condition.
type GenericCondition[T any] struct {
	*Condition
}

// And completes the current condition with a simple AND clause
func (c GenericCondition[T]) And() GenericConditionStart[T] {
	return GenericConditionStart[T]{
		ConditionStart: c.Condition.And(),
	}
}

// Or completes the current condition with a simple OR clause
func (c GenericCondition[T]) Or() GenericConditionStart[T] {
	return GenericConditionStart[T]{
		ConditionStart: c.Condition.Or(),
	}
}

// AndNot completes the current condition with a simple AND NOT clause
func (c GenericCondition[T]) AndNot() GenericConditionStart[T] {
	return GenericConditionStart[T]{
		ConditionStart: c.Condition.AndNot(),
	}
}

// OrNot completes the current condition with a simple OR NOT clause
func (c GenericCondition[T]) OrNot() GenericConditionStart[T] {
	return GenericConditionStart[T]{
		ConditionStart: c.Condition.OrNot(),
	}
}

// AndCond completes the current condition with the given cond as an AND clause
func (c GenericCondition[T]) AndCond(cond GenericCondition[T]) GenericCondition[T] {
	return GenericCondition[T]{
		Condition: c.Condition.AndCond(cond.Condition),
	}
}

// OrCond completes the current condition with the given cond as an OR clause
func (c GenericCondition[T]) OrCond(cond GenericCondition[T]) GenericCondition[T] {
	return GenericCondition[T]{
		Condition: c.Condition.OrCond(cond.Condition),
	}
}

// GenericConditionStart is a generic type for a condition start.
type GenericConditionStart[T any] struct {
	*ConditionStart
}

// Field adds a field to the condition
func (cs GenericConditionStart[T]) Field(name FieldName) *ConditionField {
	return cs.ConditionStart.Field(name)
}
