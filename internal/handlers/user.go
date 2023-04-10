package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) Profile(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "profile.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		helpers.ServerError(w, err)
	}
	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.IsEmail("email")

	userID, err := strconv.Atoi(form.Get("user_id"))
	if err != nil {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}

	if form.Has("password_new") {
		_, _, _, err = m.DB.Authenticate(form.Get("email"), form.Get("password_old"))
		if err != nil {
			form.Errors.Add("password_old", err.Error())
		} else {
			form.ValidPassword("password_new")
			form.IsDuplicate("password_new", "password_confirm", "Passwords don't match")
		}
	}

	file, handler, err := r.FormFile("profile_img")
	if err != nil {
		if !os.IsNotExist(err) {
			helpers.ServerError(w, err)
		}
	}

	fmt.Println(userID)

	if file != nil {
		defer file.Close()
		if forms.ValidExtension(handler.Filename, "png", "jpg", "jpeg") {
			form.Errors.Add("profile_img", "Only .png .jpg .jpeg files allowed")
		}
		if handler.Size > 1024*1024*0.5 {
			form.Errors.Add("profile_img", "Files larger than 500KB are not supported")
		}
		re, err := regexp.Compile(`\.\w+$`)
		if err != nil {
			helpers.ServerError(w, err)
		}
		extension := re.FindString(handler.Filename)

		tempFile, err := os.CreateTemp("./static/images/users", "*"+extension)
		if err != nil {
			form.Errors.Add("profile_img", err.Error())
		}
		defer tempFile.Close()
		imageID := strings.Split(tempFile.Name(), "/")[4]

		fmt.Println(imageID)

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			form.Errors.Add("post_image", err.Error())
		}

		tempFile.Write(fileBytes)
	}
	render.Template(w, r, "profile.page.tmpl", &models.TemplateData{
		Form: form,
	})
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	if helpers.IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Handles logging in the user
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, _, accessLevel, err := m.DB.Authenticate(email, password)

	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := m.DB.GetUser(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "first_name", user.FirstName)
	m.App.Session.Put(r.Context(), "last_name", user.LastName)
	m.App.Session.Put(r.Context(), "email", user.Email)
	m.App.Session.Put(r.Context(), "img", user.ImgID)
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	m.App.Session.Put(r.Context(), "access_level", accessLevel)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
