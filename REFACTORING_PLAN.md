# Refactoring and Modernization Plan

This document outlines the plan to refactor the Hexya ERP framework from its current state to a more modern, maintainable, and efficient architecture using modern Go features, primarily generics.

## 1. Introduce a Generic ORM Core

The current ORM relies on a code generator to create type-specific wrappers (`UserSet`, `PostSet`, etc.) around a generic `interface{}`-based core. This is a pre-generics pattern that can be replaced entirely.

-   **Action:** Define a new generic `RecordSet[T any]` type in the `models` package. `T` will represent the model's data struct (e.g., `models.User`). This will be the foundation for replacing all generated `...Set` types.
-   **Action:** Create a generic, fluent query builder, `Condition[T any]`, to replace the existing query DSL. This will allow for type-safe queries directly on model fields (e.g., `db.Query[User]().Where(u => u.Name).Equals("John")`).
-   **Action:** Adapt the model registry to handle these new generic types.

## 2. Refactor a Pilot Model

To prove the new design and create a template for migration, we will start with a single, simple model.

-   **Action:** Select a pilot model (e.g., `Post` or `Tag`).
-   **Action:** Manually refactor it to use the new generic ORM. This involves changing its definition to work with `RecordSet[T]` and updating all its usages and tests to use the new generic query builder and recordset methods.

## 3. Deprecate and Remove the Code Generator

With the new generic ORM proven, the code generator becomes obsolete.

-   **Action:** Remove the `hexya generate` command and all associated code from `cmd/generate.go` and `src/tools/generate/`.
-   **Action:** Update CI/build scripts (like `.travis.yml`) to remove the generation step. This simplifies the development workflow significantly.

## 4. Incrementally Refactor All Remaining Models

With a clear path forward, the remaining models can be migrated.

-   **Action:** Methodically go through the remaining models and refactor them one by one to use the new generic ORM, following the pattern established with the pilot model.

## 5. Update Dependencies

Once the core refactoring is complete and the codebase is stable, we can safely update external dependencies.

-   **Action:** Update the dependencies in `go.mod` to their latest stable versions to bring in security patches, performance improvements, and new features.

## 6. Comprehensive Testing

Testing is critical throughout this process to ensure no functionality is broken.

-   **Action:** At each stage, run the entire existing test suite using `./run_tests.sh`.
-   **Action:** Fix any tests that fail due to the refactoring.
-   **Action:** Add new tests specifically for the new generic ORM components to ensure they are robust and correct.
