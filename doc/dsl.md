# Hexya DSL Technical Documentation

Hexya is a Go-based framework for building enterprise resource planning (ERP) applications. It offers a Domain-Specific Language (DSL) tailored for defining models, fields, methods, views, actions, and routes in a structured and efficient manner. This documentation provides a detailed overview of Hexya's DSL, including code examples, field types, attributes, keywords, and best practices.

---

## Table of Contents

1. [Introduction](#introduction)
2. [Models](#models)
   - [Model Declaration](#model-declaration)
   - [Field Types](#field-types)
   - [Field Attributes](#field-attributes)
   - [Relationship Fields](#relationship-fields)
3. [Methods](#methods)
   - [Method Declaration](#method-declaration)
   - [Special Methods](#special-methods)
4. [Views](#views)
   - [View Types](#view-types)
   - [View Components](#view-components)
5. [Actions and Menus](#actions-and-menus)
   - [Defining Actions](#defining-actions)
   - [Creating Menus](#creating-menus)
6. [Routes and Controllers](#routes-and-controllers)
   - [Route Declaration](#route-declaration)
   - [Middleware and Groups](#middleware-and-groups)
7. [Search and Filters](#search-and-filters)
   - [Query Builder](#query-builder)
   - [Domain Expressions](#domain-expressions)
8. [Extending Modules](#extending-modules)
9. [Full Reference](#full-reference)
10. [Examples](#examples)
11. [Conclusion](#conclusion)

---

## 1. Introduction

Hexya's DSL leverages Go's strong typing and concurrency features to provide a robust platform for ERP application development. It simplifies the process by allowing developers to define business logic, data models, and UI components in a coherent and type-safe manner.

---

## 2. Models

### Model Declaration

Models represent database tables. They are declared using `models.NewModel("ModelName")`. Each model should be defined within an `init()` function.

```go
func init() {
    models.NewModel("User")
}
```

### Field Types

Hexya provides various field types to define the structure of models:

- `fields.Char`: Character strings.
- `fields.Text`: Multi-line text.
- `fields.Integer`: Integer numbers.
- `fields.Float`: Floating-point numbers.
- `fields.Boolean`: Boolean values.
- `fields.Date`: Date values.
- `fields.DateTime`: Date and time values.
- `fields.Binary`: Binary data.
- `fields.HTML`: HTML content.
- `fields.Selection`: Fields with a fixed set of choices.
- `fields.Many2One`: Many-to-one relationship.
- `fields.One2Many`: One-to-many relationship.
- `fields.Many2Many`: Many-to-many relationship.

### Field Attributes

Each field can have attributes to control its behavior:

- **String**: Display label for the field.
- **Help**: Help text for the field.
- **Required**: If true, the field is mandatory.
- **Default**: Default value for the field.
- **Readonly**: If true, the field is read-only.
- **Unique**: If true, the field value must be unique.
- **Index**: If true, creates an index on the field.
- **Translate**: If true, the field is translatable.
- **Size**: Maximum size for string fields.

**Example:**

```go
h.User().AddFields(map[string]models.FieldDefinition{
    "Name": fields.Char{
        String:   "Full Name",
        Required: true,
        Help:     "The full name of the user",
    },
})
```

### Relationship Fields

#### Many2One

Defines a many-to-one relationship with another model.

**Attributes**:

- **RelationModel**: The related model.
- **OnDelete**: Behavior on deletion ("cascade", "set null", etc.).

**Example:**

```go
h.User().AddFields(map[string]models.FieldDefinition{
    "Company": fields.Many2One{
        RelationModel: h.Company(),
        String:        "Company",
        OnDelete:      "cascade",
    },
})
```

#### One2Many

Defines a one-to-many relationship.

**Attributes**:

- **RelationModel**: The related model.
- **ReverseFK**: The field in the related model that points back.

**Example:**

```go
h.Company().AddFields(map[string]models.FieldDefinition{
    "Employees": fields.One2Many{
        RelationModel: h.User(),
        ReverseFK:     "Company",
        String:        "Employees",
    },
})
```

#### Many2Many

Defines a many-to-many relationship.

**Attributes**:

- **RelationModel**: The related model.
- **M2MLinkModelName**: The name of the link table.
- **M2MOurField**: The field in the link table that refers to this model.
- **M2MTheirField**: The field in the link table that refers to the related model.

**Example:**

```go
h.User().AddFields(map[string]models.FieldDefinition{
    "Roles": fields.Many2Many{
        RelationModel:    h.Role(),
        M2MLinkModelName: "UserRoleRel",
        M2MOurField:      "UserID",
        M2MTheirField:    "RoleID",
        String:           "Roles",
    },
})
```

---

## 3. Methods

### Method Declaration

Methods are used to define business logic and are attached to models. They are defined as functions with specific signatures and registered using `NewMethod`.

**Method Signature**:

```go
func methodName(rs m.ModelSet, args...) returnType {
    // method body
}
```

**Example**:

```go
func user_Activate(rs m.UserSet) {
    rs.SetIsActive(true)
}

func init() {
    h.User().NewMethod("Activate", user_Activate)
}
```

### Special Methods

Hexya supports special methods:

- **Default methods**: Provide default values for fields.
- **Compute methods**: Compute field values dynamically.
- **Onchange methods**: React to field changes in the UI.

#### Default Method Example

```go
func user_DefaultIsActive(rs m.UserSet) bool {
    return true
}

func init() {
    h.User().AddFields(map[string]models.FieldDefinition{
        "IsActive": fields.Boolean{
            DefaultFunc: user_DefaultIsActive,
        },
    })
}
```

#### Compute Method Example

```go
func user_ComputeFullName(rs m.UserSet) string {
    return rs.FirstName() + " " + rs.LastName()
}

func init() {
    h.User().AddFields(map[string]models.FieldDefinition{
        "FullName": fields.Char{
            Compute: user_ComputeFullName,
            String:  "Full Name",
        },
    })
}
```

---

## 4. Views

### View Types

Hexya supports different view types defined in XML:

- **Form View**: For creating and editing single records.
- **Tree (List) View**: Displays multiple records in a list.
- **Kanban View**: Visual representation of records.
- **Graph View**: For analytical purposes.

### View Components

- `<field>`: Represents a model field.
- `<group>`: Groups fields together.
- `<notebook>`: Creates tabs within forms.
- `<page>`: A tab within a notebook.
- `<button>`: Action buttons.

**Example Form View**:

```xml
<record id="view_user_form" model="ir.ui.view">
    <field name="name">user.form</field>
    <field name="model">User</field>
    <field name="arch" type="xml">
        <form string="User">
            <sheet>
                <group>
                    <field name="first_name"/>
                    <field name="last_name"/>
                </group>
                <group>
                    <field name="email"/>
                    <field name="is_active"/>
                </group>
            </sheet>
            <footer>
                <button name="activate" type="object" string="Activate" class="btn-primary"/>
            </footer>
        </form>
    </field>
</record>
```

---

## 5. Actions and Menus

### Defining Actions

Actions link views to the user interface. They are defined as records in XML.

**Example**:

```xml
<record id="action_user_form" model="ir.actions.act_window">
    <field name="name">User Form</field>
    <field name="res_model">User</field>
    <field name="view_mode">form</field>
    <field name="view_id" ref="view_user_form"/>
</record>
```

### Creating Menus

Menus provide navigation paths to actions.

**Example**:

```xml
<menuitem id="menu_root" name="My App"/>
<menuitem id="menu_user" name="Users" parent="menu_root" sequence="10"/>
<menuitem id="menu_user_form" name="Create User" parent="menu_user" action="action_user_form" sequence="10"/>
```

---

## 6. Routes and Controllers

### Route Declaration

Custom routes handle HTTP requests and are defined using the controllers package.

**Example**:

```go
func userCreateHandler(c *controllers.Context) {
    var data map[string]interface{}
    if err := c.BindJSON(&data); err != nil {
        c.JSONResponse(map[string]string{"error":

 "Invalid input"})
        return
    }
    user := h.User().Create(c.Env(), data)
    c.JSONResponse(user)
}

func init() {
    g := controllers.Registry.AddGroup("/api/v1")
    g.AddController("POST", "/users", userCreateHandler)
}
```

### Middleware and Groups

Middleware functions can be applied to route groups for tasks like authentication.

**Example**:

```go
func authMiddleware(c *controllers.Context) {
    token := c.GetHeader("Authorization")
    if !validateToken(token) {
        c.JSONResponse(map[string]string{"error": "Unauthorized"}, 401)
        c.Abort()
    }
}

func init() {
    g := controllers.Registry.AddGroup("/api/v1")
    g.Use(authMiddleware)
    g.AddController("GET", "/users", userListHandler)
}
```

---

## 7. Search and Filters

### Query Builder

Hexya provides a type-safe query builder for constructing database queries.

**Example**:

```go
activeUsers := h.User().Search(env, q.User().IsActive().Equals(true))
```

### Domain Expressions

Domain expressions allow complex query conditions.

**Operators**:

- **Equals(value)**: Equal to value.
- **NotEquals(value)**: Not equal to value.
- **GreaterThan(value)**: Greater than value.
- **LessThan(value)**: Less than value.
- **In(values...)**: Value is in values.
- **NotIn(values...)**: Value is not in values.
- **Like(pattern)**: Value matches pattern.
- **ILike(pattern)**: Case-insensitive match.
- **And(conditions...)**: Logical AND.
- **Or(conditions...)**: Logical OR.
- **Not(condition)**: Logical NOT.

**Example**:

```go
admins := h.User().Search(env, q.User().Groups().Name().Equals("Admin").And(q.User().IsActive().Equals(true)))
```

---

## 8. Extending Modules

Hexya's modular architecture allows extending existing models, fields, and methods without altering the base code.

### Adding Fields

```go
func init() {
    h.User().AddFields(map[string]models.FieldDefinition{
        "MiddleName": fields.Char{String: "Middle Name"},
    })
}
```

### Overriding Methods

```go
func user_GetDisplayName(rs m.UserSet) string {
    return rs.FirstName() + " " + rs.MiddleName() + " " + rs.LastName()
}

func init() {
    h.User().Methods().GetDisplayName().Extend(user_GetDisplayName)
}
```

---

## 9. Full Reference

### Field Types

- `fields.Char`
- `fields.Text`
- `fields.Integer`
- `fields.Float`
- `fields.Boolean`
- `fields.Date`
- `fields.DateTime`
- `fields.Binary`
- `fields.HTML`
- `fields.Selection`
- `fields.Many2One`
- `fields.One2Many`
- `fields.Many2Many`

### Field Attributes

- **String**: string - Display label.
- **Help**: string - Help text.
- **Required**: bool - Mandatory field.
- **Default**: interface{} - Default value.
- **Readonly**: bool - Read-only field.
- **Unique**: bool - Unique value constraint.
- **Index**: bool - Database index.
- **Translate**: bool - Translatable field.
- **Size**: int - Maximum size.

### Model Methods

- `NewMethod(name string, method interface{})`
- `Methods(): Access to method registry.`
- `AddFields(fields map[string]models.FieldDefinition)`

### Query Operators

- `Equals(value)`
- `NotEquals(value)`
- `GreaterThan(value)`
- `LessThan(value)`
- `In(values...)`
- `NotIn(values...)`
- `Like(pattern)`
- `ILike(pattern)`
- `And(conditions...)`
- `Or(conditions...)`
- `Not(condition)`

---

## 10. Examples

### Complete Model Definition

```go
func init() {
    models.NewModel("Article")
    h.Article().AddFields(map[string]models.FieldDefinition{
        "Title": fields.Char{
            String:   "Title",
            Required: true,
        },
        "Content": fields.Text{
            String: "Content",
        },
        "Author": fields.Many2One{
            RelationModel: h.User(),
            String:        "Author",
        },
        "Tags": fields.Many2Many{
            RelationModel:    h.Tag(),
            M2MLinkModelName: "ArticleTagRel",
            M2MOurField:      "ArticleID",
            M2MTheirField:    "TagID",
            String:           "Tags",
        },
        "PublishedAt": fields.DateTime{
            String: "Published At",
        },
    })
}
```

### Custom Search Query

```go
func recentArticles(env models.Environment) m.ArticleSet {
    thirtyDaysAgo := dates.Now().AddDate(0, 0, -30)
    return h.Article().Search(env,
        q.Article().PublishedAt().GreaterThan(thirtyDaysAgo).
        And(q.Article().Author().IsActive().Equals(true)))
}
```

### Controller with Authentication Middleware

```go
func articleListHandler(c *controllers.Context) {
    articles := recentArticles(c.Env())
    result := articles.Collect(func(a m.ArticleSet) interface{} {
        return map[string]interface{}{
            "title":        a.Title(),
            "author":       a.Author().Name(),
            "published_at": a.PublishedAt(),
        }
    })
    c.JSONResponse(result)
}

func authMiddleware(c *controllers.Context) {
    token := c.GetHeader("Authorization")
    if !validateToken(token) {
        c.JSONResponse(map[string]string{"error": "Unauthorized"}, 401)
        c.Abort()
    }
}

func init() {
    g := controllers.Registry.AddGroup("/api")
    g.Use(authMiddleware)
    g.AddController("GET", "/articles", articleListHandler)
}
```

---

## 11. Conclusion

Hexya's DSL provides a powerful and flexible way to develop ERP applications in Go. By leveraging its type-safe constructs and modular design, developers can create scalable and maintainable business applications with ease.

For more information, refer to Hexya's official documentation and GitHub repository.
