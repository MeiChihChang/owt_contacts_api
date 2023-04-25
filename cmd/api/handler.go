package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MeiChihChang/owt_contacts_api/internal/models"
	"github.com/go-chi/chi"
	"github.com/go-passwd/validator"
)

// Authenticate  authenticate a user with email & password
// @Summary      Authentication
// @Description  authenticate a user with email & password
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  RefreshToken
// @Router       /authenticate [post]
func (app *application) Authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.DB.GetContactByEmail(requestPayload.Email)
	log.Println(err)
	log.Println(user)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:    user.ID,
		Email: user.Email,
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

// InsertContact create a new contact with first_name, last_name, full_name, email, password, address, mobile, token
// @Summary      InsertContact
// @Description  create a new contact with first_name, last_name, full_name, email, password, address, mobile, token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  contact_id
// @Router       /contact/new [put]
func (app *application) InsertContact(w http.ResponseWriter, r *http.Request) {
	var contact models.Contact
	err := app.readJSON(w, r, &contact)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	log.Printf("contact : %v", contact)

	// do checkings
	message := app.checkContact(contact.FirstName, contact.LastName, contact.FullName, contact.Address, contact.Mobile, contact.Email, contact.Password)
	if len(message) > 0 {
		log.Printf("err %v", message)
		resp := JSONResponse{
			Error:   true,
			Message: message,
		}

		app.writeJSON(w, http.StatusNotAcceptable, resp)
		return
	}

	contact.CreatedAt = time.Now()
	contact.UpdatedAt = time.Now()

	newID, err := app.DB.InsertContact(contact)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	words := strings.Fields(contact.SkillIDstring)

	for _, w := range words {
		intVar, _ := strconv.Atoi(w)
		contact.SkillIDs = append(contact.SkillIDs, intVar)
	}

	log.Printf("id %v skills %v", newID, contact.SkillIDs)

	err = app.DB.UpdateContactSkills(newID, contact.SkillIDs)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("%v", newID),
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// AllContacts   list all contacts with token
// @Summary      AllContacts
// @Description  list all contacts with token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  [] contact
// @Router       /contact/all [get]
func (app *application) AllContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := app.DB.AllContacts()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, contacts)
}

// GetContact    get a contact with id, token
// @Summary      GetContact
// @Description  get a contact with id, token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  contact
// @Router       /contact/{id} [get]
func (app *application) GetContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	contactID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	contact, err := app.DB.GetContactByID(contactID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, contact)
}

// UpdateContact update a contact by id with first_name, last_name, full_name, email, password, address, mobile & []skillids, token
// @Summary      UpdateContact
// @Description  update a contact by id with content and skills, token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  Success
// @Router       /contact/update/{id} [patch]
func (app *application) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
	if claims["id"] != id {
		resp := JSONResponse{
			Error:   true,
			Message: "User cant change others contact",
		}

		app.writeJSON(w, http.StatusNotAcceptable, resp)
		return
	}

	var payload models.Contact

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println(payload)

	contact, err := app.DB.GetContactByID(id)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	contact.ID = id
	contact.FirstName = payload.FirstName
	contact.LastName = payload.LastName
	contact.FullName = payload.FullName
	contact.Address = payload.Address
	contact.Mobile = payload.Mobile
	contact.Email = payload.Email
	contact.UpdatedAt = time.Now()

	log.Printf("contact : %v", contact)

	// do checkings
	message := app.checkContact(contact.FirstName, contact.LastName, contact.FullName, contact.Address, contact.Mobile, contact.Email, contact.Password)
	if len(message) > 0 {
		log.Printf("err %v", message)
		resp := JSONResponse{
			Error:   true,
			Message: message,
		}

		app.writeJSON(w, http.StatusNotAcceptable, resp)
		return
	}

	err = app.DB.UpdateContact(*contact)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	words := strings.Fields(payload.SkillIDstring)

	for _, w := range words {
		intVar, _ := strconv.Atoi(w)
		payload.SkillIDs = append(payload.SkillIDs, intVar)
	}

	err = app.DB.UpdateContactSkills(contact.ID, payload.SkillIDs)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "contact updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// DeleteContact delete a contact with id, token
