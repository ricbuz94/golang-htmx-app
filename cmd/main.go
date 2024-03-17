package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Contact struct {
	Name  string
	Email string
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func newContact(name string, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

func (d *Data) delContact(email string) bool {
	for i, contact := range d.Contacts {
		if contact.Email == email {
			d.Contacts = append(d.Contacts[:i], d.Contacts[i+1:]...)
			return true
		}
	}
	return false
}

func (d *Data) emailAlreadyExists(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("riccardo", "riccardo@email.com"),
			newContact("pietro", "pietro@email.com"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

type ErrorResponse struct {
	ErrorStatus  int
	ErrorTitle   string
	ErrorMessage string
}

func main() {
	e := echo.New()
	e.Static("/", "static")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.Logger())

	page := newPage()
	e.Renderer = NewTemplates()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", page)
	})

	e.POST("/contact", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.emailAlreadyExists(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email già utilizzata"

			return c.Render(http.StatusUnprocessableEntity, "form", formData)
		}

		page.Data.Contacts = append(page.Data.Contacts, newContact(name, email))

		c.Response().Header().Set("Hx-Reswap", "beforeend")
		c.Response().Header().Set("Hx-Retarget", "#contacts")

		c.Render(http.StatusOK, "form", newFormData())
		return c.Render(http.StatusOK, "contact", page.Data.Contacts[len(page.Data.Contacts)-1])
	})

	e.DELETE("/contacts", func(c echo.Context) error {
		email := c.QueryParam("email")
		e.StdLogger.Printf("passed email: %v", email)

		if !page.Data.emailAlreadyExists(email) {
			return c.Render(http.StatusNotFound, "error-popup", ErrorResponse{
				ErrorStatus:  http.StatusNotFound,
				ErrorTitle:   "Contatto non trovato",
				ErrorMessage: "Non è stato trovato alcun account da rimuovere con questa email: " + email + ".",
			})
		}

		isDeleted := page.Data.delContact(email)

		if isDeleted {
			c.Response().Header().Set("Hx-Reswap", "delete")
			return c.NoContent(http.StatusOK)
		}
		return c.NoContent(http.StatusInternalServerError)
	})

	e.Logger.Fatal(e.Start(":6969"))
}