// @Summary      DeleteContact
// @Description  delete a contact with id, token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  Success
// @Router       /contact/{id} [get]
func (app *application) DeleteContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteContact(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "contact deleted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// InsertSkill   create a new skill with name, level, token
// @Summary      InsertSkill
// @Description  create a new skill with name, level, token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  skill_id
// @Router       /skill/new [put]
func (app *application) InsertSkill(w http.ResponseWriter, r *http.Request) {
	var skill models.Skill

	err := app.readJSON(w, r, &skill)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println(skill)

	// do checkings
	message := app.checkSkillNameLevel(skill.Name, skill.Level)
	if len(message) > 0 {
		resp := JSONResponse{
			Error:   true,
			Message: message,
		}

		app.writeJSON(w, http.StatusNotAcceptable, resp)
		return
	}

	skill.CreatedAt = time.Now()
	skill.UpdatedAt = time.Now()

	newID, err := app.DB.InsertSkill(skill)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("%v", newID),
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// AllSkills     list all skills with token
// @Summary      AllSkills
// @Description  list all skills with token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  [] skills
// @Router       /skill/all [get]
func (app *application) AllSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := app.DB.AllSkills()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, skills)
}

// GetSkill      get a skill with id, token
// @Summary      GetSkill
// @Description  get a skill with id, token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  skill
// @Router       /skill/{id} [get]
func (app *application) GetSkill(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	skillID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	skill, err := app.DB.GetSkillByID(skillID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, skill)
}

// UpdateSkill   update a skill by id with name, level, token
// @Summary      UpdateSkill
// @Description  update a skill by id with with name, level, token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  Success
// @Router       /skill/update/{id} [patch]
func (app *application) UpdateSkill(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	skillID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload models.Skill

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	message := app.checkSkillNameLevel(payload.Name, payload.Level)
	if len(message) > 0 {
		resp := JSONResponse{
			Error:   true,
			Message: message,
		}

		app.writeJSON(w, http.StatusNotAcceptable, resp)
		return
	}

	skill, err := app.DB.GetSkillByID(skillID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	skill.Name = payload.Name
	skill.Level = payload.Level
	skill.UpdatedAt = time.Now()

	err = app.DB.UpdateSkill(*skill)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "skill updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// DeleteSkill   delete a skill with id, token
// @Summary      DeleteSkill
// @Description  delete a skill with id, token
// @Accept       json
// @Produce      json
// @Tags         Tools
// @Security     JWTToken
// @Success      200  {object}  Success
// @Router       /skill/{id} [delete]
func (app *application) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteSkill(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "skill deleted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) checkSkillNameLevel(name string, level int) string {
	if len(name) > models.MAX_NAME_LENGHT {
		return "name length is too big"
	}

	if level < 0 || level >= len(models.LEVELS) {
		return "skill level is not at corret range"
	}

	return ""
}

func (app *application) checkContact(first_name string, last_name string, full_name string, address string, mobile string, email string, password string) string {
	// check name
	if len(first_name) > models.MAX_NAME_LENGHT || len(last_name) > models.MAX_NAME_LENGHT || len(full_name) > 2*models.MAX_NAME_LENGHT {
		return "name length is too big"
	}

	if !(strings.Contains(full_name, first_name) && strings.Contains(full_name, last_name)) {
		return "full name should contain first name & last name"
	}

	// check address
	if len(address) <= 0 || len(address) >= 255 {
		return "address is too short or long"
	}

	// check phone number
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	if !re.MatchString(mobile) {
		return "mobile number is incorrect"
	}

	// check email
	_, err := mail.ParseAddress(email)
	if err != nil {
		return "email is incorrect"
	}

	// check password
	if len(password) > 0 {
		passwordValidator := validator.New(validator.MinLength(5, errors.New("too short")), validator.MaxLength(12, errors.New("too long")))
		err = passwordValidator.Validate(password)
		if err != nil {
			return fmt.Sprintf("err : %v", err)
		}
	}

	return ""
}
